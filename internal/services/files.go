package services

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/exp/slices"
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

func (f *Files) ExtractRecipes(fileHeaders []*multipart.FileHeader) models.Recipes {
	recipes := make(models.Recipes, 0)
	var wg sync.WaitGroup

	for _, file := range fileHeaders {
		wg.Add(1)
		f := file
		go func() {
			defer wg.Done()
			content := f.Header.Get("Content-Type")
			if strings.Contains(content, "zip") {
				recipes = append(recipes, processZip(f)...)
			} else if strings.Contains(content, "json") {
				recipes = append(recipes, *processJSON(f))
			}
		}()
	}

	wg.Wait()
	return recipes
}

func processZip(file *multipart.FileHeader) models.Recipes {
	recipes := make(models.Recipes, 0)

	f, err := file.Open()
	if err != nil {
		log.Println(err)
		return recipes
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	fileSize, err := io.Copy(buf, f)
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
			f, err := file.Open()
			if err != nil {
				continue
			}

			if file.FileInfo().Size() < 1<<12 {
				f.Close()
				continue
			}

			exe, err := os.Executable()
			if err != nil {
				continue
			}

			imageUUID = uuid.New()

			var out *os.File
			if out, err = os.Create(filepath.Join(filepath.Dir(exe), "data", "images", imageUUID.String()+".jpg")); err != nil {
				f.Close()
				continue
			}

			io.Copy(out, f)
			out.Close()
			f.Close()
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
