package scraper

import (
	"time"

	"github.com/reaper47/recipya/internal/constants"
	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeBettyCrocker(root *html.Node) (rs models.RecipeSchema, err error) {
	rs, err = scrapeLdJSON(root)

	actualLayout := "02/01/2006"
	expectedLayout := constants.BasicTimeLayout

	dp, err2 := time.Parse(actualLayout, rs.DatePublished)
	if err2 == nil {
		rs.DatePublished = dp.Format(expectedLayout)
	}

	dc, err2 := time.Parse(actualLayout, rs.DateCreated)
	if err2 == nil {
		rs.DateCreated = dc.Format(expectedLayout)
	}

	dm, err2 := time.Parse(actualLayout, rs.DateModified)
	if err2 == nil {
		rs.DateModified = dm.Format(expectedLayout)
	}

	return rs, err
}
