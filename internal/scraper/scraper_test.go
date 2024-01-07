package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-cmp/cmp"
	"github.com/reaper47/recipya/internal/models"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type testcase struct {
	name string
	in   string
	want models.RecipeSchema
}

func test(t *testing.T, tc testcase) {
	defer func() {
		err := recover()
		if err != nil {
			t.Fatalf("panic while testing %s: %s", tc.name, err)
		}
	}()

	// Uncomment the two lines below to refresh the HTML files to test against any changes the sites have brought.
	/*updateHTMLFile(t, tc.name, tc.in)
	return*/

	actual := testFile(t, tc.name, tc.in)

	if !cmp.Equal(actual, tc.want) {
		t.Logf(cmp.Diff(actual, tc.want))
		t.Fatal()
	}
}

func testFile(t *testing.T, name, url string) models.RecipeSchema {
	t.Helper()
	host, _, _ := strings.Cut(name, ".")
	_, fileName, _, _ := runtime.Caller(0)
	f, err := os.Open(filepath.Join(path.Dir(fileName), "testdata", host+".html"))
	if err != nil {
		t.Fatalf("%s open file: %s", name, err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Fatalf("%s could not parse HTML: %s", host, err)
	}

	actual, err := scrapeWebsite(doc, getHost(url))
	if err != nil {
		t.Fatal(err)
	}

	if actual.URL == "" {
		actual.URL = url
	}

	actual.AtContext = atContext
	return actual
}

/*func updateHTMLFile(t *testing.T, name, url string) {
	t.Helper()

	c := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := c.Get(url)
	if err != nil {
		t.Log("could not fetch url")
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Logf("got status code %d", res.StatusCode)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Log(err)
		return
	}

	host, _, _ := strings.Cut(name, ".")
	_, fileName, _, _ := runtime.Caller(0)
	filePath := filepath.Join(path.Dir(fileName), "testdata", host+".html")
	err = os.WriteFile(filePath, body, os.ModePerm)
	if err != nil {
		t.Log(err)
		return
	}
}

type mockFiles struct {
	exportHitCount      int
	extractRecipesFunc  func(fileHeaders []*multipart.FileHeader) models.Recipes
	ReadTempFileFunc    func(name string) ([]byte, error)
	uploadImageHitCount int
	uploadImageFunc     func(rc io.ReadCloser) (uuid.UUID, error)
}

func (m *mockFiles) ExportCookbook(cookbook models.Cookbook, fileType models.FileType) (string, error) {
	m.exportHitCount++
	return cookbook.Title + fileType.Ext(), nil
}

func (m *mockFiles) ExportRecipes(recipes models.Recipes, _ models.FileType) (string, error) {
	var s string
	for _, recipe := range recipes {
		s += recipe.Name + "-"
	}
	m.exportHitCount++
	return s, nil
}

func (m *mockFiles) ExtractRecipes(fileHeaders []*multipart.FileHeader) models.Recipes {
	if m.extractRecipesFunc != nil {
		return m.extractRecipesFunc(fileHeaders)
	}
	return models.Recipes{}
}

func (m *mockFiles) ReadTempFile(name string) ([]byte, error) {
	if m.ReadTempFileFunc != nil {
		return m.ReadTempFileFunc(name)
	}
	return []byte(name), nil
}

func (m *mockFiles) UploadImage(rc io.ReadCloser) (uuid.UUID, error) {
	if m.uploadImageFunc != nil {
		return m.uploadImageFunc(rc)
	}
	m.uploadImageHitCount++
	return uuid.New(), nil
}
*/
