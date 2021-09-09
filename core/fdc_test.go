package core

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/reaper47/recipya/api"
	"github.com/reaper47/recipya/config"
	"github.com/reaper47/recipya/mock"
	"github.com/reaper47/recipya/model"
)

func init() {
	api.Client = &mock.MockClient{}
}

func TestFdc(t *testing.T) {
	t.Run(
		"FetchNutrientsInfo yields nutrition facts for a recipe",
		test_FetchNutrientsInfo_MultipleIngredients,
	)
	t.Run(
		"FetchNutrientsInfo yields an empty nutrition facts for no ingredients",
		test_FetchNutrientsInfo_NoIngredients,
	)
	t.Run(
		"FetchNutrientsInfo yields an empty nutrition facts for an undefined API key",
		test_FetchNutrientsInfo_UndefinedApiKey,
	)
	t.Run(
		"FetchNutrientsInfo yields an empty nutrition facts for an empty FDC search result",
		test_FetchNutrientsInfo_NoResults,
	)
}

func test_FetchNutrientsInfo_MultipleIngredients(t *testing.T) {
	config.Config.FdcApiKey = "hello"
	mock.GetGetFunc = func(string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(fdcData()))),
		}, nil
	}

	actual := FetchNutrientsInfo([]string{"kiwi", "kiwi"})

	expected := &model.NutritionSet{
		FdcIDs:       []int64{1898469, 1898469},
		Calories:     "706 kcal",
		Carbohydrate: "173.4g",
		Fat:          "0.0g",
		SaturatedFat: "0.0g",
		Cholesterol:  "0.0mg",
		Protein:      "0.0g",
		Sodium:       "280.0mg",
		Fiber:        "6.6g",
		Sugar:        "133.4g",
	}
	if !cmp.Equal(actual, expected) {
		t.Errorf("the food struct are not equal: %v", cmp.Diff(actual, expected))
	}
}

func test_FetchNutrientsInfo_NoIngredients(t *testing.T) {
	config.Config.FdcApiKey = "hello"
	mock.GetGetFunc = func(string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(fdcData()))),
		}, nil
	}

	actual := FetchNutrientsInfo([]string{})

	expected := &model.NutritionSet{}
	if !cmp.Equal(actual, expected) {
		t.Errorf("the food struct are not equal: %v", cmp.Diff(actual, expected))
	}
}

func test_FetchNutrientsInfo_UndefinedApiKey(t *testing.T) {
	mock.GetGetFunc = func(string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(fdcData()))),
		}, nil
	}

	actual := FetchNutrientsInfo([]string{})

	expected := &model.NutritionSet{}
	if !cmp.Equal(actual, expected) {
		t.Errorf("the food struct are not equal: %v", cmp.Diff(actual, expected))
	}
}

func test_FetchNutrientsInfo_NoResults(t *testing.T) {
	config.Config.FdcApiKey = "hello"
	data := `{"totalHits":0,"foods":[]}`
	mock.GetGetFunc = func(string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(data))),
		}, nil
	}

	actual := FetchNutrientsInfo([]string{"xyzt"})

	expected := &model.NutritionSet{}
	if !cmp.Equal(actual, expected) {
		t.Errorf("the food struct are not equal: %v", cmp.Diff(actual, expected))
	}
}

func fdcData() string {
	return `{"totalHits":724,"foods":[{"fdcId":1898469,"description":"KIWI","lowercaseDescription":"kiwi","dataType":"Branded","publishedDate":"2021-07-29","brandOwner":"RegalHealthFoodInternational,Inc.","brandName":"REGALGOURMETSNACKS","ingredients":"KIWI,SUGAR,CITRICACID,SULPHURDIOXIDE,FD&CYELLOW#5,ANDFD&CBLUE#1.","marketCountry":"UnitedStates","foodCategory":"WholesomeSnacks","score":1116.3937,"foodNutrients":[{"nutrientId":1003,"nutrientName":"Protein","nutrientNumber":"203","unitName":"G","derivationCode":"LCCS","value":0.0},{"nutrientId":1004,"nutrientName":"Totallipid(fat)","nutrientNumber":"204","unitName":"G","derivationCode":"LCCD","value":0.0},{"nutrientId":1005,"nutrientName":"Carbohydrate,bydifference","nutrientNumber":"205","unitName":"G","derivationCode":"LCCS","value":86.7},{"nutrientId":1008,"nutrientName":"Energy","nutrientNumber":"208","unitName":"KCAL","derivationCode":"LCCS","value":353},{"nutrientId":2000,"nutrientName":"Sugars,totalincludingNLEA","nutrientNumber":"269","unitName":"G","derivationCode":"LCCS","value":66.7},{"nutrientId":1079,"nutrientName":"Fiber,totaldietary","nutrientNumber":"291","unitName":"G","derivationCode":"LCCS","value":3.3},{"nutrientId":1087,"nutrientName":"Calcium,Ca","nutrientNumber":"301","unitName":"MG","derivationCode":"LCCD","value":67.0},{"nutrientId":1089,"nutrientName":"Iron,Fe","nutrientNumber":"303","unitName":"MG","derivationCode":"LCCD","value":1.2},{"nutrientId":1093,"nutrientName":"Sodium,Na","nutrientNumber":"307","unitName":"MG","derivationCode":"LCCS","value":140},{"nutrientId":1104,"nutrientName":"VitaminA,IU","nutrientNumber":"318","unitName":"IU","derivationCode":"LCCD","value":0.0},{"nutrientId":1162,"nutrientName":"VitaminC,totalascorbicacid","nutrientNumber":"401","unitName":"MG","derivationCode":"LCCD","value":4.0},{"nutrientId":1253,"nutrientName":"Cholesterol","nutrientNumber":"601","unitName":"MG","derivationCode":"LCCD","value":0.0},{"nutrientId":1257,"nutrientName":"Fattyacids,totaltrans","nutrientNumber":"605","unitName":"G","derivationCode":"LCCD","value":0.0},{"nutrientId":1258,"nutrientName":"Fattyacids,totalsaturated","nutrientNumber":"606","unitName":"G","derivationCode":"LCCD","value":0.0}]}]}`
}
