package templates

import (
	"github.com/reaper47/recipya/internal/models"
)

// Data holds general template data.
type Data struct {
	HideSidebar  bool
	HideGap      bool
	HeaderData   HeaderData
	IsViewRecipe bool

	RecipesData RecipesData
	RecipeData  RecipeData
	Categories  []string

	FormErrorData FormErrorData
}

// HeaderData holds data for the header.
type HeaderData struct {
	Hide              bool
	IsUnauthenticated bool
	AvatarInitials    string
}

// FormErrorData holds errors related to forms.
type FormErrorData struct {
	Username, Email, Password string
}

func (m FormErrorData) IsEmpty() bool {
	return m.Username == "" && m.Email == "" && m.Password == ""
}

// IndexData holds data to pass on to the index template.
type RecipesData struct {
	Recipes    []models.Recipe
	Pagination Pagination
}

// RecipeData holds data to pass to the recipe templates.
type RecipeData struct {
	Recipe           models.Recipe
	HideEditControls bool
}

// Pagination holds pagination data for templates with pagination.
type Pagination struct {
	Left, Middle, Right  []int
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
	if numPages <= 10 {
		for i := 1; i <= numPages; i++ {
			p.Left = append(p.Left, i)
		}
	} else {
		if page <= numPages/2 {
			if page > 4 {
				p.Middle = append(p.Middle, page-2, page-1, page, page+1, page+2)
				p.Left = append(p.Left, 1)
			} else {
				for i := 1; i <= page; i++ {
					p.Left = append(p.Left, i)
				}
				p.Left = append(p.Left, page+1, page+2, page+3)
			}

			p.Right = append(p.Right, numPages)
		} else {
			if page < numPages-3 {
				p.Middle = append(p.Middle, page-2, page-1, page, page+1, page+2)
				p.Right = append(p.Right, numPages)
			} else {
				for i := page - 3; i <= numPages; i++ {
					p.Right = append(p.Right, i)
				}
			}

			p.Left = append(p.Left, 1)
		}
	}
}
