package server_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/server"
	"github.com/reaper47/recipya/internal/services"
)

func TestHandlers_Recipes(t *testing.T) {
	srv := newServerTest()

	const uri = "/recipes"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("user has no recipes", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri)

		want := []string{
			`<title hx-swap-oob="true">Recipes | Recipya</title>`,
			`<p class="pb-2">Your recipe collection looks a bit empty at the moment.</p><p>Why not start adding recipes by clicking the <a class="underline font-semibold cursor-pointer" hx-get="/recipes/add" hx-target="#content" hx-push-url="true">Add recipe</a> button at the top?</p>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("user has recipes", func(t *testing.T) {
		srv.Repository = &mockRepository{
			RecipesRegistered: map[int64]models.Recipes{
				1: {
					{ID: 1, Name: "One", Category: "Lunch", Description: "Recipe one"},
					{ID: 2, Name: "Two", Category: "Soup", Description: "Recipe two"},
					{ID: 3, Name: "Three", Category: "Dinner", Description: "Recipe three"},
				},
			},
		}
		defer func() {
			srv.Repository = &mockRepository{}
		}()

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri)

		got := getBodyHTML(rr)
		assertStringsInHTML(t, got, []string{
			`<title hx-swap-oob="true">Recipes | Recipya</title>`,
			`<form class="w-72 flex md:w-96" hx-get="/recipes/search" hx-vals="{"page": 1}" hx-target="#list-recipes" hx-push-url="true" hx-trigger="submit, change target:.sort-option"><div class="w-full"><label class="input input-bordered input-sm flex justify-between px-0 gap-2 z-20"><button type="button" id="search_shortcut" class="pl-2" popovertarget="search_help" _="on click toggle .hidden on #search_help"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 self-center" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg></button> <input id="search_recipes" class="w-full" type="search" name="q" placeholder="Search for recipes..." value="" _="on keyup if event.target.value !== '' then remove .md:block from #search_shortcut else add .md:block to #search_shortcut then if (event.key is not 'Delete' and not event.key.startsWith('Arrow')) then send submit to closest <form/> then end end"> <button type="submit" class="px-2 btn btn-sm btn-primary"><svg class="w-4 h-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"></path></svg><span class="sr-only">Search</span></button></label></div><div class="dropdown dropdown-left ml-1"><div tabindex="0" role="button" class="btn btn-sm p-1"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12"></path></svg></div><div tabindex="0" class="dropdown-content z-10 menu menu-sm p-2 shadow bg-base-200 w-52 sm:menu-md prose"><h4>Sort</h4><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Default</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="default" checked></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Name:<br>A to Z</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="a-z"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Name:<br>Z to A</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="z-a"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Date created:<br>Newest to oldest</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="new-old"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Date created:<br>Oldest to newest</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="old-new"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Random</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="random"></label></div></div></div></form>`,
			`<div class="hidden absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 items-center justify-center text-white select-none rounded-t-lg sm:flex">`,
			`<img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the One recipe">`,
			`<img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Two recipe">`,
			`<img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Three recipe">`,
		})
		assertStringsNotInHTML(t, got, []string{`Your recipe collection looks a bit empty at the moment.`})
	})
}

func TestHandlers_Recipes_New(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/add"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	testcases := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
	}{
		{name: "is logged in Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody},
		{name: "is logged in no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, uri)

			want := []string{
				`<title hx-swap-oob="true">Add Recipe | Recipya</title>`,
				`<img class="object-cover w-full h-40 rounded-t-xl" src="/static/img/recipes/new/import.webp" alt="Writing on a piece of paper with a traditional pen.">`,
				`<button class="underline" hx-get="/recipes/supported-websites" hx-target="#search-results" onclick="supported_websites_dialog.showModal()">supported</button>`,
				`<dialog id="websites_dialog" class="modal"><div class="modal-box"><form method="dialog"><button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button></form><h3 class="font-bold text-lg">Fetch recipes from websites</h3><form class="py-4" hx-post="/recipes/add/website" hx-swap="none" _="on submit call websites_dialog.close() then set me.querySelector('textarea').value to ''"><div class="grid mb-4"><label class="form-control"><div class="label"><span class="label-text">Enter one or more URLs, each on a new line.</span></div><textarea class="textarea textarea-bordered whitespace-pre-line" placeholder="URL 1URL 2URL 3URL 4etc..." name="urls" rows="5"></textarea></label></div><button class="btn btn-block btn-primary btn-sm">Submit</button></form></div></dialog>`,
				`<dialog id="supported_websites_dialog" class="modal"><div class="modal-box h-2/3"><form method="dialog"><button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button></form><h3 class="mb-1"><label><input type="search" placeholder="Search a website" class="input input-bordered input-sm w-11/12" _="on input show <tbody>tr/> in next <table/> when its textContent.toLowerCase() contains my value.toLowerCase()"></label></h3><div class="overflow-x-auto"><table class="table table-zebra table-sm"><thead><tr class="text-center"><th class="py-1">Number</th><th class="py-1">Website</th></tr></thead> <tbody id="search-results"></tbody></table></div></div></dialog>`,
				`<dialog id="supported_apps_import_dialog" class="modal"><div class="modal-box h-2/3"><form method="dialog"><button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button></form><h3 class="mb-1"><label><input type="search" placeholder="Search an application" class="input input-bordered input-sm w-11/12" _="on input show <tbody>tr/> in next <table/> when its textContent.toLowerCase() contains my value.toLowerCase()"></label></h3><div class="overflow-x-auto"><table class="table table-zebra table-sm"><thead><tr class="text-center"><th class="py-1">Number</th><th class="py-1">Application</th></tr></thead> <tbody id="application-results"></tbody></table></div></div></dialog>`,
				`<dialog id="add_ocr_dialog" class="modal"><div class="modal-box"><form method="dialog"><button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button></form><h3 class="font-bold text-lg">Scan Recipe</h3><form class="py-4" hx-post="/recipes/add/ocr" hx-encoding="multipart/form-data" hx-indicator="#fullscreen-loader" hx-swap="none" _="on submit add_ocr_dialog.close()"><div class="grid mb-4"><label for="add-ocr-files-input" class="text-sm font-medium mb-1">Select your recipe's images ordered by page or a recipe document in the PDF format.</label> <input id="add-ocr-files-input" type="file" name="files" accept=".jpg, .jpeg, .png, .bmp, .tiff, .heif, .pdf" multiple class="p-2 border border-gray-300 rounded-lg shadow focus:ring-2 focus:ring-purple-600 dark:bg-gray-900 dark:border-none"></div><button class="btn btn-block btn-primary btn-sm">Submit</button></form></div></dialog>`,
				`<dialog id="import_recipes_dialog" class="modal"><div class="modal-box"><form method="dialog"><button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button></form><h3 class="font-bold text-lg">Import Recipes</h3><form class="py-4" hx-post="/recipes/add/import" enctype="multipart/form-data" hx-indicator="#fullscreen-loader" hx-swap="none"><div class="grid mb-4"><label for="import-dialog-file" class="text-sm font-semibold mb-1">Choose files in the .json, .txt, .zip or other application format.</label> <input id="import-dialog-file" type="file" name="files" accept=".cml,.crumb,.json,.mxp,.paprikarecipes,.txt,.zip" multiple class="p-2 border border-gray-300 rounded-lg shadow focus:ring-2 focus:ring-purple-600 dark:bg-gray-900 dark:border-none"></div><button type="submit" class="btn btn-block btn-primary btn-sm" onclick="import_recipes_dialog.close()">Submit</button></form></div></dialog>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}

func TestHandlers_Recipes_AddImport(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	uri := ts.URL + "/recipes/add/import"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("payload too big", func(t *testing.T) {
		b := bytes.NewBuffer(make([]byte, 130<<20))
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formData, strings.NewReader(b.String()))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertWebsocket(t, c, 3, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Could not parse the uploaded files.","title":"Form Error"}}`)
	})

	t.Run("error parsing files", func(t *testing.T) {
		contentType, body := createMultipartForm(map[string][]string{"files": {"file1"}})
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertWebsocket(t, c, 3, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Could not retrieve the files or the directory from the form.","title":"Form Error"}}`)
	})

	t.Run("valid request", func(t *testing.T) {
		contentType, body := createMultipartForm(map[string][]string{"files": {"file1.jpg"}})
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusAccepted)
	})
}

func TestHandlers_Recipes_AddManual(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	uri := ts.URL + "/recipes/add/manual"

	repo := &mockRepository{
		categories:        map[int64][]string{1: {"breakfast", "lunch", "dinner"}},
		RecipesRegistered: make(map[int64]models.Recipes),
		UsersRegistered: []models.User{
			{ID: 1, Email: "test@example.com"},
		},
		UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
	}

	resetRepo := func() int {
		repo = &mockRepository{
			categories:        map[int64][]string{1: {"breakfast", "lunch", "dinner"}},
			RecipesRegistered: make(map[int64]models.Recipes),
			UsersRegistered: []models.User{
				{ID: 1, Email: "test@example.com"},
			},
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		srv.Repository = repo
		return len(repo.RecipesRegistered)
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	testcases := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
	}{
		{name: "is logged in Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody},
		{name: "is logged in no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, uri)

			assertStringsInHTML(t, getBodyHTML(rr), []string{
				`<title hx-swap-oob="true">Add Manual | Recipya</title>`,
				`<form class="card-body" style="padding: 0" enctype="multipart/form-data" hx-post="/recipes/add/manual" hx-indicator="#fullscreen-loader">`,
				`<input required type="text" name="title" placeholder="Title of the recipe*" autocomplete="off" class="input w-full btn-ghost text-center">`,
				`<img src="" alt="" class="object-cover mb-2 w-full max-h-[39rem]"> <span class="grid gap-1 max-w-sm" style="margin: auto auto 0.25rem;"><div class="mr-1"><input type="file" accept="image/*,video/*" name="images" class="file-input file-input-sm file-input-bordered w-full max-w-sm" _="on dragover or dragenter halt the event then set the target's style.background to 'lightgray' on dragleave or drop set the target's style.background to '' on drop or change make an FileReader called reader then if event.dataTransfer get event.dataTransfer.files[0] else get event.target.files[0] end then if it.type.startsWith('video') put`,
				`<input type="number" min="1" name="yield" value="1" class="input input-bordered input-sm w-24 md:w-20 lg:w-24">`,
				`<input type="text" list="categories" name="category" class="input input-bordered input-sm w-48 md:w-36 lg:w-48" placeholder="Breakfast" autocomplete="off" value=""> <datalist id="categories"><option>breakfast</option><option>lunch</option><option>dinner</option></datalist>`,
				`<textarea name="description" placeholder="This Thai curry chicken will make you drool." class="textarea w-full h-full resize-none"></textarea>`,
				`<div class="grid grid-flow-col col-span-6 py-1 md:grid-cols-2 md:row-span-1"><div class="flex justify-self-center items-center gap-1 cursor-default" title="Prep time"><svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="24px" height="14px" viewBox="0 0 23 14" version="1.1"><defs><linearGradient id="linear0" gradientUnits="userSpaceOnUse" x1="-125.300003" y1="85.900002" x2="-64.599998" y2="85.900002" gradientTransform="matrix(0.000000000000000013,-0.225806,0.219048,0.000000000000000014,-7.447619,-14.451613)"><stop offset="0.1" style="stop-color:rgb(67.058824%,23.921569%,8.235294%);stop-opacity:1;"></stop> <stop offset="0.5" style="stop-color:rgb(78.431373%,51.372549%,30.588235%);stop-opacity:1;"></stop> <stop offset="0.8" style="stop-color:rgb(90.196078%,77.254902%,53.333333%);stop-opacity:1;"></stop> <stop offset="1" style="stop-color:rgb(94.509804%,87.45098%,62.352941%);stop-opacity:1;"></stop></linearGradient></defs> <g id="surface1"><path style=" stroke:none;fill-rule:evenodd;fill:rgb(95.294118%,82.745099%,64.705884%);fill-opacity:1;" d="M 0 8.128906 L 0.21875 6.324219 C 0.4375 5.644531 0.65625 5.195312 1.3125 4.96875 L 8.542969 2.484375 L 8.542969 1.804688 L 8.980469 0.675781 L 9.855469 0.453125 L 10.734375 0.453125 L 10.953125 0.675781 L 12.265625 1.128906 C 13.003906 1 13.761719 1.078125 14.457031 1.355469 C 15.125 1.453125 15.785156 1.601562 16.429688 1.804688 L 17.523438 1.804688 L 18.617188 0.902344 L 20.589844 0.453125 C 21.246094 0.675781 21.6875 0.902344 21.90625 1.355469 C 22.34375 1.804688 22.5625 2.710938 22.34375 4.066406 L 22.125 4.515625 C 22.5625 4.742188 22.78125 5.195312 22.78125 5.644531 L 22.78125 7.675781 L 22.125 8.804688 L 21.027344 9.484375 L 17.304688 11.289062 L 16.210938 11.742188 L 12.921875 13.324219 L 12.046875 13.773438 L 11.390625 14 L 10.078125 14 L 8.542969 13.546875 C 6.269531 12.441406 4.007812 11.308594 1.753906 10.160156 L 0.65625 9.257812 C 0.21875 9.03125 0 8.582031 0 8.128906 Z M 0 8.128906 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:url(#linear0);" d="M 1.3125 4.742188 L 8.542969 2.03125 L 8.542969 1.582031 L 8.980469 0.453125 C 9.199219 0.226562 9.636719 0 9.855469 0.226562 L 10.953125 0.226562 L 12.265625 0.902344 C 13.003906 0.773438 13.761719 0.851562 14.457031 1.128906 L 15.769531 1.355469 C 16.308594 1.65625 16.933594 1.738281 17.523438 1.582031 L 17.523438 1.355469 C 17.742188 1.128906 18.179688 0.675781 18.617188 0.675781 C 19.277344 0.226562 19.933594 0.226562 20.589844 0.226562 C 21.027344 0.453125 21.6875 0.675781 21.90625 1.128906 C 22.34375 1.582031 22.5625 2.484375 22.34375 3.839844 L 22.125 4.289062 C 22.5625 4.515625 22.78125 4.96875 22.78125 5.417969 L 22.78125 6.546875 L 22.5625 7.453125 L 22.125 8.582031 L 21.027344 9.257812 L 17.085938 11.289062 L 16.210938 11.515625 L 12.921875 13.097656 L 12.046875 13.546875 L 11.390625 13.773438 L 9.855469 13.773438 L 8.324219 13.324219 C 6.121094 12.21875 3.929688 11.089844 1.753906 9.933594 L 1.535156 9.933594 L 0.4375 9.03125 L 0 7.902344 C 0 7.292969 0.0742188 6.6875 0.21875 6.097656 C 0.21875 5.417969 0.65625 4.96875 1.3125 4.742188 Z M 1.3125 4.742188 "></path> <path style="fill:none;stroke-width:0.3;stroke-linecap:butt;stroke-linejoin:miter;stroke:rgb(95.294118%,82.745099%,64.705884%);stroke-opacity:1;stroke-miterlimit:4;" d="M 5.991848 21.001116 L 39.999151 8.995536 L 39.00051 6.00279 L 40.997792 2.006696 L 46.008832 0 L 47.007473 0 L 49.004755 1.003348 L 50.003397 1.003348 L 55.995245 3.996094 C 59.579654 2.923549 63.413723 2.923549 66.998132 3.996094 L 73.007812 6.00279 C 75.272588 6.763951 77.733526 6.763951 79.998302 6.00279 L 84.991508 2.006696 C 88.005265 1.003348 91.001189 0 93.997113 1.003348 C 96.993037 1.003348 99.008152 2.006696 101.005435 3.996094 C 103.002717 6.00279 103.002717 9.998884 102.004076 16.001674 L 102.004076 18.008371 L 104.001359 23.007812 L 105 28.007254 L 104.001359 33.006696 L 101.005435 37.00279 L 95.994395 40.998884 L 78.99966 49.008371 L 75.005095 50.997768 L 74.006454 50.997768 L 58.991168 58.003906 L 54.996603 59.993304 L 52.000679 59.993304 L 51.002038 60.996652 L 46.008832 60.996652 L 39.00051 59.007254 C 28.60394 54.128906 18.278702 49.129464 8.006963 43.991629 L 8.006963 43.00558 L 2.995924 39.995536 L 0 33.992746 C 0 31.294085 0.338825 28.612723 0.998641 26.000558 C 1.997283 23.993862 2.995924 22.004464 5.991848 21.001116 Z M 5.991848 21.001116 " transform="matrix(0.219048,0,0,0.225806,0,0)"></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(25.882354%,9.411765%,1.568628%);fill-opacity:1;" d="M 1.3125 8.128906 L 1.09375 7.675781 L 1.3125 6.097656 L 1.753906 5.644531 L 9.855469 2.933594 L 9.636719 2.257812 C 9.710938 1.882812 9.785156 1.503906 9.855469 1.128906 L 10.078125 1.128906 L 10.515625 1.355469 L 10.734375 1.355469 L 12.265625 2.03125 C 12.914062 1.859375 13.589844 1.859375 14.238281 2.03125 C 14.910156 2.207031 15.570312 2.433594 16.210938 2.710938 L 16.648438 2.710938 L 17.960938 2.484375 L 18.398438 2.03125 L 19.058594 1.582031 L 20.371094 1.355469 L 21.246094 1.804688 L 21.246094 3.839844 L 21.027344 4.066406 L 20.371094 4.515625 L 20.589844 4.515625 L 21.464844 4.96875 L 21.90625 5.417969 L 21.90625 6.324219 L 21.6875 7 L 21.464844 7.675781 L 20.589844 8.128906 C 19.933594 8.582031 18.839844 9.257812 16.867188 9.933594 L 12.484375 11.96875 L 11.828125 12.417969 C 11.617188 12.515625 11.394531 12.589844 11.171875 12.644531 L 10.296875 12.644531 L 8.980469 12.195312 C 6.703125 11.09375 4.441406 9.964844 2.191406 8.804688 Z M 1.3125 8.128906 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(65.490198%,51.372552%,26.274511%);fill-opacity:1;" d="M 6.351562 5.195312 C 8.921875 4.820312 11.476562 4.371094 14.019531 3.839844 L 15.332031 4.289062 L 16.429688 4.515625 C 17.304688 4.515625 17.742188 4.289062 17.960938 4.066406 L 18.179688 3.613281 L 18.617188 3.160156 L 19.933594 2.484375 L 21.027344 2.03125 L 21.027344 3.839844 L 20.808594 4.066406 C 20.371094 4.289062 19.933594 4.515625 19.277344 4.289062 L 17.960938 4.742188 L 19.496094 4.96875 L 21.464844 5.644531 L 21.464844 7.226562 L 21.246094 7.453125 L 20.371094 8.128906 C 17.839844 9.550781 15.203125 10.757812 12.484375 11.742188 L 11.171875 12.417969 C 10.515625 12.417969 9.636719 12.417969 8.980469 11.96875 C 6.324219 10.59375 3.695312 9.164062 1.09375 7.675781 L 1.3125 7 L 1.535156 6.324219 Z M 6.351562 5.195312 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(56.078434%,44.313726%,22.352941%);fill-opacity:1;" d="M 4.382812 7.902344 L 3.503906 7.902344 L 3.285156 9.257812 L 3.066406 9.03125 L 3.285156 7.675781 L 2.847656 7.226562 L 2.628906 7.453125 L 2.410156 8.804688 L 2.191406 8.582031 L 1.972656 8.582031 L 2.410156 7.226562 L 1.972656 7 L 1.753906 8.355469 L 1.3125 8.128906 L 1.535156 6.546875 L 6.351562 6.324219 L 10.078125 8.128906 L 10.953125 10.839844 L 11.171875 12.417969 L 10.296875 12.417969 L 10.515625 10.839844 L 9.417969 10.613281 L 9.199219 12.195312 L 8.542969 11.96875 L 8.761719 10.386719 L 8.105469 10.160156 L 7.886719 11.515625 L 7.449219 11.289062 L 7.449219 9.710938 L 7.230469 9.03125 L 6.570312 9.257812 L 6.132812 10.613281 L 5.476562 10.386719 L 5.914062 9.03125 L 4.820312 8.128906 L 4.601562 8.355469 L 4.382812 9.710938 L 4.160156 9.484375 Z M 4.382812 7.902344 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(52.941179%,41.176471%,18.82353%);fill-opacity:1;" d="M 19.933594 2.484375 L 19.933594 2.710938 C 18.839844 3.160156 18.617188 3.839844 19.496094 4.289062 L 18.617188 4.289062 L 17.960938 4.742188 L 19.496094 4.96875 L 21.464844 5.644531 L 21.464844 7.226562 L 21.246094 7.453125 L 20.371094 8.128906 C 17.839844 9.550781 15.203125 10.757812 12.484375 11.742188 L 11.828125 12.195312 L 10.515625 12.417969 L 10.734375 12.195312 L 10.734375 9.710938 L 10.515625 8.804688 C 13.609375 6.628906 16.75 4.519531 19.933594 2.484375 Z M 19.933594 2.484375 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(47.450981%,36.470589%,16.862746%);fill-opacity:1;" d="M 21.027344 7.453125 C 21.464844 7 21.464844 6.546875 21.464844 5.644531 L 21.464844 7.226562 L 21.246094 7.453125 L 20.371094 8.128906 C 17.839844 9.550781 15.203125 10.757812 12.484375 11.742188 L 11.828125 12.195312 L 10.515625 12.417969 L 10.734375 12.195312 L 11.171875 12.195312 L 12.046875 11.96875 C 15.042969 10.46875 18.035156 8.960938 21.027344 7.453125 Z M 21.027344 7.453125 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(87.450981%,71.764708%,40.784314%);fill-opacity:1;" d="M 1.753906 6.097656 C 5.445312 4.722656 9.167969 3.441406 12.921875 2.257812 L 14.238281 2.484375 L 15.550781 2.933594 L 16.648438 3.160156 C 17.523438 3.160156 17.960938 2.933594 18.179688 2.710938 L 18.617188 2.257812 L 19.058594 2.03125 C 19.714844 1.582031 20.152344 1.582031 20.808594 1.804688 C 21.246094 2.03125 21.246094 2.484375 20.808594 2.710938 L 19.496094 3.160156 L 18.179688 3.613281 L 19.058594 4.289062 C 19.855469 4.75 20.660156 5.203125 21.464844 5.644531 C 21.6875 5.871094 21.246094 6.324219 20.589844 6.773438 C 17.578125 8.257812 14.511719 9.613281 11.390625 10.839844 C 10.734375 11.066406 9.855469 10.839844 8.980469 10.613281 C 6.472656 9.3125 3.988281 7.957031 1.535156 6.546875 L 1.535156 6.097656 Z M 1.753906 6.097656 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(79.215688%,64.313728%,33.333334%);fill-opacity:1;" d="M 2.628906 7.226562 L 2.410156 7.226562 L 4.820312 6.097656 L 6.351562 4.742188 L 8.105469 3.839844 C 9.554688 3.546875 11.015625 3.320312 12.484375 3.160156 C 12.921875 3.160156 13.582031 2.933594 14.019531 2.484375 L 14.457031 2.484375 C 12.945312 3.640625 11.078125 4.203125 9.199219 4.066406 C 8.105469 4.066406 7.230469 4.289062 6.570312 4.742188 L 5.039062 6.097656 C 4.601562 6.546875 3.722656 7 2.628906 7.226562 Z M 2.628906 7.226562 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(79.215688%,64.313728%,33.333334%);fill-opacity:1;" d="M 5.257812 7 C 4.601562 7 3.941406 7.226562 3.503906 7.675781 L 3.285156 7.675781 C 4.5625 6.878906 5.976562 6.34375 7.449219 6.097656 L 10.296875 5.417969 L 11.609375 4.742188 L 12.484375 4.066406 L 14.675781 2.484375 L 14.894531 2.710938 L 12.921875 4.289062 L 10.953125 5.644531 C 9.855469 6.324219 7.886719 6.773438 5.257812 7 Z M 1.972656 6.324219 L 3.066406 5.871094 L 4.160156 5.195312 L 6.132812 4.515625 L 5.476562 4.96875 L 3.722656 6.097656 L 2.191406 6.773438 L 1.972656 6.773438 L 1.535156 6.546875 Z M 9.855469 3.160156 L 11.609375 2.710938 L 13.363281 2.257812 L 13.582031 2.257812 C 13.144531 2.710938 12.484375 2.933594 11.828125 2.933594 Z M 18.617188 7.675781 C 19.277344 6.546875 20.152344 5.871094 21.246094 5.417969 L 21.464844 5.644531 C 20.808594 5.871094 20.152344 6.324219 19.714844 7 Z M 9.199219 7.902344 C 10.078125 7.902344 10.953125 7.675781 12.046875 7 L 14.019531 5.417969 L 15.550781 4.066406 C 16.210938 3.613281 16.648438 3.160156 17.304688 3.160156 L 18.179688 2.710938 L 20.589844 1.804688 L 20.808594 1.804688 L 20.589844 2.03125 L 18.398438 2.933594 L 16.210938 3.839844 L 14.457031 5.417969 C 13.800781 6.324219 13.144531 6.773438 12.703125 7 C 12.265625 7.453125 11.390625 7.902344 10.078125 8.128906 L 7.886719 8.582031 L 6.570312 9.257812 L 5.914062 8.804688 L 7.230469 8.355469 Z M 16.867188 6.097656 C 17.304688 5.195312 17.960938 4.742188 18.839844 4.289062 L 19.496094 4.515625 L 17.742188 5.871094 L 15.992188 7.902344 C 15.113281 8.582031 14.457031 9.03125 13.582031 9.257812 L 10.953125 9.710938 L 9.417969 10.613281 L 8.761719 10.386719 L 10.515625 9.484375 L 12.703125 9.03125 C 14.382812 8.582031 15.851562 7.542969 16.867188 6.097656 Z M 16.867188 6.097656 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(79.215688%,64.313728%,33.333334%);fill-opacity:1;" d="M 6.789062 7.675781 C 5.914062 7.675781 5.257812 7.902344 4.601562 8.355469 L 4.160156 8.128906 C 4.820312 7.675781 5.914062 7.226562 7.230469 7.226562 L 10.953125 6.097656 C 12.046875 5.644531 12.921875 4.96875 13.582031 4.289062 L 15.332031 2.710938 L 15.992188 2.933594 L 14.019531 4.515625 L 12.265625 5.871094 L 9.636719 7.226562 Z M 20.152344 4.742188 L 20.589844 4.96875 L 18.839844 6.324219 C 18.179688 6.773438 17.742188 7.453125 17.304688 8.355469 C 14.960938 9.164062 12.625 9.992188 10.296875 10.839844 L 12.703125 9.933594 L 14.894531 9.03125 C 15.992188 8.582031 17.304688 7.453125 18.617188 5.644531 Z M 13.144531 7.453125 C 14.671875 6.238281 16.203125 5.035156 17.742188 3.839844 C 17.960938 3.386719 19.058594 2.933594 21.027344 2.03125 L 21.027344 2.484375 C 20.402344 2.570312 19.804688 2.800781 19.277344 3.160156 L 18.179688 3.613281 L 18.398438 3.839844 L 15.769531 6.324219 L 13.582031 8.128906 L 10.515625 8.804688 C 9.417969 9.03125 8.761719 9.484375 8.105469 9.933594 L 7.449219 9.710938 C 9.289062 8.8125 11.191406 8.058594 13.144531 7.453125 Z M 6.570312 5.871094 L 6.570312 5.644531 L 6.789062 5.195312 L 7.230469 4.515625 L 8.761719 4.289062 L 10.515625 4.289062 C 10.734375 4.515625 10.734375 4.742188 10.296875 4.96875 L 8.761719 5.644531 L 7.449219 5.871094 L 7.449219 5.195312 L 8.324219 4.742188 L 9.199219 4.515625 L 9.417969 4.742188 L 8.761719 5.195312 L 8.980469 4.96875 L 8.324219 4.96875 L 8.105469 5.195312 C 7.886719 5.417969 8.105469 5.417969 8.324219 5.417969 L 9.199219 5.195312 L 9.855469 4.742188 L 9.855469 4.515625 L 8.761719 4.515625 L 7.449219 4.96875 L 6.789062 5.417969 Z M 6.570312 5.871094 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(41.960785%,25.098041%,14.117648%);fill-opacity:1;" d="M 9.855469 3.160156 L 11.171875 2.710938 L 12.265625 3.839844 L 14.238281 5.417969 L 15.550781 7 L 16.210938 8.582031 L 15.992188 9.484375 L 15.332031 9.710938 L 14.894531 9.710938 L 14.894531 9.257812 L 15.113281 9.03125 L 15.332031 8.582031 L 14.675781 7.675781 L 14.457031 7.453125 L 14.457031 7.226562 L 14.238281 6.773438 L 14.019531 6.773438 C 13.304688 7.003906 12.570312 7.15625 11.828125 7.226562 L 11.171875 7.226562 L 10.734375 6.097656 L 10.078125 4.289062 Z M 9.855469 3.160156 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(47.843137%,47.843137%,47.843137%);fill-opacity:1;" d="M 11.171875 7 L 11.390625 5.871094 L 12.265625 4.96875 L 13.582031 4.96875 L 14.894531 5.195312 L 15.113281 5.871094 L 14.894531 6.097656 L 12.921875 6.773438 L 11.609375 7 Z M 11.171875 7 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(59.215689%,59.215689%,59.215689%);fill-opacity:1;" d="M 12.265625 6.097656 L 12.046875 6.546875 L 12.265625 7 L 11.171875 7 L 11.390625 6.324219 Z M 12.265625 6.097656 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(92.54902%,92.54902%,92.54902%);fill-opacity:1;" d="M 9.855469 1.582031 L 10.734375 1.804688 L 12.921875 2.933594 C 13.75 3.667969 14.488281 4.5 15.113281 5.417969 L 14.894531 5.417969 L 14.019531 4.289062 L 12.484375 2.710938 L 10.734375 1.804688 Z M 9.855469 1.582031 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(78.039217%,78.039217%,78.039217%);fill-opacity:1;" d="M 9.855469 1.582031 L 10.078125 1.582031 L 10.515625 1.804688 C 10.609375 3.429688 11.140625 4.992188 12.046875 6.324219 L 11.171875 7 L 10.953125 6.546875 C 10.421875 5.246094 10.054688 3.882812 9.855469 2.484375 Z M 9.855469 1.582031 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(89.411765%,89.411765%,89.411765%);fill-opacity:1;" d="M 10.078125 3.386719 L 10.078125 2.933594 L 10.734375 2.933594 L 10.734375 3.160156 Z M 9.855469 2.03125 L 10.515625 2.03125 L 10.515625 2.710938 L 9.855469 2.710938 Z M 11.828125 6.097656 L 11.171875 6.773438 L 10.953125 6.324219 L 11.609375 5.871094 Z M 10.296875 4.515625 L 10.078125 3.839844 L 10.734375 3.613281 L 10.953125 4.066406 Z M 11.390625 5.195312 L 10.734375 5.644531 L 10.515625 4.96875 L 10.953125 4.515625 Z M 11.390625 5.195312 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(56.470591%,56.470591%,56.470591%);fill-opacity:1;" d="M 14.238281 6.546875 L 14.019531 6.324219 L 14.019531 5.871094 L 14.675781 5.871094 L 15.113281 6.097656 L 15.113281 6.324219 L 14.894531 6.546875 Z M 14.238281 6.546875 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(70.19608%,70.19608%,70.19608%);fill-opacity:1;" d="M 14.019531 5.871094 L 14.238281 5.644531 L 14.894531 5.644531 L 15.113281 5.871094 L 15.113281 6.097656 L 14.238281 6.097656 Z M 14.019531 5.871094 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(55.686277%,35.686275%,17.254902%);fill-opacity:1;" d="M 14.238281 6.773438 L 14.238281 6.097656 L 14.894531 6.097656 L 15.550781 6.773438 L 16.429688 8.128906 L 16.429688 8.582031 L 15.769531 9.484375 C 15.113281 9.710938 14.675781 9.710938 14.894531 9.257812 L 15.113281 8.582031 L 15.113281 8.128906 L 14.894531 7.675781 L 14.675781 7.453125 L 14.457031 6.773438 Z M 14.238281 6.773438 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(43.137255%,27.058825%,13.725491%);fill-opacity:1;" d="M 15.769531 8.355469 L 16.210938 7.675781 L 16.429688 8.128906 L 16.429688 8.582031 C 16.429688 9.03125 16.210938 9.484375 15.769531 9.484375 L 14.894531 9.484375 L 14.894531 9.03125 L 15.113281 8.582031 L 15.332031 8.355469 Z M 15.769531 8.355469 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(64.313728%,44.705883%,26.666668%);fill-opacity:1;" d="M 14.675781 6.324219 L 14.457031 6.097656 L 14.457031 5.871094 L 15.113281 5.871094 L 15.769531 6.097656 L 16.429688 7.226562 L 16.429688 8.582031 C 15.992188 8.804688 15.550781 9.03125 15.113281 8.804688 L 15.113281 8.355469 L 15.550781 7.902344 L 15.332031 7.453125 L 14.894531 7.226562 L 14.675781 6.773438 Z M 14.675781 6.324219 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(76.078433%,52.156866%,30.980393%);fill-opacity:1;" d="M 15.113281 8.128906 L 15.550781 7.902344 L 15.332031 8.355469 L 15.332031 8.804688 L 15.113281 8.804688 Z M 14.457031 5.871094 L 14.894531 6.546875 L 15.332031 7.453125 L 15.113281 7.453125 L 14.894531 6.773438 L 14.894531 6.546875 L 14.675781 6.546875 L 14.238281 6.097656 Z M 16.429688 7.675781 L 16.429688 7.453125 C 16.648438 7.902344 16.648438 8.355469 16.210938 8.582031 Z M 16.429688 7.675781 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(76.078433%,52.156866%,30.980393%);fill-opacity:1;" d="M 14.894531 6.097656 L 15.332031 6.546875 L 15.550781 6.773438 L 15.550781 7 L 15.769531 7.453125 L 15.769531 8.128906 L 15.550781 8.804688 L 15.550781 7 L 14.675781 5.871094 Z M 14.894531 6.097656 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(76.078433%,52.156866%,30.980393%);fill-opacity:1;" d="M 15.550781 6.324219 L 15.992188 6.546875 L 16.210938 7.226562 L 16.210938 8.128906 L 15.992188 8.804688 L 15.992188 6.546875 C 15.695312 6.324219 15.402344 6.101562 15.113281 5.871094 L 15.332031 5.871094 Z M 15.550781 6.324219 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(76.078433%,52.156866%,30.980393%);fill-opacity:1;" d="M 15.113281 5.871094 L 15.332031 6.324219 L 15.550781 6.773438 L 15.769531 7.226562 L 15.992188 7.453125 L 15.992188 8.128906 L 15.769531 8.804688 L 15.769531 7 L 15.550781 6.773438 L 15.332031 6.324219 L 14.894531 5.871094 Z M 15.332031 5.871094 L 15.769531 6.324219 L 15.992188 6.546875 L 15.550781 6.324219 Z M 15.332031 5.871094 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(43.137255%,27.058825%,13.725491%);fill-opacity:1;" d="M 16.210938 8.355469 C 15.992188 8.582031 15.769531 8.804688 15.550781 8.582031 L 15.332031 8.355469 L 15.992188 8.128906 Z M 16.210938 8.355469 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(69.411767%,69.411767%,69.411767%);fill-opacity:1;" d="M 16.210938 8.355469 L 15.992188 8.582031 L 15.332031 8.582031 L 15.769531 8.355469 Z M 16.210938 8.355469 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(49.019608%,49.019608%,49.019608%);fill-opacity:1;" d="M 16.210938 8.355469 L 15.992188 8.582031 L 15.332031 8.582031 L 15.769531 8.582031 L 15.992188 8.355469 Z M 16.210938 8.355469 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(76.078433%,52.156866%,30.980393%);fill-opacity:1;" d="M 15.992188 6.773438 L 15.769531 6.773438 L 15.992188 7 L 15.992188 6.773438 L 15.992188 7 L 15.769531 7 L 15.769531 6.773438 Z M 15.992188 6.773438 "></path></g></svg><label><input type="text" name="time-preparation" value="00:15:00" class="input input-bordered input-xs max-w-24 html-duration-picker"></label></div><div class="flex justify-self-center items-center gap-1 cursor-default" title="Cooking time"><svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="24px" height="23px" viewBox="0 0 24 23" version="1.1"><g id="surface1"><path style=" stroke:none;fill-rule:nonzero;fill:rgb(62.745098%,64.705882%,65.882353%);fill-opacity:1;" d="M 4.636719 10.984375 L 4.636719 20.417969 C 4.636719 21.007812 5.078125 21.527344 5.667969 21.527344 L 18.257812 21.527344 C 18.847656 21.527344 19.363281 21.007812 19.363281 20.417969 L 19.363281 10.984375 Z M 4.636719 10.984375 "></path> <path style=" stroke:none;fill-rule:nonzero;fill:rgb(76.862745%,76.862745%,76.862745%);fill-opacity:1;" d="M 19.289062 9.953125 C 18.992188 8.699219 17.964844 7.8125 16.710938 7.8125 L 7.214844 7.8125 C 5.964844 7.8125 4.933594 8.699219 4.710938 9.953125 Z M 19.289062 9.953125 "></path> <path style=" stroke:none;fill-rule:nonzero;fill:rgb(54.509804%,69.411765%,81.960784%);fill-opacity:1;" d="M 16.710938 7.8125 L 13.914062 7.8125 C 14.945312 7.8125 15.828125 8.476562 16.269531 9.363281 C 16.417969 9.730469 16.785156 9.953125 17.226562 9.953125 L 19.289062 9.953125 C 18.992188 8.699219 17.964844 7.8125 16.710938 7.8125 "></path> <path style=" stroke:none;fill-rule:nonzero;fill:rgb(0.392157%,26.666667%,38.823529%);fill-opacity:1;" d="M 22.160156 9.953125 L 20.320312 9.953125 C 20.097656 8.183594 18.550781 6.78125 16.710938 6.78125 L 15.535156 6.78125 L 15.3125 5.898438 C 15.09375 5.160156 14.503906 4.644531 13.765625 4.644531 L 11.191406 4.644531 C 10.453125 4.644531 9.792969 5.160156 9.644531 5.898438 L 9.421875 6.78125 L 7.214844 6.78125 C 5.449219 6.78125 3.902344 8.183594 3.605469 9.953125 L 1.765625 9.953125 C 1.03125 9.953125 0.441406 10.542969 0.441406 11.277344 C 0.441406 11.5 0.441406 11.722656 0.589844 11.941406 L 1.25 13.269531 C 1.546875 13.785156 2.0625 14.152344 2.648438 14.152344 L 3.605469 14.152344 L 3.605469 20.417969 C 3.605469 21.597656 4.492188 22.558594 5.667969 22.558594 L 18.257812 22.558594 C 19.4375 22.558594 20.394531 21.597656 20.394531 20.417969 L 20.394531 14.152344 L 21.351562 14.152344 C 21.9375 14.152344 22.453125 13.859375 22.75 13.269531 L 23.410156 11.867188 C 23.558594 11.648438 23.558594 11.5 23.558594 11.277344 C 23.558594 10.542969 22.894531 9.953125 22.160156 9.953125 M 3.605469 13.121094 L 2.648438 13.121094 C 2.429688 13.121094 2.28125 12.972656 2.136719 12.753906 L 1.472656 11.5 L 1.472656 11.277344 C 1.472656 11.132812 1.621094 10.984375 1.765625 10.984375 L 3.605469 10.984375 Z M 10.75 6.117188 C 10.75 5.972656 10.96875 5.75 11.265625 5.75 L 13.839844 5.75 C 14.0625 5.75 14.28125 5.898438 14.355469 6.117188 L 14.503906 6.78125 L 10.601562 6.78125 Z M 7.214844 7.8125 L 16.710938 7.8125 C 17.964844 7.8125 18.992188 8.699219 19.289062 9.953125 L 4.710938 9.953125 C 4.933594 8.699219 5.964844 7.8125 7.214844 7.8125 M 19.363281 13.636719 L 19.363281 20.417969 C 19.363281 21.007812 18.847656 21.527344 18.257812 21.527344 L 5.667969 21.527344 C 5.078125 21.527344 4.636719 21.007812 4.636719 20.417969 L 4.636719 10.984375 L 19.363281 10.984375 Z M 22.453125 11.5 L 21.71875 12.828125 C 21.644531 12.972656 21.496094 13.121094 21.277344 13.121094 L 20.394531 13.121094 L 20.394531 11.058594 L 22.160156 11.058594 C 22.308594 11.058594 22.453125 11.207031 22.453125 11.351562 L 22.453125 11.5 "></path> <path style=" stroke:none;fill-rule:nonzero;fill:rgb(92.54902%,94.117647%,97.647059%);fill-opacity:1;" d="M 7.804688 17.324219 C 8.089844 17.324219 8.320312 17.554688 8.320312 17.839844 C 8.320312 18.125 8.089844 18.355469 7.804688 18.355469 C 7.519531 18.355469 7.289062 18.125 7.289062 17.839844 C 7.289062 17.554688 7.519531 17.324219 7.804688 17.324219 M 7.804688 16.808594 C 7.4375 16.808594 7.214844 16.585938 7.214844 16.21875 L 7.214844 13.121094 C 7.214844 12.753906 7.4375 12.53125 7.804688 12.53125 C 8.097656 12.53125 8.320312 12.753906 8.320312 13.121094 L 8.320312 16.21875 C 8.320312 16.585938 8.097656 16.808594 7.804688 16.808594 "></path> <path style=" stroke:none;fill-rule:nonzero;fill:rgb(54.509804%,69.411765%,81.960784%);fill-opacity:1;" d="M 16.195312 12.015625 L 16.195312 19.902344 C 16.195312 20.492188 15.679688 21.007812 15.09375 21.007812 L 6.699219 21.007812 C 6.183594 21.007812 5.667969 20.492188 5.667969 19.902344 L 5.667969 12.015625 C 5.667969 11.5 5.226562 10.984375 4.636719 10.984375 L 4.636719 20.417969 C 4.636719 21.007812 5.078125 21.527344 5.667969 21.527344 L 18.257812 21.527344 C 18.847656 21.527344 19.363281 21.007812 19.363281 20.417969 L 19.363281 10.984375 L 17.226562 10.984375 C 16.636719 10.984375 16.195312 11.5 16.195312 12.015625 M 8.246094 7.8125 L 7.214844 7.8125 C 5.964844 7.8125 4.933594 8.699219 4.710938 9.953125 L 4.933594 9.953125 C 5.375 9.953125 5.742188 9.730469 5.890625 9.363281 C 6.332031 8.476562 7.214844 7.8125 8.246094 7.8125 "></path> <path style=" stroke:none;fill-rule:nonzero;fill:rgb(0.392157%,26.666667%,38.823529%);fill-opacity:1;" d="M 18.773438 5.75 C 18.625 5.75 18.480469 5.675781 18.40625 5.527344 C 18.183594 5.308594 18.257812 4.9375 18.40625 4.792969 C 18.699219 4.644531 18.773438 4.421875 18.773438 4.128906 C 18.773438 3.90625 18.699219 3.6875 18.480469 3.539062 C 18.257812 3.316406 18.257812 3.023438 18.40625 2.800781 C 18.625 2.582031 18.992188 2.507812 19.140625 2.726562 C 19.582031 3.097656 19.878906 3.613281 19.878906 4.128906 C 19.878906 4.71875 19.582031 5.234375 19.140625 5.601562 L 18.773438 5.75 M 18.773438 1.03125 C 19.058594 1.03125 19.289062 1.261719 19.289062 1.546875 C 19.289062 1.832031 19.058594 2.0625 18.773438 2.0625 C 18.488281 2.0625 18.257812 1.832031 18.257812 1.546875 C 18.257812 1.261719 18.488281 1.03125 18.773438 1.03125 M 16.710938 5.75 L 16.269531 5.527344 C 16.050781 5.308594 16.121094 4.9375 16.34375 4.792969 C 16.5625 4.644531 16.710938 4.421875 16.710938 4.128906 C 16.710938 3.90625 16.5625 3.6875 16.417969 3.539062 C 15.902344 3.097656 15.679688 2.652344 15.679688 2.0625 C 15.679688 1.472656 15.902344 1.03125 16.417969 0.589844 C 16.5625 0.367188 16.933594 0.441406 17.152344 0.664062 C 17.300781 0.8125 17.300781 1.179688 17.078125 1.402344 C 16.785156 1.546875 16.710938 1.769531 16.710938 2.0625 C 16.710938 2.285156 16.785156 2.507812 17.007812 2.652344 C 17.519531 3.097656 17.742188 3.613281 17.742188 4.128906 C 17.742188 4.71875 17.519531 5.234375 17.007812 5.601562 L 16.710938 5.75 "></path></g></svg><label><input type="text" name="time-cooking" value="00:30:00" class="input input-bordered input-xs max-w-24 html-duration-picker"></label></div></div>`,
				`<table class="table table-zebra table-xs"><thead><tr><th>Nutrition<br>(per 100g)</th><th>Amount</th></tr></thead> <tbody><tr><td>Calories</td><td><label><input type="text" name="calories" autocomplete="off" placeholder="368kcal" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Total carbs</td><td><label><input type="text" name="total-carbohydrates" autocomplete="off" placeholder="35g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Sugars</td><td><label><input type="text" name="sugars" autocomplete="off" placeholder="3g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Protein</td><td><label><input type="text" name="protein" autocomplete="off" placeholder="21g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Total fat</td><td><label><input type="text" name="total-fat" autocomplete="off" placeholder="15g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Saturated fat</td><td><label><input type="text" name="saturated-fat" autocomplete="off" placeholder="1.8g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Unsaturated fat</td><td><label><input type="text" name="unsaturated-fat" autocomplete="off" placeholder="1.8g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Trans fat</td><td><label><input type="text" name="trans-fat" autocomplete="off" placeholder="1.8g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Cholesterol</td><td><label><input type="text" name="cholesterol" autocomplete="off" placeholder="1.1mg" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Sodium</td><td><label><input type="text" name="sodium" autocomplete="off" placeholder="100mg" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Fiber</td><td><label><input type="text" name="fiber" autocomplete="off" placeholder="8g" class="input input-bordered input-xs max-w-24"></label></td></tr></tbody></table>`,
				`<ol id="tools-list" class="pl-4 list-decimal"><li class="pb-2"><div class="grid grid-flow-col items-center"><label><input type="text" name="tools" placeholder="1 frying pan" class="input input-bordered input-sm w-full" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)"></label><div class="ml-2"><button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: Enter" onclick="addItem(event)">+</button> <button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" _="on click if (closest <ol/>).childElementCount > 1 remove closest <li/> else set input to (closest <li/>).querySelector('input') then set input.value to '' then input.focus()">-</button><div class="inline-block h-4 cursor-move handle ml-2"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path></svg></div></div></div></li></ol>`,
				`<ol id="ingredients-list" class="pl-4 list-decimal"><li class="pb-2"><div class="grid grid-flow-col items-center"><label><input required type="text" name="ingredients" value="" placeholder="1 cup of chopped onions" class="input input-bordered input-sm w-full" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)"></label><div class="ml-2"><button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: Enter" onclick="addItem(event)">+</button> <button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" _="on click if (closest <ol/>).childElementCount > 1 remove closest <li/> else set input to (closest <li/>).querySelector('input') then set input.value to '' then input.focus()">-</button><div class="inline-block h-4 cursor-move handle ml-2"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path></svg></div></div></div></li></ol></div><div class="col-span-6 px-6 py-2 border-gray-700 md:rounded-bl-none md:col-span-4"><h2 class="font-semibold text-center pb-2"><span class="underline">Instructions</span> <sup class="text-red-600">*</sup></h2><ol id="instructions-list" class="grid list-decimal"><li class="pt-2 md:pl-0"><div class="flex"><label class="w-11/12"><textarea required name="instructions" rows="3" class="textarea textarea-bordered w-full" placeholder="Mix all ingredients together" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)"></textarea></label><div class="grid ml-2"><button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: CTRL + Enter" onclick="addItem(event)">+</button> <button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" _="on click if (closest <ol/>).childElementCount > 1 remove closest <li/> else set input to (closest <li/>).querySelector('textarea') then set input.value to '' then input.focus()">-</button><div class="h-4 cursor-move handle grid place-content-center"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path></svg></div></div></div></li></ol>`,
				`<div class="col-span-6 px-6 py-2 border-gray-700 md:rounded-bl-none md:col-span-4"><h2 class="font-semibold text-center pb-2"><span class="underline">Instructions</span> <sup class="text-red-600">*</sup></h2>`,
				`<ol id="instructions-list" class="grid list-decimal"><li class="pt-2 md:pl-0"><div class="flex"><label class="w-11/12"><textarea required name="instructions" rows="3" class="textarea textarea-bordered w-full" placeholder="Mix all ingredients together" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)"></textarea></label><div class="grid ml-2"><button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: CTRL + Enter" onclick="addItem(event)">+</button> <button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" _="on click if (closest <ol/>).childElementCount > 1 remove closest <li/> else set input to (closest <li/>).querySelector('textarea') then set input.value to '' then input.focus()">-</button><div class="h-4 cursor-move handle grid place-content-center"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path></svg></div></div></div></li></ol>`,
				`<button class="btn btn-primary btn-block btn-sm">Submit</button>`,
			})
		})
	}

	t.Run("missing source defaults to unknown", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusCreated)
		if repo.RecipesRegistered[1][0].URL != "Unknown" {
			t.Fatalf("got source %q; want 'unknown'", repo.RecipesRegistered[1][0].URL)
		}
	})

	t.Run("missing category defaults to uncategorized", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusCreated)
		if repo.RecipesRegistered[1][0].Category != "uncategorized" {
			t.Fatalf("got category %q; want 'uncategorized'", repo.RecipesRegistered[1][0].Category)
		}
	})

	t.Run("missing yield defaults to 1", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusCreated)
		if repo.RecipesRegistered[1][0].Yield != 1 {
			t.Fatalf("got yield %d; want 1", repo.RecipesRegistered[1][0].Yield)
		}
	})

	t.Run("can only be one category", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"category":     {"dinner,lunch"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusCreated)
		if repo.RecipesRegistered[1][0].Category != "dinner" {
			t.Fatalf("got category %s; want dinner", repo.RecipesRegistered[1][0].Category)
		}
	})

	t.Run("subcategories are possible", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"category":     {"beverages:cocktails:vodka"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusCreated)
		if repo.RecipesRegistered[1][0].Category != "beverages:cocktails:vodka" {
			t.Fatalf("got category %s; want beverages:cocktails:vodka", repo.RecipesRegistered[1][0].Category)
		}
	})

	t.Run("submit recipe", func(t *testing.T) {
		_ = resetRepo()

		contentType, body := createMultipartForm(map[string][]string{
			"title":               {"Salsa"},
			"images":              {"eggs.jpg"},
			"category":            {"appetizers"},
			"source":              {"Mommy"},
			"description":         {"The best"},
			"calories":            {"666"},
			"total-carbohydrates": {"31g"},
			"sugars":              {"0.1mg"},
			"protein":             {"5g"},
			"total-fat":           {"0g"},
			"saturated-fat":       {"0g"},
			"unsaturated-fat":     {"1g"},
			"trans-fat":           {"2g"},
			"cholesterol":         {"256mg"},
			"sodium":              {"777mg"},
			"fiber":               {"2g"},
			"time-preparation":    {"00:15:30"},
			"time-cooking":        {"00:30:15"},
			"tools":               {"wok", "3 pan"},
			"ingredients":         {"ing1", "ing2"},
			"instructions":        {"ins1", "ins2"},
			"keywords":            {"meow"},
			"yield":               {"4"},
		})
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusCreated)
		id := int64(len(repo.RecipesRegistered))
		gotRecipe := repo.RecipesRegistered[1][id-1]
		want := models.Recipe{
			Category:     "appetizers",
			Description:  "The best",
			ID:           1,
			Images:       gotRecipe.Images,
			Ingredients:  []string{"ing1", "ing2"},
			Instructions: []string{"ins1", "ins2"},
			Keywords:     []string{"meow"},
			Name:         "Salsa",
			Nutrition: models.Nutrition{
				Calories:           "666",
				Cholesterol:        "256mg",
				Fiber:              "2g",
				Protein:            "5g",
				SaturatedFat:       "0g",
				Sodium:             "777mg",
				Sugars:             "0.1mg",
				TotalCarbohydrates: "31g",
				TotalFat:           "0g",
				UnsaturatedFat:     "1g",
				TransFat:           "2g",
			},
			Times: models.Times{
				Prep:  15*time.Minute + 30*time.Second,
				Cook:  30*time.Minute + 15*time.Second,
				Total: 45*time.Minute + 45*time.Second,
			},
			Tools: []models.HowToItem{
				{Type: "HowToTool", Text: "wok", Quantity: 1},
				{Type: "HowToTool", Text: "pan", Quantity: 3},
			},
			URL:   "Mommy",
			Yield: 4,
		}
		if len(gotRecipe.Images) == 0 {
			t.Fatal("got no images when want images")
		}
		if !cmp.Equal(want, gotRecipe) {
			t.Log(cmp.Diff(want, gotRecipe))
			t.Fail()
		}
		assertHeader(t, rr, "HX-Redirect", "/recipes/"+strconv.FormatInt(id, 10))
	})
}

func TestHandlers_Recipes_AddOCR(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	originalIntegrations := srv.Integrations
	originalRepo := srv.Repository
	app.Config.Integrations.AzureDI.Key = "chicken"
	app.Config.Integrations.AzureDI.Endpoint = "kyiv"

	uri := ts.URL + "/recipes/add/ocr"

	sendReq := func(files string) *httptest.ResponseRecorder {
		contentType, body := createMultipartForm(map[string][]string{"files": {files}})
		return sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("warn user when feature disabled", func(t *testing.T) {
		original := app.Config
		app.Config.Integrations.AzureDI.Key = ""
		app.Config.Integrations.AzureDI.Endpoint = ""
		defer func() {
			app.Config = original
		}()

		rr := sendReq("hello.jpg")

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-warning","message":"Please consult the docs to enable OCR.","title":"Feature Disabled"}}`)
	})

	t.Run("files must not be empty", func(t *testing.T) {
		rr := sendReq("")

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Could not retrieve the image from the form.","title":"Form Error"}}`)
	})

	t.Run("processing OCR failed", func(t *testing.T) {
		srv.Integrations = &mockIntegrations{
			processImageOCRFunc: func(_ []io.Reader) (models.Recipes, error) {
				return models.Recipes{}, errors.New("error")
			},
		}
		defer func() {
			srv.Integrations = originalIntegrations
		}()

		rr := sendReq("hello.jpg")

		assertStatus(t, rr.Code, http.StatusAccepted)
		assertWebsocket(t, c, 3, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Could not process OCR.","title":"Integrations Error"}}`)
	})

	t.Run("adding processed recipe failed", func(t *testing.T) {
		srv.Repository = &mockRepository{
			AddRecipesFunc: func(_ models.Recipes, _ int64, _ chan models.Progress) ([]int64, []models.ReportLog, error) {
				return nil, nil, errors.New("oops")
			},
			UsersRegistered:        originalRepo.Users(),
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendReq("hello.jpg")

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Recipes could not be added.","title":"Database Error"}}`
		assertWebsocket(t, c, 3, want)
	})

	t.Run("valid request", func(t *testing.T) {
		repo := &mockRepository{
			RecipesRegistered:      make(map[int64]models.Recipes),
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		srv.Repository = repo
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendReq("hello.jpg")

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"View /recipes/1","background":"alert-info","message":"Recipe scanned and uploaded.","title":"Operation Successful"}}`
		assertWebsocket(t, c, 3, want)
		if len(repo.RecipesRegistered[1]) != 1 && repo.RecipesRegistered[1][0].ID != 1 {
			t.Fatal("expected the recipe to be added")
		}
	})
}

func TestHandlers_Recipes_AddWebsite(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	originalRepo := srv.Repository
	originalScraper := srv.Scraper
	originalFiles := srv.Files

	uri := ts.URL + "/recipes/add/website"

	prepare := func() *mockRepository {
		repo := &mockRepository{
			RecipesRegistered:      map[int64]models.Recipes{1: make(models.Recipes, 0)},
			Reports:                map[int64][]models.Report{1: make([]models.Report, 0)},
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		srv.Repository = repo
		return repo
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("no input", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("websites="))

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("no valid URLs", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("urls=I am a pig\noink oink"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"No valid URLs found.","title":"Request Error"}}`)
	})

	t.Run("add one valid URL from unsupported websites", func(t *testing.T) {
		repo := prepare()
		srv.Files = &mockFiles{}
		srv.Scraper = &mockScraper{
			scraperFunc: func(_ string, _ services.FilesService) (models.RecipeSchema, error) {
				return models.RecipeSchema{}, errors.New("unsupported website")
			},
		}
		defer func() {
			srv.Repository = originalRepo
			srv.Scraper = originalScraper
			srv.Files = originalFiles
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("urls=https://www.example.com"))

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"View /reports?view=latest","background":"alert-error","message":"Fetching the recipe failed.","title":"Operation Failed"}}`
		assertWebsocket(t, c, 4, want)
		if len(repo.Reports[1]) != 1 {
			t.Fatalf("got reports %v but want one report added", repo.Reports[1])
		}
	})

	t.Run("add one valid URL from supported websites", func(t *testing.T) {
		repo := prepare()
		defer func() {
			srv.Repository = repo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("urls=https://www.example.com"))

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"View /recipes/1","background":"alert-info","message":"Recipe has been added to your collection.","title":"Operation Successful"}}`
		assertWebsocket(t, c, 4, want)
		if len(repo.Reports[1]) != 1 {
			t.Fatalf("got reports %v but want one report added", repo.Reports[1])
		}
		if len(repo.RecipesRegistered[1]) != 1 {
			t.Fatal("expected 3 recipes")
		}
	})

	t.Run("add duplicates", func(t *testing.T) {
		repo := prepare()
		defer func() {
			srv.Repository = repo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("urls=https://www.example.com\nhttps://www.example.com"))

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"View /recipes/1","background":"alert-info","message":"Recipe has been added to your collection.","title":"Operation Successful"}}`
		assertWebsocket(t, c, 4, want)
		if len(repo.Reports[1]) != 1 {
			t.Fatalf("got reports %v but want one report added", repo.Reports[1])
		}
		if len(repo.RecipesRegistered[1]) != 1 {
			t.Fatal("expected 3 recipes")
		}
	})

	t.Run("add many valid URLs from supported websites", func(t *testing.T) {
		repo := prepare()
		defer func() {
			srv.Repository = repo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("urls=https://www.example.com\nhttps://www.hello.com\nhttp://helloiam.bob.com\njesus.com"))

		assertStatus(t, rr.Code, http.StatusAccepted)
		assertWebsocket(t, c, 6, `{"type":"toast","fileName":"","data":"","toast":{"action":"View /reports?view=latest","background":"alert-info","message":"Fetched 3 recipes. 0 skipped","title":"Operation Successful"}}`)
		if len(repo.Reports[1]) != 1 {
			t.Fatalf("got reports %v but want one report added", repo.Reports[1])
		}
		if len(repo.RecipesRegistered[1]) != 3 {
			t.Fatal("expected 3 recipes")
		}
	})
}

func TestHandlers_Recipes_Delete(t *testing.T) {
	repo := &mockRepository{
		RecipesRegistered: map[int64]models.Recipes{1: make(models.Recipes, 0)},
		UsersRegistered:   []models.User{{ID: 1, Email: "test@example.com"}},
	}
	srv := newServerTest()
	srv.Repository = repo

	uri := "/recipes"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodDelete, uri+"/5")
	})

	t.Run("cannot delete recipe that does not exist", func(t *testing.T) {
		numRecipesBefore := len(repo.RecipesRegistered)
		defer func() {
			srv.Repository = repo
		}()

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodDelete, uri+"/5")

		assertStatus(t, rr.Code, http.StatusNoContent)
		if numRecipesBefore != len(repo.RecipesRegistered) {
			t.Fail()
		}
	})

	t.Run("can delete user's recipe", func(t *testing.T) {
		_, _, _ = srv.Repository.AddRecipes(models.Recipes{{ID: 1, Name: "Chicken"}}, 1, nil)

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodDelete, uri+"/1")

		assertStatus(t, rr.Code, http.StatusNoContent)
		assertHeader(t, rr, "HX-Redirect", "/")
	})
}

func TestHandlers_Recipes_Edit(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	baseRecipe := models.Recipe{
		Category:     "american",
		CreatedAt:    time.Now(),
		Cuisine:      "indonesian",
		Description:  "A delicious recipe!",
		ID:           1,
		Images:       []uuid.UUID{uuid.New()},
		Ingredients:  []string{"ing1", "ing2", "ing3"},
		Instructions: []string{"ins1", "ins2", "ins3"},
		Keywords:     []string{"chicken", "big", "marinade"},
		Name:         "Chicken Jersey",
		Nutrition: models.Nutrition{
			Calories:           "354",
			Cholesterol:        "1",
			Fiber:              "2",
			Protein:            "3",
			SaturatedFat:       "4",
			Sodium:             "5",
			Sugars:             "6",
			TotalCarbohydrates: "7",
			TotalFat:           "8",
			UnsaturatedFat:     "9",
			TransFat:           "10",
		},
		Times: models.Times{
			Prep:  30 * time.Minute,
			Cook:  1 * time.Hour,
			Total: 1*time.Hour + 30*time.Minute,
		},
		Tools: []models.HowToItem{
			models.NewHowToTool("spoons"),
			models.NewHowToTool("drill"),
		},
		UpdatedAt: time.Now(),
		URL:       "https://example.com/recipes/yummy",
		Yield:     12,
	}

	xc, _ := srv.Repository.Categories(1)
	repo := &mockRepository{
		categories:        map[int64][]string{1: xc},
		RecipesRegistered: map[int64]models.Recipes{1: {baseRecipe}},
	}
	srv.Repository = repo
	originalRepo := srv.Repository

	resetRepo := func() int {
		repo = &mockRepository{
			categories:             map[int64][]string{1: xc},
			RecipesRegistered:      map[int64]models.Recipes{1: {baseRecipe}},
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		srv.Repository = repo
		return len(repo.RecipesRegistered)
	}

	uri := ts.URL + "/recipes/%d/edit"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, fmt.Sprintf(uri, 5))
		assertMustBeLoggedIn(t, srv, http.MethodPut, fmt.Sprintf(uri, 5))
	})

	t.Run("error fetching recipe", func(t *testing.T) {
		srv.Repository = &mockRepository{
			RecipeFunc: func(id, userID int64) (*models.Recipe, error) {
				return nil, errors.New("oops")
			},
			RecipesRegistered: map[int64]models.Recipes{1: {baseRecipe}},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, fmt.Sprintf(uri, 1))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to retrieve recipe.","title":"Database Error"}}`)
	})

	t.Run("successful request", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, fmt.Sprintf(uri, 1))

		assertStringsInHTML(t, getBodyHTML(rr), []string{
			`<title hx-swap-oob="true">Edit Chicken Jersey | Recipya</title>`,
			`<input required type="text" name="title" placeholder="Title of the recipe*" autocomplete="off" class="input w-full btn-ghost text-center" value="Chicken Jersey">`,
			`<div id="media-container" class="grid grid-flow-col grid-cols-7 w-full text-center border-gray-700 md:grid-cols-6 md:col-span-3 md:border-r"><div class="buttons-container flex flex-col gap-1 col-span-2 md:col-span-1 p-1"><button id="media-button-1" type="button" class="btn btn-sm btn-ghost btn-active" onclick="switchMedia(event)">Media 1</button> <button id="add-media-button" type="button" class="btn btn-sm btn-ghost" onclick="addMedia(event)"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" width="24px" height="24px" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle> <line x1="12" y1="8" x2="12" y2="16"></line> <line x1="8" y1="12" x2="16" y2="12"></line></svg>Add</button></div><div id="media" class="col-span-5"><label id="media-1" class=""><img alt="" class="object-cover mb-2 w-full max-h-[39rem]" src=""> <span class="grid gap-1 max-w-sm" style="margin: auto auto 0.25rem;"><div class="mr-1 hidden"><input type="file" accept="image/*,video/*" name="images" class="file-input file-input-sm file-input-bordered w-full max-w-sm" value="/data/images/` + baseRecipe.Images[0].String() + `.webp" _="on dragover or dragenter halt the event then set the target's style.background to 'lightgray' on dragleave or drop set the target's style.background to '' on drop or change make an FileReader called reader then if event.dataTransfer get event.dataTransfer.files[0] else get event.target.files[0] end then if it.type.startsWith('video')`,
			`after previous <img/> then add .hidden to previous <img/> else set {src: window.URL.createObjectURL(it)} on previous <img/> end then remove .hidden from me.parentElement.parentElement.querySelectorAll('button') then add .hidden to the parentElement of me"><div class="divider">OR</div><span class="hidden input-error"></span><div class="flex"><input type="url" placeholder="Enter the URL of an image" class="input input-bordered input-sm w-full max-w-sm mr-1"> <button type="button" class="btn btn-sm" hx-get="/fetch" hx-vals="js:{url: event.target.previousElementSibling.value}" hx-swap="none" _="on htmx:afterRequest if event.detail.successful then set a to first in event.target.parentElement.parentElement.children then call updateMediaFromFetch(a, event.detail.xhr.responseURL) end">Fetch</button></div><div _="on load if not navigator.clipboard hide me"><div class="divider">OR</div><button type="button" class="btn btn-sm" onclick="pasteImage(event)">Paste copied image</button></div></div><button type="button" class="btn btn-sm btn-error btn-outline hidden" onclick="deleteMedia(event)">Delete</button></span></label> </div>`,
			`<input type="text" list="categories" name="category" class="input input-bordered input-sm w-48 md:w-36 lg:w-48" placeholder="Breakfast" autocomplete="off" value="american"> <datalist id="categories"><option>breakfast</option><option>lunch</option><option>dinner</option></datalist>`,
			`<input type="number" min="1" name="yield" class="input input-bordered input-sm w-24 md:w-20 lg:w-24" value="12">`,
			`<input type="text" placeholder="Source" name="source" class="input input-bordered input-sm md:w-28 lg:w-40 xl:w-44" value="https://example.com/recipes/yummy"`,
			`<textarea name="description" placeholder="This Thai curry chicken will make you drool." class="textarea w-full h-full resize-none">A delicious recipe!</textarea>`,
			`<div class="flex justify-self-center items-center gap-1 cursor-default" title="Prep time"><svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="24px" height="14px" viewBox="0 0 23 14" version="1.1"><defs><linearGradient id="linear0" gradientUnits="userSpaceOnUse" x1="-125.300003" y1="85.900002" x2="-64.599998" y2="85.900002" gradientTransform="matrix(0.000000000000000013,-0.225806,0.219048,0.000000000000000014,-7.447619,-14.451613)"><stop offset="0.1" style="stop-color:rgb(67.058824%,23.921569%,8.235294%);stop-opacity:1;"></stop> <stop offset="0.5" style="stop-color:rgb(78.431373%,51.372549%,30.588235%);stop-opacity:1;"></stop> <stop offset="0.8" style="stop-color:rgb(90.196078%,77.254902%,53.333333%);stop-opacity:1;"></stop> <stop offset="1" style="stop-color:rgb(94.509804%,87.45098%,62.352941%);stop-opacity:1;"></stop></linearGradient></defs> <g id="surface1"><path style=" stroke:none;fill-rule:evenodd;fill:rgb(95.294118%,82.745099%,64.705884%);fill-opacity:1;" d="M 0 8.128906 L 0.21875 6.324219 C 0.4375 5.644531 0.65625 5.195312 1.3125 4.96875 L 8.542969 2.484375 L 8.542969 1.804688 L 8.980469 0.675781 L 9.855469 0.453125 L 10.734375 0.453125 L 10.953125 0.675781 L 12.265625 1.128906 C 13.003906 1 13.761719 1.078125 14.457031 1.355469 C 15.125 1.453125 15.785156 1.601562 16.429688 1.804688 L 17.523438 1.804688 L 18.617188 0.902344 L 20.589844 0.453125 C 21.246094 0.675781 21.6875 0.902344 21.90625 1.355469 C 22.34375 1.804688 22.5625 2.710938 22.34375 4.066406 L 22.125 4.515625 C 22.5625 4.742188 22.78125 5.195312 22.78125 5.644531 L 22.78125 7.675781 L 22.125 8.804688 L 21.027344 9.484375 L 17.304688 11.289062 L 16.210938 11.742188 L 12.921875 13.324219 L 12.046875 13.773438 L 11.390625 14 L 10.078125 14 L 8.542969 13.546875 C 6.269531 12.441406 4.007812 11.308594 1.753906 10.160156 L 0.65625 9.257812 C 0.21875 9.03125 0 8.582031 0 8.128906 Z M 0 8.128906 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:url(#linear0);" d="M 1.3125 4.742188 L 8.542969 2.03125 L 8.542969 1.582031 L 8.980469 0.453125 C 9.199219 0.226562 9.636719 0 9.855469 0.226562 L 10.953125 0.226562 L 12.265625 0.902344 C 13.003906 0.773438 13.761719 0.851562 14.457031 1.128906 L 15.769531 1.355469 C 16.308594 1.65625 16.933594 1.738281 17.523438 1.582031 L 17.523438 1.355469 C 17.742188 1.128906 18.179688 0.675781 18.617188 0.675781 C 19.277344 0.226562 19.933594 0.226562 20.589844 0.226562 C 21.027344 0.453125 21.6875 0.675781 21.90625 1.128906 C 22.34375 1.582031 22.5625 2.484375 22.34375 3.839844 L 22.125 4.289062 C 22.5625 4.515625 22.78125 4.96875 22.78125 5.417969 L 22.78125 6.546875 L 22.5625 7.453125 L 22.125 8.582031 L 21.027344 9.257812 L 17.085938 11.289062 L 16.210938 11.515625 L 12.921875 13.097656 L 12.046875 13.546875 L 11.390625 13.773438 L 9.855469 13.773438 L 8.324219 13.324219 C 6.121094 12.21875 3.929688 11.089844 1.753906 9.933594 L 1.535156 9.933594 L 0.4375 9.03125 L 0 7.902344 C 0 7.292969 0.0742188 6.6875 0.21875 6.097656 C 0.21875 5.417969 0.65625 4.96875 1.3125 4.742188 Z M 1.3125 4.742188 "></path> <path style="fill:none;stroke-width:0.3;stroke-linecap:butt;stroke-linejoin:miter;stroke:rgb(95.294118%,82.745099%,64.705884%);stroke-opacity:1;stroke-miterlimit:4;" d="M 5.991848 21.001116 L 39.999151 8.995536 L 39.00051 6.00279 L 40.997792 2.006696 L 46.008832 0 L 47.007473 0 L 49.004755 1.003348 L 50.003397 1.003348 L 55.995245 3.996094 C 59.579654 2.923549 63.413723 2.923549 66.998132 3.996094 L 73.007812 6.00279 C 75.272588 6.763951 77.733526 6.763951 79.998302 6.00279 L 84.991508 2.006696 C 88.005265 1.003348 91.001189 0 93.997113 1.003348 C 96.993037 1.003348 99.008152 2.006696 101.005435 3.996094 C 103.002717 6.00279 103.002717 9.998884 102.004076 16.001674 L 102.004076 18.008371 L 104.001359 23.007812 L 105 28.007254 L 104.001359 33.006696 L 101.005435 37.00279 L 95.994395 40.998884 L 78.99966 49.008371 L 75.005095 50.997768 L 74.006454 50.997768 L 58.991168 58.003906 L 54.996603 59.993304 L 52.000679 59.993304 L 51.002038 60.996652 L 46.008832 60.996652 L 39.00051 59.007254 C 28.60394 54.128906 18.278702 49.129464 8.006963 43.991629 L 8.006963 43.00558 L 2.995924 39.995536 L 0 33.992746 C 0 31.294085 0.338825 28.612723 0.998641 26.000558 C 1.997283 23.993862 2.995924 22.004464 5.991848 21.001116 Z M 5.991848 21.001116 " transform="matrix(0.219048,0,0,0.225806,0,0)"></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(25.882354%,9.411765%,1.568628%);fill-opacity:1;" d="M 1.3125 8.128906 L 1.09375 7.675781 L 1.3125 6.097656 L 1.753906 5.644531 L 9.855469 2.933594 L 9.636719 2.257812 C 9.710938 1.882812 9.785156 1.503906 9.855469 1.128906 L 10.078125 1.128906 L 10.515625 1.355469 L 10.734375 1.355469 L 12.265625 2.03125 C 12.914062 1.859375 13.589844 1.859375 14.238281 2.03125 C 14.910156 2.207031 15.570312 2.433594 16.210938 2.710938 L 16.648438 2.710938 L 17.960938 2.484375 L 18.398438 2.03125 L 19.058594 1.582031 L 20.371094 1.355469 L 21.246094 1.804688 L 21.246094 3.839844 L 21.027344 4.066406 L 20.371094 4.515625 L 20.589844 4.515625 L 21.464844 4.96875 L 21.90625 5.417969 L 21.90625 6.324219 L 21.6875 7 L 21.464844 7.675781 L 20.589844 8.128906 C 19.933594 8.582031 18.839844 9.257812 16.867188 9.933594 L 12.484375 11.96875 L 11.828125 12.417969 C 11.617188 12.515625 11.394531 12.589844 11.171875 12.644531 L 10.296875 12.644531 L 8.980469 12.195312 C 6.703125 11.09375 4.441406 9.964844 2.191406 8.804688 Z M 1.3125 8.128906 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(65.490198%,51.372552%,26.274511%);fill-opacity:1;" d="M 6.351562 5.195312 C 8.921875 4.820312 11.476562 4.371094 14.019531 3.839844 L 15.332031 4.289062 L 16.429688 4.515625 C 17.304688 4.515625 17.742188 4.289062 17.960938 4.066406 L 18.179688 3.613281 L 18.617188 3.160156 L 19.933594 2.484375 L 21.027344 2.03125 L 21.027344 3.839844 L 20.808594 4.066406 C 20.371094 4.289062 19.933594 4.515625 19.277344 4.289062 L 17.960938 4.742188 L 19.496094 4.96875 L 21.464844 5.644531 L 21.464844 7.226562 L 21.246094 7.453125 L 20.371094 8.128906 C 17.839844 9.550781 15.203125 10.757812 12.484375 11.742188 L 11.171875 12.417969 C 10.515625 12.417969 9.636719 12.417969 8.980469 11.96875 C 6.324219 10.59375 3.695312 9.164062 1.09375 7.675781 L 1.3125 7 L 1.535156 6.324219 Z M 6.351562 5.195312 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(56.078434%,44.313726%,22.352941%);fill-opacity:1;" d="M 4.382812 7.902344 L 3.503906 7.902344 L 3.285156 9.257812 L 3.066406 9.03125 L 3.285156 7.675781 L 2.847656 7.226562 L 2.628906 7.453125 L 2.410156 8.804688 L 2.191406 8.582031 L 1.972656 8.582031 L 2.410156 7.226562 L 1.972656 7 L 1.753906 8.355469 L 1.3125 8.128906 L 1.535156 6.546875 L 6.351562 6.324219 L 10.078125 8.128906 L 10.953125 10.839844 L 11.171875 12.417969 L 10.296875 12.417969 L 10.515625 10.839844 L 9.417969 10.613281 L 9.199219 12.195312 L 8.542969 11.96875 L 8.761719 10.386719 L 8.105469 10.160156 L 7.886719 11.515625 L 7.449219 11.289062 L 7.449219 9.710938 L 7.230469 9.03125 L 6.570312 9.257812 L 6.132812 10.613281 L 5.476562 10.386719 L 5.914062 9.03125 L 4.820312 8.128906 L 4.601562 8.355469 L 4.382812 9.710938 L 4.160156 9.484375 Z M 4.382812 7.902344 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(52.941179%,41.176471%,18.82353%);fill-opacity:1;" d="M 19.933594 2.484375 L 19.933594 2.710938 C 18.839844 3.160156 18.617188 3.839844 19.496094 4.289062 L 18.617188 4.289062 L 17.960938 4.742188 L 19.496094 4.96875 L 21.464844 5.644531 L 21.464844 7.226562 L 21.246094 7.453125 L 20.371094 8.128906 C 17.839844 9.550781 15.203125 10.757812 12.484375 11.742188 L 11.828125 12.195312 L 10.515625 12.417969 L 10.734375 12.195312 L 10.734375 9.710938 L 10.515625 8.804688 C 13.609375 6.628906 16.75 4.519531 19.933594 2.484375 Z M 19.933594 2.484375 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(47.450981%,36.470589%,16.862746%);fill-opacity:1;" d="M 21.027344 7.453125 C 21.464844 7 21.464844 6.546875 21.464844 5.644531 L 21.464844 7.226562 L 21.246094 7.453125 L 20.371094 8.128906 C 17.839844 9.550781 15.203125 10.757812 12.484375 11.742188 L 11.828125 12.195312 L 10.515625 12.417969 L 10.734375 12.195312 L 11.171875 12.195312 L 12.046875 11.96875 C 15.042969 10.46875 18.035156 8.960938 21.027344 7.453125 Z M 21.027344 7.453125 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(87.450981%,71.764708%,40.784314%);fill-opacity:1;" d="M 1.753906 6.097656 C 5.445312 4.722656 9.167969 3.441406 12.921875 2.257812 L 14.238281 2.484375 L 15.550781 2.933594 L 16.648438 3.160156 C 17.523438 3.160156 17.960938 2.933594 18.179688 2.710938 L 18.617188 2.257812 L 19.058594 2.03125 C 19.714844 1.582031 20.152344 1.582031 20.808594 1.804688 C 21.246094 2.03125 21.246094 2.484375 20.808594 2.710938 L 19.496094 3.160156 L 18.179688 3.613281 L 19.058594 4.289062 C 19.855469 4.75 20.660156 5.203125 21.464844 5.644531 C 21.6875 5.871094 21.246094 6.324219 20.589844 6.773438 C 17.578125 8.257812 14.511719 9.613281 11.390625 10.839844 C 10.734375 11.066406 9.855469 10.839844 8.980469 10.613281 C 6.472656 9.3125 3.988281 7.957031 1.535156 6.546875 L 1.535156 6.097656 Z M 1.753906 6.097656 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(79.215688%,64.313728%,33.333334%);fill-opacity:1;" d="M 2.628906 7.226562 L 2.410156 7.226562 L 4.820312 6.097656 L 6.351562 4.742188 L 8.105469 3.839844 C 9.554688 3.546875 11.015625 3.320312 12.484375 3.160156 C 12.921875 3.160156 13.582031 2.933594 14.019531 2.484375 L 14.457031 2.484375 C 12.945312 3.640625 11.078125 4.203125 9.199219 4.066406 C 8.105469 4.066406 7.230469 4.289062 6.570312 4.742188 L 5.039062 6.097656 C 4.601562 6.546875 3.722656 7 2.628906 7.226562 Z M 2.628906 7.226562 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(79.215688%,64.313728%,33.333334%);fill-opacity:1;" d="M 5.257812 7 C 4.601562 7 3.941406 7.226562 3.503906 7.675781 L 3.285156 7.675781 C 4.5625 6.878906 5.976562 6.34375 7.449219 6.097656 L 10.296875 5.417969 L 11.609375 4.742188 L 12.484375 4.066406 L 14.675781 2.484375 L 14.894531 2.710938 L 12.921875 4.289062 L 10.953125 5.644531 C 9.855469 6.324219 7.886719 6.773438 5.257812 7 Z M 1.972656 6.324219 L 3.066406 5.871094 L 4.160156 5.195312 L 6.132812 4.515625 L 5.476562 4.96875 L 3.722656 6.097656 L 2.191406 6.773438 L 1.972656 6.773438 L 1.535156 6.546875 Z M 9.855469 3.160156 L 11.609375 2.710938 L 13.363281 2.257812 L 13.582031 2.257812 C 13.144531 2.710938 12.484375 2.933594 11.828125 2.933594 Z M 18.617188 7.675781 C 19.277344 6.546875 20.152344 5.871094 21.246094 5.417969 L 21.464844 5.644531 C 20.808594 5.871094 20.152344 6.324219 19.714844 7 Z M 9.199219 7.902344 C 10.078125 7.902344 10.953125 7.675781 12.046875 7 L 14.019531 5.417969 L 15.550781 4.066406 C 16.210938 3.613281 16.648438 3.160156 17.304688 3.160156 L 18.179688 2.710938 L 20.589844 1.804688 L 20.808594 1.804688 L 20.589844 2.03125 L 18.398438 2.933594 L 16.210938 3.839844 L 14.457031 5.417969 C 13.800781 6.324219 13.144531 6.773438 12.703125 7 C 12.265625 7.453125 11.390625 7.902344 10.078125 8.128906 L 7.886719 8.582031 L 6.570312 9.257812 L 5.914062 8.804688 L 7.230469 8.355469 Z M 16.867188 6.097656 C 17.304688 5.195312 17.960938 4.742188 18.839844 4.289062 L 19.496094 4.515625 L 17.742188 5.871094 L 15.992188 7.902344 C 15.113281 8.582031 14.457031 9.03125 13.582031 9.257812 L 10.953125 9.710938 L 9.417969 10.613281 L 8.761719 10.386719 L 10.515625 9.484375 L 12.703125 9.03125 C 14.382812 8.582031 15.851562 7.542969 16.867188 6.097656 Z M 16.867188 6.097656 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(79.215688%,64.313728%,33.333334%);fill-opacity:1;" d="M 6.789062 7.675781 C 5.914062 7.675781 5.257812 7.902344 4.601562 8.355469 L 4.160156 8.128906 C 4.820312 7.675781 5.914062 7.226562 7.230469 7.226562 L 10.953125 6.097656 C 12.046875 5.644531 12.921875 4.96875 13.582031 4.289062 L 15.332031 2.710938 L 15.992188 2.933594 L 14.019531 4.515625 L 12.265625 5.871094 L 9.636719 7.226562 Z M 20.152344 4.742188 L 20.589844 4.96875 L 18.839844 6.324219 C 18.179688 6.773438 17.742188 7.453125 17.304688 8.355469 C 14.960938 9.164062 12.625 9.992188 10.296875 10.839844 L 12.703125 9.933594 L 14.894531 9.03125 C 15.992188 8.582031 17.304688 7.453125 18.617188 5.644531 Z M 13.144531 7.453125 C 14.671875 6.238281 16.203125 5.035156 17.742188 3.839844 C 17.960938 3.386719 19.058594 2.933594 21.027344 2.03125 L 21.027344 2.484375 C 20.402344 2.570312 19.804688 2.800781 19.277344 3.160156 L 18.179688 3.613281 L 18.398438 3.839844 L 15.769531 6.324219 L 13.582031 8.128906 L 10.515625 8.804688 C 9.417969 9.03125 8.761719 9.484375 8.105469 9.933594 L 7.449219 9.710938 C 9.289062 8.8125 11.191406 8.058594 13.144531 7.453125 Z M 6.570312 5.871094 L 6.570312 5.644531 L 6.789062 5.195312 L 7.230469 4.515625 L 8.761719 4.289062 L 10.515625 4.289062 C 10.734375 4.515625 10.734375 4.742188 10.296875 4.96875 L 8.761719 5.644531 L 7.449219 5.871094 L 7.449219 5.195312 L 8.324219 4.742188 L 9.199219 4.515625 L 9.417969 4.742188 L 8.761719 5.195312 L 8.980469 4.96875 L 8.324219 4.96875 L 8.105469 5.195312 C 7.886719 5.417969 8.105469 5.417969 8.324219 5.417969 L 9.199219 5.195312 L 9.855469 4.742188 L 9.855469 4.515625 L 8.761719 4.515625 L 7.449219 4.96875 L 6.789062 5.417969 Z M 6.570312 5.871094 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(41.960785%,25.098041%,14.117648%);fill-opacity:1;" d="M 9.855469 3.160156 L 11.171875 2.710938 L 12.265625 3.839844 L 14.238281 5.417969 L 15.550781 7 L 16.210938 8.582031 L 15.992188 9.484375 L 15.332031 9.710938 L 14.894531 9.710938 L 14.894531 9.257812 L 15.113281 9.03125 L 15.332031 8.582031 L 14.675781 7.675781 L 14.457031 7.453125 L 14.457031 7.226562 L 14.238281 6.773438 L 14.019531 6.773438 C 13.304688 7.003906 12.570312 7.15625 11.828125 7.226562 L 11.171875 7.226562 L 10.734375 6.097656 L 10.078125 4.289062 Z M 9.855469 3.160156 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(47.843137%,47.843137%,47.843137%);fill-opacity:1;" d="M 11.171875 7 L 11.390625 5.871094 L 12.265625 4.96875 L 13.582031 4.96875 L 14.894531 5.195312 L 15.113281 5.871094 L 14.894531 6.097656 L 12.921875 6.773438 L 11.609375 7 Z M 11.171875 7 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(59.215689%,59.215689%,59.215689%);fill-opacity:1;" d="M 12.265625 6.097656 L 12.046875 6.546875 L 12.265625 7 L 11.171875 7 L 11.390625 6.324219 Z M 12.265625 6.097656 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(92.54902%,92.54902%,92.54902%);fill-opacity:1;" d="M 9.855469 1.582031 L 10.734375 1.804688 L 12.921875 2.933594 C 13.75 3.667969 14.488281 4.5 15.113281 5.417969 L 14.894531 5.417969 L 14.019531 4.289062 L 12.484375 2.710938 L 10.734375 1.804688 Z M 9.855469 1.582031 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(78.039217%,78.039217%,78.039217%);fill-opacity:1;" d="M 9.855469 1.582031 L 10.078125 1.582031 L 10.515625 1.804688 C 10.609375 3.429688 11.140625 4.992188 12.046875 6.324219 L 11.171875 7 L 10.953125 6.546875 C 10.421875 5.246094 10.054688 3.882812 9.855469 2.484375 Z M 9.855469 1.582031 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(89.411765%,89.411765%,89.411765%);fill-opacity:1;" d="M 10.078125 3.386719 L 10.078125 2.933594 L 10.734375 2.933594 L 10.734375 3.160156 Z M 9.855469 2.03125 L 10.515625 2.03125 L 10.515625 2.710938 L 9.855469 2.710938 Z M 11.828125 6.097656 L 11.171875 6.773438 L 10.953125 6.324219 L 11.609375 5.871094 Z M 10.296875 4.515625 L 10.078125 3.839844 L 10.734375 3.613281 L 10.953125 4.066406 Z M 11.390625 5.195312 L 10.734375 5.644531 L 10.515625 4.96875 L 10.953125 4.515625 Z M 11.390625 5.195312 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(56.470591%,56.470591%,56.470591%);fill-opacity:1;" d="M 14.238281 6.546875 L 14.019531 6.324219 L 14.019531 5.871094 L 14.675781 5.871094 L 15.113281 6.097656 L 15.113281 6.324219 L 14.894531 6.546875 Z M 14.238281 6.546875 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(70.19608%,70.19608%,70.19608%);fill-opacity:1;" d="M 14.019531 5.871094 L 14.238281 5.644531 L 14.894531 5.644531 L 15.113281 5.871094 L 15.113281 6.097656 L 14.238281 6.097656 Z M 14.019531 5.871094 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(55.686277%,35.686275%,17.254902%);fill-opacity:1;" d="M 14.238281 6.773438 L 14.238281 6.097656 L 14.894531 6.097656 L 15.550781 6.773438 L 16.429688 8.128906 L 16.429688 8.582031 L 15.769531 9.484375 C 15.113281 9.710938 14.675781 9.710938 14.894531 9.257812 L 15.113281 8.582031 L 15.113281 8.128906 L 14.894531 7.675781 L 14.675781 7.453125 L 14.457031 6.773438 Z M 14.238281 6.773438 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(43.137255%,27.058825%,13.725491%);fill-opacity:1;" d="M 15.769531 8.355469 L 16.210938 7.675781 L 16.429688 8.128906 L 16.429688 8.582031 C 16.429688 9.03125 16.210938 9.484375 15.769531 9.484375 L 14.894531 9.484375 L 14.894531 9.03125 L 15.113281 8.582031 L 15.332031 8.355469 Z M 15.769531 8.355469 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(64.313728%,44.705883%,26.666668%);fill-opacity:1;" d="M 14.675781 6.324219 L 14.457031 6.097656 L 14.457031 5.871094 L 15.113281 5.871094 L 15.769531 6.097656 L 16.429688 7.226562 L 16.429688 8.582031 C 15.992188 8.804688 15.550781 9.03125 15.113281 8.804688 L 15.113281 8.355469 L 15.550781 7.902344 L 15.332031 7.453125 L 14.894531 7.226562 L 14.675781 6.773438 Z M 14.675781 6.324219 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(76.078433%,52.156866%,30.980393%);fill-opacity:1;" d="M 15.113281 8.128906 L 15.550781 7.902344 L 15.332031 8.355469 L 15.332031 8.804688 L 15.113281 8.804688 Z M 14.457031 5.871094 L 14.894531 6.546875 L 15.332031 7.453125 L 15.113281 7.453125 L 14.894531 6.773438 L 14.894531 6.546875 L 14.675781 6.546875 L 14.238281 6.097656 Z M 16.429688 7.675781 L 16.429688 7.453125 C 16.648438 7.902344 16.648438 8.355469 16.210938 8.582031 Z M 16.429688 7.675781 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(76.078433%,52.156866%,30.980393%);fill-opacity:1;" d="M 14.894531 6.097656 L 15.332031 6.546875 L 15.550781 6.773438 L 15.550781 7 L 15.769531 7.453125 L 15.769531 8.128906 L 15.550781 8.804688 L 15.550781 7 L 14.675781 5.871094 Z M 14.894531 6.097656 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(76.078433%,52.156866%,30.980393%);fill-opacity:1;" d="M 15.550781 6.324219 L 15.992188 6.546875 L 16.210938 7.226562 L 16.210938 8.128906 L 15.992188 8.804688 L 15.992188 6.546875 C 15.695312 6.324219 15.402344 6.101562 15.113281 5.871094 L 15.332031 5.871094 Z M 15.550781 6.324219 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(76.078433%,52.156866%,30.980393%);fill-opacity:1;" d="M 15.113281 5.871094 L 15.332031 6.324219 L 15.550781 6.773438 L 15.769531 7.226562 L 15.992188 7.453125 L 15.992188 8.128906 L 15.769531 8.804688 L 15.769531 7 L 15.550781 6.773438 L 15.332031 6.324219 L 14.894531 5.871094 Z M 15.332031 5.871094 L 15.769531 6.324219 L 15.992188 6.546875 L 15.550781 6.324219 Z M 15.332031 5.871094 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(43.137255%,27.058825%,13.725491%);fill-opacity:1;" d="M 16.210938 8.355469 C 15.992188 8.582031 15.769531 8.804688 15.550781 8.582031 L 15.332031 8.355469 L 15.992188 8.128906 Z M 16.210938 8.355469 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(69.411767%,69.411767%,69.411767%);fill-opacity:1;" d="M 16.210938 8.355469 L 15.992188 8.582031 L 15.332031 8.582031 L 15.769531 8.355469 Z M 16.210938 8.355469 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(49.019608%,49.019608%,49.019608%);fill-opacity:1;" d="M 16.210938 8.355469 L 15.992188 8.582031 L 15.332031 8.582031 L 15.769531 8.582031 L 15.992188 8.355469 Z M 16.210938 8.355469 "></path> <path style=" stroke:none;fill-rule:evenodd;fill:rgb(76.078433%,52.156866%,30.980393%);fill-opacity:1;" d="M 15.992188 6.773438 L 15.769531 6.773438 L 15.992188 7 L 15.992188 6.773438 L 15.992188 7 L 15.769531 7 L 15.769531 6.773438 Z M 15.992188 6.773438 "></path></g></svg><label><input type="text" name="time-preparation" class="input input-bordered input-xs max-w-24 html-duration-picker" value="00:30:00"></label></div><div class="flex justify-self-center items-center gap-1 cursor-default" title="Cooking time"><svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="24px" height="23px" viewBox="0 0 24 23" version="1.1"><g id="surface1"><path style=" stroke:none;fill-rule:nonzero;fill:rgb(62.745098%,64.705882%,65.882353%);fill-opacity:1;" d="M 4.636719 10.984375 L 4.636719 20.417969 C 4.636719 21.007812 5.078125 21.527344 5.667969 21.527344 L 18.257812 21.527344 C 18.847656 21.527344 19.363281 21.007812 19.363281 20.417969 L 19.363281 10.984375 Z M 4.636719 10.984375 "></path> <path style=" stroke:none;fill-rule:nonzero;fill:rgb(76.862745%,76.862745%,76.862745%);fill-opacity:1;" d="M 19.289062 9.953125 C 18.992188 8.699219 17.964844 7.8125 16.710938 7.8125 L 7.214844 7.8125 C 5.964844 7.8125 4.933594 8.699219 4.710938 9.953125 Z M 19.289062 9.953125 "></path> <path style=" stroke:none;fill-rule:nonzero;fill:rgb(54.509804%,69.411765%,81.960784%);fill-opacity:1;" d="M 16.710938 7.8125 L 13.914062 7.8125 C 14.945312 7.8125 15.828125 8.476562 16.269531 9.363281 C 16.417969 9.730469 16.785156 9.953125 17.226562 9.953125 L 19.289062 9.953125 C 18.992188 8.699219 17.964844 7.8125 16.710938 7.8125 "></path> <path style=" stroke:none;fill-rule:nonzero;fill:rgb(0.392157%,26.666667%,38.823529%);fill-opacity:1;" d="M 22.160156 9.953125 L 20.320312 9.953125 C 20.097656 8.183594 18.550781 6.78125 16.710938 6.78125 L 15.535156 6.78125 L 15.3125 5.898438 C 15.09375 5.160156 14.503906 4.644531 13.765625 4.644531 L 11.191406 4.644531 C 10.453125 4.644531 9.792969 5.160156 9.644531 5.898438 L 9.421875 6.78125 L 7.214844 6.78125 C 5.449219 6.78125 3.902344 8.183594 3.605469 9.953125 L 1.765625 9.953125 C 1.03125 9.953125 0.441406 10.542969 0.441406 11.277344 C 0.441406 11.5 0.441406 11.722656 0.589844 11.941406 L 1.25 13.269531 C 1.546875 13.785156 2.0625 14.152344 2.648438 14.152344 L 3.605469 14.152344 L 3.605469 20.417969 C 3.605469 21.597656 4.492188 22.558594 5.667969 22.558594 L 18.257812 22.558594 C 19.4375 22.558594 20.394531 21.597656 20.394531 20.417969 L 20.394531 14.152344 L 21.351562 14.152344 C 21.9375 14.152344 22.453125 13.859375 22.75 13.269531 L 23.410156 11.867188 C 23.558594 11.648438 23.558594 11.5 23.558594 11.277344 C 23.558594 10.542969 22.894531 9.953125 22.160156 9.953125 M 3.605469 13.121094 L 2.648438 13.121094 C 2.429688 13.121094 2.28125 12.972656 2.136719 12.753906 L 1.472656 11.5 L 1.472656 11.277344 C 1.472656 11.132812 1.621094 10.984375 1.765625 10.984375 L 3.605469 10.984375 Z M 10.75 6.117188 C 10.75 5.972656 10.96875 5.75 11.265625 5.75 L 13.839844 5.75 C 14.0625 5.75 14.28125 5.898438 14.355469 6.117188 L 14.503906 6.78125 L 10.601562 6.78125 Z M 7.214844 7.8125 L 16.710938 7.8125 C 17.964844 7.8125 18.992188 8.699219 19.289062 9.953125 L 4.710938 9.953125 C 4.933594 8.699219 5.964844 7.8125 7.214844 7.8125 M 19.363281 13.636719 L 19.363281 20.417969 C 19.363281 21.007812 18.847656 21.527344 18.257812 21.527344 L 5.667969 21.527344 C 5.078125 21.527344 4.636719 21.007812 4.636719 20.417969 L 4.636719 10.984375 L 19.363281 10.984375 Z M 22.453125 11.5 L 21.71875 12.828125 C 21.644531 12.972656 21.496094 13.121094 21.277344 13.121094 L 20.394531 13.121094 L 20.394531 11.058594 L 22.160156 11.058594 C 22.308594 11.058594 22.453125 11.207031 22.453125 11.351562 L 22.453125 11.5 "></path> <path style=" stroke:none;fill-rule:nonzero;fill:rgb(92.54902%,94.117647%,97.647059%);fill-opacity:1;" d="M 7.804688 17.324219 C 8.089844 17.324219 8.320312 17.554688 8.320312 17.839844 C 8.320312 18.125 8.089844 18.355469 7.804688 18.355469 C 7.519531 18.355469 7.289062 18.125 7.289062 17.839844 C 7.289062 17.554688 7.519531 17.324219 7.804688 17.324219 M 7.804688 16.808594 C 7.4375 16.808594 7.214844 16.585938 7.214844 16.21875 L 7.214844 13.121094 C 7.214844 12.753906 7.4375 12.53125 7.804688 12.53125 C 8.097656 12.53125 8.320312 12.753906 8.320312 13.121094 L 8.320312 16.21875 C 8.320312 16.585938 8.097656 16.808594 7.804688 16.808594 "></path> <path style=" stroke:none;fill-rule:nonzero;fill:rgb(54.509804%,69.411765%,81.960784%);fill-opacity:1;" d="M 16.195312 12.015625 L 16.195312 19.902344 C 16.195312 20.492188 15.679688 21.007812 15.09375 21.007812 L 6.699219 21.007812 C 6.183594 21.007812 5.667969 20.492188 5.667969 19.902344 L 5.667969 12.015625 C 5.667969 11.5 5.226562 10.984375 4.636719 10.984375 L 4.636719 20.417969 C 4.636719 21.007812 5.078125 21.527344 5.667969 21.527344 L 18.257812 21.527344 C 18.847656 21.527344 19.363281 21.007812 19.363281 20.417969 L 19.363281 10.984375 L 17.226562 10.984375 C 16.636719 10.984375 16.195312 11.5 16.195312 12.015625 M 8.246094 7.8125 L 7.214844 7.8125 C 5.964844 7.8125 4.933594 8.699219 4.710938 9.953125 L 4.933594 9.953125 C 5.375 9.953125 5.742188 9.730469 5.890625 9.363281 C 6.332031 8.476562 7.214844 7.8125 8.246094 7.8125 "></path> <path style=" stroke:none;fill-rule:nonzero;fill:rgb(0.392157%,26.666667%,38.823529%);fill-opacity:1;" d="M 18.773438 5.75 C 18.625 5.75 18.480469 5.675781 18.40625 5.527344 C 18.183594 5.308594 18.257812 4.9375 18.40625 4.792969 C 18.699219 4.644531 18.773438 4.421875 18.773438 4.128906 C 18.773438 3.90625 18.699219 3.6875 18.480469 3.539062 C 18.257812 3.316406 18.257812 3.023438 18.40625 2.800781 C 18.625 2.582031 18.992188 2.507812 19.140625 2.726562 C 19.582031 3.097656 19.878906 3.613281 19.878906 4.128906 C 19.878906 4.71875 19.582031 5.234375 19.140625 5.601562 L 18.773438 5.75 M 18.773438 1.03125 C 19.058594 1.03125 19.289062 1.261719 19.289062 1.546875 C 19.289062 1.832031 19.058594 2.0625 18.773438 2.0625 C 18.488281 2.0625 18.257812 1.832031 18.257812 1.546875 C 18.257812 1.261719 18.488281 1.03125 18.773438 1.03125 M 16.710938 5.75 L 16.269531 5.527344 C 16.050781 5.308594 16.121094 4.9375 16.34375 4.792969 C 16.5625 4.644531 16.710938 4.421875 16.710938 4.128906 C 16.710938 3.90625 16.5625 3.6875 16.417969 3.539062 C 15.902344 3.097656 15.679688 2.652344 15.679688 2.0625 C 15.679688 1.472656 15.902344 1.03125 16.417969 0.589844 C 16.5625 0.367188 16.933594 0.441406 17.152344 0.664062 C 17.300781 0.8125 17.300781 1.179688 17.078125 1.402344 C 16.785156 1.546875 16.710938 1.769531 16.710938 2.0625 C 16.710938 2.285156 16.785156 2.507812 17.007812 2.652344 C 17.519531 3.097656 17.742188 3.613281 17.742188 4.128906 C 17.742188 4.71875 17.519531 5.234375 17.007812 5.601562 L 16.710938 5.75 "></path></g></svg><label><input type="text" name="time-cooking" class="input input-bordered input-xs max-w-24 html-duration-picker" value="01:00:00"></label></div></div>`,
			`<tbody><tr><td>Calories</td><td><label><input type="text" name="calories" autocomplete="off" placeholder="kcal" class="input input-bordered input-xs max-w-24" value="354"></label></td></tr><tr><td>Total carbs</td><td><label><input type="text" name="total-carbohydrates" autocomplete="off" placeholder="g" class="input input-bordered input-xs max-w-24" value="7"></label></td></tr><tr><td>Sugars</td><td><label><input type="text" name="sugars" autocomplete="off" placeholder="g" class="input input-bordered input-xs max-w-24" value="6"></label></td></tr><tr><td>Protein</td><td><label><input type="text" name="protein" autocomplete="off" placeholder="g" class="input input-bordered input-xs max-w-24" value="3"></label></td></tr><tr><td>Total fat</td><td><label><input type="text" name="total-fat" autocomplete="off" placeholder="g" class="input input-bordered input-xs max-w-24" value="8"></label></td></tr><tr><td>Saturated fat</td><td><label><input type="text" name="saturated-fat" autocomplete="off" placeholder="g" class="input input-bordered input-xs max-w-24" value="4"></label></td></tr><tr><td>Unsaturated fat</td><td><label><input type="text" name="unsaturated-fat" autocomplete="off" placeholder="g" class="input input-bordered input-xs max-w-24" value="9"></label></td></tr><tr><td>Trans fat</td><td><label><input type="text" name="trans-fat" autocomplete="off" placeholder="g" class="input input-bordered input-xs max-w-24" value="10"></label></td></tr><tr><td>Cholesterol</td><td><label><input type="text" name="cholesterol" autocomplete="off" placeholder="mg" class="input input-bordered input-xs max-w-24" value="1"></label></td></tr><tr><td>Sodium</td><td><label><input type="text" name="sodium" autocomplete="off" placeholder="mg" class="input input-bordered input-xs max-w-24" value="5"></label></td></tr><tr><td>Fiber</td><td><label><input type="text" name="fiber" autocomplete="off" placeholder="g" class="input input-bordered input-xs max-w-24" value="2"></label></td></tr></tbody>`,
			`<input type="text" name="tools" placeholder="1 frying pan" class="input input-bordered input-sm w-full" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)">`,
			`<input type="text" name="time-preparation" class="input input-bordered input-xs max-w-24 html-duration-picker" value="00:30:00">`,
			`<input required type="text" name="ingredients" value="ing1" placeholder="1 cup of chopped onions" class="input input-bordered input-sm w-full" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)">`,
			`<textarea required name="instructions" rows="3" class="textarea textarea-bordered w-full" placeholder="Mix all ingredients together" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)">ins1</textarea>`,
		})
	})

	t.Run("no updated image", func(t *testing.T) {
		files := &mockFiles{}
		srv.Files = files
		srv.Repository = &mockRepository{RecipesRegistered: map[int64]models.Recipes{1: {baseRecipe}}}
		contentType, body := createMultipartForm(map[string][]string{})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		got, _ := srv.Repository.Recipe(baseRecipe.ID, 1)
		if !slices.Equal(got.Images, baseRecipe.Images) {
			t.Fatal("image should not have been updated")
		}
		if files.uploadImageHitCount > 0 {
			t.Fatal("must not have uploaded image")
		}
	})

	t.Run("updated image", func(t *testing.T) {
		files := &mockFiles{}
		srv.Files = files
		srv.Repository = &mockRepository{RecipesRegistered: map[int64]models.Recipes{1: {baseRecipe}}}
		contentType, body := createMultipartForm(map[string][]string{"images": {"eggs.jpg"}})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		got, _ := srv.Repository.Recipe(baseRecipe.ID, 1)
		isImagesEqual := slices.Equal(got.Images, baseRecipe.Images)
		isVideosEqual := slices.Equal(got.Videos, baseRecipe.Videos)
		if (isImagesEqual && !isVideosEqual) || (!isImagesEqual && isVideosEqual) {
			t.Fatal("image should have been updated")
		}
		if files.uploadImageHitCount == 0 {
			t.Fatal("must have uploaded image")
		}
	})

	t.Run("update every field", func(t *testing.T) {
		files := &mockFiles{}
		srv.Files = files
		srv.Repository = &mockRepository{RecipesRegistered: map[int64]models.Recipes{1: {baseRecipe}}}
		fields := map[string][]string{
			"title":               {"Salsa"},
			"image":               {"jesus.jpg"},
			"category":            {"appetizers"},
			"source":              {"Mommy"},
			"description":         {"The best"},
			"calories":            {"666"},
			"total-carbohydrates": {"31g"},
			"sugars":              {"0.1mg"},
			"protein":             {"5g"},
			"total-fat":           {"24g"},
			"saturated-fat":       {"58g"},
			"unsaturated-fat":     {"1g"},
			"trans-fat":           {"2g"},
			"cholesterol":         {"256mg"},
			"sodium":              {"777mg"},
			"fiber":               {"2g"},
			"time-preparation":    {"00:15:30"},
			"time-cooking":        {"00:30:15"},
			"ingredients":         {"cheese", "avocado"},
			"instructions":        {"mix", "eat"},
			"keywords":            {"chicken", "big", "marinade"},
			"yield":               {"4"},
		}
		contentType, body := createMultipartForm(fields)

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		assertHeader(t, rr, "HX-Redirect", "/recipes/1")
		got, _ := srv.Repository.Recipe(baseRecipe.ID, 1)
		want := models.Recipe{
			Category:     "appetizers",
			CreatedAt:    baseRecipe.CreatedAt,
			Description:  "The best",
			ID:           baseRecipe.ID,
			Images:       got.Images,
			Ingredients:  []string{"cheese", "avocado"},
			Instructions: []string{"mix", "eat"},
			Keywords:     []string{"chicken", "big", "marinade"},
			Name:         "Salsa",
			Nutrition: models.Nutrition{
				Calories:           "666",
				Cholesterol:        "256mg",
				Fiber:              "2g",
				Protein:            "5g",
				SaturatedFat:       "58g",
				Sodium:             "777mg",
				Sugars:             "0.1mg",
				TotalCarbohydrates: "31g",
				TotalFat:           "24g",
				UnsaturatedFat:     "1g",
				TransFat:           "2g",
			},
			Times: models.Times{
				Prep:  15*time.Minute + 30*time.Second,
				Cook:  30*time.Minute + 15*time.Second,
				Total: 45*time.Minute + 45*time.Second,
			},
			Tools:     baseRecipe.Tools,
			UpdatedAt: baseRecipe.UpdatedAt,
			URL:       "Mommy",
			Yield:     4,
		}
		if !cmp.Equal(*got, want) {
			t.Log(cmp.Diff(*got, want))
			t.Fail()
		}
	})

	t.Run("missing source defaults to unknown", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		if repo.RecipesRegistered[1][0].URL != "Unknown" {
			t.Fatalf("got source %q; want 'unknown'", repo.RecipesRegistered[1][0].URL)
		}
	})

	t.Run("missing category defaults to uncategorized", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		if repo.RecipesRegistered[1][0].Category != "uncategorized" {
			t.Fatalf("got category %q; want 'uncategorized'", repo.RecipesRegistered[1][0].Category)
		}
	})

	t.Run("missing yield defaults to 1", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		if repo.RecipesRegistered[1][0].Yield != 1 {
			t.Fatalf("got yield %d; want 1", repo.RecipesRegistered[1][0].Yield)
		}
	})

	t.Run("can only be one category", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"category":     {"dinner,lunch"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		if repo.RecipesRegistered[1][0].Category != "dinner" {
			t.Fatalf("got category %s; want dinner", repo.RecipesRegistered[1][0].Category)
		}
	})

	t.Run("subcategories are possible", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"category":     {"beverages:cocktails:vodka"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		if repo.RecipesRegistered[1][0].Category != "beverages:cocktails:vodka" {
			t.Fatalf("got category %s; want beverages:cocktails:vodka", repo.RecipesRegistered[1][0].Category)
		}
	})
}

func TestHandlers_Recipes_Scale(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	originalRepo := srv.Repository

	uri := ts.URL + "/recipes/1/scale"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	yieldTestcases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "yield query parameter must be present",
			in:   "",
			want: "No yield in the query.",
		},
		{
			name: "yield query parameter must be greater than zero",
			in:   "-1",
			want: "Yield must be greater than zero.",
		},
		{
			name: "yield query parameter must be greater than zero",
			in:   "0",
			want: "Yield must be greater than zero.",
		},
	}
	for _, tc := range yieldTestcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?yield="+tc.in)

			assertStatus(t, rr.Code, http.StatusBadRequest)
			assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"`+tc.want+`","title":"General Error"}}`)
		})
	}

	t.Run("cannot find recipe in database", func(t *testing.T) {
		srv.Repository = &mockRepository{
			RecipesRegistered: map[int64]models.Recipes{1: make(models.Recipes, 0)},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?yield=6")

		assertStatus(t, rr.Code, http.StatusNotFound)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Recipe not found.","title":"General Error"}}`)
	})

	t.Run("valid request double yield", func(t *testing.T) {
		srv.Repository = &mockRepository{
			RecipesRegistered: map[int64]models.Recipes{
				1: []models.Recipe{
					{
						ID:   1,
						Name: "American Jerky",
						Ingredients: []string{
							"2lb chicken",
							"1/2 cup bread loaf",
							"½ tbsp beef broth",
							"7 1/2 cups flour",
							"2 big apples",
							"Lots of big apples",
							"2.5 slices of bacon",
							"2 1/3 cans of bamboo sticks",
							"1½can of tomato paste",
							"6 ¾ peanut butter jars",
							"7.5mL of whiskey",
							"2 tsp lemon juice",
							"Ground ginger",
							"3 Large or 4 medium ripe Hass avocados",
							"1/4-1/2 teaspoon salt plus more for seasoning",
							"1/2 fresh pineapple, cored and cut into 1 1/2-inch pieces",
							"Un sac de chips de 1kg",
							"Two 15-ounce can Goya beans",
						},
						Instructions: []string{},
						Yield:        4,
					},
				},
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?yield=8")

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">4 lb chicken</span></label>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">1 cup bread loaf</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">1 tbsp beef broth</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">15 cups flour</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">4 big apples</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">Lots of big apples</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">5 slices of bacon</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">4 2/3 cans of bamboo sticks</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">3 can of tomato paste</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">13 1/2 peanut butter jars</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">15 mL of whiskey</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">1 1/3 tbsp lemon juice</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">Ground ginger</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">6 Large or 8 medium ripe Hass avocados</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">1/2 tsp salt plus more for seasoning</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">1 fresh pineapple, cored and cut into 3-inch pieces</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">Un sac de chips de 1kg</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">4 15-ounce can Goya beans</span>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Recipes_Categories(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	baseRecipes := map[int64]models.Recipes{
		1: {
			{ID: 1, Name: "Chinese Firmware", Category: "breakfast"},
			{ID: 2, Name: "Lovely Canada", Category: "lunch"},
			{ID: 3, Name: "Lovely Ukraine", Category: "dinner"},
			{ID: 4, Name: "Space Disco", Category: "snack"},
			{ID: 4, Name: "Maple Pancakes", Category: "breakfast"},
		},
	}
	srv.Repository = &mockRepository{
		categories:        map[int64][]string{1: {"uncategorized", "chicken"}},
		RecipesRegistered: baseRecipes,
	}
	originalRepo := srv.Repository

	uri := ts.URL + "/recipes/categories"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("delete category updates recipes", func(t *testing.T) {
		d := url.Values{}
		d.Add("category", "breakfast")

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri+"?"+d.Encode(), formHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		for _, r := range srv.Repository.RecipesAll(1) {
			if r.Category == "breakfast" {
				t.Fatal("expected recipes with category 'breakfast' to be updated")
			}
		}
	})

	t.Run("failed to delete category", func(t *testing.T) {
		srv.Repository = &mockRepository{
			DeleteCategoryFunc: func(_ string, _ int64) error {
				return errors.New("that's some bad hat harry")
			},
			RecipesRegistered: baseRecipes,
		}
		defer func() {
			srv.Repository = originalRepo
		}()
		d := url.Values{}
		d.Add("category", "breakfast")

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri+"?"+d.Encode(), formHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to delete category.","title":"General Error"}}`)
	})

	t.Run("cannot delete uncategorized category", func(t *testing.T) {
		d := url.Values{}
		d.Add("category", "uncategorized")

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri+"?"+d.Encode(), formHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to delete category.","title":"General Error"}}`)
	})

	t.Run("add new category", func(t *testing.T) {
		categoriesBefore, _ := srv.Repository.Categories(1)
		d := url.Values{}
		d.Add("category", "uber chicken")

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader(d.Encode()))

		assertStatus(t, rr.Code, http.StatusCreated)
		categoriesAfter, _ := srv.Repository.Categories(1)
		if len(categoriesBefore) == len(categoriesAfter) {
			t.Fatal("categories should have been updated")
		}
		assertStringsInHTML(t, getBodyHTML(rr), []string{
			`<div class="badge badge-outline p-3 pr-0"><form class="inline-flex" hx-delete="/recipes/categories" hx-target="closest <div/>" hx-swap="delete"><input type="hidden" name="category" value="uber chicken"> <span class="select-none">uber chicken</span> <button type="submit" class="btn btn-xs btn-ghost">X</button></form></div>`,
			`<div class="badge badge-outline p-3 pr-0"><form class="inline-flex" hx-post="/recipes/categories" hx-target="closest <div/>" hx-swap="outerHTML"><label class="form-control"><input required type="text" placeholder="New category" class="input input-ghost input-xs w-[16ch] focus:outline-none" name="category" autocomplete="off"></label> <button class="btn btn-xs btn-ghost">&#10003;</button></form></div>`,
		})
	})

	t.Run("cannot add existing category", func(t *testing.T) {
		categoriesBefore, _ := srv.Repository.Categories(1)
		d := url.Values{}
		d.Add("category", "uncategorized")

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader(d.Encode()))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		categoriesAfter, _ := srv.Repository.Categories(1)
		if len(categoriesBefore) != len(categoriesAfter) {
			t.Fatal("categories should not have been updated")
		}
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to add category.","title":"General Error"}}`)
	})

	t.Run("failed to add category", func(t *testing.T) {
		srv.Repository = &mockRepository{
			AddRecipeCategoryFunc: func(_ string, _ int64) error {
				return errors.New("that's very bad")
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()
		categoriesBefore, _ := srv.Repository.Categories(1)
		d := url.Values{}
		d.Add("category", "uber chicken")

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader(d.Encode()))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		categoriesAfter, _ := srv.Repository.Categories(1)
		if len(categoriesBefore) != len(categoriesAfter) {
			t.Fatal("categories should not have been updated")
		}
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to add category.","title":"General Error"}}`)
	})
}

func TestHandlers_Recipes_Search(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	srv.Repository = &mockRepository{
		RecipesRegistered: map[int64]models.Recipes{
			1: {
				{ID: 1, Name: "Chinese Firmware"},
				{ID: 2, Name: "Lovely Canada"},
				{ID: 3, Name: "Lovely Ukraine"},
			},
		},
	}

	uri := ts.URL + "/recipes/search"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("page below 1 returns recipes from first page", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?q=lovely&page=0&mode=name")

		assertStatus(t, rr.Code, http.StatusOK)
		assertStringsInHTML(t, getBodyHTML(rr), []string{
			`<button class="join-item btn btn-disabled">«</button>`,
			`<!-- Left Section --><button aria-current="page" class="join-item btn btn-active">1</button>`,
			`<!-- Right Section --><button class="join-item btn btn-disabled">»</button></div><div class="text-center"><p class="text-sm">Showing <span class="font-medium">1</span> to <span class="font-medium">2</span> of <span id="search-count" class="font-medium">2</span> results</p></div>`,
		})
	})

	t.Run("empty query redirects to recipes index", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?q=&page=0&sort=a-z")

		assertStatus(t, rr.Code, http.StatusOK)
		assertHeader(t, rr, "HX-Retarget", "#content")
	})

	t.Run("no results", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?q=kool-aid&mode=name")

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<div class="grid place-content-center text-sm text-center h-3/5 md:text-base"><p>No results found.</p></div>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	searches := []struct {
		query string
		want  []string
	}{
		{
			query: "lov",
			want: []string{
				`<section class="card-side sm:card card-compact card-bordered bg-base-100 shadow-lg indicator w-full"><span class="hidden sm:block"><span class="badge badge-primary select-none cursor-pointer badge-sm p-2 m-1 sm:badge-md sm:m-0 hover:bg-neutral indicator-item indicator-center" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put "cat:" into #search_recipes.value"></span></span><figure class="relative cursor-pointer" hx-get="/recipes/2" hx-target="#content" hx-push-url="true" hx-trigger="mousedown" hx-swap="innerHTML show:window:top transition:true"><img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Canada recipe"><div class="hidden absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 items-center justify-center text-white select-none rounded-t-lg sm:flex"><p class="p-2 text-sm"></p></div></figure><div class="card-body justify-between"><h2 class="sm:font-semibold sm:w-[25ch] sm:break-words sm:min-h-14 sm:min-h-28">Lovely Canada</h2><div class="sm:max-h-14 sm:overflow-y-auto sm:content-end"><div class="flex flex-col flex-wrap overflow-x-auto max-h-12 pb-2 sm:pb-0 sm:max-h-none sm:flex-auto sm:flex-row"><span class="sm:hidden"><span class="badge badge-primary select-none cursor-pointer badge-sm p-2 m-1 sm:badge-md sm:m-0 hover:bg-neutral" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put "cat:" into #search_recipes.value"></span></span> </div></div><div class="card-actions flex-col-reverse h-fit"><button class="btn btn-block btn-xs btn-outline sm:btn-sm" hx-get="/recipes/2" hx-target="#content" hx-trigger="mousedown" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true">View</button></div></div></section>`,
				`<section class="card-side sm:card card-compact card-bordered bg-base-100 shadow-lg indicator w-full"><span class="hidden sm:block"><span class="badge badge-primary select-none cursor-pointer badge-sm p-2 m-1 sm:badge-md sm:m-0 hover:bg-neutral indicator-item indicator-center" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put "cat:" into #search_recipes.value"></span></span><figure class="relative cursor-pointer" hx-get="/recipes/3" hx-target="#content" hx-push-url="true" hx-trigger="mousedown" hx-swap="innerHTML show:window:top transition:true"><img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Ukraine recipe"><div class="hidden absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 items-center justify-center text-white select-none rounded-t-lg sm:flex"><p class="p-2 text-sm"></p></div></figure><div class="card-body justify-between"><h2 class="sm:font-semibold sm:w-[25ch] sm:break-words sm:min-h-14 sm:min-h-28">Lovely Ukraine</h2><div class="sm:max-h-14 sm:overflow-y-auto sm:content-end"><div class="flex flex-col flex-wrap overflow-x-auto max-h-12 pb-2 sm:pb-0 sm:max-h-none sm:flex-auto sm:flex-row"><span class="sm:hidden"><span class="badge badge-primary select-none cursor-pointer badge-sm p-2 m-1 sm:badge-md sm:m-0 hover:bg-neutral" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put "cat:" into #search_recipes.value"></span></span> </div></div><div class="card-actions flex-col-reverse h-fit"><button class="btn btn-block btn-xs btn-outline sm:btn-sm" hx-get="/recipes/3" hx-target="#content" hx-trigger="mousedown" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true">View</button></div></div></section>`,
			},
		},
		{
			query: "chi",
			want: []string{
				`<section class="card-side sm:card card-compact card-bordered bg-base-100 shadow-lg indicator w-full"><span class="hidden sm:block"><span class="badge badge-primary select-none cursor-pointer badge-sm p-2 m-1 sm:badge-md sm:m-0 hover:bg-neutral indicator-item indicator-center" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put "cat:" into #search_recipes.value"></span></span><figure class="relative cursor-pointer" hx-get="/recipes/1" hx-target="#content" hx-push-url="true" hx-trigger="mousedown" hx-swap="innerHTML show:window:top transition:true"><img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Chinese Firmware recipe"><div class="hidden absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 items-center justify-center text-white select-none rounded-t-lg sm:flex"><p class="p-2 text-sm"></p></div></figure><div class="card-body justify-between"><h2 class="sm:font-semibold sm:w-[25ch] sm:break-words sm:min-h-14 sm:min-h-28">Chinese Firmware</h2><div class="sm:max-h-14 sm:overflow-y-auto sm:content-end"><div class="flex flex-col flex-wrap overflow-x-auto max-h-12 pb-2 sm:pb-0 sm:max-h-none sm:flex-auto sm:flex-row"><span class="sm:hidden"><span class="badge badge-primary select-none cursor-pointer badge-sm p-2 m-1 sm:badge-md sm:m-0 hover:bg-neutral" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put "cat:" into #search_recipes.value"></span></span> </div></div><div class="card-actions flex-col-reverse h-fit"><button class="btn btn-block btn-xs btn-outline sm:btn-sm" hx-get="/recipes/1" hx-target="#content" hx-trigger="mousedown" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true">View</button></div></div></section>`,
			},
		},
		{
			query: "lovely",
			want: []string{
				`<section class="card-side sm:card card-compact card-bordered bg-base-100 shadow-lg indicator w-full"><span class="hidden sm:block"><span class="badge badge-primary select-none cursor-pointer badge-sm p-2 m-1 sm:badge-md sm:m-0 hover:bg-neutral indicator-item indicator-center" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put "cat:" into #search_recipes.value"></span></span><figure class="relative cursor-pointer" hx-get="/recipes/2" hx-target="#content" hx-push-url="true" hx-trigger="mousedown" hx-swap="innerHTML show:window:top transition:true"><img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Canada recipe"><div class="hidden absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 items-center justify-center text-white select-none rounded-t-lg sm:flex"><p class="p-2 text-sm"></p></div></figure><div class="card-body justify-between"><h2 class="sm:font-semibold sm:w-[25ch] sm:break-words sm:min-h-14 sm:min-h-28">Lovely Canada</h2><div class="sm:max-h-14 sm:overflow-y-auto sm:content-end"><div class="flex flex-col flex-wrap overflow-x-auto max-h-12 pb-2 sm:pb-0 sm:max-h-none sm:flex-auto sm:flex-row"><span class="sm:hidden"><span class="badge badge-primary select-none cursor-pointer badge-sm p-2 m-1 sm:badge-md sm:m-0 hover:bg-neutral" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put "cat:" into #search_recipes.value"></span></span> </div></div><div class="card-actions flex-col-reverse h-fit"><button class="btn btn-block btn-xs btn-outline sm:btn-sm" hx-get="/recipes/2" hx-target="#content" hx-trigger="mousedown" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true">View</button></div></div></section>`,
				`<section class="card-side sm:card card-compact card-bordered bg-base-100 shadow-lg indicator w-full"><span class="hidden sm:block"><span class="badge badge-primary select-none cursor-pointer badge-sm p-2 m-1 sm:badge-md sm:m-0 hover:bg-neutral indicator-item indicator-center" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put "cat:" into #search_recipes.value"></span></span><figure class="relative cursor-pointer" hx-get="/recipes/3" hx-target="#content" hx-push-url="true" hx-trigger="mousedown" hx-swap="innerHTML show:window:top transition:true"><img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Ukraine recipe"><div class="hidden absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 items-center justify-center text-white select-none rounded-t-lg sm:flex"><p class="p-2 text-sm"></p></div></figure><div class="card-body justify-between"><h2 class="sm:font-semibold sm:w-[25ch] sm:break-words sm:min-h-14 sm:min-h-28">Lovely Ukraine</h2><div class="sm:max-h-14 sm:overflow-y-auto sm:content-end"><div class="flex flex-col flex-wrap overflow-x-auto max-h-12 pb-2 sm:pb-0 sm:max-h-none sm:flex-auto sm:flex-row"><span class="sm:hidden"><span class="badge badge-primary select-none cursor-pointer badge-sm p-2 m-1 sm:badge-md sm:m-0 hover:bg-neutral" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put "cat:" into #search_recipes.value"></span></span> </div></div><div class="card-actions flex-col-reverse h-fit"><button class="btn btn-block btn-xs btn-outline sm:btn-sm" hx-get="/recipes/3" hx-target="#content" hx-trigger="mousedown" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true">View</button></div></div></section>`,
			},
		},
	}
	for _, tc := range searches {
		t.Run("results for "+tc.query, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?mode=name&q="+tc.query, formHeader, strings.NewReader("id=1&page=1&q="+tc.query))

			assertStatus(t, rr.Code, http.StatusOK)
			assertStringsInHTML(t, getBodyHTML(rr), tc.want)
		})
	}
}

func TestHandlers_Recipes_Share(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	uri := func(id int64) string {
		return fmt.Sprintf("%s/recipes/%d/share", ts.URL, id)
	}

	app.Config.Server.URL = "https://www.recipya.com"
	recipe := models.Recipe{
		Category:     "American",
		Description:  "This is the most delicious recipe!",
		ID:           1,
		Images:       []uuid.UUID{uuid.New()},
		Ingredients:  []string{"Ing1", "Ing2", "Ing3"},
		Instructions: []string{"Ins1", "Ins2", "Ins3"},
		Name:         "Chicken Jersey",
		Nutrition: models.Nutrition{
			Calories:           "500",
			Cholesterol:        "1mg",
			Fiber:              "2g",
			Protein:            "3g",
			SaturatedFat:       "4g",
			Sodium:             "5mg",
			Sugars:             "6g",
			TotalCarbohydrates: "7g",
			TotalFat:           "8g",
			UnsaturatedFat:     "9g",
			TransFat:           "10g",
		},
		Times: models.Times{
			Prep:  5 * time.Minute,
			Cook:  1*time.Hour + 5*time.Minute,
			Total: 1*time.Hour + 10*time.Minute,
		},
		URL:   "https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/",
		Yield: 2,
	}
	_, _, _ = srv.Repository.AddRecipes(models.Recipes{recipe}, 1, nil)
	link, _ := srv.Repository.AddShareLink(models.Share{RecipeID: 1, CookbookID: -1, UserID: 1})

	t.Run("create valid share link", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodPost, uri(1))

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<input type="url" value="` + ts.URL + `/r/33320755-82f9-47e5-bb0a-d1b55cbd3f7b" class="input input-bordered w-full" readonly="readonly"></label>`,
			`<button id="copy_button" class="btn btn-neutral" title="Copy to clipboard" onClick="__templ_copyToClipboard`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("create invalid share link", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodPost, uri(10))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to create share link.","title":"General Error"}}`)
	})

	t.Run("access share link anonymous", func(t *testing.T) {
		rr := sendRequestNoBody(srv, http.MethodGet, link)

		assertStatus(t, rr.Code, http.StatusOK)
		body := getBodyHTML(rr)
		assertStringsInHTML(t, body, []string{
			`<title hx-swap-oob="true">` + recipe.Name + " | Recipya</title>",
			`<button title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release() then add @d=`,
			`<a href="/auth/login" class="btn btn-ghost">Log In</a> <a href="/auth/register" class="btn btn-ghost">Sign Up</a>`,
			`<span class="text-center pb-2 print:w-full" itemprop="name">Chicken Jersey</span>`,
			`<button class="mr-2 hidden sm:block" title="Print recipe" _="on click print()">`,
			`<img id="output" style="object-fit: cover" alt="Image of the recipe" class="w-full max-h-80 md:max-h-[34rem]" src="/static/img/recipes/placeholder.webp">`,
			`<div class="badge badge-primary badge-outline">American</div>`,
			`<p class="text-sm text-center">2 servings</p>`,
			`<a class="btn btn-sm btn-outline no-underline print:hidden" href="https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/" target="_blank">Source</a><p class="hidden print:block print:whitespace-nowrap print:overflow-hidden print:text-ellipsis print:max-w-xs">https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/</p>`,
			`<textarea class="textarea w-full h-full resize-none" readonly>This is the most delicious recipe!</textarea>`,
			`<p class="text-xs">Per 100g: calories 500 kcal; total carbohydrates 7 g; sugar 6 g; protein 3 g; total fat 8 g; saturated fat 4 g; unsaturated fat 9 g; trans fat 10 g; cholesterol 1 mg; sodium 5 mg; fiber 2 g</p>`,
			`<div class="grid grid-flow-col border-gray-700 col-span-6 py-1 md:border-y md:grid-cols-3 md:row-span-1 print:border-none"><div class="flex justify-self-center items-center gap-1 cursor-default" title="Prep time">`,
			`<time datetime="PT05M">5m</time></div><div class="flex justify-self-center items-center gap-1 cursor-default" title="Cooking time">`,
			`<time datetime="PT1H05M">1h05m</time></div><div class="flex justify-self-center items-center gap-1 cursor-default" title="Total time">`,
			`<time datetime="PT1H10M">1h10m</time></div></div>`,
			`<table class="table table-zebra table-xs print:hidden"><thead><tr><th>Nutrition (per 100g)</th><th>Amount</th></tr></thead> <tbody><tr><td>Calories:</td><td>500 kcal</td></tr><tr><td>Total carbs:</td><td>7 g</td></tr><tr><td>Sugars:</td><td>6 g</td></tr><tr><td>Protein:</td><td>3 g</td></tr><tr><td>Total fat:</td><td>8 g</td></tr><tr><td>Saturated fat:</td><td>4 g</td></tr><tr><td>Unsaturated fat:</td><td>9 g</td></tr><tr><td>Trans fat:</td><td>10 g</td></tr><tr><td>Cholesterol:</td><td>1 mg</td></tr><tr><td>Sodium:</td><td>5 mg</td></tr><tr><td>Fiber:</td><td>2 g</td></tr></tbody></table>`,
			`<div id="ingredients-instructions-container" class="grid text-sm md:grid-flow-col md:col-span-6"><div class="col-span-6 border-gray-700 px-4 py-2 border-y md:col-span-2 md:border-r md:border-y-0 print:hidden"><h2 class="font-semibold text-center underline pb-1">Ingredients</h2><ul><li class="form-control hover:bg-gray-100 dark:hover:bg-gray-700"><label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">Ing1</span></label></li><li class="form-control hover:bg-gray-100 dark:hover:bg-gray-700"><label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">Ing2</span></label></li><li class="form-control hover:bg-gray-100 dark:hover:bg-gray-700"><label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">Ing3</span></label></li></ul></div><div class="col-span-6 px-8 py-2 border-gray-700 md:rounded-bl-none md:col-span-4 print:hidden"><h2 class="font-semibold text-center underline pb-1">Instructions</h2><ol class="grid list-decimal"><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins1</span></li><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins2</span></li><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins3</span></li></ol></div></div><div class="hidden print:grid col-span-6 ml-2 my-1"><h1 class="text-sm print:mb-1"><b>Ingredients</b></h1><ol class="col-span-6 w-full print:mb-2" style="column-count: 1"><li class="text-sm"><label><input type="checkbox"></label> <span class="pl-2">Ing1</span></li><li class="text-sm"><label><input type="checkbox"></label> <span class="pl-2">Ing2</span></li><li class="text-sm"><label><input type="checkbox"></label> <span class="pl-2">Ing3</span></li></ol></div><div class="hidden col-span-5 overflow-visible print:inline"><h1 class="text-sm print:ml-2 print:mb-1"><b>Instructions</b></h1><ol class="col-span-6 list-decimal w-full ml-6"><li class="print:mr-4"><span class="text-sm whitespace-pre-line">Ins1</span></li><li class="print:mr-4"><span class="text-sm whitespace-pre-line">Ins2</span></li><li class="print:mr-4"><span class="text-sm whitespace-pre-line">Ins3</span></li></ol></div>`,
			`<h2 class="font-semibold text-center underline pb-1">Instructions</h2><ol class="grid list-decimal"><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins1</span></li><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins2</span></li><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins3</span></li></ol>`,
		})
		assertStringsNotInHTML(t, body, []string{
			`id="share-dialog"`,
		})
	})

	testcases := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
	}{
		{name: "other user Hx-Request", sendFunc: sendHxRequestAsLoggedInOtherNoBody},
		{name: "other user no Hx-Request", sendFunc: sendRequestAsLoggedInOtherNoBody},
	}
	for _, tc := range testcases {
		t.Run("access share link logged in "+tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, link)

			assertStatus(t, rr.Code, http.StatusOK)
			want := []string{
				`<title hx-swap-oob="true">` + recipe.Name + " | Recipya</title>",
				`<button title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release() then add @d=`,
				`<button class="mr-2" title="Add recipe to collection" hx-get="/recipes/1/share/add" hx-push-url="true">`,
			}
			notWant := []string{
				`id="share-dialog"`,
				`title="Share recipe"`,
			}
			body := getBodyHTML(rr)
			assertStringsInHTML(t, body, want)
			assertStringsNotInHTML(t, body, notWant)
		})
	}

	testcases2 := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
	}{
		{name: "host user Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody},
		{name: "host user no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody},
	}
	for _, tc := range testcases2 {
		t.Run("access share link logged "+tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, link)

			assertStatus(t, rr.Code, http.StatusOK)
			want := []string{
				`<title hx-swap-oob="true">` + recipe.Name + " | Recipya</title>",
				`<button title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release() then add @d=`,
				`<button class="mr-2 hidden sm:block" title="Print recipe" _="on click print()">`,
				`<button class="mr-2 hidden sm:block" hx-delete="/recipes/1" hx-swap="none" title="Delete recipe" hx-confirm="Are you sure you wish to delete this recipe?" hx-indicator="#fullscreen-loader">`,
			}
			notWant := []string{
				`id="share-dialog-result`,
				`title="Add recipe to collection`,
				`title="Share recipe`,
				`title="Add recipe to collection`,
			}
			body := getBodyHTML(rr)
			assertStringsInHTML(t, body, want)
			assertStringsNotInHTML(t, body, notWant)
		})
	}
}

func TestHandlers_Recipes_Duplicate(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	uri := func(id int64) string {
		return fmt.Sprintf("%s/recipes/%d/duplicate", ts.URL, id)
	}

	originalRepo := srv.Repository

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri(1))
	})

	t.Run("recipe does not exist", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri(10))

		assertStatus(t, rr.Code, http.StatusNotFound)
	})

	t.Run("recipe exists", func(t *testing.T) {
		recipe := models.Recipe{
			Category:     "American",
			Description:  "This is the most delicious recipe!",
			ID:           1,
			Images:       []uuid.UUID{uuid.New()},
			Ingredients:  []string{"Ing1", "Ing2", "Ing3"},
			Instructions: []string{"Ins1", "Ins2", "Ins3"},
			Keywords:     []string{"green sauce", "sheet pan meatballs"},
			Name:         "Chicken Jersey",
			Times: models.Times{
				Prep:  5 * time.Minute,
				Cook:  1*time.Hour + 5*time.Minute,
				Total: 1*time.Hour + 10*time.Minute,
			},
			URL:   "https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/",
			Yield: 2,
		}
		_, _, _ = srv.Repository.AddRecipes(models.Recipes{recipe}, 1, nil)
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri(1))

		assertStatus(t, rr.Code, http.StatusOK)
		ingredients := ""
		for _, ing := range recipe.Ingredients {
			ingredients += `<li class="pb-2"><div class="grid grid-flow-col items-center"><label><input required type="text" name="ingredients" value="` + ing + `" placeholder="1 cup of chopped onions" class="input input-bordered input-sm w-full" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)"></label><div class="ml-2"><button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: Enter" onclick="addItem(event)">+</button> <button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" _="on click if (closest <ol/>).childElementCount > 1 remove closest <li/> else set input to (closest <li/>).querySelector('input') then set input.value to '' then input.focus()">-</button><div class="inline-block h-4 cursor-move handle ml-2"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path></svg></div></div></div></li>`
		}
		instructions := ""
		for _, ins := range recipe.Instructions {
			instructions += `<li class="pt-2 md:pl-0"><div class="flex"><label class="w-11/12"><textarea required name="instructions" rows="3" class="textarea textarea-bordered w-full" placeholder="Mix all ingredients together" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)">` + ins + `</textarea></label><div class="grid ml-2"><button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: CTRL + Enter" onclick="addItem(event)">+</button> <button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" _="on click if (closest <ol/>).childElementCount > 1 remove closest <li/> else set input to (closest <li/>).querySelector('textarea') then set input.value to '' then input.focus()">-</button><div class="h-4 cursor-move handle grid place-content-center"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path></svg></div></div></div></li>`
		}
		assertStringsInHTML(t, getBodyHTML(rr), []string{
			`<title hx-swap-oob="true">Add Manual | Recipya</title>`,
			`<form class="card-body" style="padding: 0" enctype="multipart/form-data" hx-post="/recipes/add/manual" hx-indicator="#fullscreen-loader">`,
			`<input required type="text" name="title" placeholder="Title of the recipe*" autocomplete="off" class="input w-full btn-ghost text-center" value="` + recipe.Name + ` (copy)">`,
			`<div class="badge badge-sm badge-neutral p-3 pr-0"><input type="hidden" name="keywords" value="green sauce"> <span class="select-none">green sauce</span> <button type="button" class="btn btn-xs btn-ghost" _="on click remove closest <div/>">X</button></div><div class="badge badge-sm badge-neutral p-3 pr-0"><input type="hidden" name="keywords" value="sheet pan meatballs"> <span class="select-none">sheet pan meatballs</span> <button type="button" class="btn btn-xs btn-ghost" _="on click remove closest <div/>">X</button></div><div id="hidden_keyword" class="hidden badge badge-sm badge-neutral p-3 pr-0"><input type="hidden" name="keywords" value=""> <span class="select-none"></span> <button type="button" class="btn btn-xs btn-ghost" _="on click remove closest <div/>">X</button></div><div id="empty_keyword" class="badge badge-sm badge-neutral badge-outline p-3 pr-0" _="on keydown if event.key is 'Enter' halt the event then addKeyword(event)"><label><input id="new_keyword" type="text" placeholder="New keyword" class="input input-ghost input-xs w-[16ch] focus:outline-none" autocomplete="off" list="keywords"> <datalist id="keywords"><option>big</option></datalist></label> <button type="button" class="btn btn-xs btn-ghost" _="on click addKeyword(event)">&#10003;</button></div></div>`,
			`<input type="number" min="1" name="yield" value="` + strconv.FormatInt(int64(recipe.Yield), 10) + `" class="input input-bordered input-sm w-24 md:w-20 lg:w-24">`,
			`<input type="text" list="categories" name="category" class="input input-bordered input-sm w-48 md:w-36 lg:w-48" placeholder="Breakfast" autocomplete="off" value="` + recipe.Category + `"> <datalist id="categories"><option>breakfast</option><option>lunch</option><option>dinner</option></datalist>`,
			`<textarea name="description" placeholder="This Thai curry chicken will make you drool." class="textarea w-full h-full resize-none">` + recipe.Description + `</textarea>`,
			`<label><input type="text" name="time-preparation" value="00:05:00" class="input input-bordered input-xs max-w-24 html-duration-picker"></label>`,
			`<label><input type="text" name="time-cooking" value="01:05:00" class="input input-bordered input-xs max-w-24 html-duration-picker"></label>`,
			`<ol id="ingredients-list" class="pl-4 list-decimal">` + ingredients + `</ol></div>`,
			`<div class="col-span-6 px-6 py-2 border-gray-700 md:rounded-bl-none md:col-span-4"><h2 class="font-semibold text-center pb-2"><span class="underline">Instructions</span> <sup class="text-red-600">*</sup></h2><ol id="instructions-list" class="grid list-decimal">` + instructions + `</ol>`,
			`<div class="col-span-6 px-6 py-2 border-gray-700 md:rounded-bl-none md:col-span-4"><h2 class="font-semibold text-center pb-2"><span class="underline">Instructions</span> <sup class="text-red-600">*</sup></h2>`,
			`<ol id="instructions-list" class="grid list-decimal">` + instructions + `</ol>`,
			`<button class="btn btn-primary btn-block btn-sm">Submit</button>`,
		})
	})
}

func TestHandlers_Recipes_ShareAdd(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	uri := func(id int64) string {
		return fmt.Sprintf("%s/recipes/%d/share/add", ts.URL, id)
	}

	originalRepo := srv.Repository

	app.Config.Server.URL = "https://www.recipya.com"
	recipe := models.Recipe{
		Category:     "American",
		Description:  "This is the most delicious recipe!",
		ID:           1,
		Images:       []uuid.UUID{uuid.New()},
		Ingredients:  []string{"Ing1", "Ing2", "Ing3"},
		Instructions: []string{"Ins1", "Ins2", "Ins3"},
		Name:         "Chicken Jersey",
		Times: models.Times{
			Prep:  5 * time.Minute,
			Cook:  1*time.Hour + 5*time.Minute,
			Total: 1*time.Hour + 10*time.Minute,
		},
		URL:   "https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/",
		Yield: 2,
	}
	_, _, _ = srv.Repository.AddRecipes(models.Recipes{recipe}, 1, nil)

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri(1))
	})

	t.Run("recipe ID in path must be positive", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri(-1))

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})

	testcases := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
		isHxReq  bool
	}{
		{name: "is logged in Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody, isHxReq: true},
		{name: "is logged in no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody, isHxReq: false},
	}
	for _, tc := range testcases {
		t.Run("failed to add shared recipe "+tc.name, func(t *testing.T) {
			srv.Repository = &mockRepository{
				AddShareRecipeFunc: func(_, _ int64) (int64, error) {
					return 0, errors.New("failed to add shared recipe")
				},
			}
			defer func() {
				srv.Repository = originalRepo
			}()

			rr := tc.sendFunc(srv, http.MethodGet, uri(1))

			if tc.isHxReq {
				assertStatus(t, rr.Code, http.StatusInternalServerError)
				assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to add shared recipe to user's collection.","title":"General Error"}}`)
			} else {
				assertStatus(t, rr.Code, http.StatusSeeOther)
				assertHeader(t, rr, "Location", "/recipes")
			}
		})
	}

	testcases2 := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
		isHxReq  bool
	}{
		{name: "is logged in Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody, isHxReq: true},
		{name: "is logged in no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody, isHxReq: false},
	}
	for _, tc := range testcases2 {
		t.Run("recipe exists in collection"+tc.name, func(t *testing.T) {
			_, _, _ = srv.Repository.AddRecipes(models.Recipes{recipe}, 1, nil)
			defer func() {
				srv.Repository = originalRepo
			}()

			rr := tc.sendFunc(srv, http.MethodGet, uri(1))

			if tc.isHxReq {
				assertStatus(t, rr.Code, http.StatusInternalServerError)
				assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to add shared recipe to user's collection.","title":"General Error"}}`)
			} else {
				assertStatus(t, rr.Code, http.StatusSeeOther)
				assertHeader(t, rr, "Location", "/recipes")
			}
		})
	}

	testcases3 := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
		isHxReq  bool
	}{
		{name: "is logged in Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody, isHxReq: true},
		{name: "is logged in no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody, isHxReq: false},
	}
	for _, tc := range testcases3 {
		t.Run("valid request"+tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, uri(8))

			if tc.isHxReq {
				assertStatus(t, rr.Code, http.StatusOK)
				assertHeader(t, rr, "HX-Redirect", "/recipes/2")
			} else {
				assertStatus(t, rr.Code, http.StatusSeeOther)
				assertHeader(t, rr, "Location", "/recipes/2")
			}
		})
	}
}

func TestHandlers_Recipes_SupportedApplications(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/supported-applications"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("returns list of applications to logged in user", func(t *testing.T) {
		var want []string
		applications := [][]string{
			{"AccuChef", "https://www.accuchef.com"},
			{"ChefTap", "https://cheftap.com"},
			{"Crouton", "https://crouton.app"},
			{"Easy Recipe Deluxe", "https://easy-recipe-deluxe.software.informer.com"},
			{"Kalorio", "https://www.kalorio.de"},
			{"MasterCook", "https://www.mastercook.com"},
			{"Paprika", "https://www.paprikaapp.com"},
			{"Recipe Keeper", "https://recipekeeperonline.com"},
			{"RecipeSage", "https://recipesage.com"},
			{"Saffron", "https://www.mysaffronapp.com"},
		}
		for i, a := range applications {
			want = append(want, `<tr class="border text-center"><td class="border dark:border-gray-800">`+strconv.Itoa(i+1)+`</td><td class="border py-1 dark:border-gray-800"><a class="underline" href="`+a[1]+`" target="_blank">`+a[0]+`</a></td></tr>`)
		}

		rr := sendRequestAsLoggedInNoBody(srv, http.MethodGet, uri)

		assertStatus(t, rr.Code, http.StatusOK)
		assertHeader(t, rr, "Content-Type", "text/html")
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Recipes_SupportedWebsites(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/supported-websites"

	website1 := `<tr class="border text-center"><td class="border dark:border-gray-800">1</td><td class="border py-1 dark:border-gray-800"><a class="underline" href="https://101cookbooks.com" target="_blank">101cookbooks.com</a></td></tr>`
	website2 := `<tr class="border text-center"><td class="border dark:border-gray-800">2</td><td class="border py-1 dark:border-gray-800"><a class="underline" href="https://www.afghankitchenrecipes.com" target="_blank">afghankitchenrecipes.com</a></td></tr>`

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("returns list of websites to logged in user", func(t *testing.T) {
		rr := sendRequestAsLoggedInNoBody(srv, http.MethodGet, uri)

		assertStatus(t, rr.Code, http.StatusOK)
		assertHeader(t, rr, "Content-Type", "text/html")
		want := []string{website1, website2}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Recipes_View(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri+"/999")
	})

	t.Run("recipe is not in user collection", func(t *testing.T) {
		rr := sendRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"/999")

		assertStatus(t, rr.Code, http.StatusNotFound)
		assertStringsInHTML(t, getBodyHTML(rr), []string{
			`<title hx-swap-oob="true">Page Not Found | Recipya</title>`,
			"Page Not Found",
			"The page you requested to view is not found. Please go back to the main page.",
		})
	})

	testcases := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
	}{
		{name: "logged in Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody},
		{name: "logged in no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody},
	}
	for _, tc := range testcases {
		t.Run("recipe is in user's collection when "+tc.name, func(t *testing.T) {
			image, _ := uuid.Parse("e81ba735-a4af-4c66-8c17-2f2ccc1b1a95")
			r := models.Recipe{
				Category:     "American",
				Description:  "This is the most delicious recipe!",
				ID:           1,
				Images:       []uuid.UUID{image},
				Ingredients:  []string{"Ing1", "Ing2", "Ing3"},
				Instructions: []string{"Ins1", "Ins2", "Ins3"},
				Name:         "Chicken Jersey",
				Nutrition: models.Nutrition{
					Calories:           "500",
					Cholesterol:        "1mg",
					Fiber:              "2g",
					Protein:            "3g",
					SaturatedFat:       "4g",
					Sodium:             "5mg",
					Sugars:             "6g",
					TotalCarbohydrates: "7g",
					TotalFat:           "8g",
					UnsaturatedFat:     "9g",
					TransFat:           "10g",
				},
				Times: models.Times{
					Prep:  5 * time.Minute,
					Cook:  1*time.Hour + 5*time.Minute,
					Total: 1*time.Hour + 10*time.Minute,
				},
				URL:   "https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/",
				Yield: 2,
			}
			_, _, _ = srv.Repository.AddRecipes(models.Recipes{r}, 1, nil)

			rr := tc.sendFunc(srv, http.MethodGet, uri+"/1")

			assertStatus(t, rr.Code, http.StatusOK)
			assertStringsInHTML(t, getBodyHTML(rr), []string{
				`<title hx-swap-oob="true">` + r.Name + " | Recipya</title>",
				`<button title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release() then add @d='M 12.276 18.55 v -0.748 a 4.79 4.79 0 0 1 1.463 -3.458 a 5.763 5.763 0 0 0 1.804 -4.21 a 5.821 5.821 0 0 0 -6.475 -5.778 c -2.779 0.307 -4.99 2.65 -5.146 5.448 a 5.82 5.82 0 0 0 1.757 4.503 a 4.906 4.906 0 0 1 1.5 3.495 v 0.747 a 1.44 1.44 0 0 0 1.44 1.439 h 2.218 a 1.44 1.44 0 0 0 1.44 -1.439 z m -1.058 0 c 0 0.209 -0.17 0.38 -0.38 0.38 h -2.22 c -0.21 0 -0.38 -0.171 -0.38 -0.38 v -0.748 c 0 -1.58 -0.664 -3.13 -1.822 -4.254 A 4.762 4.762 0 0 1 4.98 9.863 c 0.127 -2.289 1.935 -4.204 4.205 -4.455 a 4.762 4.762 0 0 1 5.3 4.727 a 4.714 4.714 0 0 1 -1.474 3.443 a 5.853 5.853 0 0 0 -1.791 4.225 v 0.746 z M 11.45 20.51 H 8.006 a 0.397 0.397 0 1 0 0 0.795 h 3.444 a 0.397 0.397 0 1 0 0 -0.794 z M 11.847 22.162 a 0.397 0.397 0 0 0 -0.397 -0.397 H 8.006 a 0.397 0.397 0 1 0 0 0.794 h 3.444 c 0.22 0 0.397 -0.178 0.397 -0.397 z z z z z z z z M 10.986 23.416 H 8.867 a 0.397 0.397 0 1 0 0 0.794 h 1.722 c 0.22 0 0.397 -0.178 0.397 -0.397 z' to #icon-bulb else call initWakeLock() then add @d='M12.276 18.55v-.748a4.79 4.79 0 0 1 1.463-3.458 5.763 5.763 0 0 0 1.804-4.21 5.821 5.821 0 0 0-6.475-5.778c-2.779.307-4.99 2.65-5.146 5.448a5.82 5.82 0 0 0 1.757 4.503 4.906 4.906 0 0 1 1.5 3.495v.747a1.44 1.44 0 0 0 1.44 1.439h2.218a1.44 1.44 0 0 0 1.44-1.439zm-1.058 0c0 .209-.17.38-.38.38h-2.22c-.21 0-.38-.171-.38-.38v-.748c0-1.58-.664-3.13-1.822-4.254A4.762 4.762 0 0 1 4.98 9.863c.127-2.289 1.935-4.204 4.205-4.455a4.762 4.762 0 0 1 5.3 4.727 4.714 4.714 0 0 1-1.474 3.443 5.853 5.853 0 0 0-1.791 4.225v.746zM11.45 20.51H8.006a.397.397 0 1 0 0 .795h3.444a.397.397 0 1 0 0-.794zM11.847 22.162a.397.397 0 0 0-.397-.397H8.006a.397.397 0 1 0 0 .794h3.444c.22 0 .397-.178.397-.397zM.397 10.125h2.287a.397.397 0 1 0 0-.794H.397a.397.397 0 1 0 0 .794zM19.456 9.728a.397.397 0 0 0-.397-.397h-2.287a.397.397 0 1 0 0 .794h2.287c.22 0 .397-.178.397-.397zM9.331.397v2.287a.397.397 0 1 0 .794 0V.397a.397.397 0 1 0-.794 0zM16.045 2.85 14.43 4.465a.397.397 0 1 0 .561.561l1.617-1.617a.397.397 0 1 0-.562-.56zM5.027 14.429a.397.397 0 0 0-.56 0l-1.618 1.616a.397.397 0 1 0 .562.562l1.617-1.617a.397.397 0 0 0 0-.561zM4.466 5.027a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.56L3.41 2.848a.397.397 0 1 0-.561.561zM16.045 16.607a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.562L14.99 14.43a.397.397 0 1 0-.561.56zM10.986 23.416a.397.397 0 0 0-.397-.397H8.867a.397.397 0 1 0 0 .794h1.722c.22 0 .397-.178.397-.397z' to #icon-bulb end"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" width="24px" height="24px" stroke="currentColor"><path id="icon-bulb" d="M12.276 18.55v-.748a4.79 4.79 0 0 1 1.463-3.458 5.763 5.763 0 0 0 1.804-4.21 5.821 5.821 0 0 0-6.475-5.778c-2.779.307-4.99 2.65-5.146 5.448a5.82 5.82 0 0 0 1.757 4.503 4.906 4.906 0 0 1 1.5 3.495v.747a1.44 1.44 0 0 0 1.44 1.439h2.218a1.44 1.44 0 0 0 1.44-1.439zm-1.058 0c0 .209-.17.38-.38.38h-2.22c-.21 0-.38-.171-.38-.38v-.748c0-1.58-.664-3.13-1.822-4.254A4.762 4.762 0 0 1 4.98 9.863c.127-2.289 1.935-4.204 4.205-4.455a4.762 4.762 0 0 1 5.3 4.727 4.714 4.714 0 0 1-1.474 3.443 5.853 5.853 0 0 0-1.791 4.225v.746zM11.45 20.51H8.006a.397.397 0 1 0 0 .795h3.444a.397.397 0 1 0 0-.794zM11.847 22.162a.397.397 0 0 0-.397-.397H8.006a.397.397 0 1 0 0 .794h3.444c.22 0 .397-.178.397-.397zM.397 10.125h2.287a.397.397 0 1 0 0-.794H.397a.397.397 0 1 0 0 .794zM19.456 9.728a.397.397 0 0 0-.397-.397h-2.287a.397.397 0 1 0 0 .794h2.287c.22 0 .397-.178.397-.397zM9.331.397v2.287a.397.397 0 1 0 .794 0V.397a.397.397 0 1 0-.794 0zM16.045 2.85 14.43 4.465a.397.397 0 1 0 .561.561l1.617-1.617a.397.397 0 1 0-.562-.56zM5.027 14.429a.397.397 0 0 0-.56 0l-1.618 1.616a.397.397 0 1 0 .562.562l1.617-1.617a.397.397 0 0 0 0-.561zM4.466 5.027a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.56L3.41 2.848a.397.397 0 1 0-.561.561zM16.045 16.607a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.562L14.99 14.43a.397.397 0 1 0-.561.56zM10.986 23.416a.397.397 0 0 0-.397-.397H8.867a.397.397 0 1 0 0 .794h1.722c.22 0 .397-.178.397-.397z"></path></svg></button>`,
				`<button class="ml-2 hidden sm:block" title="Edit recipe" hx-get="/recipes/1/edit" hx-push-url="true" hx-target="#content" hx-swap="innerHTML transition:true"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg></button>`,
				`<span class="text-center pb-2 print:w-full" itemprop="name">Chicken Jersey</span>`,
				`<button class="mr-2 hidden sm:block" title="Share recipe" hx-post="/recipes/1/share" hx-target="#share-dialog-result" _="on htmx:afterRequest from me if event.detail.successful if navigator.canShare set name to document.querySelector('[itemprop=name]').textContent then set data to {title: name, text: name, url: document.querySelector('#share-dialog-result input').value} then call navigator.share(data) else call share_dialog.showModal() end">`,
				`<button class="mr-2 hidden sm:block" title="Print recipe" _="on click print()">`,
				`<button class="mr-2 hidden sm:block" hx-delete="/recipes/1" hx-swap="none" title="Delete recipe" hx-confirm="Are you sure you wish to delete this recipe?" hx-indicator="#fullscreen-loader">`,
				`<img id="output" style="object-fit: cover" alt="Image of the recipe" class="w-full max-h-80 md:max-h-[34rem]" src="/static/img/recipes/placeholder.webp">`,
				`<div class="badge badge-primary badge-outline">American</div>`,
				`<button class="mr-2 hidden sm:block" title="Share recipe" hx-post="/recipes/1/share" hx-target="#share-dialog-result" _="on htmx:afterRequest from me if event.detail.successful if navigator.canShare set name to document.querySelector('[itemprop=name]').textContent then set data to {title: name, text: name, url: document.querySelector('#share-dialog-result input').value} then call navigator.share(data) else call share_dialog.showModal() end"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" width="24px" height="24px" stroke="currentColor">`,
				`<form autocomplete="off" _="on submit halt the event" class="print:hidden"><label class="form-control w-full"><div class="label p-0"><span class="label-text">Servings</span></div><input id="yield" type="number" min="1" name="yield" value="2" class="input input-bordered input-sm w-24" hx-get="/recipes/1/scale" hx-trigger="input" hx-target="#ingredients-instructions-container"></label></form>`,
				`<a class="btn btn-sm btn-outline no-underline print:hidden" href="https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/" target="_blank">Source</a>`,
				`<textarea class="textarea w-full h-full resize-none" readonly>This is the most delicious recipe!</textarea>`,
				`<p class="text-xs">Per 100g: calories 500 kcal; total carbohydrates 7 g; sugar 6 g; protein 3 g; total fat 8 g; saturated fat 4 g; unsaturated fat 9 g; trans fat 10 g; cholesterol 1 mg; sodium 5 mg; fiber 2 g</p>`,
				`<div class="grid grid-flow-col border-gray-700 col-span-6 py-1 md:border-y md:grid-cols-3 md:row-span-1 print:border-none"><div class="flex justify-self-center items-center gap-1 cursor-default" title="Prep time">`,
				`<table class="table table-zebra table-xs print:hidden"><thead><tr><th>Nutrition (per 100g)</th><th>Amount</th></tr></thead> <tbody><tr><td>Calories:</td><td>500 kcal</td></tr><tr><td>Total carbs:</td><td>7 g</td></tr><tr><td>Sugars:</td><td>6 g</td></tr><tr><td>Protein:</td><td>3 g</td></tr><tr><td>Total fat:</td><td>8 g</td></tr><tr><td>Saturated fat:</td><td>4 g</td></tr><tr><td>Unsaturated fat:</td><td>9 g</td></tr><tr><td>Trans fat:</td><td>10 g</td></tr><tr><td>Cholesterol:</td><td>1 mg</td></tr><tr><td>Sodium:</td><td>5 mg</td></tr><tr><td>Fiber:</td><td>2 g</td></tr></tbody></table>`,
			})
		})
	}

	anImage1 := uuid.New()
	anImage2 := uuid.New()
	aVideo1 := uuid.New()
	aVideo2 := uuid.New()

	f1, _ := os.Create(filepath.Join(app.ImagesDir, anImage1.String()+app.ImageExt))
	f2, _ := os.Create(filepath.Join(app.ImagesDir, anImage2.String()+app.ImageExt))
	f3, _ := os.Create(filepath.Join(app.VideosDir, aVideo1.String()+app.VideoExt))
	f4, _ := os.Create(filepath.Join(app.VideosDir, aVideo2.String()+app.VideoExt))
	_, _ = f3.Write([]byte("video 1"))
	_, _ = f4.Write([]byte("video 2"))
	f1.Close()
	f2.Close()
	f3.Close()
	f4.Close()

	testcases2 := []struct {
		name   string
		images []uuid.UUID
		videos []uuid.UUID
		want   []string
	}{
		{
			name:   "0 image + 0 video",
			images: nil,
			videos: nil,
			want: []string{
				`<img style="object-fit: cover" alt="Image of the recipe" class="w-full max-h-80 md:max-h-[34rem]" src="/static/img/recipes/placeholder.webp">`,
			},
		},
		{
			name:   "1 image + 0 video",
			images: []uuid.UUID{anImage1},
			videos: nil,
			want: []string{
				`<img id="output" style="object-fit: cover" alt="Image of the recipe" class="w-full max-h-80 md:max-h-[34rem]" src="/data/images/` + anImage1.String() + `.webp">`,
			},
		},
		{
			name:   "2 images + 0 video",
			images: []uuid.UUID{anImage1, anImage2},
			videos: nil,
			want: []string{
				`<img style="object-fit: cover" alt="Image of the recipe" class="w-full max-h-80 md:max-h-[34rem]" src="/data/images/` + anImage1.String() + `.webp"><div class="absolute flex justify-between transform -translate-y-1/2 left-5 right-5 top-1/2"><a class="btn btn-circle" href="#media-1">❮</a> <a class="btn btn-circle" href="#media-1">❯</a></div></div><div id="media-1" class="carousel-item relative w-full">`,
				`<img style="object-fit: cover" alt="Image of the recipe" class="w-full max-h-80 md:max-h-[34rem]" src="/data/images/` + anImage2.String() + `.webp"><div class="absolute flex justify-between transform -translate-y-1/2 left-5 right-5 top-1/2"><a class="btn btn-circle" href="#media-0">❮</a> <a class="btn btn-circle" href="#media-0">❯</a></div></div></div></div>`,
			},
		},
		{
			name:   "0 image + 1 video",
			images: nil,
			videos: []uuid.UUID{aVideo1},
			want: []string{
				`<video controls preload="metadata" src="/data/videos/` + aVideo1.String() + `.webm" type="video/webm"></video>`,
			},
		},
		{
			name:   "0 image + 2 video",
			images: nil,
			videos: []uuid.UUID{aVideo1, aVideo2},
			want: []string{
				`<video controls preload="metadata" src="/data/videos/` + aVideo1.String() + `.webm" type="video/webm"></video><div class="absolute flex justify-between transform -translate-y-1/2 left-5 right-5 top-1/2"><a class="btn btn-circle" href="#media--1">❮</a> <a class="btn btn-circle" href="#media-1">❯</a></div></div><div id="media-1" class="carousel-item relative w-full">`,
				`<div id="media-1" class="carousel-item relative w-full"><video controls preload="metadata" src="/data/videos/` + aVideo2.String() + `.webm" type="video/webm"></video>`,
			},
		},
		{
			name:   "2 images + 2 videos",
			images: []uuid.UUID{anImage1, anImage2},
			videos: []uuid.UUID{aVideo1, aVideo2},
			want: []string{
				`<img style="object-fit: cover" alt="Image of the recipe" class="w-full max-h-80 md:max-h-[34rem]" src="/data/images/` + anImage1.String() + `.webp"><div class="absolute flex justify-between transform -translate-y-1/2 left-5 right-5 top-1/2"><a class="btn btn-circle" href="#media-3">❮</a> <a class="btn btn-circle" href="#media-1">❯</a></div></div><div id="media-1" class="carousel-item relative w-full">`,
				`<img style="object-fit: cover" alt="Image of the recipe" class="w-full max-h-80 md:max-h-[34rem]" src="/data/images/` + anImage2.String() + `.webp"><div class="absolute flex justify-between transform -translate-y-1/2 left-5 right-5 top-1/2"><a class="btn btn-circle" href="#media-0">❮</a> <a class="btn btn-circle" href="#media-2">❯</a></div></div><div id="media-2" class="carousel-item relative w-full">`,
				`<video controls preload="metadata" src="/data/videos/` + aVideo1.String() + `.webm" type="video/webm"></video><div class="absolute flex justify-between transform -translate-y-1/2 left-5 right-5 top-1/2"><a class="btn btn-circle" href="#media-1">❮</a> <a class="btn btn-circle" href="#media-3">❯</a></div></div><div id="media-3" class="carousel-item relative w-full">`,
				`<video controls preload="metadata" src="/data/videos/` + aVideo2.String() + `.webm" type="video/webm"></video>`,
			},
		},
	}
	for _, tc := range testcases2 {
		t.Run("media "+tc.name, func(t *testing.T) {
			videos := make([]models.VideoObject, 0, len(tc.videos))
			for _, v := range tc.videos {
				videos = append(videos, models.VideoObject{ID: v})
			}

			srv.Repository = &mockRepository{
				RecipesRegistered: map[int64]models.Recipes{
					1: {
						models.Recipe{
							Category:     "American",
							ID:           1,
							Images:       tc.images,
							Ingredients:  []string{"Ing1", "Ing2", "Ing3"},
							Instructions: []string{"Ins1", "Ins2", "Ins3"},
							Name:         "Chicken Jersey",
							Times: models.Times{
								Prep:  5 * time.Minute,
								Cook:  1*time.Hour + 5*time.Minute,
								Total: 1*time.Hour + 10*time.Minute,
							},
							URL:    "https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/",
							Videos: videos,
							Yield:  2,
						},
					},
				},
			}

			rr := sendRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"/1")

			assertStatus(t, rr.Code, http.StatusOK)
			assertStringsInHTML(t, getBodyHTML(rr), tc.want)
		})
	}

	for _, file := range []*os.File{f1, f2, f3, f4} {
		os.Remove(file.Name())
	}
}
