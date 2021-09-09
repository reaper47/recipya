package core

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/reaper47/recipya/api"
	"github.com/reaper47/recipya/config"
	"github.com/reaper47/recipya/model"
)

type extractInfo struct {
	unit string
	sum  float64
}

func (e *extractInfo) format() string {
	return strconv.FormatFloat(e.sum, 'f', 1, 32) + strings.ToLower(e.unit)
}

// FetchNutrientsInfo gets the nutrition facts of a product based on its ingredients.
func FetchNutrientsInfo(foodList []string) *model.NutritionSet {
	apiKey := config.Config.FdcApiKey
	if apiKey == "" || len(foodList) == 0 {
		return &model.NutritionSet{}
	}
	q := "pageSize=1"
	baseURL := "https://api.nal.usda.gov/fdc/v1/foods/search?" + q + "&api_key=" + apiKey

	foods := getFoods(foodList, baseURL)
	if len(foods) == 0 {
		return &model.NutritionSet{}
	}
	return getNutritionFacts(foods)
}

func getFoods(foodList []string, baseURL string) []*model.Foods {
	var urls []string
	for _, food := range foodList {
		urls = append(urls, baseURL+"&query="+food)
	}

	var wg sync.WaitGroup
	c := make(chan *model.HttpResponse, len(urls))
	for _, url := range urls {
		wg.Add(1)
		go api.GetAsync(url, c, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	var foods []*model.Foods
	for msg := range c {
		body, err := ioutil.ReadAll(msg.Response.Body)
		if err != nil {
			continue
		}
		msg.Response.Body.Close()

		var f model.Foods
		err = json.Unmarshal([]byte(body), &f)
		if err != nil || len(f.List) == 0 {
			continue
		}
		foods = append(foods, &f)
	}
	return foods
}

func getNutritionFacts(foods []*model.Foods) *model.NutritionSet {
	var fdcIDs []int64
	nutrients := map[int64]*extractInfo{
		model.FDCID_CALORIES:      {},
		model.FDCID_CARBOHYDRATE:  {},
		model.FDCID_FAT:           {},
		model.FDCID_SATURATED_FAT: {},
		model.FDCID_CHOLESTEROL:   {},
		model.FDCID_PROTEIN:       {},
		model.FDCID_SODIUM:        {},
		model.FDCID_FIBER:         {},
		model.FDCID_SUGAR:         {},
	}

	for _, foodList := range foods {
		if len(foodList.List) == 0 {
			continue
		}
		food := foodList.List[0]
		fdcIDs = append(fdcIDs, food.Id)

		for _, n := range food.Nutrients {
			if n.IsOfInterest() {
				if n.Amount > 0 {
					nutrients[n.Id].sum += n.Amount
				} else if n.Value > 0 {
					nutrients[n.Id].sum += n.Value
				}

				if nutrients[n.Id].unit == "" {
					nutrients[n.Id].unit = n.Unit
				}
			}
		}
	}

	cal := nutrients[model.FDCID_CALORIES]
	return &model.NutritionSet{
		FdcIDs:       fdcIDs,
		Calories:     strconv.FormatInt(int64(cal.sum), 10) + " " + strings.ToLower(cal.unit),
		Carbohydrate: nutrients[model.FDCID_CARBOHYDRATE].format(),
		Fat:          nutrients[model.FDCID_FAT].format(),
		SaturatedFat: nutrients[model.FDCID_SATURATED_FAT].format(),
		Cholesterol:  nutrients[model.FDCID_CHOLESTEROL].format(),
		Protein:      nutrients[model.FDCID_PROTEIN].format(),
		Sodium:       nutrients[model.FDCID_SODIUM].format(),
		Fiber:        nutrients[model.FDCID_FIBER].format(),
		Sugar:        nutrients[model.FDCID_SUGAR].format(),
	}
}
