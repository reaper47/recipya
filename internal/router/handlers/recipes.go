package handlers

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/reaper47/recipya/internal/config"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"golang.org/x/image/draw"
)

// RecipesAdd handles the GET /recipes/new URI.
func RecipesAdd(wr http.ResponseWriter, req *http.Request) {
	err := templates.Render(wr, "recipes-new.gohtml", nil)
	if err != nil {
		log.Println(err)
	}
}

// Recipe handles the /recipes/{id:[0-9]+} endpoint.
func Recipe(wr http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)
	if err != nil {
		showErrorPage(wr, "The id is not specified:", err)
		return
	}

	switch req.Method {
	case http.MethodGet:
		handleGetRecipe(wr, id)
	case http.MethodDelete:
		handleDeleteRecipe(wr, req, id)
	}
}

func handleGetRecipe(wr http.ResponseWriter, id int64) {
	r, err := config.App().Repo.GetRecipe(id)
	if err != nil {
		showErrorPage(wr, "Could not retrieve the recipe:", err)
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(r)
	if err != nil {
		showErrorPage(wr, "Could not send the recipe:", err)
		return
	}
	wr.Write(j)
}

func handleDeleteRecipe(wr http.ResponseWriter, req *http.Request, id int64) {
	err := config.App().Repo.DeleteRecipe(id)
	if err != nil {
		showErrorPage(wr, "Could not delete the recipe:", err)
		return
	}
	http.Redirect(wr, req, "/", http.StatusSeeOther)
}

// GetRecipesNewManual handles the GET /recipes/new/manual URI.
func GetRecipesNewManual(wr http.ResponseWriter, req *http.Request) {
	err := templates.Render(wr, "recipes-new-manual.gohtml", nil)
	if err != nil {
		log.Println(err)
	}
}

// PostRecipesNewManual handles the POST /recipes/new/manual URI.
func PostRecipesNewManual(wr http.ResponseWriter, req *http.Request) {
	yield, err := strconv.ParseInt(req.FormValue("yields"), 10, 16)
	if err != nil {
		showErrorPage(wr, "An error occured when retrieving the yield count:", err)
		return
	}

	prepTime, err := timeToDuration(req, "time-preparation")
	if err != nil {
		showErrorPage(wr, "An error occured when parsing the preparation time:", err)
		return
	}

	cookTime, err := timeToDuration(req, "time-cooking")
	if err != nil {
		showErrorPage(wr, "An error occured when parsing the cooking time:", err)
		return
	}

	file, _, err := req.FormFile("image")
	var imageUUID uuid.UUID
	if err == nil {
		defer file.Close()
		imageUUID, err = saveImage(file, "img")
		if err != nil {
			showErrorPage(wr, "An error occured when saving the image:", err)
			return
		}
	}

	r := models.Recipe{
		Name:        req.FormValue("title"),
		Description: req.FormValue("description"),
		Image:       imageUUID,
		Url:         req.FormValue("source"),
		Yield:       int16(yield),
		Category:    req.FormValue("category"),
		Times: models.Times{
			Prep: prepTime,
			Cook: cookTime,
		},
		Ingredients:  getFormItems(req, "ingredient"),
		Instructions: getFormItems(req, "instruction"),
		Nutrition: models.Nutrition{
			Calories:           req.FormValue("calories"),
			TotalCarbohydrates: req.FormValue("total-carbohydrates"),
			Sugars:             req.FormValue("sugars"),
			Protein:            req.FormValue("protein"),
			TotalFat:           req.FormValue("total-fat"),
			SaturatedFat:       req.FormValue("saturated-fat"),
			Cholesterol:        req.FormValue("cholesterol"),
			Sodium:             req.FormValue("sodium"),
			Fiber:              req.FormValue("fiber"),
		},
	}

	id, err := config.App().Repo.InsertNewRecipe(r)
	if err != nil {
		showErrorPage(wr, "An error occured when inserting the recipe:", err)
		return
	}

	fmt.Println(id)
	wr.WriteHeader(http.StatusCreated)
}

func timeToDuration(req *http.Request, field string) (time.Duration, error) {
	t := strings.SplitN(req.FormValue(field), ":", 3)
	return time.ParseDuration(t[0] + "h" + t[1] + "m" + t[2] + "s")
}

func getFormItems(req *http.Request, field string) []string {
	items := []string{}
	i := 1
	for {
		item := req.FormValue(field + "-" + strconv.Itoa(i))
		if item == "" {
			break
		}
		items = append(items, item)
		i++
	}
	return items
}

func saveImage(file multipart.File, dir string) (uuid.UUID, error) {
	uuid := uuid.New()

	tmp, err := os.Create("data/" + dir + "/" + uuid.String())
	if err != nil {
		return uuid, err
	}
	defer tmp.Close()

	img, err := compressImage(file)
	if err != nil {
		return uuid, err
	}

	err = jpeg.Encode(tmp, img, &jpeg.Options{Quality: 50})
	if err != nil {
		return uuid, err
	}

	_, err = io.Copy(tmp, file)
	if err != nil {
		return uuid, err
	}
	return uuid, nil
}

func compressImage(f multipart.File) (*image.RGBA, error) {
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	x := img.Bounds().Max.X
	if x > 2160 {
		x /= 2
	}

	y := img.Bounds().Max.Y
	if y > 1440 {
		y /= 2
	}

	dst := image.NewRGBA(image.Rect(0, 0, x, y))
	draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
	return dst, nil
}
