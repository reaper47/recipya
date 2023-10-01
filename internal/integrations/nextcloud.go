package integrations

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"io"
	"net/http"
	"sync"
)

const (
	baseURLNextcloud = "/apps/cookbook/api/v1"
)

// NextcloudImport imports recipes from a Nextcloud instance.
func NextcloudImport(client *http.Client, baseURL, username, password string, uploadImageFunc func(rc io.ReadCloser) (uuid.UUID, error)) (*models.Recipes, error) {
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	header := fmt.Sprintf("Basic %s", auth)

	recipesURL := fmt.Sprintf("%s%s/recipes", baseURL, baseURLNextcloud)
	resRecipes, err := sendBasicAuthRequest(client, recipesURL, header)
	if err != nil {
		return nil, err
	}
	defer resRecipes.Body.Close()

	var allRecipes []models.NextcloudRecipes
	err = json.NewDecoder(resRecipes.Body).Decode(&allRecipes)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	recipes := make(models.Recipes, len(allRecipes))
	for i, r := range allRecipes {
		wg.Add(1)
		go func(i int, r models.NextcloudRecipes, authHeader string) {
			defer wg.Done()

			url := fmt.Sprintf("%s%s/recipes/%d", baseURL, baseURLNextcloud, r.ID)
			res, err := sendBasicAuthRequest(client, url, header)
			if err != nil {
				return
			}
			defer res.Body.Close()

			var rs models.RecipeSchema
			err = json.NewDecoder(res.Body).Decode(&rs)
			if err != nil {
				return
			}

			recipe, err := rs.Recipe()
			if err != nil {
				return
			}

			recipes[i] = *recipe

			url = fmt.Sprintf("%s/image?size=thumb", url)
			imageRes, err := sendBasicAuthRequest(client, url, header)
			if err != nil {
				return
			}
			defer imageRes.Body.Close()

			imageUUID, err := uploadImageFunc(imageRes.Body)
			if err != nil {
				return
			}

			recipes[i].Image = imageUUID
		}(i, r, header)
	}
	wg.Wait()

	return &recipes, nil
}

func sendBasicAuthRequest(client *http.Client, url string, auth string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", auth)
	return client.Do(req)
}
