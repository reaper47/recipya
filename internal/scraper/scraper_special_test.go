package scraper_test

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/scraper"
	"github.com/reaper47/recipya/internal/services"
	"io"
	"net/http"
	"testing"
)

func TestScraper_Bergamot(t *testing.T) {
	s := scraper.NewScraper(&mockHTTPClient{DoFunc: func(req *http.Request) (*http.Response, error) {
		json := `{"id":210338,"shortId":"mIB4jYQtZU1A97","userId":585,"userFavorite":0,"sourceId":10,"sourceUrl":"https://www.elle.fr/Elle-a-Table/Recettes-de-cuisine/Soupe-miso-aux-oignons-nouveaux-tofu-et-saumon-emiette-4188650","lang":"","title":"Soupe miso aux oignons nouveaux, tofu et saumon émietté","description":"La soupe miso enrichie de saumon.","userNote":null,"ingredients":[{"data":["100 g de saumon frais","6 oignons nouveaux","30 g d'algues wakame séchées","70 g de pâte de miso blanc","quelques cives","300 g de tofu soyeux","1 cuillère(s) à soupe d'huile de sésame","1 cuillère(s) à soupe de graines de sésame"]}],"instructions":[{"data":["Dans une poêle bien chaude, faites cuire le saumon côté peau pendant 5 mn, puis laissez-le refroidir avant de l’émietter.","Dans une casserole, versez 1,5l d’eau, la moitié des oignons lavés et coupés en deux dans la hauteur, et les algues, puis portez à ébullition, réduisez ensuite le feu et laissez mijoter pendant 20 mn. Filtrez et ajoutez le miso, mélangez soigneusement.","Ajoutez le reste des oignons coupés en quatre, les cives lavées et émincées, le tofu coupé en dés et les miettes de saumon. Arrosez d’huile de sésame et parsemez de graines de sésame. Dégustez bien chaud."]}],"time":{"prepTime":20,"cookTime":null,"totalTime":50},"nutrition":{},"servings":4,"createdAt":"2024-01-16T16:16:24.000Z","updatedAt":"2024-01-16T16:16:24.000Z","deletedAt":null,"photos":[{"id":198989,"recipeId":210338,"reference":"210338PZAE79ER","order":0,"status":"uploaded","isUserUploaded":0,"sourceUrl":"https://resize.elle.fr/portrait_1280/var/plain_site/storage/images/elle-a-table/recettes-de-cuisine/soupe-miso-aux-oignons-nouveaux-tofu-et-saumon-emiette-4188650/101348896-2-fre-FR/Soupe-miso-aux-oignons-nouveaux-tofu-et-saumon-emiette.jpg","filenameExtension":"jpg","createdAt":"2024-01-16T16:16:24.000Z","updatedAt":"2024-01-16T16:16:24.000Z","deletedAt":null,"photoUrl":"https://aihkimhfpo.cloudimg.io/v7/foodbox/210338PZAE79ER.jpg?w=1280","photoThumbUrl":"https://aihkimhfpo.cloudimg.io/v7/foodbox/210338PZAE79ER.jpg?w=600&h=338"}],"sourceDomain":"elle.fr"}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(json))),
		}, nil
	}})
	f := services.NewFilesService()

	got, err := s.Scrape("https://dashboard.bergamot.app/shared/mIB4jYQtZU1A97", f)
	if err != nil {
		return
	}

	want := models.RecipeSchema{
		AtContext:     "https://schema.org",
		AtType:        models.SchemaType{Value: "Recipe"},
		CookTime:      "PT30M",
		CookingMethod: models.CookingMethod{},
		Cuisine:       models.Cuisine{},
		DateCreated:   "2024-01-16 16:16:24 +0000 UTC",
		DateModified:  "2024-01-16 16:16:24 +0000 UTC",
		DatePublished: "2024-01-16 16:16:24 +0000 UTC",
		Description:   models.Description{Value: "La soupe miso enrichie de saumon."},
		Keywords:      models.Keywords{},
		Image: models.Image{
			Value: "https://aihkimhfpo.cloudimg.io/v7/foodbox/210338PZAE79ER.jpg?w=600&h=338",
		},
		Ingredients: models.Ingredients{
			Values: []string{
				"100 g de saumon frais",
				"6 oignons nouveaux",
				"30 g d'algues wakame séchées",
				"70 g de pâte de miso blanc",
				"quelques cives",
				"300 g de tofu soyeux",
				"1 cuillère(s) à soupe d'huile de sésame",
				"1 cuillère(s) à soupe de graines de sésame",
			},
		},
		Instructions: models.Instructions{
			Values: []string{
				"Dans une poêle bien chaude, faites cuire le saumon côté peau pendant 5 mn, puis laissez-le refroidir avant de l’émietter.",
				"Dans une casserole, versez 1,5l d’eau, la moitié des oignons lavés et coupés en deux dans la hauteur, et les algues, puis portez à ébullition, réduisez ensuite le feu et laissez mijoter pendant 20 mn. Filtrez et ajoutez le miso, mélangez soigneusement.",
				"Ajoutez le reste des oignons coupés en quatre, les cives lavées et émincées, le tofu coupé en dés et les miettes de saumon. Arrosez d’huile de sésame et parsemez de graines de sésame. Dégustez bien chaud.",
			},
		},
		Name:            "Soupe miso aux oignons nouveaux, tofu et saumon émietté",
		NutritionSchema: models.NutritionSchema{},
		PrepTime:        "PT20M",
		Tools:           models.Tools{},
		Yield:           models.Yield{Value: 4},
		URL:             "https://dashboard.bergamot.app/shared/mIB4jYQtZU1A97",
	}
	if !cmp.Equal(got, want) {
		t.Logf(cmp.Diff(got, want))
		t.Fatal()
	}
}

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return nil, nil
}
