package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/reaper47/recipya/internal/config"
	"github.com/reaper47/recipya/internal/constants"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/repository"
	"github.com/reaper47/recipya/internal/templates"
	"golang.org/x/image/draw"
)

// RecipesAdd handles the GET /recipes/new URI.
func RecipesAdd(w http.ResponseWriter, req *http.Request) {
	err := templates.Render(w, "recipe-new.gohtml", templates.Data{
		HeaderData: templates.HeaderData{
			AvatarInitials: getSession(req).UserInitials,
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
	r := config.App().Repo.Recipe(id)
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
	r := config.App().Repo.Recipe(id)
	if r.ID == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	s := getSession(req)
	err := templates.Render(w, "recipe-edit.gohtml", templates.Data{
		Categories: config.App().Repo.Categories(s.UserID),
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

	if err = config.App().Repo.UpdateRecipe(r); err != nil {
		showErrorPage(w, "An error occured when updating the recipe:", err)
		return
	}
	http.Redirect(w, req, "/recipes/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

// GetRecipesNewManual handles the GET /recipes/new/manual URI.
func GetRecipesNewManual(w http.ResponseWriter, req *http.Request) {
	s := getSession(req)
	err := templates.Render(w, "recipe-new-manual.gohtml", templates.Data{
		Categories: config.App().Repo.Categories(s.UserID),
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
		showErrorPage(w, "Could not retrieve the recipe from the form.", err)
		return
	}

	id, err := config.App().Repo.InsertNewRecipe(r, getSession(req).UserID)
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

	times, err := models.NewTimes(req.FormValue("time-preparation"), req.FormValue("time-cooking"))
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
		Name:         req.FormValue("title"),
		Description:  req.FormValue("description"),
		Image:        imageUUID,
		Url:          req.FormValue("source"),
		Yield:        int16(yield),
		Category:     strings.TrimSpace(strings.ToLower(req.FormValue("category"))),
		Times:        times,
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

func getFormItems(req *http.Request, field string) []string {
	itemMap := make(map[string]bool)

	items := []string{}
	i := 1
	for {
		item := strings.ToLower(req.FormValue(field + "-" + strconv.Itoa(i)))
		if item == "" {
			break
		}

		if _, found := itemMap[item]; !found {
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

// Categories handles the /recipes/categories endpoint.
func Categories(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		handlePostCategories(w, req)
	}
}

func handlePostCategories(w http.ResponseWriter, req *http.Request) {
	j := make(map[string]string)
	err := json.NewDecoder(req.Body).Decode(&j)
	if err != nil {
		writeJson(w, "Could not decode categories JSON.", http.StatusInternalServerError)
		return
	}

	c, ok := j["category"]
	if !ok {
		writeJson(w, "JSON does not contain the key 'category'.", http.StatusBadRequest)
		return
	}

	s := getSession(req)
	err = config.App().Repo.InsertCategory(c, s.UserID)
	if err != nil {
		writeJson(w, "Could not insert the category - "+c+".", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ImportRecipes handles the POST /recipes/import endpoint.
func ImportRecipes(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(120 << 20)
	if err != nil {
		showErrorPage(w, "Could not parse the uploaded files.", err)
		return
	}

	files, filesOk := req.MultipartForm.File["files"]
	if !filesOk {
		showErrorPage(w, "Could not retrieve the files or the directory from the form.", nil)
		return
	}

	userID := getSession(req).UserID
	var wg sync.WaitGroup
	for i, file := range files {
		wg.Add(1)

		f := file
		i := i

		go func() {
			defer wg.Done()
			importRecipe(f, i, userID)
		}()
	}
	wg.Wait()

	http.Redirect(w, req, "/recipes", http.StatusSeeOther)
}

func importRecipe(file *multipart.FileHeader, fnumber int, userID int64) {
	content := file.Header.Get("Content-Type")
	if strings.Contains(content, "zip") {
		processZip(file, userID)
	} else if strings.Contains(content, "json") {
		processJSON(file, userID)
	}
}

func processZip(file *multipart.FileHeader, userID int64) {
	f, err := file.Open()
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	fsize, err := io.Copy(buf, f)
	if err != nil {
		log.Println(err)
		return
	}

	z, err := zip.NewReader(bytes.NewReader(buf.Bytes()), fsize)
	if err != nil {
		log.Println(err)
		return
	}

	for _, file := range z.File {
		if filepath.Ext(file.Name) == ".json" {
			f, err := file.Open()
			if err != nil {
				log.Println(err)
				f.Close()
				continue
			}

			err = extractRecipe(f, userID)
			if err != nil {
				log.Println(err)
			}
			f.Close()
		}
	}
}

func processJSON(file *multipart.FileHeader, userID int64) {
	f, err := file.Open()
	if err != nil {
		log.Printf("error opening file %s: '%s'\n", file.Filename, err)
		return
	}
	defer f.Close()

	err = extractRecipe(f, userID)
	if err != nil {
		log.Printf("could not extract %s - %s", file.Filename, err)
		return
	}
}

func extractRecipe(rd io.Reader, userID int64) error {
	buf, err := ioutil.ReadAll(rd)
	if err != nil {
		log.Println(err)
		return err
	}

	var rs models.RecipeSchema
	err = json.Unmarshal(buf, &rs)
	if err != nil {
		log.Println(err)
		return err
	}

	r, err := rs.ToRecipe()
	if err != nil {
		log.Println("ToRecipe err - ", err)
		return err
	}

	_, err = config.App().Repo.InsertNewRecipe(r, userID)
	return err
}

// ExportRecipes handles the POST /settings/export endpoint.
func ExportRecipes(w http.ResponseWriter, req *http.Request) {
	buf := new(bytes.Buffer)
	z := zip.NewWriter(buf)

	recipes, err := config.App().Repo.Recipes(getSession(req).UserID, -1)
	if err != nil {
		showErrorPage(w, "Could not retrieve user recipes.", err)
		return
	}

	for _, r := range recipes {
		j, err := json.MarshalIndent(r, "", "\t")
		if err != nil {
			log.Printf("Could not marshal recipe %#v\r", r)
			continue
		}

		zw, err := z.Create(r.Category + "/" + r.Name + ".json")
		if err != nil {
			log.Printf("Could not create file for %s\r", j)
			continue
		}

		_, err = zw.Write(j)
		if err != nil {
			log.Printf("Could not write file %s\r", j)
			continue
		}
	}

	err = z.Close()
	if err != nil {
		log.Printf("Could not close zip writer: '%s'\n", err)
	}

	w.Header().Set(constants.HeaderContentType, constants.ApplicationZip)
	w.Header().Set(constants.HeaderContentDisposition, `attachment; filename="recipya-recipes.zip"`)
	w.Write(buf.Bytes())
}
