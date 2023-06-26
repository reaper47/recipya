package services

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/reaper47/recipya/internal/models"
	"io"
	"log"
	"mime/multipart"
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

	for _, file := range z.File {
		if filepath.Ext(file.Name) == ".json" {
			f, err := file.Open()
			if err != nil {
				log.Println(err)
				_ = f.Close()
				continue
			}

			r, err := extractRecipe(f)
			if err != nil {
				log.Println(err)
			}
			f.Close()
			recipes = append(recipes, *r)
		}
	}
	return recipes
}

func processJSON(file *multipart.FileHeader) *models.Recipe {
	f, err := file.Open()
	if err != nil {
		log.Printf("error opening file %s: %s", file.Filename, err.Error())
		return nil
	}
	defer f.Close()

	r, err := extractRecipe(f)
	if err != nil {
		log.Printf("could not extract '%s': %s", file.Filename, err.Error())
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
		return nil, fmt.Errorf("extract recipe: %s", err)
	}

	r, err := rs.Recipe()
	if err != nil {
		return nil, fmt.Errorf("rs.Recipe() err: %s", err)
	}
	return r, err
}
