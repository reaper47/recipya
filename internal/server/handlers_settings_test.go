package server_test

import (
	"errors"
	"fmt"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services"
	"github.com/reaper47/recipya/internal/units"
	"maps"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestHandlers_Settings(t *testing.T) {
	srv := newServerTest()

	uri := "/settings"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("profile tab not displayed when autologin", func(t *testing.T) {
		app.Config.Server.IsAutologin = true
		defer func() {
			app.Config.Server.IsAutologin = false
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<button class="px-2 bg-gray-300 hover:bg-gray-300 dark:hover:bg-gray-800" hx-get="/settings/tabs/recipes" hx-target="#settings-tab-content" role="tab" aria-selected="false" aria-controls="tab-content" _="on click remove .bg-gray-300 .dark:bg-gray-800 from <div[role='tablist'] button/> then add .bg-gray-300 .dark:bg-gray-800"> Recipes </button>`,
			`<p class="text-end font-semibold select-none">Export data:<br><span class="font-light text-sm">Download your recipes in the selected file format.</span></p>`,
			`<form method="dialog" hx-post="/integrations/import/nextcloud" hx-swap="none" onsubmit="document.querySelector('#integrations-nextcloud-dialog').close()"><div class="block"><label for="integrations-nextcloud-dialog-url" class="font-medium">Nextcloud URL</label><input id="integrations-nextcloud-dialog-url" type="url" name="url" placeholder="https://nextcloud.mydomain.com" class="text-input"></div><div class="block mt-3"><label for="integrations-nextcloud-dialog-username" class="font-medium">Username</label><input id="integrations-nextcloud-dialog-username" type="text" name="username" placeholder="Enter your Nextcloud username" class="text-input"></div><div class="block mt-3"><label for="integrations-nextcloud-dialog-password" class="font-medium">Password</label><input id="integrations-nextcloud-dialog-password" type="password" name="password" placeholder="Enter your Nextcloud password" class="text-input"></div><button class="mt-3 w-full rounded-lg bg-indigo-600 px-4 py-2 text-lg font-semibold tracking-wide text-white hover:bg-green-600"> Import </button></form>`,
			`<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 ml-1" fill="black" viewBox="0 0 24 24" stroke="currentColor"><path d="M16 11v5H2v-5H0v5a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-5z"/><path d="m9 14 5-6h-4V0H8v8H4z"/></svg>`,
			`<label for="systems" class="text-end font-semibold">Measurement system:</label>`,
			`<select id="systems" name="system" hx-post="/settings/measurement-system" hx-swap="none" class="h-fit w-fit bg-gray-50 border border-gray-300 rounded-lg p-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-gray-200 dark:focus:ring-blue-500 dark:focus:border-blue-500"><option value="imperial" >imperial</option><option value="metric" selected>metric</option></select>`,
			`<input type="checkbox" name="convert" id="convert" class="w-fit h-fit mt-1" hx-post="/settings/convert-automatically" hx-trigger="click">`,
			`<div class="grid grid-cols-2 gap-4 mb-2"><label for="integrations" class="text-end font-semibold">Integrations:<br><span class="font-light text-sm">Import recipes from the selected solution.</span></label><div class="grid gap-1 grid-flow-col w-fit h-fit"><label><select id="integrations" name="integrations" hx-post="/settings/measurement-system" hx-swap="none" class="h-fit w-fit bg-gray-50 border border-gray-300 rounded-lg p-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-gray-200 dark:focus:ring-blue-500 dark:focus:border-blue-500"><option value="nextcloud" selected>Nextcloud</option></select></label><button class="bg-white border border-gray-300 rounded-lg py-1 px-2 hover:bg-gray-100 focus:ring-2 focus:ring-gray-200 dark:bg-gray-800 dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:border-gray-600 dark:focus:ring-gray-700" onmousedown="document.querySelector('#integrations-nextcloud-dialog').showModal()"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" class="bi bi-cloud-arrow-down" viewBox="0 0 16 16"><path fill-rule="evenodd" d="M7.646 10.854a.5.5 0 0 0 .708 0l2-2a.5.5 0 0 0-.708-.708L8.5 9.293V5.5a.5.5 0 0 0-1 0v3.793L6.354 8.146a.5.5 0 1 0-.708.708l2 2z"/><path d="M4.406 3.342A5.53 5.53 0 0 1 8 2c2.69 0 4.923 2 5.166 4.579C14.758 6.804 16 8.137 16 9.773 16 11.569 14.502 13 12.687 13H3.781C1.708 13 0 11.366 0 9.318c0-1.763 1.266-3.223 2.942-3.593.143-.863.698-1.723 1.464-2.383zm.653.757c-.757.653-1.153 1.44-1.153 2.056v.448l-.445.049C2.064 6.805 1 7.952 1 9.318 1 10.785 2.23 12 3.781 12h8.906C13.98 12 15 10.988 15 9.773c0-1.216-1.02-2.228-2.313-2.228h-.5v-.5C12.188 4.825 10.328 3 8 3a4.53 4.53 0 0 0-2.941 1.1z"/></svg></button></div></div>`,
			`<dialog id="integrations-nextcloud-dialog" class="p-4 dark:bg-gray-600 rounded-lg dark:text-gray-200"><h3 class="flex justify-between mb-4"><span class="flex font-semibold underline">Import from Nextcloud</span><button type="button" class="hover:bg-gray-800 hover:text-white px-2 border rounded" onclick="document.querySelector('#integrations-nextcloud-dialog').close()"> X </button></h3><form method="dialog" hx-post="/integrations/import/nextcloud" hx-swap="none" onsubmit="document.querySelector('#integrations-nextcloud-dialog').close()"><div class="block"><label for="integrations-nextcloud-dialog-url" class="font-medium">Nextcloud URL</label><input id="integrations-nextcloud-dialog-url" type="url" name="url" placeholder="https://nextcloud.mydomain.com" class="text-input"></div><div class="block mt-3"><label for="integrations-nextcloud-dialog-username" class="font-medium">Username</label><input id="integrations-nextcloud-dialog-username" type="text" name="username" placeholder="Enter your Nextcloud username" class="text-input"></div><div class="block mt-3"><label for="integrations-nextcloud-dialog-password" class="font-medium">Password</label><input id="integrations-nextcloud-dialog-password" type="password" name="password" placeholder="Enter your Nextcloud password" class="text-input"></div><button class="mt-3 w-full rounded-lg bg-indigo-600 px-4 py-2 text-lg font-semibold tracking-wide text-white hover:bg-green-600"> Import </button></form></dialog>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
		notWant := []string{
			`<button class="px-2 bg-gray-300 hover:bg-gray-300 dark:bg-gray-800 dark:hover:bg-gray-800" hx-get="/settings/tabs/profile" hx-target="#settings-tab-content" role="tab" aria-selected="false" aria-controls="tab-content" _="on click remove .bg-gray-300 .dark:bg-gray-800 from <div[role='tablist'] button/> then add .bg-gray-300 .dark:bg-gray-800"> Profile </button>`,
			`<p class="grid justify-end font-semibold">Change password:</p>`,
			`<form class="h-fit w-fit border p-4 rounded-lg dark:border-none dark:bg-gray-600" hx-post="/auth/change-password" hx-indicator="#fullscreen-loader" hx-swap="none">`,
			`<input class="text-input" id="password-current" name="password-current" placeholder="Enter current password" required type="password"/>`,
			`<input class="text-input" id="password-new" name="password-new" placeholder="Enter new password" required type="password"/>`,
			`<input class="text-input" id="password-confirm" name="password-confirm" placeholder="Retype new password" required type="password"/>`,
			`<button type="submit" class="w-full p-2 font-semibold text-white bg-blue-500 rounded-lg hover:bg-blue-800"> Update </button>`,
		}
		assertStringsNotInHTML(t, getBodyHTML(rr), notWant)
	})

	t.Run("display settings", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">Settings | Recipya</title>`,
			`<button class="px-2 bg-gray-300 hover:bg-gray-300 dark:bg-gray-800 dark:hover:bg-gray-800" hx-get="/settings/tabs/profile" hx-target="#settings-tab-content" role="tab" aria-selected="false" aria-controls="tab-content" _="on click remove .bg-gray-300 .dark:bg-gray-800 from <div[role='tablist'] button/> then add .bg-gray-300 .dark:bg-gray-800"> Profile </button>`,
			`<button class="px-2 hover:bg-gray-300 dark:hover:bg-gray-800" hx-get="/settings/tabs/recipes" hx-target="#settings-tab-content" role="tab" aria-selected="false" aria-controls="tab-content" _="on click remove .bg-gray-300 .dark:bg-gray-800 from <div[role='tablist'] button/> then add .bg-gray-300 .dark:bg-gray-800"> Recipes </button>`,
			`<button class="px-2 hover:bg-gray-300 dark:hover:bg-gray-800" hx-get="/settings/tabs/advanced" hx-target="#settings-tab-content" role="tab" aria-selected="false" aria-controls="tab-content" _="on click remove .bg-gray-300 .dark:bg-gray-800 from <div[role='tablist'] button/> then add .bg-gray-300 .dark:bg-gray-800"> Advanced </button>`,
			`<div id="settings-tab-content" role="tabpanel" class="text-sm md:text-base p-4 auto-rows-min">`,
			`<p class="grid justify-end font-semibold">Change password:</p>`,
			`<form class="h-fit w-fit border p-4 rounded-lg dark:border-none dark:bg-gray-600" hx-post="/auth/change-password" hx-indicator="#fullscreen-loader" hx-swap="none">`,
			`<input class="text-input" id="password-current" name="password-current" placeholder="Enter current password" required type="password"/>`,
			`<input class="text-input" id="password-new" name="password-new" placeholder="Enter new password" required type="password"/>`,
			`<input class="text-input" id="password-confirm" name="password-confirm" placeholder="Retype new password" required type="password"/>`,
			`<button type="submit" class="w-full p-2 font-semibold text-white bg-blue-500 rounded-lg hover:bg-blue-800"> Update </button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Settings_BackupsRestore(t *testing.T) {
	srv := newServerTest()
	originalFiles := srv.Files
	originalRepo := srv.Repository

	uri := "/settings/backups/restore"

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
			message := fmt.Sprintf("{\"showToast\":\"{\\\"message\\\":\\\"%s is an invalid backup.\\\",\\\"backgroundColor\\\":\\\"bg-red-500\\\"}\"}", after)
			assertHeader(t, rr, "HX-Trigger", message)
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Failed to backup current data.\",\"backgroundColor\":\"bg-red-500\"}"}`)
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Failed to extract backup.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("restore backup failed", func(t *testing.T) {
		srv.Repository = &mockRepository{
			restoreUserBackupFunc: func(_ *models.UserBackup) error {
				return errors.New("restore failed")
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("date=2006-01-02"))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Failed to restore backup.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("valid request", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("date=2006-01-02"))

		assertStatus(t, rr.Code, http.StatusOK)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Backup restored successfully.\",\"backgroundColor\":\"bg-blue-500\"}"}`)
	})
}

func TestHandlers_Settings_CalculateNutrition(t *testing.T) {
	srv := newServerTest()
	srv.Repository = &mockRepository{
		UserSettingsRegistered: map[int64]*models.UserSettings{
			1: {CalculateNutritionFact: false}},
	}

	uri := "/settings/calculate-nutrition"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("error updating the setting", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInOther(srv, http.MethodPost, uri, formHeader, strings.NewReader("calculate-nutrition=off"))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Failed to set setting.\",\"backgroundColor\":\"bg-red-500\"}"}`)
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

func TestHandlers_Settings_ConvertAutomatically(t *testing.T) {
	srv := newServerTest()
	srv.Repository = &mockRepository{
		UserSettingsRegistered: map[int64]*models.UserSettings{
			1: {
				ConvertAutomatically: false,
				MeasurementSystem:    units.ImperialSystem,
			},
		},
	}

	uri := "/settings/convert-automatically"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("error updating the setting", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInOther(srv, http.MethodPost, uri, formHeader, strings.NewReader("convert=off"))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Failed to set setting.\",\"backgroundColor\":\"bg-red-500\"}"}`)
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

func TestHandlers_Settings_MeasurementSystems(t *testing.T) {
	srv := newServerTest()

	uri := "/settings/measurement-system"

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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"System already set to imperial.\",\"backgroundColor\":\"bg-orange-500\"}"}`)
	})

	t.Run("system does not exist", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("system=peanuts"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Measurement system does not exist.\",\"backgroundColor\":\"bg-red-500\"}"}`)
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Error switching units system.\",\"backgroundColor\":\"bg-red-500\"}"}`)
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
						"Preheat the oven to 177 °C (177 °C).",
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
					Description: "Preheat the oven to 177 °C (177 °C). " +
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
					Description: "Preheat the oven to 351 °F (351 °F). " +
						"Stir in flour, chocolate chips, and walnuts. " +
						"Drop spoonfuls of dough 1.18 inches apart onto ungreased baking sheets. " +
						"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
					Ingredients: []string{
						"1 cup butter, softened",
						"2 eggs",
						"2 tsp vanilla extract",
						"1 tsp baking soda",
						"1.5 pints all-purpose flour",
						"2 cups semisweet chocolate chips",
					},
					Instructions: []string{
						"Preheat the oven to 351 °F (351 °F).",
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

func TestHandlers_Settings_Recipes_ExportSchema(t *testing.T) {
	srv, ts, c := createWSServer()
	defer func() {
		_ = c.Close()
	}()

	originalRepo := srv.Repository

	f := &mockFiles{}
	srv.Files = f

	uri := ts.URL + "/settings/export/recipes"

	validExportTypes := []string{"json", "pdf"}

	t.Run("must be logged in", func(t *testing.T) {
		for _, q := range validExportTypes {
			assertMustBeLoggedIn(t, srv, http.MethodGet, uri+"?type="+q)
		}
	})

	t.Run("lost socket connection", func(t *testing.T) {
		brokers := maps.Clone(srv.Brokers)
		srv.Brokers = nil
		defer func() {
			srv.Brokers = brokers
		}()

		rr := sendRequestAsLoggedIn(srv, http.MethodGet, "/settings/export/recipes", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Connection lost. Please reload page.\",\"backgroundColor\":\"bg-orange-500\"}"}`)
	})

	t.Run("invalid file type", func(t *testing.T) {
		_ = c.SetReadDeadline(time.Now().Add(2 * time.Second))

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Invalid export file format.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("no export if no recipes", func(t *testing.T) {
		for _, q := range validExportTypes {
			t.Run(q, func(t *testing.T) {
				_ = c.SetReadDeadline(time.Now().Add(2 * time.Second))
				originalHitCount := f.exportHitCount

				rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?type="+q, noHeader, nil)

				assertStatus(t, rr.Code, http.StatusAccepted)
				want := `{"type":"toast","fileName":"","data":"{\"message\":\"No recipes in database.\",\"background\":\"bg-orange-500\"}"}`
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
				_ = c.SetReadDeadline(time.Now().Add(2 * time.Second))
				originalHitCount := f.exportHitCount

				rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?type="+q, noHeader, nil)

				assertStatus(t, rr.Code, http.StatusAccepted)
				want := `{"type":"file","fileName":"recipes_` + q + `.zip","data":"Q2hpY2tlbi1KZXJzZXkt"}`
				assertWebsocket(t, c, 3, want)
				if f.exportHitCount != originalHitCount+1 {
					t.Fatalf("expected the export function to be called")
				}
			})
		}
	})
}

func TestHandlers_Settings_TabsAdvanced(t *testing.T) {
	srv := newServerTest()

	uri := "/settings/tabs/advanced"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("successful request", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<p class="text-end font-semibold select-none">Restore backup:</p>`,
			`<form class="grid gap-1 grid-flow-col w-fit" hx-post="/settings/backups/restore" hx-include="select[name='date']" hx-swap="none" hx-indicator="#fullscreen-loader" hx-confirm="Continue with this backup? Today's data will be backed up if not already done.">`,
			`<label><select required id="file-type" name="date" class="bg-gray-50 border border-gray-300 rounded-lg p-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-gray-200 dark:focus:ring-blue-500 dark:focus:border-blue-500"></select></label>`,
			`<button class="bg-white border border-gray-300 rounded-lg py-2 px-4 justify-start hover:bg-gray-100 h-fit focus:ring-2 focus:ring-gray-200 dark:bg-gray-800 dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:border-gray-600 dark:focus:ring-gray-700"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M15.59 14.37a6 6 0 0 1-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 0 0 6.16-12.12A14.98 14.98 0 0 0 9.631 8.41m5.96 5.96a14.926 14.926 0 0 1-5.841 2.58m-.119-8.54a6 6 0 0 0-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 0 0-2.58 5.84m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 0 1-2.448-2.448 14.9 14.9 0 0 1 .06-.312m-2.24 2.39a4.493 4.493 0 0 0-1.757 4.306 4.493 4.493 0 0 0 4.306-1.758M16.5 9a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z"/></svg></button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Settings_TabsProfile(t *testing.T) {
	srv := newServerTest()

	uri := "/settings/tabs/profile"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("successful request", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<p class="grid justify-end font-semibold">Change password:</p>`,
			`<form class="h-fit w-fit border p-4 rounded-lg dark:border-none dark:bg-gray-600" hx-post="/auth/change-password" hx-indicator="#fullscreen-loader" hx-swap="none">`,
			`<input class="text-input" id="password-current" name="password-current" placeholder="Enter current password" required type="password"/>`,
			`<input class="text-input" id="password-new" name="password-new" placeholder="Enter new password" required type="password"/>`,
			`<input class="text-input" id="password-confirm" name="password-confirm" placeholder="Retype new password" required type="password"/>`,
			`<button type="submit" class="w-full p-2 font-semibold text-white bg-blue-500 rounded-lg hover:bg-blue-800"> Update </button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Settings_TabsRecipes(t *testing.T) {
	srv := newServerTest()

	uri := "/settings/tabs/recipes"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("error on getting units systems", func(t *testing.T) {
		srv.Repository = &mockRepository{
			MeasurementSystemsFunc: func(userID int64) ([]units.System, models.UserSettings, error) {
				return nil, models.UserSettings{}, errors.New("error fetching systems")
			},
		}
		defer func() {
			srv.Repository = &mockRepository{}
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Error fetching units systems.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("successful request", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<p class="text-end font-semibold select-none">Export data:<br><span class="font-light text-sm">Download your recipes in the selected file format.</span></p>`,
			`<form method="dialog" hx-post="/integrations/import/nextcloud" hx-swap="none" onsubmit="document.querySelector('#integrations-nextcloud-dialog').close()"><div class="block"><label for="integrations-nextcloud-dialog-url" class="font-medium">Nextcloud URL</label><input id="integrations-nextcloud-dialog-url" type="url" name="url" placeholder="https://nextcloud.mydomain.com" class="text-input"></div><div class="block mt-3"><label for="integrations-nextcloud-dialog-username" class="font-medium">Username</label><input id="integrations-nextcloud-dialog-username" type="text" name="username" placeholder="Enter your Nextcloud username" class="text-input"></div><div class="block mt-3"><label for="integrations-nextcloud-dialog-password" class="font-medium">Password</label><input id="integrations-nextcloud-dialog-password" type="password" name="password" placeholder="Enter your Nextcloud password" class="text-input"></div><button class="mt-3 w-full rounded-lg bg-indigo-600 px-4 py-2 text-lg font-semibold tracking-wide text-white hover:bg-green-600"> Import </button></form>`,
			`<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 ml-1" fill="black" viewBox="0 0 24 24" stroke="currentColor"><path d="M16 11v5H2v-5H0v5a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-5z"/><path d="m9 14 5-6h-4V0H8v8H4z"/></svg>`,
			`<label for="systems" class="text-end font-semibold">Measurement system:</label>`,
			`<select id="systems" name="system" hx-post="/settings/measurement-system" hx-swap="none" class="h-fit w-fit bg-gray-50 border border-gray-300 rounded-lg p-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-gray-200 dark:focus:ring-blue-500 dark:focus:border-blue-500"><option value="imperial" >imperial</option><option value="metric" selected>metric</option></select>`,
			`<input type="checkbox" name="convert" id="convert" class="w-fit h-fit mt-1" hx-post="/settings/convert-automatically" hx-trigger="click">`,
			`<div class="grid grid-cols-2 gap-4 mb-2"><label for="integrations" class="text-end font-semibold">Integrations:<br><span class="font-light text-sm">Import recipes from the selected solution.</span></label><div class="grid gap-1 grid-flow-col w-fit h-fit"><label><select id="integrations" name="integrations" hx-post="/settings/measurement-system" hx-swap="none" class="h-fit w-fit bg-gray-50 border border-gray-300 rounded-lg p-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-gray-200 dark:focus:ring-blue-500 dark:focus:border-blue-500"><option value="nextcloud" selected>Nextcloud</option></select></label><button class="bg-white border border-gray-300 rounded-lg py-1 px-2 hover:bg-gray-100 focus:ring-2 focus:ring-gray-200 dark:bg-gray-800 dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:border-gray-600 dark:focus:ring-gray-700" onmousedown="document.querySelector('#integrations-nextcloud-dialog').showModal()"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" class="bi bi-cloud-arrow-down" viewBox="0 0 16 16"><path fill-rule="evenodd" d="M7.646 10.854a.5.5 0 0 0 .708 0l2-2a.5.5 0 0 0-.708-.708L8.5 9.293V5.5a.5.5 0 0 0-1 0v3.793L6.354 8.146a.5.5 0 1 0-.708.708l2 2z"/><path d="M4.406 3.342A5.53 5.53 0 0 1 8 2c2.69 0 4.923 2 5.166 4.579C14.758 6.804 16 8.137 16 9.773 16 11.569 14.502 13 12.687 13H3.781C1.708 13 0 11.366 0 9.318c0-1.763 1.266-3.223 2.942-3.593.143-.863.698-1.723 1.464-2.383zm.653.757c-.757.653-1.153 1.44-1.153 2.056v.448l-.445.049C2.064 6.805 1 7.952 1 9.318 1 10.785 2.23 12 3.781 12h8.906C13.98 12 15 10.988 15 9.773c0-1.216-1.02-2.228-2.313-2.228h-.5v-.5C12.188 4.825 10.328 3 8 3a4.53 4.53 0 0 0-2.941 1.1z"/></svg></button></div></div>`,
			`<dialog id="integrations-nextcloud-dialog" class="p-4 dark:bg-gray-600 rounded-lg dark:text-gray-200"><h3 class="flex justify-between mb-4"><span class="flex font-semibold underline">Import from Nextcloud</span><button type="button" class="hover:bg-gray-800 hover:text-white px-2 border rounded" onclick="document.querySelector('#integrations-nextcloud-dialog').close()"> X </button></h3><form method="dialog" hx-post="/integrations/import/nextcloud" hx-swap="none" onsubmit="document.querySelector('#integrations-nextcloud-dialog').close()"><div class="block"><label for="integrations-nextcloud-dialog-url" class="font-medium">Nextcloud URL</label><input id="integrations-nextcloud-dialog-url" type="url" name="url" placeholder="https://nextcloud.mydomain.com" class="text-input"></div><div class="block mt-3"><label for="integrations-nextcloud-dialog-username" class="font-medium">Username</label><input id="integrations-nextcloud-dialog-username" type="text" name="username" placeholder="Enter your Nextcloud username" class="text-input"></div><div class="block mt-3"><label for="integrations-nextcloud-dialog-password" class="font-medium">Password</label><input id="integrations-nextcloud-dialog-password" type="password" name="password" placeholder="Enter your Nextcloud password" class="text-input"></div><button class="mt-3 w-full rounded-lg bg-indigo-600 px-4 py-2 text-lg font-semibold tracking-wide text-white hover:bg-green-600"> Import </button></form></dialog>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}
