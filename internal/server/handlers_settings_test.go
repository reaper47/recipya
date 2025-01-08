package server_test

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services"
	"github.com/reaper47/recipya/internal/units"
)

func TestHandlers_Settings(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	var (
		uri                = ts.URL + "/settings"
		originalRepo       = srv.Repository
		measurementErrFunc = func(_ int64) ([]units.System, models.UserSettings, error) {
			return nil, models.UserSettings{}, errors.New("kerch bridge on fire... Your defence is terrified")
		}
	)

	app.Config = app.ConfigFile{}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("error fetch measurement systems with htmx", func(t *testing.T) {
		app.Config.Server.IsAutologin = true
		srv.Repository = &mockRepository{
			MeasurementSystemsFunc: measurementErrFunc,
		}
		defer func() {
			app.Config.Server.IsAutologin = false
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Error fetching unit systems: kerch bridge on fire... Your defence is terrified","title":"Database Error"}}`
		assertWebsocket(t, c, 1, want)
	})

	t.Run("error fetch measurement systems without htmx", func(t *testing.T) {
		app.Config.Server.IsAutologin = true
		srv.Repository = &mockRepository{
			MeasurementSystemsFunc: measurementErrFunc,
		}
		defer func() {
			app.Config.Server.IsAutologin = false
			srv.Repository = originalRepo
		}()

		rr := sendRequestAsLoggedInNoBody(srv, http.MethodGet, uri)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertStringsInHTML(t, getBodyHTML(rr), []string{"Error fetching unit systems: kerch bridge on fire... Your defence is terrified"})
	})

	t.Run("server and connections settings not displayed when not admin", func(t *testing.T) {
		srv.Repository = &mockRepository{
			categories: map[int64][]string{2: {"breakfast"}},
			UsersRegistered: []models.User{
				{ID: 1, Email: "admin@admin.com"},
				{ID: 2, Email: "yay@nay.com"},
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedInOtherNoBody(srv, http.MethodGet, uri)

		assertStatus(t, rr.Code, http.StatusOK)
		body := getBodyHTML(rr)
		assertStringsInHTML(t, body, []string{
			`<div id="settings_recipes"`,
			`<div id="settings_data"`,
			`<div id="settings_account"`,
			`<div id="settings_about"`,
		})
		assertStringsNotInHTML(t, body, []string{
			`<a class="setting-tab" _="on click add .hidden to the children of #settings_blocks then remove .hidden from #settings_server">`,
			`<a class="setting-tab" _="on click add .hidden to the children of #settings_blocks then remove .hidden from #settings_connections">`,
			`<div id="settings_server"`,
			`<div id="settings_connections"`,
		})
	})

	t.Run("demo sees fake connections data", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInOtherNoBody(srv, http.MethodGet, uri)

		assertStringsNotInHTML(t, getBodyHTML(rr), []string{
			`<input name="email.from" type="text" placeholder="SendGrid email" value="demo@demo.com" autocomplete="off" class="input input-bordered input-sm w-full">`,
			`<input name="email.apikey" type="text" placeholder="API key" value="demo" autocomplete="off" class="input input-bordered input-sm w-full">`,
			`<input name="integrations.ocr.key" type="text" placeholder="Resource key 1" value="demo" autocomplete="off" class="input input-bordered input-sm w-full">`,
			`<input name="integrations.ocr.url" type="url" placeholder="Vision endpoint URL" value="https://www.example.com" autocomplete="off" class="input input-bordered input-sm w-full">`,
		})
	})

	t.Run("application update available", func(t *testing.T) {
		app.Info.IsUpdateAvailable = true
		defer func() {
			app.Info.IsUpdateAvailable = false
		}()

		rr := sendRequestAsLoggedInNoBody(srv, http.MethodGet, uri)

		assertStatus(t, rr.Code, http.StatusOK)
		assertStringsInHTML(t, getBodyHTML(rr), []string{
			`<div><p class="font-semibold">Recipya Version</p><p class="text-sm mt-2">v1.3.0 (update available)</p><p class="text-xs">Last checked: 0001-01-01<br>Last updated: 0001-01-01<br><br>Read the <a class="link" href="https://recipya.musicavis.ca/about/changelog/v1.3.0" target="_blank">release notes</a></p></div><div class="flex flex-row self-start"><button class="btn btn-sm" hx-get="/update" hx-swap="none" hx-indicator="#fullscreen-loader">Update</button></div></div>`,
		})
	})

	t.Run("display settings", func(t *testing.T) {
		xc, _ := srv.Repository.Categories(1)
		srv.Repository = &mockRepository{
			categories: map[int64][]string{1: xc},
			RecipesRegistered: map[int64]models.Recipes{
				1: {
					{ID: 1, Name: "Chinese Firmware", Category: "breakfast"},
					{ID: 2, Name: "Lovely Canada", Category: "lunch"},
					{ID: 3, Name: "Lovely Ukraine", Category: "dinner"},
					{ID: 4, Name: "Space Disco", Category: "snack"},
					{ID: 5, Name: "Maple Pancakes", Category: "breakfast"},
					{ID: 6, Name: "Maple Pancakes", Category: "uncategorized"},
				},
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<div class="flex flex-col menu-sm sm:flex-row sm:menu-md">`,
			`<ul class="menu menu-horizontal pt-0 flex-nowrap overflow-x-auto sm:w-56 sm:overflow-x-clip sm:flex-wrap sm:menu sm:menu-vertical" _="on click remove .active from .setting-tab then add .active to closest <a/> to event.target">`,
			`<li><a class="setting-tab active" _="on click add .hidden to the children of #settings_blocks then remove .hidden from #settings_recipes"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="m21 7.5-2.25-1.313M21 7.5v2.25m0-2.25-2.25 1.313M3 7.5l2.25-1.313M3 7.5l2.25 1.313M3 7.5v2.25m9 3 2.25-1.313M12 12.75l-2.25-1.313M12 12.75V15m0 6.75 2.25-1.313M12 21.75V19.5m0 2.25-2.25-1.313m0-16.875L12 2.25l2.25 1.313M21 14.25v2.25l-2.25 1.313m-13.5 0L3 16.5v-2.25"></path></svg>Recipes</a></li>`,
			`<li><a class="setting-tab" _="on click add .hidden to the children of #settings_blocks then remove .hidden from #settings_connections"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15a4.5 4.5 0 0 0 4.5 4.5H18a3.75 3.75 0 0 0 1.332-7.257 3 3 0 0 0-3.758-3.848 5.25 5.25 0 0 0-10.233 2.33A4.502 4.502 0 0 0 2.25 15Z"></path></svg>Connections</a></li>`,
			`<li><a class="setting-tab" _="on click add .hidden to the children of #settings_blocks then remove .hidden from #settings_data"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M20.25 6.375c0 2.278-3.694 4.125-8.25 4.125S3.75 8.653 3.75 6.375m16.5 0c0-2.278-3.694-4.125-8.25-4.125S3.75 4.097 3.75 6.375m16.5 0v11.25c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125V6.375m16.5 0v3.75m-16.5-3.75v3.75m16.5 0v3.75C20.25 16.153 16.556 18 12 18s-8.25-1.847-8.25-4.125v-3.75m16.5 0c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125"></path></svg>Data</a></li>`,
			`<li><a class="setting-tab" _="on click add .hidden to the children of #settings_blocks then remove .hidden from #settings_server"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M21.75 17.25v-.228a4.5 4.5 0 0 0-.12-1.03l-2.268-9.64a3.375 3.375 0 0 0-3.285-2.602H7.923a3.375 3.375 0 0 0-3.285 2.602l-2.268 9.64a4.5 4.5 0 0 0-.12 1.03v.228m19.5 0a3 3 0 0 1-3 3H5.25a3 3 0 0 1-3-3m19.5 0a3 3 0 0 0-3-3H5.25a3 3 0 0 0-3 3m16.5 0h.008v.008h-.008v-.008Zm-3 0h.008v.008h-.008v-.008Z"></path></svg>Server</a></li>`,
			`<li><a class="setting-tab" _="on click add .hidden to the children of #settings_blocks then remove .hidden from #settings_account"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"></path></svg>Account</a></li>`,
			`<li><a class="setting-tab" _="on click add .hidden to the children of #settings_blocks then remove .hidden from #settings_about"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z"></path></svg>About</a></li></ul>`,
			`<div id="settings_blocks" class="w-full md:h-[26rem] md:max-h-[26rem]" style="padding-right: 1rem">`,
			`<div id="settings_recipes" class="p-3 md:p-0 md:pr-4 md:max-h-96 overflow-y-auto"><div class="flex justify-between items-center text-sm"><details class="w-full"><summary class="font-semibold cursor-default">Categories</summary><div class="flex flex-wrap gap-2 p-2"><div class="badge badge-outline p-3 pr-0"><form class="inline-flex" hx-delete="/recipes/categories" hx-target="closest <div/>" hx-swap="delete"><input type="hidden" name="category" value="breakfast"> <span class="select-none">breakfast</span> <button type="submit" class="btn btn-xs btn-ghost">X</button></form></div><div class="badge badge-outline p-3 pr-0"><form class="inline-flex" hx-delete="/recipes/categories" hx-target="closest <div/>" hx-swap="delete"><input type="hidden" name="category" value="lunch"> <span class="select-none">lunch</span> <button type="submit" class="btn btn-xs btn-ghost">X</button></form></div><div class="badge badge-outline p-3 pr-0"><form class="inline-flex" hx-delete="/recipes/categories" hx-target="closest <div/>" hx-swap="delete"><input type="hidden" name="category" value="dinner"> <span class="select-none">dinner</span> <button type="submit" class="btn btn-xs btn-ghost">X</button></form></div><div class="badge badge-outline p-3 pr-0"><form class="inline-flex" hx-post="/recipes/categories" hx-target="closest <div/>" hx-swap="outerHTML"><label class="form-control"><input required type="text" placeholder="New category" class="input input-ghost input-xs w-[16ch] focus:outline-none" name="category" autocomplete="off"></label> <button class="btn btn-xs btn-ghost">&#10003;</button></form></div></div></details></div><div class="divider m-0"></div><div class="flex justify-between items-center text-sm"><label for="settings_recipes_measurement_system" class="font-semibold">Measurement system</label> <select id="settings_recipes_measurement_system" name="system" class="w-fit select select-bordered select-sm" hx-post="/settings/measurement-system" hx-swap="none"><option value="imperial">imperial</option><option value="metric" selected>metric</option></select></div><div class="flex justify-between items-center text-sm mt-2"><label for="settings_recipes_convert"><span class="font-semibold">Convert automatically</span><br><span class="text-xs">Convert new recipes to your preferred measurement system.</span></label> <input type="checkbox" name="convert" id="settings_recipes_convert" class="checkbox" hx-post="/settings/convert-automatically" hx-trigger="click"></div><div class="divider m-0"></div><div class="flex justify-between items-center text-sm mt-2"><label for="settings_recipes_calc_nutrition"><span class="font-semibold">Calculate nutrition facts</span><br><span class="text-xs block max-w-[45ch]">Calculate the nutrition facts automatically when adding a recipe. The processing will be done in the background.</span></label> <input id="settings_recipes_calc_nutrition" type="checkbox" name="calculate-nutrition" class="checkbox" hx-post="/settings/calculate-nutrition" hx-trigger="click"></div><div class="divider m-0"></div><div class="flex justify-between items-center text-sm"><details class="w-full"><summary class="font-semibold cursor-default">Placeholders</summary><div class="flex flex-wrap gap-2 p-2 flex-row"><div class="max-w-60"><p class="text-center mb-1 font-medium underline">Recipe</p><form hx-post="/placeholder" hx-encoding="multipart/form-data" hx-swap="none" _="on htmx:afterRequest call reloadImg('/data/images/Placeholders/placeholder.recipe.webp')"><img src="/data/images/Placeholders/placeholder.recipe.webp" alt="Recipe placeholder" class="w-60 h-60"> <input type="hidden" name="name" value="recipe"> <input type="file" name="images" class="file-input file-input-bordered file-input-sm max-w-60 mt-1"> <button class="btn btn-neutral btn-sm btn-block my-1">Update</button></form><button class="btn btn-error btn-sm btn-block" hx-post="/placeholder/restore" hx-vals="js:{t: "recipe"}" hx-swap="none" _="on htmx:afterRequest call reloadImg('/data/images/Placeholders/placeholder.recipe.webp')">Restore original</button></div><div class="max-w-60"><p class="text-center mb-1 font-medium underline">Cookbook</p><form hx-post="/placeholder" hx-encoding="multipart/form-data" hx-swap="none" _="on htmx:afterRequest call reloadImg('/data/images/Placeholders/placeholder.cookbook.webp')"><img src="/data/images/Placeholders/placeholder.cookbook.webp" alt="Cookbook placeholder" class="w-60 h-60"> <input type="hidden" name="name" value="cookbook"> <input type="file" name="images" class="file-input file-input-bordered file-input-sm max-w-60 mt-1"> <button class="btn btn-neutral btn-sm btn-block my-1">Update</button></form><button class="btn btn-error btn-sm btn-block" hx-post="/placeholder/restore" hx-vals="js:{name: "cookbook"}" hx-swap="none" _="on htmx:afterRequest call reloadImg('/data/images/Placeholders/placeholder.cookbook.webp')">Restore original</button></div></div></details></div>`,
			`<div id="settings_connections" class="p-3 overflow-y-auto max-h-96 hidden md:p-0 md:pr-4"><div class="flex justify-between items-center text-sm"><details class="w-full"><summary class="font-semibold cursor-default">Twilio SendGrid<br><span class="text-xs font-normal">This connection is used to send emails.</span></summary><form class="grid w-full" hx-put="/settings/config" hx-swap="none"><label class="form-control w-full"><span class="label"><span class="label-text text-sm">From</span></span> <input name="email.from" type="text" placeholder="SendGrid email" value="" autocomplete="off" class="input input-bordered input-sm w-full"></label> <label class="form-control w-full"><span class="label"><span class="label-text text-sm">SendGrid API key</span></span> <input name="email.apikey" type="text" placeholder="API key" value="" autocomplete="off" class="input input-bordered input-sm w-full"></label> <button class="btn btn-sm mt-2">Update</button></form></details> <button type="button" title="Test connection" class="btn btn-xs float-right self-baseline" hx-get="/integrations/test-connection?api=sg" hx-swap="none"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99"></path></svg></button></div><div class="divider m-0"></div><div class="flex justify-between items-center text-sm"><details class="w-full"><summary class="font-semibold cursor-default">Azure AI Document Intelligence<br><span class="text-xs font-normal">This connection is used to digitize recipe images.</span></summary><form class="grid w-full" hx-put="/settings/config" hx-swap="none"><label class="form-control w-full"><span class="label"><span class="label-text text-sm">Resource key</span></span> <input name="integrations.ocr.key" type="text" placeholder="Resource key 1" value="" autocomplete="off" class="input input-bordered input-sm w-full"></label> <label class="form-control w-full"><span class="label"><span class="label-text text-sm">Endpoint</span></span> <input name="integrations.ocr.url" type="url" placeholder="Vision endpoint URL" value="" autocomplete="off" class="input input-bordered input-sm w-full"></label> <button class="btn btn-sm mt-2">Update</button></form></details> <button type="button" title="Test connection" class="btn btn-xs float-right self-baseline" hx-get="/integrations/test-connection?api=azure-di" hx-swap="none"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99"></path></svg></button></div></div>`,
			`<div id="settings_server" class="hidden p-3 md:p-0 md:pr-4 md:max-h-96"><div class="flex justify-between items-center text-sm"><form class="grid w-full" hx-put="/settings/config" hx-swap="none"><p class="font-semibold">Configuration</p><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Autologin</span> <input name="server.autologin" type="checkbox" class="checkbox"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">No signups</span> <input name="server.noSignups" type="checkbox" class="checkbox"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Is production</span> <input name="server.production" type="checkbox" class="checkbox"></label></div><button class="btn btn-sm mt-2">Update</button></form></div></div>`,
			`<div id="settings_data" class="hidden p-3 md:p-0 md:pr-4"><div class="flex justify-between items-center text-sm"><details class="w-full"><summary class="font-semibold cursor-default">Import data<br><span class="text-xs font-normal">Import from Mealie, Tandoor, Nextcloud, etc.</span></summary><form class="flex flex-col text-sm" hx-post="/integrations/import" hx-swap="none"><label class="form-control w-full"><span class="label"><span class="label-text text-sm">Solution</span></span> <select name="integration" class="w-fit select select-bordered select-sm"><option value="mealie" selected>Mealie</option> <option value="nextcloud">Nextcloud</option> <option value="tandoor">Tandoor</option></select></label> <label class="form-control w-full"><span class="label"><span class="label-text text-sm">Base URL</span></span> <input type="url" name="url" placeholder="https://instance.mydomain.com" class="input input-bordered input-sm w-full" required></label> <label class="form-control w-full"><span class="label"><span class="label-text text-sm">Username</span></span> <input type="text" name="username" placeholder="Enter your username" class="input input-bordered input-sm w-full" required></label> <label class="form-control w-full"><span class="label"><span class="label-text text-sm">Password</span></span> <input type="password" name="password" placeholder="Enter your password" class="input input-bordered input-sm w-full" required></label> <button class="btn btn-sm mt-2"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" class="bi bi-cloud-arrow-down" viewBox="0 0 16 16"><path fill-rule="evenodd" d="M7.646 10.854a.5.5 0 0 0 .708 0l2-2a.5.5 0 0 0-.708-.708L8.5 9.293V5.5a.5.5 0 0 0-1 0v3.793L6.354 8.146a.5.5 0 1 0-.708.708l2 2z"></path> <path d="M4.406 3.342A5.53 5.53 0 0 1 8 2c2.69 0 4.923 2 5.166 4.579C14.758 6.804 16 8.137 16 9.773 16 11.569 14.502 13 12.687 13H3.781C1.708 13 0 11.366 0 9.318c0-1.763 1.266-3.223 2.942-3.593.143-.863.698-1.723 1.464-2.383zm.653.757c-.757.653-1.153 1.44-1.153 2.056v.448l-.445.049C2.064 6.805 1 7.952 1 9.318 1 10.785 2.23 12 3.781 12h8.906C13.98 12 15 10.988 15 9.773c0-1.216-1.02-2.228-2.313-2.228h-.5v-.5C12.188 4.825 10.328 3 8 3a4.53 4.53 0 0 0-2.941 1.1z"></path></svg>Import</button></form></details></div><div class="divider m-0"></div><div class="flex justify-between items-center text-sm"><div><p class="font-semibold">Export data</p><p class="text-xs">Download your data in the selected file format.</p></div><form class="grid gap-1 grid-flow-col w-fit" hx-get="/settings/export/recipes" hx-include="select[name='type']" hx-swap="none"><label class="form-control w-full max-w-xs"><select required id="file-type" name="type" class="w-fit select select-bordered select-sm"><optgroup label="Recipes"><option value="json" selected>JSON</option> <option value="pdf">PDF</option></optgroup></select></label> <button class="btn btn-outline btn-sm"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 ml-1" fill="black" viewBox="0 0 24 24" stroke="currentColor"><path d="M16 11v5H2v-5H0v5a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-5z"></path> <path d="m9 14 5-6h-4V0H8v8H4z"></path></svg></button></form></div></div>`,
			`<div id="settings_account" class="hidden p-3 md:p-0 md:pr-4 md:max-h-96"><div><div class="flex justify-between items-center text-sm"><div><p class="font-semibold">Theme</p><p class="font-normal text-sm">Select your preferred theme.</p></div><div id="themes_palette" class="dropdown dropdown-end hidden z-30 [@supports(color:oklch(0%_0_0))]:block" _="on load call themeChange(document.querySelector('#theme_palette'))"><div tabindex="0" role="button" class="btn btn-ghost"><svg width="20" height="20" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="h-5 w-5 stroke-current md:hidden"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"></path></svg> <span id="theme_name" class="hidden font-normal md:inline" _="on load set theme to localStorage.getItem('theme') then if not theme put 'system' into me else put theme into me">Theme</span> <svg width="12px" height="12px" class="hidden h-2 w-2 fill-current opacity-60 sm:inline-block" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2048 2048"><path d="M1799 349l242 241-1017 1017L7 590l242-241 775 775 775-775z"></path></svg></div><div tabindex="0" class="dropdown-content bg-base-200 text-base-content rounded-box top-px h-[28.6rem] max-h-[calc(100vh-10rem)] w-56 overflow-y-auto border border-white/5 shadow-2xl outline outline-1 outline-black/5 mt-16"><div class="grid grid-cols-1 gap-3 p-3"><button class="outline-base-content text-stbbcgoodfood.com/recipesart outline-offset-4 [&amp;_svg]:visible" data-act-class="[&amp;_svg]:visible" data-set-theme="" _="on click put 'system' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme=""><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">system</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="light" _="on click put 'light' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="light"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg id="light_checkmark" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">light</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="dark" _="on click put 'dark' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="dark"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">dark</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="cupcake" _="on click put 'cupcake' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="cupcake"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">cupcake</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="bumblebee" _="on click put 'bumblebee' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="bumblebee"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">bumblebee</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="emerald" _="on click put 'emerald' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="emerald"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">emerald</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="corporate" _="on click put 'corporate' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="corporate"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">corporate</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="synthwave" _="on click put 'synthwave' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="synthwave"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">synthwave</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="retro" _="on click put 'retro' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="retro"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">retro</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="cyberpunk" _="on click put 'cyberpunk' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="cyberpunk"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">cyberpunk</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="valentine" _="on click put 'valentine' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="valentine"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">valentine</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="halloween" _="on click put 'halloween' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="halloween"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">halloween</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="garden" _="on click put 'garden' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="garden"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">garden</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="forest" _="on click put 'forest' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="forest"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">forest</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="aqua" _="on click put 'aqua' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="aqua"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">aqua</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="lofi" _="on click put 'lofi' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="lofi"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">lofi</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="pastel" _="on click put 'pastel' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="pastel"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">pastel</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="fantasy" _="on click put 'fantasy' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="fantasy"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">fantasy</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="wireframe" _="on click put 'wireframe' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="wireframe"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">wireframe</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="black" _="on click put 'black' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="black"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">black</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="luxury" _="on click put 'luxury' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="luxury"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">luxury</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="dracula" _="on click put 'dracula' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="dracula"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">dracula</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="cmyk" _="on click put 'cmyk' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="cmyk"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">cmyk</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="autumn" _="on click put 'autumn' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="autumn"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">autumn</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="business" _="on click put 'business' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="business"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">business</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="acid" _="on click put 'acid' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="acid"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">acid</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="lemonade" _="on click put 'lemonade' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="lemonade"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">lemonade</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="night" _="on click put 'night' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="night"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">night</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="coffee" _="on click put 'coffee' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="coffee"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">coffee</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="winter" _="on click put 'winter' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="winter"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">winter</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="dim" _="on click put 'dim' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="dim"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">dim</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="nord" _="on click put 'nord' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="nord"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">nord</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <button class="outline-base-content text-start outline-offset-4" data-act-class="[&amp;_svg]:visible" data-set-theme="sunset" _="on click put 'sunset' into #theme_name"><span class="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans" data-theme="sunset"><span class="grid grid-cols-5 grid-rows-3"><span class="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="invisible h-3 w-3 shrink-0"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path></svg> <span class="flex-grow text-sm">sunset</span> <span class="flex h-full shrink-0 flex-wrap gap-1"><span class="bg-primary rounded-badge w-2"></span> <span class="bg-secondary rounded-badge w-2"></span> <span class="bg-accent rounded-badge w-2"></span> <span class="bg-neutral rounded-badge w-2"></span></span></span></span></span></button> <a class="outline-base-content overflow-hidden rounded-lg text-center" href="/theme-generator/"><p class="px-2 text-xs">Credits to DaisyUI for this list</p></a></div></div></div></div></div><div class="divider m-0"></div><div class="flex justify-between items-center text-sm"><details class="w-full"><summary class="font-semibold cursor-default">Change password</summary><form class="flex flex-col text-sm" hx-post="/auth/change-password" hx-indicator="#fullscreen-loader" hx-swap="none"><label class="form-control w-full"><span class="label"><span class="label-text text-sm">Current password</span></span> <input type="password" placeholder="Enter current password" class="input input-bordered input-sm w-full" name="password-current" required></label> <label class="form-control w-full"><span class="label"><span class="label-text text-sm">New password</span></span> <input type="password" placeholder="Enter new password" class="input input-bordered input-sm w-full" name="password-new" required></label> <label class="form-control w-full"><span class="label"><span class="label-text text-sm">Confirm password</span></span> <input type="password" placeholder="Retype new password" class="input input-bordered input-sm w-full" name="password-confirm" required></label> <button class="btn btn-sm mt-2">Update password</button></form></details></div><div class="divider m-0"></div><div><div class="flex justify-between items-center text-sm"><div><p class="font-semibold">Delete Account</p><p class="font-normal text-sm">This will delete all your data.</p></div><button type="submit" class="btn btn-sm" hx-delete="/auth/user" hx-confirm="Are you sure you want to delete your account? This action is irreversible.">Delete</button></div></div></div>`,
			`<div id="settings_about" class="p-3 md:p-0 md:pr-4 hidden"><div><div class="flex justify-between items-center text-sm"><div><p class="font-semibold">Recipya Version</p><p class="text-sm mt-2">v1.3.0 (latest)</p><p class="text-xs">Last checked: 0001-01-01<br>Last updated: 0001-01-01<br><br>Read the <a class="link" href="https://recipya.musicavis.ca/about/changelog/v1.3.0" target="_blank">release notes</a></p></div><div class="flex flex-row self-start"><img id="settings_about_update_check" class="htmx-indicator mr-1" src="/static/img/bars.svg" alt="Checking..."> <button class="btn btn-sm" hx-get="/update/check" hx-target="#settings_about" hx-swap="outerHTML" hx-indicator="#settings_about_update_check">Check for updates</button></div></div></div><div class="divider m-0"></div><div class="flex space-x-1"><a href="https://app.element.io/#/room/#recipya:matrix.org"><img alt="Support" src="https://img.shields.io/badge/Element-Recipya-blue?logo=element&amp;logoColor=white"></a> <a href="https://github.com/reaper47/recipya" target="_blank"><img alt="Github Repo" src="https://img.shields.io/github/stars/reaper47/recipya?style=social&amp;label=Star on Github"></a></div></div>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Settings_BackupsRestore(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	originalFiles := srv.Files
	originalRepo := srv.Repository

	uri := ts.URL + "/settings/backups/restore"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	testcases := []struct {
		name string
		in   string
	}{
		{name: "empty body", in: ""},
		{name: "date is empty", in: "date="},
		{name: "date is invalid", in: "date=01/02/1999"},
		{name: "date contains letters", in: "date=01 January 2024"},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader(tc.in))

			assertStatus(t, rr.Code, http.StatusBadRequest)
			_, after, _ := strings.Cut(tc.in, "=")
			message := fmt.Sprintf(`{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"%s is an invalid backup.","title":"Form Error"}}`, after)
			assertWebsocket(t, c, 1, message)
		})
	}

	t.Run("backup user data failed", func(t *testing.T) {
		srv.Files = &mockFiles{
			backupUserDataFunc: func(_ services.RepositoryService, _ int64) error {
				return errors.New("could not backup data")
			},
		}
		defer func() {
			srv.Files = originalFiles
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("date=2006-01-02"))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to backup current data.","title":"Files Error"}}`)
	})

	t.Run("extract user backup failed", func(t *testing.T) {
		srv.Files = &mockFiles{
			extractUserBackupFunc: func(date string, userID int64) (*models.UserBackup, error) {
				return nil, errors.New("backup failed")
			},
		}
		defer func() {
			srv.Files = originalFiles
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("date=2006-01-02"))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to extract backup.","title":"Files Error"}}`)
	})

	t.Run("restore backup failed", func(t *testing.T) {
		srv.Repository = &mockRepository{
			RestoreUserBackupFunc: func(_ *models.UserBackup) error {
				return errors.New("restore failed")
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("date=2006-01-02"))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to restore backup.","title":"Database Error"}}`)
	})

	t.Run("valid request", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("date=2006-01-02"))

		assertStatus(t, rr.Code, http.StatusOK)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-info","message":"","title":"Backup restored successfully."}}`)
	})
}

func TestHandlers_Settings_CalculateNutrition(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	srv.Repository = &mockRepository{
		UserSettingsRegistered: map[int64]*models.UserSettings{
			1: {CalculateNutritionFact: false},
		},
	}

	uri := ts.URL + "/settings/calculate-nutrition"
	originalRepo := srv.Repository

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("error updating the setting", func(t *testing.T) {
		srv.Repository = &mockRepository{
			UpdateCalculateNutritionFunc: func(userID int64, isEnabled bool) error {
				return errors.New("whoops")
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("calculate-nutrition=off"))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to set setting.","title":"Database Error"}}`
		assertWebsocket(t, c, 1, want)
	})

	t.Run("unchecked does not convert new recipes", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("convert=off"))

		assertStatus(t, rr.Code, http.StatusNoContent)
	})

	t.Run("checked converts new recipes", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("convert=on"))

		assertStatus(t, rr.Code, http.StatusNoContent)
	})
}

func TestHandlers_Settings_Config(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	uri := ts.URL + "/settings/config"
	fields := "email.from=test%40gmail.com&email.apikey=GHJ&integrations.ocr.key=JGKL&integrations.ocr.url=https%3A%2F%2Fwww.google.com&server.autologin=on&server.noSignups=on&server.production=on"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPut, uri)
	})

	t.Run("only admin may update config", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInOther(srv, http.MethodPut, uri, formHeader, strings.NewReader(fields))

		assertStatus(t, rr.Code, http.StatusForbidden)
		assertStringsInHTML(t, getBodyHTML(rr), []string{"Access denied: You are not an admin."})
	})

	t.Run("demo cannot update config", func(t *testing.T) {
		app.Config.Server.IsDemo = true
		defer func() {
			app.Config.Server.IsDemo = false
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, uri, formHeader, strings.NewReader(fields))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to update configuration.","title":"Database Error"}}`)
	})

	t.Run("update all fields of the config", func(t *testing.T) {
		_ = os.Setenv("RECIPYA_IS_TEST", "true")
		original := app.Config
		app.Config.Server.URL = "http://localhost"
		app.Config.Server.Port = 8078
		defer func() {
			os.Clearenv()
			app.Config = original
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, uri, formHeader, strings.NewReader(fields))

		want := app.ConfigFile{
			Email: app.ConfigEmail{
				From:           "test@gmail.com",
				SendGridAPIKey: "GHJ",
			},
			Integrations: app.ConfigIntegrations{
				AzureDI: app.AzureDI{
					Key:      "JGKL",
					Endpoint: "https://www.google.com",
				},
			},
			Server: app.ConfigServer{
				IsAutologin:  true,
				IsDemo:       false,
				IsNoSignups:  true,
				IsProduction: true,
				Port:         8078,
				URL:          "http://localhost",
			},
		}
		assertStatus(t, rr.Code, http.StatusNoContent)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-info","message":"Configuration updated.","title":"Operation Successful"}}`)
		if !cmp.Equal(app.Config, want) {
			t.Log(cmp.Diff(app.Config, want))
			t.Fail()
		}
	})
}

func TestHandlers_Settings_ConvertAutomatically(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	srv.Repository = &mockRepository{
		UserSettingsRegistered: map[int64]*models.UserSettings{
			1: {
				ConvertAutomatically: false,
				MeasurementSystem:    units.ImperialSystem,
			},
			2: {
				ConvertAutomatically: false,
				MeasurementSystem:    units.ImperialSystem,
			},
		},
	}

	uri := ts.URL + "/settings/convert-automatically"
	originalRepo := srv.Repository

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("error updating the setting", func(t *testing.T) {
		srv.Repository = &mockRepository{
			UpdateConvertMeasurementSystemFunc: func(_ int64, _ bool) error {
				return errors.New("muga ftw")
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("convert=off"))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to set setting.","title":"Database Error"}}`)
	})

	t.Run("unchecked does not convert new recipes", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("convert=off"))

		assertStatus(t, rr.Code, http.StatusNoContent)
	})

	t.Run("checked converts new recipes", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("convert=on"))

		assertStatus(t, rr.Code, http.StatusNoContent)
	})
}

func TestHandlers_Settings_Recipes_ExportSchema(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	f := &mockFiles{}
	srv.Files = f
	originalRepo := srv.Repository

	uri := ts.URL + "/settings/export/recipes"
	validExportTypes := []string{"json", "pdf"}

	t.Run("must be logged in", func(t *testing.T) {
		for _, q := range validExportTypes {
			assertMustBeLoggedIn(t, srv, http.MethodGet, uri+"?type="+q)
		}
	})

	t.Run("lost socket connection", func(t *testing.T) {
		brokers := srv.Brokers.Clone()
		srv.Brokers = nil
		defer func() {
			srv.Brokers = brokers
		}()

		rr := sendRequestAsLoggedInNoBody(srv, http.MethodGet, "/settings/export/recipes")

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-warning\",\"message\":\"Connection lost. Please reload page.\",\"title\":\"Websocket\"}"}`)
	})

	t.Run("invalid file type", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri)

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Invalid export file format.","title":"Files Error"}}`)
	})

	t.Run("no export if no recipes", func(t *testing.T) {
		for _, q := range validExportTypes {
			t.Run(q, func(t *testing.T) {
				originalHitCount := f.exportHitCount

				rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?type="+q)

				assertStatus(t, rr.Code, http.StatusAccepted)
				want := `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-warning","message":"","title":"No recipes in database."}}`
				assertWebsocket(t, c, 3, want)
				if originalHitCount != f.exportHitCount {
					t.Fatalf("expected the export function not to be called")
				}
			})
		}
	})

	t.Run("have recipes", func(t *testing.T) {
		repo := &mockRepository{
			RecipesRegistered: map[int64]models.Recipes{
				1: {{ID: 1, Name: "Chicken"}, {ID: 3, Name: "Jersey"}},
				2: {{ID: 2, Name: "BBQ"}},
			},
			UsersRegistered: srv.Repository.Users(),
		}
		srv.Repository = repo
		defer func() {
			srv.Repository = originalRepo
		}()

		for _, q := range validExportTypes {
			t.Run(q, func(t *testing.T) {
				originalHitCount := f.exportHitCount

				rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?type="+q)

				assertStatus(t, rr.Code, http.StatusAccepted)
				want := `{"type":"file","fileName":"recipes_` + q + `.zip","data":"Q2hpY2tlbi1KZXJzZXkt","toast":{"action":"","background":"","message":"","title":""}}`
				assertWebsocket(t, c, 3, want)
				if f.exportHitCount != originalHitCount+1 {
					t.Fatalf("expected the export function to be called")
				}
			})
		}
	})
}

func TestHandlers_Settings_MeasurementSystems(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	uri := ts.URL + "/settings/measurement-system"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("selected system is already user's selected", func(t *testing.T) {
		srv.Repository = &mockRepository{
			MeasurementSystemsFunc: func(userID int64) ([]units.System, models.UserSettings, error) {
				return []units.System{units.ImperialSystem}, models.UserSettings{
					ConvertAutomatically: false,
					MeasurementSystem:    units.ImperialSystem,
				}, nil
			},
		}
		defer func() {
			srv.Repository = &mockRepository{}
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("system=imperial"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-warning","message":"","title":"System already set to imperial."}}`)
	})

	t.Run("system does not exist", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("system=peanuts"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Measurement system does not exist.","title":"Form Error"}}`)
	})

	t.Run("failed to switch system", func(t *testing.T) {
		srv.Repository = &mockRepository{
			SwitchMeasurementSystemFunc: func(system units.System, userID int64) error {
				return errors.New("error switching")
			},
		}
		defer func() {
			srv.Repository = &mockRepository{}
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("system=imperial"))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Error switching units system.","title":"Database Error"}}`)
	})

	testcases := []struct {
		name    string
		recipes models.Recipes
		system  units.System
		want    models.Recipes
	}{
		{
			name: "successful request imperial to metric",
			recipes: models.Recipes{
				{
					Description: "Preheat the oven to 351 °F (351 °F). " +
						"Stir in flour, chocolate chips, and walnuts. " +
						"Drop spoonfuls of dough 1.18 inches apart onto ungreased baking sheets. " +
						"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
					Ingredients: []string{
						"1 cup butter, softened",
						"2 eggs",
						"2 teaspoons vanilla extract",
						"1 teaspoon baking soda",
						"3 cups all-purpose flour",
						"2 cups semisweet chocolate chips",
					},
					Instructions: []string{
						"Preheat the oven to 350 degrees F (175 degrees C).",
						"Stir in flour, chocolate chips, and walnuts.",
						"Drop spoonfuls of dough 2 inches apart onto ungreased baking sheets.",
						"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
					},
				},
			},
			system: units.MetricSystem,
			want: models.Recipes{
				{
					Description: "Preheat the oven to 177 °C (177 °C). " +
						"Stir in flour, chocolate chips, and walnuts. " +
						"Drop spoonfuls of dough 3 cm apart onto ungreased baking sheets. " +
						"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
					Ingredients: []string{
						"2.37 dl butter, softened",
						"2 eggs",
						"10 ml vanilla extract",
						"5 ml baking soda",
						"7.1 dl all-purpose flour",
						"4.73 dl semisweet chocolate chips",
					},
					Instructions: []string{
						"Preheat the oven to 350 degrees F (175 degrees C).",
						"Stir in flour, chocolate chips, and walnuts.",
						"Drop spoonfuls of dough 5.08 cm apart onto ungreased baking sheets.",
						"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
					},
				},
			},
		},
		{
			name: "successful request metric to imperial",
			recipes: models.Recipes{
				{
					Description: "Preheat the oven to 177 °C (351 °F). " +
						"Stir in flour, chocolate chips, and walnuts. " +
						"Drop spoonfuls of dough 29.97 mm apart onto ungreased baking sheets. " +
						"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
					Ingredients: []string{
						"2.37 dl butter, softened",
						"2 eggs",
						"10 ml vanilla extract",
						"5 ml baking soda",
						"7.1 dl all-purpose flour",
						"4.73 dl semisweet chocolate chips",
					},
					Instructions: []string{
						"Preheat the oven to 177 °C (175 degrees C).",
						"Stir in flour, chocolate chips, and walnuts.",
						"Drop spoonfuls of dough 50.8 mm apart onto ungreased baking sheets.",
						"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
					},
				},
			},
			system: units.ImperialSystem,
			want: models.Recipes{
				{
					Description: "Preheat the oven to 177 °C (351 °F). Stir in flour, chocolate chips, and walnuts. Drop spoonfuls of dough 1.18 inches apart onto ungreased baking sheets. Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
					Ingredients: []string{
						"1 cup butter, softened",
						"2 eggs",
						"2 tsp vanilla extract",
						"1 tsp baking soda",
						"1.5 pints all-purpose flour",
						"2 cups semisweet chocolate chips",
					},
					Instructions: []string{
						"Preheat the oven to 351 °F (347 °F).",
						"Stir in flour, chocolate chips, and walnuts.",
						"Drop spoonfuls of dough 2 inches apart onto ungreased baking sheets.",
						"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
					},
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockRepository{
				MeasurementSystemsFunc: func(userID int64) ([]units.System, models.UserSettings, error) {
					system := units.MetricSystem
					if tc.system == units.MetricSystem {
						system = units.ImperialSystem
					}
					return []units.System{units.ImperialSystem, units.MetricSystem}, models.UserSettings{
						ConvertAutomatically: false,
						MeasurementSystem:    system,
					}, nil
				},
				RecipesRegistered: map[int64]models.Recipes{1: tc.recipes},
			}
			srv.Repository = repo

			rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("system="+tc.system.String()))

			assertStatus(t, rr.Code, http.StatusNoContent)
			if !slices.EqualFunc(repo.RecipesRegistered[1], tc.want, func(r1 models.Recipe, r2 models.Recipe) bool {
				isDescriptionEqual := r1.Description == r2.Description
				if !isDescriptionEqual {
					t.Logf("got description:\n%s\nbut want:\n%s", r1.Description, r2.Description)
				}

				isIngredientsEqual := slices.Equal(r1.Ingredients, r2.Ingredients)
				if !isIngredientsEqual {
					t.Logf("got ingredients:\n%v\nbut want:\n%v", r1.Ingredients, r2.Ingredients)
				}

				isInstructionsEqual := slices.Equal(r1.Instructions, r2.Instructions)
				if !isInstructionsEqual {
					t.Logf("got instructions:\n%v\nbut want:\n%v", r1.Instructions, r2.Instructions)
				}

				return isDescriptionEqual && isIngredientsEqual && isInstructionsEqual

			}) {
				t.Fail()
			}
		})
	}
}
