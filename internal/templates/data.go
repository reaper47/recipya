package templates

import (
	"github.com/reaper47/recipya/internal/models"
)

// IndexData holds data to pass on to the index template.
type RecipesData struct {
	Recipes    []models.Recipe
	Pagination Pagination
}

// RecipeData holds data to pass to the recipe templates.
type RecipeData struct {
	Recipe models.Recipe
}

// Pagination holds pagination data for templates with pagination.
type Pagination struct {
	// [<][1][2][3][...][87][88][89][>] =>
	Left, Right          []int
	Prev, Selected, Next int
	NumPages, NumResults int
}

// Init initializes the pagination struct.
func (p *Pagination) Init(page int) {
	numPages := p.NumResults / 12
	if p.NumResults%12 > 0 {
		numPages++
	}
	p.NumPages = numPages

	if page == 1 {
		p.Prev = 1
	} else {
		p.Prev = page - 1
	}

	if (page * 12) == p.NumResults {
		p.Next = page
	} else {
		p.Next = page + 1
	}

	p.Selected = page
	if page < numPages/2 {
		if p.Prev == p.Selected {
			p.Left = append(p.Left, p.Selected, p.Next, p.Next+1)
		} else {
			p.Left = append(p.Left, p.Prev, p.Selected, p.Next)
		}

		p.Right = append(p.Right, numPages-2, numPages-1, numPages)
	} else {
		if p.Selected == numPages {
			p.Right = append(p.Right, p.Prev-1, p.Prev, p.Selected)
		} else {
			p.Right = append(p.Right, p.Prev, p.Selected, p.Next)
		}
		p.Left = append(p.Left, 1, 2, 3)
	}
}
