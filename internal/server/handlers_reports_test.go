package server_test

import (
	"errors"
	"fmt"
	"github.com/reaper47/recipya/internal/models"
	"net/http"
	"testing"
	"time"
)

func TestHandlers_Reports(t *testing.T) {
	srv := newServerTest()
	originalRepo := srv.Repository

	uri := "/reports"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("no reports selected on load", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertStringsInHTML(t, getBodyHTML(rr), []string{
			`<title hx-swap-oob="true">Reports | Recipya</title>`,
			`<button class="active" hx-get="/reports?tab=imports" hx-target="#tab-content">Imports</button>`,
			`<ul class="col-span-1 border-r overflow-auto max-h-44 border-b md:border-b-0 sm:max-h-full dark:border-r-gray-800"></ul>`,
			`<p class="p-4 sm:p-0">No report selected. Please select a report to view its content.</p>`,
		})
	})

	t.Run("error fetching reports", func(t *testing.T) {
		srv.Repository = &mockRepository{ReportsFunc: func(userID int64) ([]models.Report, error) {
			return nil, errors.New("oops")
		}}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Failed to fetch reports.\",\"title\":\"\"}"}`)
	})

	t.Run("view latest report", func(t *testing.T) {
		srv.Repository = &mockRepository{
			Reports: map[int64][]models.Report{1: {
				{
					ID:        1,
					CreatedAt: time.Date(2020, 03, 14, 1, 6, 0, 0, time.Local),
					ExecTime:  3 * time.Second,
					Logs:      []models.ReportLog{{ID: 1}, {ID: 2}},
				},
				{
					ID:        2,
					CreatedAt: time.Date(2020, 03, 15, 4, 9, 0, 0, time.Local),
					ExecTime:  9 * time.Second,
					Logs:      []models.ReportLog{{ID: 1}},
				},
			}},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?view=latest", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		body := getBodyHTML(rr)
		t.Log(body[987:])
		assertStringsInHTML(t, body, []string{
			`<title hx-swap-oob="true">Reports | Recipya</title>`,
			`<button class="active" hx-get="/reports?tab=imports" hx-target="#tab-content">Imports</button>`,
			`<ul class="col-span-1 border-r overflow-auto max-h-44 border-b md:border-b-0 sm:max-h-full dark:border-r-gray-800"><li class="item p-2 hover:bg-slate-200 cursor-default dark:hover:bg-slate-700 bg-slate-200 dark:bg-slate-700" hx-get="/reports/1" hx-target="#report-view-pane" hx-swap="outerHTML" hx-push-url="true" _="on click remove .bg-slate-200 .dark:bg-slate-700 from .item then add .bg-slate-200 .dark:bg-slate-700"><span><b>14 Mar 20 01:06 CET</b><br><span class="text-sm">Execution time: 3s</span></span><span class="badge badge-primary float-right select-none">2</span></li><li class="item p-2 hover:bg-slate-200 cursor-default dark:hover:bg-slate-700" hx-get="/reports/2" hx-target="#report-view-pane" hx-swap="outerHTML" hx-push-url="true" _="on click remove .bg-slate-200 .dark:bg-slate-700 from .item then add .bg-slate-200 .dark:bg-slate-700"><span><b>15 Mar 20 04:09 CET</b><br><span class="text-sm">Execution time: 9s</span></span><span class="badge badge-primary float-right select-none">1</span></li></ul>`,
			`<div id="report-view-pane" class="col-span-3"><div class="overflow-auto h-[89vh]"><table class="table table-xs sm:table-md"><thead><tr><th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=id-reverse" hx-target="#report-view-pane">ID <span>▾</span></th><th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=title" hx-target="#report-view-pane">Title</th><th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=success" hx-target="#report-view-pane">Success</th><th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=error" hx-target="#report-view-pane">Error</th></tr></thead><tbody><tr class="bg-red-200 dark:bg-red-700"><th>1</th><td></td><td>X</td><td>-</td></tr><tr class="bg-red-200 dark:bg-red-700"><th>2</th><td></td><td>X</td><td>-</td></tr></tbody></table></div></div>`,
		})
	})

	t.Run("user has import reports", func(t *testing.T) {
		srv.Repository = &mockRepository{
			Reports: map[int64][]models.Report{1: {
				{
					ID:        1,
					CreatedAt: time.Date(2020, 03, 14, 1, 6, 0, 0, time.Local),
					ExecTime:  3 * time.Second,
					Logs:      []models.ReportLog{{ID: 1}, {ID: 2}},
				},
				{
					ID:        2,
					CreatedAt: time.Date(2020, 03, 15, 4, 9, 0, 0, time.Local),
					ExecTime:  9 * time.Second,
					Logs:      []models.ReportLog{{ID: 1}},
				},
			}},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertStringsInHTML(t, getBodyHTML(rr), []string{
			`<title hx-swap-oob="true">Reports | Recipya</title>`,
			`<button class="active" hx-get="/reports?tab=imports" hx-target="#tab-content">Imports</button>`,
			`<ul class="col-span-1 border-r overflow-auto max-h-44 border-b md:border-b-0 sm:max-h-full dark:border-r-gray-800"><li class="item p-2 hover:bg-slate-200 cursor-default dark:hover:bg-slate-700" hx-get="/reports/1" hx-target="#report-view-pane" hx-swap="outerHTML" hx-push-url="true" _="on click remove .bg-slate-200 .dark:bg-slate-700 from .item then add .bg-slate-200 .dark:bg-slate-700"><span><b>14 Mar 20 01:06 CET</b><br><span class="text-sm">Execution time: 3s</span></span><span class="badge badge-primary float-right select-none">2</span></li><li class="item p-2 hover:bg-slate-200 cursor-default dark:hover:bg-slate-700" hx-get="/reports/2" hx-target="#report-view-pane" hx-swap="outerHTML" hx-push-url="true" _="on click remove .bg-slate-200 .dark:bg-slate-700 from .item then add .bg-slate-200 .dark:bg-slate-700"><span><b>15 Mar 20 04:09 CET</b><br><span class="text-sm">Execution time: 9s</span></span><span class="badge badge-primary float-right select-none">1</span></li></ul>`,
			`<div id="report-view-pane" class="grid col-span-3 place-content-center"><p class="p-4 sm:p-0">No report selected. Please select a report to view its content.</p></div>`,
		})
	})
}

func TestHandlers_Reports_Report(t *testing.T) {
	srv := newServerTest()
	originalRepo := srv.Repository

	uri := func(id string) string { return fmt.Sprintf("/reports/%s", id) }

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri("1"))
	})

	t.Run("redirect if access report without htmx", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri("1"), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "Location", "/reports")
	})

	testcases := []struct {
		name string
		id   string
	}{
		{name: "id cannot be 0", id: "0"},
		{name: "id cannot be negative", id: "0"},
		{name: "id cannot be anything else", id: "bob"},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri(tc.id), noHeader, nil)

			assertStatus(t, rr.Code, http.StatusBadRequest)
			assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Report ID must be positive.\",\"title\":\"\"}"}`)
		})
	}

	t.Run("report does not exist", func(t *testing.T) {
		srv.Repository = &mockRepository{
			Reports: map[int64][]models.Report{1: make([]models.Report, 0)},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri("6"), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Failed to fetch report.\",\"title\":\"\"}"}`)
	})

	t.Run("user cannot view report of other user", func(t *testing.T) {
		srv.Repository = &mockRepository{
			Reports: map[int64][]models.Report{
				1: {{ID: 1}},
				2: {{ID: 2}},
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedInOther(srv, http.MethodGet, uri("1"), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Failed to fetch report.\",\"title\":\"\"}"}`)
	})

	sortTestcases := []struct {
		name        string
		sortType    string
		wantHeaders []string
	}{
		{
			name:     "no sort",
			sortType: "",
			wantHeaders: []string{
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=id-reverse" hx-target="#report-view-pane">ID <span>▾</span></th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=title" hx-target="#report-view-pane">Title</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=success" hx-target="#report-view-pane">Success</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=error" hx-target="#report-view-pane">Error</th>`,
			},
		},
		{
			name:     "id",
			sortType: "id",
			wantHeaders: []string{
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=id-reverse" hx-target="#report-view-pane">ID <span>▾</span></th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=title" hx-target="#report-view-pane">Title</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=success" hx-target="#report-view-pane">Success</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=error" hx-target="#report-view-pane">Error</th>`,
			},
		},
		{
			name:     "title",
			sortType: "title",
			wantHeaders: []string{
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=id" hx-target="#report-view-pane">ID</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=title-reverse" hx-target="#report-view-pane">Title <span>▾</span></th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=success" hx-target="#report-view-pane">Success</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=error" hx-target="#report-view-pane">Error</th>`,
			},
		},
		{
			name:     "success",
			sortType: "success",
			wantHeaders: []string{
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=id" hx-target="#report-view-pane">ID</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=title" hx-target="#report-view-pane">Title</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=success-reverse" hx-target="#report-view-pane">Success <span>▾</span></th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=error" hx-target="#report-view-pane">Error</th>`,
			},
		},
		{
			name:     "error",
			sortType: "error",
			wantHeaders: []string{
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=id" hx-target="#report-view-pane">ID</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=title" hx-target="#report-view-pane">Title</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=success" hx-target="#report-view-pane">Success</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=error-reverse" hx-target="#report-view-pane">Error <span>▾</span></th>`,
			},
		},
		{
			name:     "error-reverse",
			sortType: "error-reverse",
			wantHeaders: []string{
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=id" hx-target="#report-view-pane">ID</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=title" hx-target="#report-view-pane">Title</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=success" hx-target="#report-view-pane">Success</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=error" hx-target="#report-view-pane">Error <span>▴</span></th>`,
			},
		},
		{
			name:     "success-reverse",
			sortType: "success-reverse",
			wantHeaders: []string{
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=id" hx-target="#report-view-pane">ID</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=title" hx-target="#report-view-pane">Title</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=success" hx-target="#report-view-pane">Success <span>▴</span></th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=error" hx-target="#report-view-pane">Error</th>`,
			},
		},
		{
			name:     "title-reverse",
			sortType: "title-reverse",
			wantHeaders: []string{
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=id" hx-target="#report-view-pane">ID</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=title" hx-target="#report-view-pane">Title <span>▴</span></th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=success" hx-target="#report-view-pane">Success</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=error" hx-target="#report-view-pane">Error</th>`,
			},
		},
		{
			name:     "id-reverse",
			sortType: "id-reverse",
			wantHeaders: []string{
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=id" hx-target="#report-view-pane">ID <span>▴</span></th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=title" hx-target="#report-view-pane">Title</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=success" hx-target="#report-view-pane">Success</th>`,
				`<th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=error" hx-target="#report-view-pane">Error</th>`,
			},
		},
	}
	for _, tc := range sortTestcases {
		t.Run(tc.name, func(t *testing.T) {
			srv.Repository = &mockRepository{
				Reports: map[int64][]models.Report{
					1: {{ID: 1, Logs: []models.ReportLog{{ID: 1, Title: "Fried Chicken", IsSuccess: true}}}},
				},
			}
			defer func() {
				srv.Repository = originalRepo
			}()

			rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri("1")+"?sort="+tc.sortType, noHeader, nil)

			assertStatus(t, rr.Code, http.StatusOK)
			assertStringsInHTML(t, getBodyHTML(rr), tc.wantHeaders)
		})
	}

	tabsTestcases := []struct {
		name string
		tab  string
		want []string
	}{
		{
			name: "imports",
			tab:  "imports",
			want: []string{
				`<div id="report-view-pane" class="col-span-3"><div class="overflow-auto h-[89vh]">`,
				`<table class="table table-xs sm:table-md"><thead><tr><th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=id-reverse" hx-target="#report-view-pane">ID <span>▾</span></th><th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=title" hx-target="#report-view-pane">Title</th><th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=success" hx-target="#report-view-pane">Success</th><th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=error" hx-target="#report-view-pane">Error</th></tr></thead><tbody><tr class=""><th>1</th><td>Fried Chicken</td><td>&#x2713;</td><td>-</td></tr><tr class="bg-red-200 dark:bg-red-700"><th>2</th><td>Coq au vin with fries</td><td>X</td><td>Meaning of life not found.</td></tr></tbody></table></div></div>`,
			},
		},
	}
	for _, tc := range tabsTestcases {
		t.Run(tc.name, func(t *testing.T) {
			srv.Repository = &mockRepository{
				Reports: map[int64][]models.Report{
					1: {
						{
							ID:        1,
							CreatedAt: time.Date(2020, 03, 14, 1, 6, 0, 0, time.Local),
							ExecTime:  3 * time.Second,
							Logs: []models.ReportLog{
								{ID: 1, Title: "Fried Chicken", IsSuccess: true},
								{ID: 2, Title: "Coq au vin with fries", IsSuccess: false, Error: "Meaning of life not found."},
							},
						},
					},
				},
			}
			defer func() {
				srv.Repository = originalRepo
			}()

			rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri("1")+"?tab="+tc.tab, noHeader, nil)

			assertStatus(t, rr.Code, http.StatusOK)
			body := getBodyHTML(rr)
			assertStringsInHTML(t, body, tc.want)
			assertStringsNotInHTML(t, body, []string{
				`<div id="tab-content" role="tabpanel" class="w-[90vw] text-sm md:text-base p-4 auto-rows-min md:w-full">`,
				`<p>No report selected. Please select a report to view its content.</p>`,
			})
		})
	}

	t.Run("valid request", func(t *testing.T) {
		srv.Repository = &mockRepository{
			Reports: map[int64][]models.Report{
				1: {
					{
						ID:        1,
						CreatedAt: time.Date(2020, 03, 14, 1, 6, 0, 0, time.Local),
						ExecTime:  3 * time.Second,
						Logs: []models.ReportLog{
							{ID: 1, Title: "Fried Chicken", IsSuccess: true},
							{ID: 2, Title: "Coq au vin with fries", IsSuccess: false, Error: "Meaning of life not found."},
						},
					},
				},
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri("1"), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		body := getBodyHTML(rr)
		assertStringsInHTML(t, body, []string{
			`<div id="report-view-pane" class="col-span-3"><div class="overflow-auto h-[89vh]">`,
			`<table class="table table-xs sm:table-md"><thead><tr><th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=id-reverse" hx-target="#report-view-pane">ID <span>▾</span></th><th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=title" hx-target="#report-view-pane">Title</th><th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=success" hx-target="#report-view-pane">Success</th><th class="cursor-default hover:bg-blue-50 dark:hover:bg-blue-700" hx-get="?sort=error" hx-target="#report-view-pane">Error</th></tr></thead><tbody><tr class=""><th>1</th><td>Fried Chicken</td><td>&#x2713;</td><td>-</td></tr><tr class="bg-red-200 dark:bg-red-700"><th>2</th><td>Coq au vin with fries</td><td>X</td><td>Meaning of life not found.</td></tr></tbody></table>`,
		})
		assertStringsNotInHTML(t, body, []string{`<p>No report selected. Please select a report to view its content.</p>`})
	})
}
