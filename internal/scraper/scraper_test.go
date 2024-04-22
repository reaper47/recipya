package scraper_test

import (
	"bytes"
	"github.com/blang/semver"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v59/github"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/scraper"
	"github.com/reaper47/recipya/internal/services"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

const atContext = "https://schema.org"

type testcase struct {
	name string
	in   string
	want models.RecipeSchema
}

func test(t *testing.T, tc testcase) {
	t.Helper()

	defer func() {
		err := recover()
		if err != nil {
			t.Fatalf("panic while testing %s: %s", tc.name, err)
		}
	}()

	// Uncomment the two lines below to refresh the HTML files to test against any changes the sites have brought.
	// Once the HTML files have been updated, comment the lines back and run the tests.
	/*updateHTMLFile(t, tc.name, tc.in)
	return*/

	actual := testFile(t, tc.name, tc.in)

	if !cmp.Equal(actual, tc.want) {
		t.Logf(cmp.Diff(actual, tc.want))
		t.Fatal()
	}
	_, err := actual.Recipe()
	if err != nil {
		t.Fatal(err)
	}
}

func testFile(t *testing.T, name, url string) models.RecipeSchema {
	t.Helper()

	scrape := scraper.NewScraper(&mockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			host, _, _ := strings.Cut(name, ".")
			_, fileName, _, _ := runtime.Caller(0)
			f, err := os.Open(filepath.Join(path.Dir(fileName), "testdata", host+".html"))
			if err != nil {
				t.Fatalf("%s open file: %s", name, err)
			}
			defer f.Close()

			data, err := io.ReadAll(f)
			if err != nil {
				return nil, err
			}

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(data)),
			}, nil
		},
	})

	got, err := scrape.Scrape(url, &mockFiles{})
	if err != nil {
		t.Fatal(err)
	}

	if got.URL == "" {
		got.URL = url
	}

	got.AtContext = atContext
	return got
}

func updateHTMLFile(t *testing.T, name, url string) {
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

var anUploadedImage = uuid.New()

type mockFiles struct {
	exportHitCount      int
	extractRecipesFunc  func(fileHeaders []*multipart.FileHeader) models.Recipes
	ReadTempFileFunc    func(name string) ([]byte, error)
	uploadImageHitCount int
	uploadImageFunc     func(rc io.ReadCloser) (uuid.UUID, error)
}

func (m *mockFiles) BackupGlobal() error {
	panic("implement me")
}

func (m *mockFiles) Backups(userID int64) []time.Time {
	panic("implement me")
}

func (m *mockFiles) BackupUserData(repo services.RepositoryService, userID int64) error {
	panic("implement me")
}

func (m *mockFiles) BackupUsersData(repo services.RepositoryService) error {
	panic("implement me")
}

func (m *mockFiles) IsAppLatest(current semver.Version) (bool, *github.RepositoryRelease, error) {
	return false, nil, nil
}

func (m *mockFiles) ScrapeAndStoreImage(rawURL string) (uuid.UUID, error) {
	return anUploadedImage, nil
}

func (m *mockFiles) ExtractUserBackup(date string, userID int64) (*models.UserBackup, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockFiles) ExportCookbook(cookbook models.Cookbook, fileType models.FileType) (string, error) {
	m.exportHitCount++
	return cookbook.Title + fileType.Ext(), nil
}

func (m *mockFiles) ExportRecipes(recipes models.Recipes, fileType models.FileType, progress chan int) (*bytes.Buffer, error) {
	var sb strings.Builder
	for _, recipe := range recipes {
		sb.WriteString(recipe.Name + "-")
	}
	m.exportHitCount++
	return bytes.NewBufferString(sb.String()), nil
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

func (m *mockFiles) UpdateApp(_ semver.Version) error {
	return nil
}

func (m *mockFiles) UploadImage(rc io.ReadCloser) (uuid.UUID, error) {
	if m.uploadImageFunc != nil {
		return m.uploadImageFunc(rc)
	}
	m.uploadImageHitCount++
	return uuid.New(), nil
}
