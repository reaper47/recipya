package server

import (
	"cmp"
	"github.com/a-h/templ"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/web/components"
	"net/http"
	"slices"
)

func (s *Server) reportsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		reports, err := s.Repository.ReportsImport(getUserID(r))
		if err != nil {
			s.Brokers.SendToast(models.NewErrorDBToast("Failed to fetch reports."), userID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := templates.Data{
			About:           templates.NewAboutData(),
			IsAdmin:         userID == 1,
			IsAutologin:     app.Config.Server.IsAutologin,
			IsAuthenticated: true,
			IsHxRequest:     r.Header.Get("Hx-Request") == "true",
			Reports:         templates.ReportsData{Imports: reports},
		}

		var c templ.Component
		switch r.URL.Query().Get("tab") {
		case "imports":
			c = components.ReportsTabImports(data, false)
		default:
			var isHighlightFirst bool
			if r.URL.Query().Get("view") == "latest" {
				data.Reports.CurrentReport, err = s.Repository.Report(reports[0].ID, getUserID(r))
				if err != nil {
					s.Brokers.SendToast(models.NewErrorDBToast("Failed to fetch report."), userID)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				isHighlightFirst = true
			}
			c = components.ReportsIndex(data, isHighlightFirst)
		}

		_ = c.Render(r.Context(), w)
	}
}

func (s *Server) reportsReportHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isHxRequest := r.Header.Get("Hx-Request") == "true"
		if !isHxRequest {
			http.Redirect(w, r, "/reports", http.StatusSeeOther)
			return
		}

		userID := getUserID(r)

		id, err := parsePathPositiveID(r.PathValue("id"))
		if err != nil {
			s.Brokers.SendToast(models.NewErrorReqToast("Report ID must be positive."), userID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		report, err := s.Repository.Report(id, getUserID(r))
		if err != nil {
			s.Brokers.SendToast(models.NewErrorDBToast("Failed to fetch report."), userID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sortType := r.URL.Query().Get("sort")
		switch sortType {
		case "id":
			sortType = "id"
			slices.SortFunc(report, func(a, b models.ReportLog) int { return cmp.Compare(a.ID, b.ID) })
		case "id-reverse":
			sortType = "id-reverse"
			slices.SortFunc(report, func(a, b models.ReportLog) int { return cmp.Compare(b.ID, a.ID) })
		case "title":
			sortType = "title"
			slices.SortFunc(report, func(a, b models.ReportLog) int { return cmp.Compare(a.Title, b.Title) })
		case "title-reverse":
			sortType = "title-reverse"
			slices.SortFunc(report, func(a, b models.ReportLog) int { return cmp.Compare(b.Title, a.Title) })
		case "success":
			sortType = "success"
			slices.SortFunc(report, func(a, b models.ReportLog) int {
				isASuccess := 1
				if !a.IsSuccess {
					isASuccess = 0
				}

				isBSuccess := 1
				if !b.IsSuccess {
					isBSuccess = 0
				}

				return cmp.Compare(isASuccess, isBSuccess)
			})
		case "success-reverse":
			sortType = "success-reverse"
			slices.SortFunc(report, func(a, b models.ReportLog) int {
				isASuccess := 1
				if !a.IsSuccess {
					isASuccess = 0
				}

				isBSuccess := 1
				if !b.IsSuccess {
					isBSuccess = 0
				}

				return cmp.Compare(isBSuccess, isASuccess)
			})
		case "error":
			sortType = "error"
			slices.SortFunc(report, func(a, b models.ReportLog) int { return cmp.Compare(a.Error, b.Error) })
		case "error-reverse":
			sortType = "error-reverse"
			slices.SortFunc(report, func(a, b models.ReportLog) int { return cmp.Compare(b.Error, a.Error) })
		}

		_ = components.Report(templates.Data{
			About:           templates.NewAboutData(),
			IsAdmin:         userID == 1,
			IsAutologin:     app.Config.Server.IsAutologin,
			IsAuthenticated: true,
			IsHxRequest:     isHxRequest,
			Reports: templates.ReportsData{
				CurrentReport: report,
				Sort:          sortType,
			},
		}).Render(r.Context(), w)
	}
}
