package services

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"image"
	"image/jpeg"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	fontFamily    = "Arial"
	fontSizeBig   = 16
	fontSizeSmall = 9
)

// NewFilesService creates a new Files that satisfies the FilesService interface.
func NewFilesService() *Files {
	return &Files{}
}

// Files is the entity that manages the email client.
type Files struct{}

type exportData struct {
	recipeName  string
	recipeImage uuid.UUID
	data        []byte
}

// BackupDB backs up the whole database to the backup directory.
func (f *Files) BackupDB() error {
	name := fmt.Sprintf("recipya-%s.zip", time.Now().Format(time.DateOnly))
	target := filepath.Join(app.BackupPath, "all", name)

	err := os.MkdirAll(filepath.Dir(target), os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create backup dir: %q", err)
	}

	zf, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("could not create backup %q", name)
	}
	defer func() {
		_ = zf.Close()
	}()

	zw := zip.NewWriter(zf)
	defer func() {
		_ = zw.Close()
	}()

	source := filepath.Dir(app.DBBasePath)
	backupBase := filepath.Base(app.BackupPath)

	err = filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		omit := []string{"backup", "all", "fdc.db", app.RecipyaDB + "-wal", app.RecipyaDB + "-shm"}
		base := filepath.Base(filepath.Dir(path))
		if base == backupBase || base == "all" || slices.Contains(omit, info.Name()) {
			return nil
		}

		h, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		h.Method = zip.Deflate
		h.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			h.Name += "/"
		}

		w, err := zw.CreateHeader(h)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()

		_, err = io.Copy(w, f)
		return err
	})
	if err != nil {
		return fmt.Errorf("could not assemble backup %q", name)
	}

	cleanBackups(filepath.Dir(target))
	return nil
}

// Backups gets the list of backup dates sorted in descending order for the given user.
func (f *Files) Backups(userID int64) []time.Time {
	var backups []time.Time
	root := filepath.Join(app.BackupPath, strconv.FormatInt(userID, 64))
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		name := info.Name()
		ext := filepath.Ext(name)
		if ext != ".bak" {
			return nil
		}

		_, after, found := strings.Cut(strings.TrimSuffix(name, ext), ".")
		if found {
			parsed, err := time.Parse("", after)
			if err == nil {
				backups = append(backups, parsed)
			}
		}
		return nil
	})

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].After(backups[j])
	})
	return backups
}

func (f *Files) BackupUserData(repo RepositoryService) error {
	for _, user := range repo.Users() {
		err := backupUserData(repo, f, user.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func backupUserData(repo RepositoryService, files *Files, userID int64) error {
	name := fmt.Sprintf("recipya-%s.zip", time.Now().Format(time.DateOnly))
	target := filepath.Join(app.BackupPath, "users", strconv.FormatInt(userID, 10), name)

	err := os.MkdirAll(filepath.Dir(target), os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create backup dir: %q", err)
	}

	zf, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("could not create backup %q", name)
	}
	defer func() {
		_ = zf.Close()
	}()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	w, err := zw.CreateHeader(&zip.FileHeader{
		Name:   "recipes.zip",
		Method: zip.Store,
	})
	if err != nil {
		return err
	}

	recipes := repo.RecipesAll(userID)
	buf, err := files.ExportRecipes(recipes, models.JSON, nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, buf)
	return err
}

func cleanBackups(root string) {
	files, err := os.ReadDir(root)
	if err != nil {
		log.Printf("cleanBackups for %q error: %q", root, err)
		return
	}

	sort.Slice(files, func(i, j int) bool {
		aInfo, err1 := files[i].Info()
		bInfo, err2 := files[j].Info()
		if err1 != nil || err2 != nil {
			return false
		}
		return bInfo.ModTime().Before(aInfo.ModTime())
	})

	if len(files) > 10 {
		for _, file := range files[10:] {
			_ = os.Remove(filepath.Join(root, file.Name()))
		}
	}
}

// ExportRecipes creates a zip containing the recipes to export in the desired file type.
func (f *Files) ExportRecipes(recipes models.Recipes, fileType models.FileType, progress chan int) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	switch fileType {
	case models.JSON:
		for i, e := range exportRecipesJSON(recipes) {
			if progress != nil {
				progress <- i
			}

			out, err := writer.Create(e.recipeName + "/recipe" + fileType.Ext())
			if err != nil {
				return nil, err
			}

			_, err = out.Write(e.data)
			if err != nil {
				return nil, err
			}

			if e.recipeImage != uuid.Nil {
				filePath := filepath.Join(app.ImagesDir, e.recipeImage.String()+".jpg")

				_, err = os.Stat(filePath)
				if err == nil {
					out, err = writer.Create(e.recipeName + "/image.jpg")
					if err != nil {
						return nil, err
					}

					data, err := os.ReadFile(filePath)
					if err != nil {
						return nil, err
					}

					_, err = out.Write(data)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	case models.PDF:
		processed := make(map[string]struct{})
		for i, e := range exportRecipesPDF(recipes) {
			if progress != nil {
				progress <- i
			}

			name := strings.ReplaceAll(e.recipeName+fileType.Ext(), "/", "_")

			_, found := processed[name]
			if found {
				name += "_" + uuid.NewString()[:4]
			}
			processed[name] = struct{}{}

			out, err := writer.Create(name)
			if err != nil {
				return nil, err
			}

			_, err = out.Write(e.data)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, errors.New("unsupported export file type")
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func exportRecipesJSON(recipes models.Recipes) []exportData {
	data := make([]exportData, len(recipes))
	for i, r := range recipes {
		xb, err := json.Marshal(r.Schema())
		if err != nil {
			continue
		}
		data[i] = exportData{
			recipeName:  r.Name,
			recipeImage: r.Image,
			data:        xb,
		}
	}
	return data
}

func exportRecipesPDF(recipes models.Recipes) []exportData {
	data := make([]exportData, len(recipes))
	for i, r := range recipes {
		data[i] = exportData{
			recipeName:  r.Name,
			recipeImage: r.Image,
			data:        recipeToPDF(&r),
		}
	}
	return data
}

func recipeToPDF(r *models.Recipe) []byte {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetAuthor("Recipya user", false)
	pdf.SetCreator("Recipya", false)
	sanitized := strings.ToValidUTF8(r.Name, "")
	pdf.SetSubject(sanitized, true)
	pdf.SetTitle(sanitized, true)
	pdf.SetCreationDate(time.Now())
	addRecipeToPDF(pdf, r)
	return pdfToBytes(pdf, r.Name)
}

func addRecipeToPDF(pdf *gofpdf.Fpdf, r *models.Recipe) *gofpdf.Fpdf {
	viewData := templates.NewViewRecipeData(1, r, true, false)

	tr := pdf.UnicodeTranslatorFromDescriptor("")
	marginLeft, marginTop, marginRight, _ := pdf.GetMargins()
	pageWidth, pageHeight := pdf.GetPageSize()

	pdf.SetHeaderFunc(func() {
		pdf.SetFont(fontFamily, "B", fontSizeBig)
		wd := pageWidth
		pdf.SetX(marginLeft)
		pdf.MultiCell(wd-marginLeft-marginRight, 9, r.Name, "1", "C", false)
	})

	pdf.SetFooterFunc(func() {
		if pdf.PageNo() == 1 {
			return
		}
		pdf.SetY(-15)
		pdf.SetFont(fontFamily, "I", fontSizeSmall-1)
		pdf.SetTextColor(128, 128, 128)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()-1), "", 0, "C", false, 0, "")
	})

	pdf.SetFont(fontFamily, "", fontSizeSmall)
	pdf.AddPage()
	pdf.Rect(marginLeft, marginTop, pageWidth-marginLeft-marginRight, pageHeight-3*marginTop, "D")

	// Category, servings, source
	pdf.SetX(marginLeft)

	var (
		colWd   = (pageWidth - marginLeft - marginRight) / 3.
		lineHt  = 5.0
		cellGap = 2.0
	)

	type cellType struct {
		str  string
		list [][]byte
		ht   float64
	}
	var (
		cellList [3]cellType
		cell     cellType
	)

	source := r.URL
	parse, err := url.Parse(source)
	if err == nil {
		source = parse.Hostname()
	}

	cols := []string{
		r.Category,
		strconv.FormatInt(int64(r.Yield), 10) + " servings",
		"Source: " + source,
	}

	y := pdf.GetY()
	originalY := y + 9
	maxHt := lineHt
	for j := 0; j < 3; j++ {
		lines := pdf.SplitLines([]byte(cols[j]), colWd-cellGap-cellGap)
		height := float64(len(lines)) * lineHt
		if height > maxHt {
			maxHt = height
		}
		cellList[j] = cellType{
			str:  cols[j],
			list: lines,
			ht:   height,
		}
	}

	x := marginLeft
	for i := 0; i < 3; i++ {
		pdf.Rect(pdf.GetX(), y, colWd, maxHt+cellGap+cellGap, "D")
		cell = cellList[i]
		cellY := y + cellGap + (maxHt-cell.ht)/2
		for splitJ := 0; splitJ < len(cell.list); splitJ++ {
			var linkStr string
			if i == 2 && parse != nil {
				linkStr = r.URL
			}

			pdf.SetXY(x+cellGap, cellY)
			pdf.CellFormat(colWd-cellGap, lineHt, tr(string(cell.list[splitJ])), "", 0, "C", false, 0, linkStr)
			cellY += lineHt
		}
		x += colWd
	}
	y += maxHt + cellGap + cellGap

	for j := 0; j < 3; j++ {
		lines := pdf.SplitLines([]byte(cols[j]), colWd-cellGap-cellGap)
		height := float64(len(lines)) * lineHt
		if height > maxHt {
			maxHt = height
		}
		cellList[j] = cellType{
			str:  cols[j],
			list: lines,
			ht:   height,
		}
	}

	// Times
	cols = []string{
		"Prep: " + viewData.FormattedTimes.Prep,
		"Cook: " + viewData.FormattedTimes.Cook,
		"Total: " + viewData.FormattedTimes.Total,
	}
	widths := []float64{colWd + 4*cellGap, colWd - 8*cellGap, colWd + 4*cellGap}
	h := lineHt + cellGap/2
	pdf.SetXY(marginLeft, y)
	pdf.Rect(pdf.GetX(), y, widths[0], maxHt, "D")
	pdf.CellFormat(widths[0], h, cols[0], "", 0, "C", false, 0, "")
	pdf.SetX(marginLeft + widths[0])
	pdf.Rect(pdf.GetX(), y, widths[1], maxHt, "D")
	pdf.CellFormat(widths[1], h, cols[1], "", 0, "C", false, 0, "")
	pdf.SetX(marginLeft + widths[0] + widths[1])
	pdf.Rect(pdf.GetX(), y, widths[2], maxHt, "D")
	pdf.SetFont(fontFamily, "B", fontSizeSmall)
	pdf.CellFormat(widths[2], h, cols[2], "", 0, "C", false, 0, "")
	pdf.SetFont(fontFamily, "", fontSizeSmall)
	y += maxHt

	// Description
	lines := pdf.SplitLines([]byte(r.Description), 3*colWd)
	height := float64(len(lines)) * lineHt
	if height > maxHt {
		maxHt = height
	}
	cellList[0] = cellType{
		str:  r.Description,
		list: lines,
		ht:   height,
	}

	x = marginLeft
	pdf.Rect(x, y, 3*colWd, maxHt+cellGap+cellGap, "D")
	cell = cellList[0]
	cellY := y + cellGap + (maxHt-cell.ht)/2
	for splitJ := 0; splitJ < len(cell.list); splitJ++ {
		pdf.SetXY(x+cellGap, cellY)
		pdf.CellFormat(marginLeft, lineHt, tr(string(cell.list[splitJ])), "", 0, "L", false, 0, "")
		cellY += lineHt
	}
	y += maxHt + cellGap + cellGap

	pdf.SetFont(fontFamily, "", fontSizeSmall)
	pdf.SetY(originalY)

	pdf.SetY(y)
	pdf.Ln(1)
	pdf.SetX(marginLeft)
	nutrition := make([]string, 0)
	if r.Nutrition.Calories != "" {
		nutrition = append(nutrition, "Calories: "+r.Nutrition.Calories+";")
	}
	if r.Nutrition.Cholesterol != "" {
		nutrition = append(nutrition, " Cholesterol: "+r.Nutrition.Cholesterol+";")
	}
	if r.Nutrition.Fiber != "" {
		nutrition = append(nutrition, " Fiber: "+r.Nutrition.Fiber+";")
	}
	if r.Nutrition.Protein != "" {
		nutrition = append(nutrition, " Protein: "+r.Nutrition.Protein+";")
	}
	if r.Nutrition.SaturatedFat != "" {
		nutrition = append(nutrition, " Saturated fat: "+r.Nutrition.SaturatedFat+";")
	}
	if r.Nutrition.Sodium != "" {
		nutrition = append(nutrition, " Sodium: "+r.Nutrition.Sodium+";")
	}
	if r.Nutrition.Sugars != "" {
		nutrition = append(nutrition, " Sugars: "+r.Nutrition.Sugars+";")
	}
	if r.Nutrition.TotalCarbohydrates != "" {
		nutrition = append(nutrition, " Total carbohydrates: "+r.Nutrition.TotalCarbohydrates+";")
	}
	if r.Nutrition.TotalFat != "" {
		nutrition = append(nutrition, " Total fat: "+r.Nutrition.TotalFat+";")
	}
	if r.Nutrition.UnsaturatedFat != "" {
		nutrition = append(nutrition, " Unsaturated fat: "+r.Nutrition.UnsaturatedFat+";")
	}
	if len(nutrition) > 0 {
		nutrition[0] = "  " + nutrition[0]
		nutrition[len(nutrition)/2-1] += "\n"

		pdf.SetX(marginLeft + cellGap)
		pdf.SetFont(fontFamily, "B", fontSizeSmall)
		pdf.CellFormat(12, 6, "Nutrition Facts", "", 1, "L", false, 0, "")
		pdf.SetFont(fontFamily, "", fontSizeSmall)
		pdf.SetX(marginLeft)
		pdf.MultiCell(pageWidth-2*marginLeft, 5, tr(strings.Join(nutrition, " ")), "B", "1", false)
	}

	// Ingredients
	ingredientsY := pdf.GetY()
	pdf.SetX(marginLeft + cellGap)
	pdf.SetFont(fontFamily, "B", fontSizeSmall)
	pdf.CellFormat(0, 6, "Ingredients", "", 1, "L", false, 0, "")
	pdf.SetFont(fontFamily, "", fontSizeSmall)

	onNewPage := true
	for _, ing := range r.Ingredients {
		currY := pdf.GetY()
		if currY > pageHeight-3*marginTop && onNewPage {
			pdf.AddPage()
			pdf.SetX(marginLeft + cellGap)
			pdf.SetFont(fontFamily, "B", fontSizeSmall)
			pdf.CellFormat(0, 7, "Ingredients (continued)", "", 1, "L", false, 0, "")
			pdf.SetFont(fontFamily, "", fontSizeSmall)
			onNewPage = false
		}
		pdf.MultiCell(pageWidth/3-marginLeft/2, 5, tr("-> "+ing), "", "L", false)
	}

	// Instructions
	pdf.SetPage(pdf.PageNo())
	pdf.SetXY(marginLeft+pageWidth/3, ingredientsY)
	pdf.SetFont(fontFamily, "B", fontSizeSmall)
	pdf.CellFormat(0, 6, "Instructions", "", 1, "L", false, 0, "")
	pdf.SetFont(fontFamily, "", fontSizeSmall)
	pdf.SetX(marginLeft + pageWidth/3)

	_, f := pdf.GetPageSize()
	for i, ins := range r.Instructions {
		pdf.SetX(marginLeft + pageWidth/3)
		if pdf.GetY() > f-15 {
			pdf.AddPage()
			pdf.SetXY(marginLeft+pageWidth/3, 9+marginTop)
			pdf.SetPage(pdf.PageNo())
			pdf.SetFont(fontFamily, "B", fontSizeSmall)
			pdf.CellFormat(0, 7, "Instructions (continued)", "", 1, "L", false, 0, "")
			pdf.SetFont(fontFamily, "", fontSizeSmall)
			pdf.SetX(marginLeft + pageWidth/3)
		}
		pdf.MultiCell(2*pageWidth/3-2*marginRight, 5, tr(strconv.Itoa(i+1)+". "+ins), "", "L", false)
	}
	pdf.SetPage(pdf.PageNo())
	pdf.Rect(marginLeft, marginTop, pageWidth-marginLeft-marginRight, pageHeight-3*marginTop, "D")
	return pdf
}

// ExtractRecipes extracts the recipes from the HTTP files.
func (f *Files) ExtractRecipes(fileHeaders []*multipart.FileHeader) models.Recipes {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("ExtractRecipes recovered from panic for %#v file headers: %q", fileHeaders, err)
		}
	}()

	var (
		recipes models.Recipes
		wg      sync.WaitGroup
		mu      sync.Mutex
	)
	wg.Add(len(fileHeaders))

	for _, file := range fileHeaders {
		go func(fh *multipart.FileHeader) {
			defer wg.Done()
			content := fh.Header.Get("Content-Type")
			if strings.Contains(content, "zip") {
				mu.Lock()
				recipes = append(recipes, f.processZip(fh)...)
				mu.Unlock()
			} else if strings.Contains(content, "json") {
				mu.Lock()
				recipes = append(recipes, *processJSON(fh))
				mu.Unlock()
			} else if content == "application/octet-stream" {
				switch strings.ToLower(filepath.Ext(fh.Filename)) {
				case models.MXP.Ext():
					mu.Lock()
					recipes = append(recipes, processMasterCook(fh)...)
					mu.Unlock()
				}
			}
		}(file)
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
	defer func() {
		_ = openFile.Close()
	}()

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
		imageUUID    uuid.UUID
		recipeNumber int
	)

	for _, zf := range z.File {
		if imageUUID != uuid.Nil && (zf.FileInfo().IsDir() || (recipeNumber > 0 && recipes[recipeNumber-1].Image == uuid.Nil)) {
			recipes[recipeNumber-1].Image = imageUUID
			imageUUID = uuid.Nil
		}

		validImageFormats := []string{".jpg", ".jpeg", ".png"}
		if imageUUID == uuid.Nil && slices.Contains(validImageFormats, filepath.Ext(zf.Name)) {
			imageFile, err := zf.Open()
			if err != nil {
				log.Printf("Error opening image file: %q", err)
				continue
			}

			if zf.FileInfo().Size() < 1<<12 {
				_ = imageFile.Close()
				continue
			}

			imageUUID, err = f.UploadImage(imageFile)
			if err != nil {
				log.Printf("Error uploading image: %q", err)
			}
			_ = imageFile.Close()
		}

		switch strings.ToLower(filepath.Ext(zf.Name)) {
		case models.JSON.Ext():
			jsonFile, err := zf.Open()
			if err != nil {
				log.Println(err)
				continue
			}

			r, err := extractRecipe(jsonFile)
			if err != nil {
				log.Printf("could not extract %s: %q", zf.Name, err.Error())
				_ = jsonFile.Close()
				continue
			}
			r.Image = uuid.Nil

			recipes = append(recipes, *r)
			recipeNumber++
			_ = jsonFile.Close()
		case models.MXP.Ext():
			mxpFile, err := zf.Open()
			if err != nil {
				log.Println(err)
				continue
			}

			xr := models.NewRecipesFromMasterCook(mxpFile)
			if len(xr) > 0 {
				recipes = append(recipes, xr...)
				recipeNumber += len(xr)
			}

			_ = mxpFile.Close()
		}
	}

	n := len(recipes)
	if n > 0 && recipes[n-1].Image == uuid.Nil {
		recipes[n-1].Image = imageUUID
	}

	return recipes
}

func processJSON(file *multipart.FileHeader) *models.Recipe {
	f, err := file.Open()
	if err != nil {
		log.Printf("error opening file %s: %q", file.Filename, err.Error())
		return nil
	}
	defer func() {
		_ = f.Close()
	}()

	r, err := extractRecipe(f)
	if err != nil {
		log.Printf("could not extract %s: %q", file.Filename, err.Error())
		return nil
	}
	return r
}

func processMasterCook(file *multipart.FileHeader) models.Recipes {
	f, err := file.Open()
	if err != nil {
		log.Printf("error opening file %s: %q", file.Filename, err.Error())
		return nil
	}
	defer func() {
		_ = f.Close()
	}()

	return models.NewRecipesFromMasterCook(f)
}

func extractRecipe(rd io.Reader) (*models.Recipe, error) {
	buf, err := io.ReadAll(rd)
	if err != nil {
		log.Println(err)
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
	return r, err
}

// ExportCookbook exports the cookbook in the desired file type.
// It returns the name of file in the temporary directory.
func (f *Files) ExportCookbook(cookbook models.Cookbook, fileType models.FileType) (string, error) {
	buf := new(bytes.Buffer)

	var tempFileName string
	switch fileType {
	case models.PDF:
		export := exportCookbookToPDF(&cookbook)
		_, err := buf.Write(export.data)
		if err != nil {
			return "", err
		}
		tempFileName = strings.Join(strings.Split(cookbook.Title, " "), "_") + "_*.pdf"
	default:
		return "", errors.New("unsupported export file type")
	}

	out, err := os.CreateTemp("", tempFileName)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = out.Close()
	}()

	_, err = out.Write(buf.Bytes())
	if err != nil {
		return "", err
	}

	return filepath.Base(out.Name()), nil
}

func exportCookbookToPDF(cookbook *models.Cookbook) exportData {
	return exportData{
		recipeName:  cookbook.Title,
		recipeImage: cookbook.Image,
		data:        cookbookToPDF(cookbook),
	}
}

func cookbookToPDF(cookbook *models.Cookbook) []byte {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetAuthor("Recipya user", false)
	pdf.SetCreator("Recipya", false)
	sanitized := strings.ToValidUTF8(cookbook.Title, "")
	pdf.SetSubject(sanitized, true)
	pdf.SetTitle(sanitized, true)
	pdf.SetCreationDate(time.Now())

	tr := pdf.UnicodeTranslatorFromDescriptor("")
	marginLeft, marginTop, marginRight, _ := pdf.GetMargins()
	pageWidth, pageHeight := pdf.GetPageSize()

	pdf.SetFont(fontFamily, "", fontSizeSmall)
	pdf.AddPage()
	pdf.SetPage(1)
	pdf.Rect(marginLeft, marginTop, pageWidth-marginLeft-marginRight, pageHeight-3*marginTop, "D")

	pdf.SetXY(pageWidth/2-marginLeft-marginRight, pageHeight/4-marginTop)
	pdf.SetFont(fontFamily, "B", fontSizeBig)
	pdf.CellFormat(12, 6, tr(cookbook.Title), "", 1, "L", false, 0, "")

	if cookbook.Image != uuid.Nil {
		exe, err := os.Executable()
		if err != nil {
			return nil
		}

		imageFile := filepath.Join(filepath.Dir(exe), "data", "images", cookbook.Image.String()+".jpg")
		pdf.ImageOptions(imageFile, pdf.GetX()+pageWidth/2-4*marginLeft, pdf.GetY()+marginTop, 0, 0, false, gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, 0, "")
	}

	pdf.SetXY(marginLeft+3, pageHeight-2.7*marginTop)
	pdf.SetFont(fontFamily, "B", 10)
	pdf.CellFormat(12, 6, "Dominant Categories: ", "", 1, "L", false, 0, "")
	pdf.SetFont(fontFamily, "", 10)
	pdf.SetXY(marginLeft*5.2, pageHeight-2.7*marginTop)
	categories := strings.Join(cookbook.DominantCategories(5), ", ")
	pdf.CellFormat(12, 6, tr(categories), "", 1, "L", false, 0, "")

	pdf.SetXY(pageWidth-marginLeft*3.2, pageHeight-2.7*marginTop)
	pdf.SetFont(fontFamily, "B", 10)
	n := len(cookbook.Recipes)
	s := " recipe"
	if n > 1 {
		s += "s"
	}
	numRecipes := strconv.Itoa(n) + s
	pdf.CellFormat(12, 6, numRecipes, "", 1, "L", false, 0, "")
	pdf.SetFont(fontFamily, "", fontSizeSmall)

	for _, r := range cookbook.Recipes {
		addRecipeToPDF(pdf, &r)
	}
	return pdfToBytes(pdf, cookbook.Title)
}

func pdfToBytes(pdf *gofpdf.Fpdf, name string) []byte {
	buf := &bytes.Buffer{}
	err := pdf.Output(buf)
	if err != nil {
		log.Printf("could not create a pdf for %q", name)
		return []byte{}
	}
	return buf.Bytes()
}

// ReadTempFile gets the content of a file in the temporary directory.
func (f *Files) ReadTempFile(name string) ([]byte, error) {
	file := filepath.Join(os.TempDir(), name)
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	_ = os.Remove(file)
	return data, nil
}

// UploadImage uploads an image to the server.
func (f *Files) UploadImage(rc io.ReadCloser) (uuid.UUID, error) {
	img, _, err := image.Decode(rc)
	if err != nil {
		return uuid.Nil, err
	}

	imageUUID := uuid.New()
	out, err := os.Create(filepath.Join(app.ImagesDir, imageUUID.String()+".jpg"))
	if err != nil {
		return uuid.Nil, err
	}
	defer func() {
		_ = out.Close()
	}()

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width > 800 || height > 800 {
		img = imaging.Resize(img, width/2, height/2, imaging.NearestNeighbor)
	}

	err = jpeg.Encode(out, img, &jpeg.Options{Quality: 33})
	if err != nil {
		return uuid.Nil, err
	}
	return imageUUID, nil
}
