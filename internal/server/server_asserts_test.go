package server_test

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func assertCookbooksViewMode(tb testing.TB, mode models.ViewMode, got string) {
	tb.Helper()

	var want []string
	if mode == models.GridViewMode {
		want = []string{
			`<title hx-swap-oob="true">Cookbooks | Recipya</title>`,
			`<div class="p-2 hover:bg-red-600 hover:text-white" title="Display as list" hx-get="/cookbooks?view=list" hx-target="#content">`,
			`<section id="cookbook-1" class="cookbook card card-compact bg-base-100 shadow-lg indicator">`,
			`<button class="btn btn-block btn-sm btn-outline" hx-get="/cookbooks/1?page=1" hx-target="#content" hx-push-url="/cookbooks/1">Open</button>`,
			`<button class="btn btn-block btn-sm btn-outline" hx-get="/cookbooks/2?page=1" hx-target="#content" hx-push-url="/cookbooks/2">Open</button>`,
			`<button class="btn btn-block btn-sm btn-outline" hx-get="/cookbooks/3?page=1" hx-target="#content" hx-push-url="/cookbooks/3">Open</button>`,
		}
	} else {
		want = []string{
			`<title hx-swap-oob="true">Cookbooks | Recipya</title>`,
			`<div class="p-2 bg-blue-600 text-white" title="Display as list">`,
			`<div class="p-2 hover:bg-red-600 hover:text-white" title="Display as grid" hx-get="/cookbooks?view=grid" hx-target="#content">`,
			`<div id="cookbook-1" class="cookbook card card-side card-compact bg-base-100 shadow-md max-w-[30rem] border indicator dark:border-slate-600">`,
			`<button class="btn btn-outline btn-sm" hx-get="/cookbooks/1?page=1" hx-target="#content" hx-push-url="/cookbooks/1">Open</button>`,
			`<button class="btn btn-outline btn-sm" hx-get="/cookbooks/2?page=1" hx-target="#content" hx-push-url="/cookbooks/2">Open</button>`,
			`<button class="btn btn-outline btn-sm" hx-get="/cookbooks/3?page=1" hx-target="#content" hx-push-url="/cookbooks/3">Open</button>`,
		}
	}

	assertStringsInHTML(tb, got, want)
	notWant := []string{
		`Your cookbooks collection looks a bit empty at the moment.`,
		`<p>Why not start by creating a cookbook by clicking the <a class="underline font-semibold cursor-pointer" hx-post="/cookbooks" hx-prompt="Enter the name of your cookbook" hx-target="#content" hx-push-url="true">Add cookbook</a> button at the top? </p>`,
	}
	assertStringsNotInHTML(tb, got, notWant)
}

func assertHeader(tb testing.TB, rr *httptest.ResponseRecorder, key, value string) {
	tb.Helper()
	got := rr.Result().Header.Get(key)
	if got != value {
		tb.Fatalf("expected header %s to be %s but got %s", key, value, got)
	}
}

func assertImage(tb testing.TB, got, want uuid.UUID) {
	tb.Helper()
	if got != want {
		tb.Fatalf("got image %q but expected %q", got, want)
	}
}

func assertImageNotNil(tb testing.TB, got uuid.UUID) {
	tb.Helper()
	if got == uuid.Nil {
		tb.Fatal("got nil image when expected not nil")
	}
}

func assertMustBeLoggedIn(tb testing.TB, srv *server.Server, method string, uri string) {
	tb.Helper()
	rr := sendRequest(srv, method, uri, noHeader, nil)

	assertStatus(tb, rr.Code, http.StatusSeeOther)
	assertHeader(tb, rr, "Location", "/auth/login")
}

func assertStatus(tb testing.TB, got, want int) {
	tb.Helper()
	if got != want {
		tb.Fatalf("expected status %d but got %d", want, got)
	}
}

func assertStringsInHTML(tb testing.TB, bodyHTML string, wants []string) {
	tb.Helper()
	for _, want := range wants {
		if !strings.Contains(bodyHTML, want) {
			tb.Fatalf("string not found in HTML:\n%s", want)
		}
	}
}

func assertStringsNotInHTML(tb testing.TB, bodyHTML string, wants []string) {
	tb.Helper()
	for _, want := range wants {
		if strings.Contains(bodyHTML, want) {
			tb.Fatalf("string found in HTML when it should not:\n%s", want)
		}
	}
}

func assertUploadImageHitCount(tb testing.TB, got, want int) {
	tb.Helper()
	if got != want {
		tb.Fatalf("got %d images uploaded but want %d", got, want)
	}
}

func assertUserSettings(tb testing.TB, got, want *models.UserSettings) {
	tb.Helper()
	if got.ConvertAutomatically != want.ConvertAutomatically {
		tb.Fatalf("settings ConvertAutomatically got %t but want %t", got.ConvertAutomatically, want.ConvertAutomatically)
	}

	if got.CookbooksViewMode != want.CookbooksViewMode {
		tb.Fatalf("settings CookbooksViewMode got %d but want %d", got.CookbooksViewMode, want.CookbooksViewMode)
	}

	if got.MeasurementSystem != want.MeasurementSystem {
		tb.Fatalf("settings MeasurementSystem got %q but want %q", got.MeasurementSystem, want.MeasurementSystem)
	}
}

func assertWebsocket(tb testing.TB, conn *websocket.Conn, message int, want string) {
	tb.Helper()
	mt, got := readMessage(conn, message)
	got = bytes.Join(bytes.Fields(bytes.ReplaceAll(got, []byte("\r\n"), []byte(""))), []byte(" "))
	if mt != websocket.TextMessage || string(got) != want {
		tb.Errorf("got:\n%s\nbut want:\n%s", got, want)
	}
}
