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
			`<title hx-swap-oob="true">Settings | Recipya</title>`,
			`<div class="grid place-content-center md:place-content-stretch md:grid-flow-col md:h-full" style="grid-template-columns: min-content">`,
			`<div class="hidden md:grid text-sm md:text-base bg-gray-200 max-w-[6rem] mt-[1px] dark:bg-gray-600 dark:border-r dark:border-r-gray-500" role="tablist">`,
			`<button class="px-2 hover:bg-gray-300 dark:hover:bg-gray-800 bg-gray-300 dark:bg-gray-800" hx-get="/settings/tabs/recipes" hx-target="#settings-tab-content" hx-trigger="mousedown" role="tab" aria-selected="false" aria-controls="tab-content" _="on mousedown remove .bg-gray-300 .dark:bg-gray-800 from <div[role='tablist'] button/> then add .bg-gray-300 .dark:bg-gray-800">Recipes</button>`,
			`<button class="px-2 hover:bg-gray-300 dark:bg-gray-600 dark:hover:bg-gray-800" hx-get="/settings/tabs/advanced" hx-target="#settings-tab-content" hx-trigger="mousedown" role="tab" aria-selected="false" aria-controls="tab-content" _="on mousedown remove .bg-gray-300 .dark:bg-gray-800 from <div[role='tablist'] button/> then add .bg-gray-300 .dark:bg-gray-800">Advanced</button>`,
			`<div id="settings_bottom_tabs" class="btm-nav btm-nav-sm z-20 md:hidden" _="on click remove .active from <button/> in settings_bottom_tabs then add .active to event.srcElement">`,
			`<button class="active" hx-get="/settings/tabs/recipes" hx-target="#settings-tab-content">Recipes</button>`,
			`<button hx-get="/settings/tabs/advanced" hx-target="#settings-tab-content">Advanced</button>`,
			`<div id="settings-tab-content" role="tabpanel" class="w-[90vw] text-sm md:text-base p-4 auto-rows-min md:w-full">`,
			`<div class="mb-4 md:mb-2 md:grid md:grid-cols-2 md:gap-4"><p class="mb-1 font-semibold md:text-end">Export data:<br><span class="font-light text-sm">Download your recipes in the selected file format.</span></p><form class="grid gap-1 grid-flow-col w-fit" hx-get="/settings/export/recipes" hx-include="select[name='type']" hx-swap="none">`,
			`<label class="form-control w-full"><div class="label"><span class="label-text font-medium">Nextcloud URL</span></div><input type="url" name="url" placeholder="https://nextcloud.mydomain.com" class="input input-bordered w-full" required></label>`,
			`<label class="form-control w-full"><div class="label"><span class="label-text font-medium">Username</span></div><input type="text" name="username" placeholder="Enter your Nextcloud username" class="input input-bordered w-full" required></label>`,
			`<label class="form-control w-full pb-2"><div class="label"><span class="label-text font-medium">Password</span></div><input type="password" name="password" placeholder="Enter your Nextcloud password" class="input input-bordered w-full" required></label>`,
			`<button class="btn btn-block btn-primary btn-sm mt-2">Import</button></form>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
		notWant := []string{
			`hx-get="/settings/tabs/profile"`,
			`Change password:`,
			`hx-post="/auth/change-password"`,
			`Enter current password`,
			`Enter new password`,
			`Retype new password`,
			`Update`,
		}
		assertStringsNotInHTML(t, getBodyHTML(rr), notWant)
	})

	t.Run("display settings", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">Settings | Recipya</title>`,
			`<div class="grid place-content-center md:place-content-stretch md:grid-flow-col md:h-full" style="grid-template-columns: min-content">`,
			`<div class="hidden md:grid text-sm md:text-base bg-gray-200 max-w-[6rem] mt-[1px] dark:bg-gray-600 dark:border-r dark:border-r-gray-500" role="tablist">`,
			`<button class="px-2 bg-gray-300 hover:bg-gray-300 dark:bg-gray-800 dark:hover:bg-gray-800" hx-get="/settings/tabs/profile" hx-target="#settings-tab-content" hx-trigger="mousedown" role="tab" aria-selected="false" aria-controls="tab-content" _="on mousedown remove .bg-gray-300 .dark:bg-gray-800 from <div[role='tablist'] button/> then add .bg-gray-300 .dark:bg-gray-800">Profile</button>`,
			`<button class="px-2 hover:bg-gray-300 dark:hover:bg-gray-800" hx-get="/settings/tabs/recipes" hx-target="#settings-tab-content" hx-trigger="mousedown" role="tab" aria-selected="false" aria-controls="tab-content" _="on mousedown remove .bg-gray-300 .dark:bg-gray-800 from <div[role='tablist'] button/> then add .bg-gray-300 .dark:bg-gray-800">Recipes</button>`,
			`<button class="px-2 hover:bg-gray-300 dark:bg-gray-600 dark:hover:bg-gray-800" hx-get="/settings/tabs/advanced" hx-target="#settings-tab-content" hx-trigger="mousedown" role="tab" aria-selected="false" aria-controls="tab-content" _="on mousedown remove .bg-gray-300 .dark:bg-gray-800 from <div[role='tablist'] button/> then add .bg-gray-300 .dark:bg-gray-800">Advanced</button>`,
			`<div id="settings_bottom_tabs" class="btm-nav btm-nav-sm z-20 md:hidden" _="on click remove .active from <button/> in settings_bottom_tabs then add .active to event.srcElement">`,
			`<button class="active" hx-get="/settings/tabs/profile" hx-target="#settings-tab-content">Profile</button>`,
			`<button class="" hx-get="/settings/tabs/recipes" hx-target="#settings-tab-content">Recipes</button>`,
			`<button hx-get="/settings/tabs/advanced" hx-target="#settings-tab-content">Advanced</button>`,
			`<div id="settings-tab-content" role="tabpanel" class="w-[90vw] text-sm md:text-base p-4 auto-rows-min md:w-full">`,
			`<div class="mb-4 md:mb-2 md:grid md:grid-cols-2 md:gap-4"><p class="mb-1 font-semibold md:text-end">Change password:</p><div class="card card-bordered card-compact w-96 bg-base-100 max-w-xs"><div class="card-body pt-2"><form hx-post="/auth/change-password" hx-indicator="#fullscreen-loader" hx-swap="none">`,
			`<label class="form-control w-full"><div class="label"><span class="label-text">Current password?</span></div><input type="password" placeholder="Enter current password" class="input input-bordered input-sm w-full max-w-xs" name="password-current" required></label>`,
			`<label class="form-control w-full"><div class="label"><span class="label-text">New password?</span></div><input type="password" placeholder="Enter new password" class="input input-bordered input-sm w-full max-w-xs" name="password-new" required></label>`,
			`<label class="form-control w-full"><div class="label"><span class="label-text">Confirm password?</span></div><input type="password" placeholder="Retype new password" class="input input-bordered input-sm w-full max-w-xs" name="password-confirm" required></label>`,
			`<div type="submit" class="card-actions justify-end mt-2"><button class="btn btn-primary btn-block btn-sm">Update</button></div></form></div></div></div>`,
			`<div class="mb-2 grid grid-cols-2 gap-4"><p class="mb-1 font-semibold md:text-end">Delete Account:<br><span class="font-light text-sm">This will delete all your data.</span></p><button type="submit" hx-delete="/auth/user" hx-confirm="Are you sure you want to delete your account? This action is irreversible." class="btn btn-error w-28">Delete</button></div>`,
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
			message := fmt.Sprintf(`{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"%s is an invalid backup.\",\"title\":\"Form Error\"}"}`, after)
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Failed to backup current data.\",\"title\":\"Files Error\"}"}`)
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Failed to extract backup.\",\"title\":\"Files Error\"}"}`)
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Failed to restore backup.\",\"title\":\"Database Error\"}"}`)
	})

	t.Run("valid request", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("date=2006-01-02"))

		assertStatus(t, rr.Code, http.StatusOK)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-info\",\"message\":\"\",\"title\":\"Backup restored successfully.\"}"}`)
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Failed to set setting.\",\"title\":\"Database Error\"}"}`)
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Failed to set setting.\",\"title\":\"Database Error\"}"}`)
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-warning\",\"message\":\"\",\"title\":\"System already set to imperial.\"}"}`)
	})

	t.Run("system does not exist", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("system=peanuts"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Measurement system does not exist.\",\"title\":\"Form Error\"}"}`)
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Error switching units system.\",\"title\":\"Database Error\"}"}`)
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
	defer c.Close()

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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-warning\",\"message\":\"Connection lost. Please reload page.\",\"title\":\"Websocket\"}"}`)
	})

	t.Run("invalid file type", func(t *testing.T) {
		_ = c.SetReadDeadline(time.Now().Add(2 * time.Second))

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Invalid export file format.\",\"title\":\"Files Error\"}"}`)
	})

	t.Run("no export if no recipes", func(t *testing.T) {
		for _, q := range validExportTypes {
			t.Run(q, func(t *testing.T) {
				_ = c.SetReadDeadline(time.Now().Add(2 * time.Second))
				originalHitCount := f.exportHitCount

				rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?type="+q, noHeader, nil)

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
				_ = c.SetReadDeadline(time.Now().Add(2 * time.Second))
				originalHitCount := f.exportHitCount

				rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?type="+q, noHeader, nil)

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
			`<p class="mb-1 font-semibold select-none md:text-end md:mb-0">Restore backup:</p>`,
			`<form class="grid gap-1 grid-flow-col w-fit" hx-post="/settings/backups/restore" hx-include="select[name='date']" hx-swap="none" hx-indicator="#fullscreen-loader" hx-confirm="Continue with this backup? Today's data will be backed up if not already done.">`,
			`<label><select required id="file-type" name="date" class="select select-bordered select-sm"></select></label>`,
			`<button class="btn btn-sm btn-outline"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M15.59 14.37a6 6 0 0 1-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 0 0 6.16-12.12A14.98 14.98 0 0 0 9.631 8.41m5.96 5.96a14.926 14.926 0 0 1-5.841 2.58m-.119-8.54a6 6 0 0 0-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 0 0-2.58 5.84m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 0 1-2.448-2.448 14.9 14.9 0 0 1 .06-.312m-2.24 2.39a4.493 4.493 0 0 0-1.757 4.306 4.493 4.493 0 0 0 4.306-1.758M16.5 9a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z"></path></svg></button>`,
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
			`<p class="mb-1 font-semibold md:text-end">Change password:</p>`,
			`<form hx-post="/auth/change-password" hx-indicator="#fullscreen-loader" hx-swap="none">`,
			`<label class="form-control w-full"><div class="label"><span class="label-text">Current password?</span></div><input type="password" placeholder="Enter current password" class="input input-bordered input-sm w-full max-w-xs" name="password-current" required></label>`,
			`<label class="form-control w-full"><div class="label"><span class="label-text">New password?</span></div><input type="password" placeholder="Enter new password" class="input input-bordered input-sm w-full max-w-xs" name="password-new" required></label>`,
			`<label class="form-control w-full"><div class="label"><span class="label-text">Confirm password?</span></div><input type="password" placeholder="Retype new password" class="input input-bordered input-sm w-full max-w-xs" name="password-confirm" required></label>`,
			`<div type="submit" class="card-actions justify-end mt-2"><button class="btn btn-primary btn-block btn-sm">Update</button>`,
			`<div class="mb-2 grid grid-cols-2 gap-4"><p class="mb-1 font-semibold md:text-end">Delete Account:<br><span class="font-light text-sm">This will delete all your data.</span></p><button type="submit" hx-delete="/auth/user" hx-confirm="Are you sure you want to delete your account? This action is irreversible." class="btn btn-error w-28">Delete</button></div>`,
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Error fetching units systems.\",\"title\":\"Database Error\"}"}`)
	})

	t.Run("successful request", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<div class="mb-4 md:mb-2 md:grid md:grid-cols-2 md:gap-4"><p class="mb-1 font-semibold md:text-end">Export data:<br><span class="font-light text-sm">Download your recipes in the selected file format.</span></p><form class="grid gap-1 grid-flow-col w-fit" hx-get="/settings/export/recipes" hx-include="select[name='type']" hx-swap="none"><label class="form-control w-full max-w-xs"><select required id="file-type" name="type" class="w-fit select select-bordered select-sm"><optgroup label="Recipes"><option value="json" selected>JSON</option> <option value="pdf">PDF</option></optgroup></select></label>`,
			`<div class="mb-4 md:mb-2 md:grid md:grid-cols-2 md:gap-4"><p class="mb-1 font-semibold md:text-end">Measurement system:</p><label class="form-control w-full max-w-xs"><select name="system" hx-post="/settings/measurement-system" hx-swap="none" class="w-fit select select-bordered select-sm"><option value="imperial">imperial</option><option value="metric" selected>metric</option></select></label></div>`,
			`<div class="flex mb-4 md:mb-2 md:grid md:grid-cols-2 md:gap-4"><label for="convert" class="mb-1 font-semibold md:text-end">Convert automatically:<br><span class="font-light text-sm">Convert new recipes to your preferred measurement system.</span></label> <input type="checkbox" name="convert" id="convert" class="checkbox" hx-post="/settings/convert-automatically" hx-trigger="click"></div>`,
			`<div class="flex mb-4 md:mb-2 md:grid md:grid-cols-2 md:gap-4"><label for="calculate-nutrition" class="mb-1 font-semibold md:text-end">Calculate nutrition facts:<br><span class="font-light text-sm md:max-w-96 md:inline-block">Calculate the nutrition facts automatically when adding a recipe. The processing will be done in the background.</span></label> <input id="calculate-nutrition" type="checkbox" name="calculate-nutrition" class="checkbox" hx-post="/settings/calculate-nutrition" hx-trigger="click"></div>`,
			`<div class="md:grid md:grid-cols-2 md:gap-4"><label class="font-semibold md:text-end">Integrations:<br><span class="font-light text-sm">Import recipes from the selected solution.</span></label><div class="grid gap-1 grid-flow-col w-fit h-fit mt-1 md:mt-0"><label class="form-control w-full max-w-xs"><select name="integrations" class="w-fit select select-bordered select-sm"><option value="nextcloud" selected>Nextcloud</option></select></label> <button class="btn btn-outline btn-sm" onmousedown="integrations_nextcloud_dialog.showModal()">`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}
