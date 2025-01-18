package services

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/integrations"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type application int

const (
	recipeKeeper application = iota
	recipya
)

func (f *Files) processCrouton(file *multipart.FileHeader) *models.Recipe {
	fi, err := file.Open()
	if err != nil {
		slog.Error("Failed to open file", "error", err, "file", file)
		return nil
	}
	defer fi.Close()

	r := models.NewRecipeFromCrouton(fi, f.UploadImage)
	return &r
}

func (f *Files) processJSON(file *multipart.FileHeader) models.Recipes {
	fi, err := file.Open()
	if err != nil {
		slog.Error("Failed to open file", "error", err, "file", file)
		return nil
	}
	defer fi.Close()

	xr, err := f.extractJSONRecipes(fi)
	if err != nil {
		slog.Error("Could not extract file", "file", file, "error", err)
		return nil
	}
	return xr
}

func processMasterCook(file *multipart.FileHeader) models.Recipes {
	f, err := file.Open()
	if err != nil {
		slog.Error("Failed to open file", "file", file, "error", err)
		return nil
	}
	defer f.Close()

	return models.NewRecipesFromMasterCook(f)
}

func (f *Files) processPaprikaRecipes(rc io.ReadCloser, file *multipart.FileHeader) models.Recipes {
	if rc == nil {
		openedFile, err := file.Open()
		if err != nil {
			slog.Error("Failed to open file", "file", file, "error", err)
			return nil
		}
		defer openedFile.Close()
		rc = openedFile
	}

	z, err := unzipMem(rc)
	if err != nil {
		return nil
	}

	recipes := make(models.Recipes, 0, len(z.File))
	for _, zf := range z.File {
		openFile, err := zf.Open()
		if err != nil {
			slog.Error("Could not open paprika recipe", "file", zf.Name, "error", err)
			continue
		}

		data, err := unzipGzip(openFile)
		if err != nil {
			openFile.Close()
			slog.Error("Could not unzip paprika recipe", "file", zf.Name, "error", err)
			return nil
		}
		openFile.Close()

		var p models.PaprikaRecipe
		err = json.Unmarshal(data, &p)
		if err != nil {
			slog.Error("Could not unmarshal paprika recipe", "file", zf.Name, "data", string(data), "error", err)
			continue
		}

		img := uuid.Nil
		if p.PhotoData != "" {
			decode, err := base64.StdEncoding.DecodeString(p.PhotoData)
			if err == nil {
				img, err = f.UploadImage(io.NopCloser(bytes.NewReader(decode)))
				if err != nil {
					slog.Error("Failed to upload Paprika image", "file", zf.Name, "error", err)
				}
			}
		} else if p.ImageURL != "" {
			img, err = f.ScrapeAndStoreImage(p.ImageURL)
			if err != nil {
				slog.Error("Failed to fetch and upload Paprika image", "file", zf.Name, "error", err)
			}
		}

		recipes = append(recipes, p.Recipe(img))
	}

	return recipes
}

func (f *Files) processTxt(file *multipart.FileHeader) models.Recipes {
	openFile, err := file.Open()
	if err != nil {
		slog.Error("Failed to open file", "file", file, "error", err)
		return make(models.Recipes, 0)
	}
	defer openFile.Close()

	data, err := io.ReadAll(openFile)
	if err != nil {
		slog.Error("Failed to read file", "file", file.Filename, "error", err)
		return make(models.Recipes, 0)
	}

	recipe, err := models.NewRecipeFromTextFile(bytes.NewBuffer(data))
	if errors.Is(err, models.ErrIsAccuChef) {
		return models.NewRecipesFromAccuChef(bytes.NewBuffer(data))
	} else if errors.Is(err, models.ErrIsEasyRecipeDeluxe) {
		return models.NewRecipesFromEasyRecipeDeluxe(bytes.NewBuffer(data))
	}

	if err != nil {
		slog.Error("Could not create recipe from text file", "file", file.Filename, "error", err)
		return nil
	}
	return models.Recipes{recipe}
}

func (f *Files) processZip(file *multipart.FileHeader) models.Recipes {
	openFile, err := file.Open()
	if err != nil {
		slog.Error("Failed to open file", "file", file, "error", err)
		return make(models.Recipes, 0)
	}
	defer openFile.Close()

	z, err := unzipMem(openFile)
	if err != nil {
		return make(models.Recipes, 0)
	}

	switch detectApp(z) {
	case recipeKeeper:
		return f.processRecipeKeeper(z)
	default:
		return f.processRecipeFiles(z)
	}
}

func detectApp(zr *zip.Reader) application {
	files := make([]string, 0, len(zr.File))
	for _, file := range zr.File {
		files = append(files, strings.Split(file.Name, "/")[0])
	}

	files = extensions.Unique(files)
	slices.Sort(files)

	switch len(files) {
	case 2:
		if files[0] == "images" && files[1] == "recipes.html" {
			return recipeKeeper
		}
	}
	return recipya
}

func (f *Files) processRecipeFiles(zr *zip.Reader) models.Recipes {
	var (
		imageUUID    uuid.UUID
		recipeNumber int
		recipes      = make(models.Recipes, 0, len(zr.File))
	)

	for _, zf := range zr.File {
		if strings.Contains(zf.Name, "__MACOSX") {
			continue
		}

		if imageUUID != uuid.Nil && (zf.FileInfo().IsDir() || (recipeNumber > 0 && len(recipes[recipeNumber-1].Images) == 0)) {
			recipes[recipeNumber-1].Images = append(recipes[recipeNumber-1].Images, imageUUID)
			imageUUID = uuid.Nil
		}

		validImageFormats := []string{".jpg", ".jpeg", ".png"}
		if imageUUID == uuid.Nil && slices.Contains(validImageFormats, filepath.Ext(zf.Name)) {
			imageFile, err := zf.Open()
			if err != nil {
				slog.Error("Failed to open image file", "file", zf, "error", err)
				continue
			}

			if zf.FileInfo().Size() < 1<<12 {
				_ = imageFile.Close()
				continue
			}

			imageUUID, err = f.UploadImage(imageFile)
			if err != nil {
				slog.Error("Failed to upload image", "file", zf, "error", err)
			}

			_ = imageFile.Close()
			continue
		}

		openedFile, err := zf.Open()
		if err != nil {
			slog.Error("Failed to open file", "file", zf, "error", err)
			continue
		}

		switch strings.ToLower(filepath.Ext(zf.Name)) {
		case models.CML.Ext():
			xr := models.NewRecipesFromCML(openedFile, nil, f.UploadImage)
			if len(xr) > 0 {
				recipes = append(recipes, xr...)
				recipeNumber += len(xr)
			}
		case models.Crumb.Ext():
			recipes = append(recipes, models.NewRecipeFromCrouton(openedFile, f.UploadImage))
			recipeNumber++
		case models.JSON.Ext():
			xr, err := f.extractJSONRecipes(openedFile)
			if err != nil {
				_ = openedFile.Close()
				slog.Error("Failed to extract", "file", zf, "error", err)
				continue
			}

			recipes = append(recipes, xr...)
			recipeNumber += len(xr)
		case models.MXP.Ext():
			xr := models.NewRecipesFromMasterCook(openedFile)
			if len(xr) > 0 {
				recipes = append(recipes, xr...)
				recipeNumber += len(xr)
			}
		case models.Paprika.Ext():
			xr := f.processPaprikaRecipes(openedFile, nil)
			if len(xr) > 0 {
				recipes = append(recipes, xr...)
				recipeNumber += len(xr)
			}
		case models.TXT.Ext():
			recipe, err := models.NewRecipeFromTextFile(openedFile)
			if errors.Is(err, models.ErrIsAccuChef) {
				_ = openedFile.Close()
				xr := models.NewRecipesFromAccuChef(openedFile)
				recipes = append(recipes, xr...)
				recipeNumber += len(xr)
				continue
			} else if errors.Is(err, models.ErrIsEasyRecipeDeluxe) {
				_ = openedFile.Close()
				xr := models.NewRecipesFromEasyRecipeDeluxe(openedFile)
				recipes = append(recipes, xr...)
				recipeNumber += len(xr)
				continue
			} else if err != nil {
				_ = openedFile.Close()
				slog.Error("Could not create recipe from text file", "file", zf.Name, "error", err)
				continue
			}
			recipes = append(recipes, recipe)
			recipeNumber++
		case app.ImageExt:
			parts := strings.Split(zf.Name, "/")

			imagePath := filepath.Join(app.ImagesDir, parts[len(parts)-1])
			_, err = os.Stat(imagePath)
			if errors.Is(err, os.ErrNotExist) {
				dest, err := os.Create(imagePath)
				if err != nil {
					slog.Error("Failed to create image file", "file", zf, "error", err)
					continue
				}

				_, err = io.Copy(dest, openedFile)
				if err != nil {
					_ = dest.Close()
					slog.Error("Failed to copy image file", "file", zf, "error", err)
					continue
				}

				_ = dest.Close()
			}

			thumbnailPath := filepath.Join(app.ThumbnailsDir, parts[len(parts)-1])
			_, err = os.Stat(thumbnailPath)
			if errors.Is(err, os.ErrNotExist) {
				thumbnail, err := os.Create(thumbnailPath)
				if err != nil {
					slog.Error("Failed to create thumbnail file", "file", zf, "error", err)
					continue
				}

				imageFile, err := os.Open(imagePath)
				if err != nil {
					slog.Error("Failed to create image file", "file", zf, "error", err)
					continue
				}

				_, err = io.Copy(thumbnail, imageFile)
				if err != nil {
					_ = imageFile.Close()
					_ = thumbnail.Close()
					slog.Error("Failed to copy thumbnail file", "file", zf, "error", err)
					continue
				}

				_ = imageFile.Close()
				_ = thumbnail.Close()
			}

		case app.VideoExt:
			parts := strings.Split(zf.Name, "/")
			path := filepath.Join(app.VideosDir, parts[len(parts)-1])
			_, err = os.Stat(path)
			if errors.Is(err, os.ErrNotExist) {
				dest, err := os.Create(path)
				if err != nil {
					slog.Error("Failed to create video file", "file", zf, "error", err)
					continue
				}

				_, err = io.Copy(dest, openedFile)
				if err != nil {
					slog.Error("Failed to copy video file", "file", zf, "error", err)
					continue
				}
			}
		}

		_ = openedFile.Close()
	}

	n := len(recipes)
	if n > 0 && len(recipes[n-1].Images) == 0 {
		recipes[n-1].Images = append(recipes[n-1].Images, imageUUID)
	}

	return recipes
}

func (f *Files) extractJSONRecipes(rd io.Reader) (models.Recipes, error) {
	buf, err := io.ReadAll(rd)
	if err != nil {
		slog.Error("Failed to read file", "reader", rd, "error", err)
		return nil, err
	}

	if len(buf) == 0 {
		return nil, errors.New("no bytes to unmarshal")
	}

	var xrs []models.RecipeSchema

	switch buf[0] {
	case '{':
		var rs models.RecipeSchema
		err = json.Unmarshal(buf, &rs)

		if rs.Ingredients == nil && rs.Instructions == nil {
			var m integrations.MealieRecipe
			err = json.Unmarshal(buf, &m)
			xrs = append(xrs, m.Schema())
		} else {
			xrs = append(xrs, rs)
		}
	case '[':
		err = json.Unmarshal(buf, &xrs)
	default:
		err = errors.New("unexpected JSON format")
	}

	if err != nil {
		return nil, err
	}

	xr := make(models.Recipes, 0, len(xrs))
	for _, rs := range xrs {
		r, err := rs.Recipe()
		if err != nil {
			return nil, fmt.Errorf("rs.Recipe() err: %w", err)
		}

		if rs.Image != nil && rs.Image.Value != "" {
			img, err := uuid.Parse(filepath.Base(rs.Image.Value))
			if err != nil {
				img, _ = f.ScrapeAndStoreImage(rs.Image.Value)
			}

			if img != uuid.Nil {
				r.Images = []uuid.UUID{img}
			}
		}

		if r.Name == "" || len(r.Ingredients) == 0 || len(r.Instructions) == 0 {
			continue
		}
		xr = append(xr, *r)
	}

	return xr, nil
}

func (f *Files) parseJSONRecipe(rd io.Reader) (*models.Recipe, error) {
	buf, err := io.ReadAll(rd)
	if err != nil {
		slog.Error("Failed to read file", "reader", rd, "error", err)
		return nil, err
	}

	var rs models.RecipeSchema
	err = json.Unmarshal(buf, &rs)
	if err != nil {
		return nil, fmt.Errorf("extract recipe: %w", err)
	}

	r, err := rs.Recipe()
	if err != nil {
		return nil, fmt.Errorf("rs.Recipe() err: %w", err)
	}

	if rs.Image.Value != "" {
		img, err := f.ScrapeAndStoreImage(rs.Image.Value)
		if err != nil {
			slog.Error("Could not scrape and store image", "image", rs.Image.Value, "error", err)
		}

		if img != uuid.Nil {
			r.Images = []uuid.UUID{img}
		}
	}

	return r, err
}

func (f *Files) processRecipeKeeper(zr *zip.Reader) models.Recipes {
	file := zr.File[slices.IndexFunc(zr.File, func(file *zip.File) bool {
		return file.Name == "recipes.html"
	})]

	rc, err := file.Open()
	if err != nil {
		return nil
	}
	defer rc.Close()

	root, err := goquery.NewDocumentFromReader(rc)
	if err != nil {
		return nil
	}

	recipes := models.NewRecipesFromRecipeKeeper(root)

	root.Find(".recipe-details").Each(func(i int, sel *goquery.Selection) {
		src, ok := sel.Find(".recipe-photo").First().Attr("src")
		if ok {
			imgFile := zr.File[slices.IndexFunc(zr.File, func(file *zip.File) bool {
				return file.Name == src
			})]

			open, err := imgFile.Open()
			if err == nil {
				img, _ := f.UploadImage(open)
				if err != nil {
					slog.Error("Could not scrape and store image", "image", imgFile.Name, "error", err)
				}

				if img != uuid.Nil {
					recipes[i].Images = []uuid.UUID{img}
				}

				_ = open.Close()
			}
		}
	})

	return recipes
}
