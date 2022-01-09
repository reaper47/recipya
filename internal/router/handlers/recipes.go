package handlers

import (
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
	"github.com/reaper47/recipya/internal/repository"
	"github.com/reaper47/recipya/internal/templates"
	"golang.org/x/image/draw"
)

// RecipesAdd handles the GET /recipes/new URI.
func RecipesAdd(w http.ResponseWriter, req *http.Request) {
	s := getSession(req)
	err := templates.Render(w, "recipe-new.gohtml", templates.Data{
		HeaderData: templates.HeaderData{
			AvatarInitials: s.UserInitials,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

// Recipe handles the /recipes/{id:[0-9]+} endpoint.
func Recipe(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)
	if err != nil {
		showErrorPage(w, "The id is not specified:", err)
		return
	}

	switch req.Method {
	case http.MethodGet:
		_, isAuthenticated := repository.IsAuthenticated(w, req)
		handleGetRecipe(w, req, id, isAuthenticated)
	case http.MethodDelete:
		handleDeleteRecipe(w, req, id)
	}
}

func handleGetRecipe(w http.ResponseWriter, req *http.Request, id int64, isAuthenticated bool) {
	r := config.App().Repo.GetRecipe(id)
	if r.ID == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	s, _ := repository.IsAuthenticated(w, req)

	err := templates.Render(w, "recipe-view.gohtml", templates.Data{
		RecipeData:   templates.RecipeData{Recipe: r},
		IsViewRecipe: true,
		HideSidebar:  !isAuthenticated,
		HeaderData: templates.HeaderData{
			IsUnauthenticated: !isAuthenticated,
			AvatarInitials:    s.UserInitials,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func handleDeleteRecipe(w http.ResponseWriter, req *http.Request, id int64) {
	err := config.App().Repo.DeleteRecipe(id)
	if err != nil {
		showErrorPage(w, "Could not delete the recipe:", err)
		return
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// EditRecipe handles the /recipes/{id:[0-9]+}/edit endpoint
func EditRecipe(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)
	if err != nil {
		showErrorPage(w, "The id is not specified:", err)
		return
	}

	switch req.Method {
	case http.MethodGet:
		handleGetEditRecipe(w, req, id)
	case http.MethodPost:
		handlePostEditRecipe(w, req, id)
	}
}

func handleGetEditRecipe(w http.ResponseWriter, req *http.Request, id int64) {
	r := config.App().Repo.GetRecipe(id)
	if r.ID == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	s := getSession(req)
	err := templates.Render(w, "recipe-edit.gohtml", templates.Data{
		Categories: config.App().Repo.GetCategories(),
		RecipeData: templates.RecipeData{Recipe: r},
		HeaderData: templates.HeaderData{
			AvatarInitials: s.UserInitials,
		},
	})
	if err != nil {
		log.Println("EditRecipe error", err)
		return
	}
}

func handlePostEditRecipe(w http.ResponseWriter, req *http.Request, id int64) {
	r, err := getRecipeFromForm(req)
	if err != nil {
		showErrorPage(w, "Could not retrieve the recipe from the form:", err)
		return
	}
	r.ID = id

	err = config.App().Repo.UpdateRecipe(r)
	if err != nil {
		showErrorPage(w, "An error occured when updating the recipe:", err)
		return
	}
	http.Redirect(w, req, "/recipes/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

// GetRecipesNewManual handles the GET /recipes/new/manual URI.
func GetRecipesNewManual(w http.ResponseWriter, req *http.Request) {
	s := getSession(req)
	err := templates.Render(w, "recipe-new-manual.gohtml", templates.Data{
		Categories: config.App().Repo.GetCategories(),
		HeaderData: templates.HeaderData{
			AvatarInitials: s.UserInitials,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

// PostRecipesNewManual handles the POST /recipes/new/manual URI.
func PostRecipesNewManual(w http.ResponseWriter, req *http.Request) {
	r, err := getRecipeFromForm(req)
	if err != nil {
		showErrorPage(w, "Could not retrieve the recipe from the form:", err)
		return
	}

	s := getSession(req)
	id, err := config.App().Repo.InsertNewRecipe(r, s.UserID)
	if err != nil {
		showErrorPage(w, "An error occured when inserting the recipe:", err)
		return
	}
	http.Redirect(w, req, "/recipes/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

func getRecipeFromForm(req *http.Request) (models.Recipe, error) {
	yield, err := strconv.ParseInt(req.FormValue("yields"), 10, 16)
	if err != nil {
		return models.Recipe{}, err
	}

	prepTime, err := timeToDuration(req, "time-preparation")
	if err != nil {
		return models.Recipe{}, err
	}

	cookTime, err := timeToDuration(req, "time-cooking")
	if err != nil {
		return models.Recipe{}, err
	}

	file, _, err := req.FormFile("image")
	var imageUUID uuid.UUID
	if err == nil {
		defer file.Close()
		imageUUID, err = saveImage(file, "img")
		if err != nil {
			return models.Recipe{}, err
		}
	}

	return models.Recipe{
		Name:        req.FormValue("title"),
		Description: req.FormValue("description"),
		Image:       imageUUID,
		Url:         req.FormValue("source"),
		Yield:       int16(yield),
		Category:    strings.ToLower(req.FormValue("category")),
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
	}, nil
}

func timeToDuration(req *http.Request, field string) (time.Duration, error) {
	t := strings.SplitN(req.FormValue(field), ":", 3)
	return time.ParseDuration(t[0] + "h" + t[1] + "m" + t[2] + "s")
}

func getFormItems(req *http.Request, field string) []string {
	itemMap := make(map[string]bool)

	items := []string{}
	i := 1
	for {
		item := strings.ToLower(req.FormValue(field + "-" + strconv.Itoa(i)))
		if item == "" {
			break
		}

		_, found := itemMap[item]
		if !found {
			itemMap[item] = true
			items = append(items, item)
		}
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
