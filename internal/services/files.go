package services

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/exp/slices"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// NewFilesService creates a new Files that satisfies the FilesService interface.
func NewFilesService() *Files {
	return &Files{}
}

// Files is the entity that manages the email client.
type Files struct{}

func (f *Files) ExportRecipes(recipes models.Recipes) (string, error) {
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	for _, recipe := range recipes {
		data, err := json.Marshal(recipe.Schema())
		if err != nil {
			return "", err
		}

		out, err := writer.Create(recipe.Name + "/recipe.json")
		if err != nil {
			return "", err
		}

		if _, err = out.Write(data); err != nil {
			return "", err
		}

		if recipe.Image != uuid.Nil {
			fileName := recipe.Image.String() + ".jpg"
			filePath := filepath.Join(app.ImagesDir, fileName)

			if _, err := os.Stat(filePath); err == nil {
				out, err = writer.Create(recipe.Name + "/image.jpg")
				if err != nil {
					return "", err
				}

				data, err := os.ReadFile(filePath)
				if err != nil {
					return "", err
				}

				if _, err = out.Write(data); err != nil {
					return "", err
				}
			}
		}
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	out, err := os.CreateTemp("", "*")
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := out.Write(buf.Bytes()); err != nil {
		return "", err
	}

	return filepath.Base(out.Name()), nil
}

func (f *Files) ExtractRecipes(fileHeaders []*multipart.FileHeader) models.Recipes {
	recipes := make(models.Recipes, 0)
	var wg sync.WaitGroup

	for _, file := range fileHeaders {
		wg.Add(1)
		file := file
		go func() {
			defer wg.Done()
			content := file.Header.Get("Content-Type")
			if strings.Contains(content, "zip") {
				recipes = append(recipes, f.processZip(file)...)
			} else if strings.Contains(content, "json") {
				recipes = append(recipes, *processJSON(file))
			}
		}()
	}

	wg.Wait()
	return recipes
}

func (f *Files) processZip(file *multipart.FileHeader) models.Recipes {
	recipes := make(models.Recipes, 0)

	openFile, err := file.Open()
	if err != nil {
		log.Println(err)
		return recipes
	}
	defer openFile.Close()

	buf := new(bytes.Buffer)
	fileSize, err := io.Copy(buf, openFile)
	if err != nil {
		log.Println(err)
		return recipes
	}

	z, err := zip.NewReader(bytes.NewReader(buf.Bytes()), fileSize)
	if err != nil {
		log.Println(err)
		return recipes
	}

	var (
		isDir        bool
		imageUUID    uuid.UUID
		recipeNumber int
	)

	for _, file := range z.File {
		if isDir = file.FileInfo().IsDir(); isDir && imageUUID != uuid.Nil {
			recipes[recipeNumber-1].Image = imageUUID
			imageUUID = uuid.Nil
		}

		validImageFormats := []string{".jpg", ".jpeg", ".png"}
		if imageUUID == uuid.Nil && slices.Contains(validImageFormats, filepath.Ext(file.Name)) {
			zipFile, err := file.Open()
			if err != nil {
				log.Printf("Error opening image file: %q", err)
				continue
			}

			if file.FileInfo().Size() < 1<<12 {
				zipFile.Close()
				continue
			}

			imageUUID, err = f.UploadImage(zipFile)
			if err != nil {
				log.Printf("Error uploading image: %q", err)
			}
			zipFile.Close()
		}

		if filepath.Ext(file.Name) == ".json" {
			f, err := file.Open()
			if err != nil {
				log.Println(err)
				continue
			}

			r, err := extractRecipe(f)
			if err != nil {
				log.Printf("could not extract %s: %q", file.Name, err.Error())
				f.Close()
				continue
			}

			recipes = append(recipes, *r)
			recipeNumber++
			f.Close()
		}
	}
	return recipes
}

func processJSON(file *multipart.FileHeader) *models.Recipe {
	f, err := file.Open()
	if err != nil {
		log.Printf("error opening file %s: %q", file.Filename, err.Error())
		return nil
	}
	defer f.Close()

	r, err := extractRecipe(f)
	if err != nil {
		log.Printf("could not extract %s: %q", file.Filename, err.Error())
		return nil
	}
	return r
}

func extractRecipe(rd io.Reader) (*models.Recipe, error) {
	buf, err := io.ReadAll(rd)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var rs models.RecipeSchema

	if err = json.Unmarshal(buf, &rs); err != nil {
		return nil, fmt.Errorf("extract recipe: %q", err)
	}

	r, err := rs.Recipe()
	if err != nil {
		return nil, fmt.Errorf("rs.Recipe() err: %q", err)
	}
	return r, err
}

func (f *Files) ReadTempFile(name string) ([]byte, error) {
	file := filepath.Join(os.TempDir(), name)
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	os.Remove(file)
	return data, nil
}

func (f *Files) UploadImage(rc io.ReadCloser) (uuid.UUID, error) {
	img, _, err := image.Decode(rc)
	if err != nil {
		return uuid.Nil, err
	}

	imageUUID := uuid.New()
	out, err := os.Create(filepath.Join(app.ImagesDir, imageUUID.String()+".jpg"))
	if err != nil {
		return uuid.Nil, nil
	}
	defer out.Close()

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width > 800 || height > 800 {
		img = imaging.Resize(img, width/2, height/2, imaging.NearestNeighbor)
	}

	if err := jpeg.Encode(out, img, &jpeg.Options{Quality: 33}); err != nil {
		return uuid.Nil, nil
	}
	return imageUUID, nil
}
