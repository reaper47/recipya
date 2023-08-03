package server_test

import (
	"errors"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/units"
	"golang.org/x/exp/slices"
	"net/http"
	"strings"
	"testing"
)

func TestHandlers_Settings(t *testing.T) {
	srv := newServerTest()

	uri := "/settings"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("display settings", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<button class="px-2 bg-gray-300 hover:bg-gray-300 dark:bg-gray-800 dark:hover:bg-gray-800" hx-get="/settings/tabs/profile" hx-target="#settings-tab-content" role="tab" aria-selected="false" aria-controls="tab-content" _="on click remove .bg-gray-300 .dark:bg-gray-800 from <div[role='tablist'] button/> then add .bg-gray-300 .dark:bg-gray-800"> Profile </button>`,
			`<button class="px-2 hover:bg-gray-300 dark:hover:bg-gray-800" hx-get="/settings/tabs/recipes" hx-target="#settings-tab-content" role="tab" aria-selected="false" aria-controls="tab-content" _="on click remove .bg-gray-300 .dark:bg-gray-800 from <div[role='tablist'] button/> then add .bg-gray-300 .dark:bg-gray-800"> Recipes </button>`,
			`<div id="settings-tab-content" role="tabpanel" class="text-sm md:text-base p-4 auto-rows-min">`,
			`<p class="grid justify-end font-semibold">Change password:</p>`,
			`<form class="h-fit w-fit border p-4 rounded-lg dark:border-none dark:bg-gray-600" hx-post="/auth/change-password" hx-indicator="#fullscreen-loader" hx-swap="none">`,
			`<input class="w-full rounded-lg bg-gray-100 px-4 py-2 dark:bg-gray-900" id="password-current" name="password-current" placeholder="Enter current password" required type="password"/>`,
			`<input class="w-full rounded-lg bg-gray-100 px-4 py-2 dark:bg-gray-900" id="password-new" name="password-new" placeholder="Enter new password" required type="password"/>`,
			`<input class="w-full rounded-lg bg-gray-100 px-4 py-2 dark:bg-gray-900" id="password-confirm" name="password-confirm" placeholder="Retype new password" required type="password"/>`,
			`<button type="submit" class="w-full p-2 font-semibold text-white bg-blue-500 rounded-lg hover:bg-blue-800"> Update </button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
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
			MeasurementSystemsFunc: func(userID int64) ([]units.System, units.System, error) {
				return []units.System{units.ImperialSystem}, units.ImperialSystem, nil
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
				MeasurementSystemsFunc: func(userID int64) ([]units.System, units.System, error) {
					system := units.MetricSystem
					if tc.system == units.MetricSystem {
						system = units.ImperialSystem
					}
					return []units.System{units.ImperialSystem, units.MetricSystem}, system, nil
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
			`<input class="w-full rounded-lg bg-gray-100 px-4 py-2 dark:bg-gray-900" id="password-current" name="password-current" placeholder="Enter current password" required type="password"/>`,
			`<input class="w-full rounded-lg bg-gray-100 px-4 py-2 dark:bg-gray-900" id="password-new" name="password-new" placeholder="Enter new password" required type="password"/>`,
			`<input class="w-full rounded-lg bg-gray-100 px-4 py-2 dark:bg-gray-900" id="password-confirm" name="password-confirm" placeholder="Retype new password" required type="password"/>`,
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
			MeasurementSystemsFunc: func(userID int64) ([]units.System, units.System, error) {
				return nil, "", errors.New("error fetching systems")
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
			`<p class="text-end font-semibold select-none">Export data:</p>`,
			`<button type="button" hx-get="/recipes/export" hx-swap="none" class="h-fit w-fit bg-white border border-gray-300 rounded-lg py-2 px-4 justify-start hover:bg-gray-100 focus:ring-4 focus:ring-gray-200 dark:bg-gray-800 dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:border-gray-600 dark:focus:ring-gray-700">`,
			`<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 ml-1" fill="black" viewBox="0 0 24 24" stroke="currentColor"><path d="M16 11v5H2v-5H0v5a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-5z"/><path d="m9 14 5-6h-4V0H8v8H4z"/></svg>`,
			`<label for="systems" class="text-end font-semibold">Measurement system:</label>`,
			`<select id="systems" name="system" hx-post="/settings/measurement-system" hx-swap="none" class="h-fit w-fit bg-gray-50 border border-gray-300 rounded-lg p-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-gray-200 dark:focus:ring-blue-500 dark:focus:border-blue-500"><option value="imperial" >imperial</option><option value="metric" selected>metric</option></select>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}
