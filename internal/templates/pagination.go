package templates

// NewPagination initializes a pagination struct of n pages.
func NewPagination(numPages, numResults, page int) Pagination {
	if page < 0 {
		return Pagination{}
	}

	pg := Pagination{NumPages: numPages, NumResults: numResults}

	totalPages := numResults / 12
	if numResults%12 > 0 {
		totalPages++
	}
	totalPages = numPages

	if page == 1 {
		pg.Prev = 1
	} else {
		pg.Prev = page - 1
	}

	if (page * 12) == numResults {
		pg.Next = page
	} else {
		pg.Next = page + 1
	}

	pg.Selected = page
	if numPages <= 10 {
		for i := 1; i <= numPages; i++ {
			pg.Left = append(pg.Left, i)
		}
	} else {
		if page <= numPages/2 {
			if page > 4 {
				pg.Middle = append(pg.Middle, page-2, page-1, page, page+1, page+2)
				pg.Left = append(pg.Left, 1)
			} else {
				for i := 1; i <= page; i++ {
					pg.Left = append(pg.Left, i)
				}
				pg.Left = append(pg.Left, page+1, page+2, page+3)
			}

			pg.Right = append(pg.Right, numPages)
		} else {
			if page < numPages-3 {
				pg.Middle = append(pg.Middle, page-2, page-1, page, page+1, page+2)
				pg.Right = append(pg.Right, numPages)
			} else {
				for i := page - 3; i <= numPages; i++ {
					pg.Right = append(pg.Right, i)
				}
			}

			pg.Left = append(pg.Left, 1)
		}
	}

	return pg
}

// Pagination holds data related to pagination.
type Pagination struct {
	Left, Middle, Right  []int
	Prev, Selected, Next int
	NumPages, NumResults int
}
