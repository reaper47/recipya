package templates

// ResultsPerPage is the number of results to display per pagination page.
const (
	ResultsPerPage    = 15
	ResultsPerPageStr = "15"
)

// NewPagination initializes a pagination struct of n pages.
func NewPagination(page, numPages, numResults uint64, resultsPerPage uint64, url, queries string, isSwap bool) Pagination {
	p := Pagination{
		Functions:      NewFunctionsData[uint64](),
		IsHidden:       numResults == 0,
		IsSwap:         isSwap,
		NumPages:       numPages,
		NumResults:     numResults,
		ResultsPerPage: resultsPerPage,
		URL:            url,
		URLQueries:     queries,
	}

	totalPages := numResults / resultsPerPage
	if numResults%resultsPerPage > 0 {
		totalPages++
	} else {
		totalPages = numPages
	}
	numPages = totalPages
	p.NumPages = totalPages

	if page == 1 {
		p.Prev = 1
	} else {
		p.Prev = page - 1
	}

	if (page * resultsPerPage) == numResults {
		p.Next = page
	} else {
		p.Next = page + 1
	}

	p.Selected = page
	if numPages <= 10 {
		for i := uint64(1); i <= numPages; i++ {
			p.Left = append(p.Left, i)
		}
	} else {
		if page <= numPages/2 {
			if page > 4 {
				p.Middle = append(p.Middle, page-2, page-1, page, page+1, page+2)
				p.Left = append(p.Left, 1)
			} else {
				for i := uint64(1); i <= page; i++ {
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

	return p
}

// Pagination holds data related to pagination.
type Pagination struct {
	Left, Middle, Right  []uint64
	Prev, Selected, Next uint64
	NumPages, NumResults uint64
	ResultsPerPage       uint64

	Functions  FunctionsData[uint64]
	IsHidden   bool
	IsSwap     bool
	URL        string
	URLQueries string
}
