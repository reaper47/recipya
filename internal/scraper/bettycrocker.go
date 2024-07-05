package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"time"
)

func scrapeBettyCrocker(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return rs, err
	}

	actualLayout := "02/01/2006"

	dp, err := time.Parse(actualLayout, rs.DatePublished)
	if err == nil {
		rs.DatePublished = dp.Format(time.DateOnly)
	}

	dc, err := time.Parse(actualLayout, rs.DateCreated)
	if err == nil {
		rs.DateCreated = dc.Format(time.DateOnly)
	}

	dm, err := time.Parse(actualLayout, rs.DateModified)
	if err == nil {
		rs.DateModified = dm.Format(time.DateOnly)
	}

	return rs, nil
}
