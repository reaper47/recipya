import Vue from "vue";
import Vuex from "vuex";
import createPersistedState from "vuex-persistedstate";

import Recipe from "@/models/recipe";

Vue.use(Vuex);

export default new Vuex.Store({
  plugins: [createPersistedState()],
  state: {
    isLoading: false,
    recipes: [],
  },
  getters: {
    isLoading: (state) => state.isLoading,
    recipes: (state) => state.recipes,
    recipe: (state) => (id) => state.recipes.find((recipe) => recipe.id === id),
  },
  mutations: {
    IS_LOADING(state, value) {
      state.isLoading = value;
    },
    SET_RECIPES(state, recipes) {
      state.recipes.splice(
        0,
        state.recipes.length,
        ...recipes.map((item) => new Recipe(item))
      );
    },
  },
  actions: {
    async search({ commit }, { ingredientsList, limit, mode }) {
      commit("IS_LOADING", true);

      const ingredients = Array.from(new Set(ingredientsList)).join(",");
      const url = `${document.location.origin}/api/v1/search?ingredients=${ingredients}&mode=${mode}&n=${limit}`;

      // Url for developmemt with the Go server running
      //const url = `http://localhost:3001/api/v1/search?ingredients=${ingredients}&mode=${mode}&n=${limit}`;

      const res = await fetch(url);
      const data = await res.json();
      if (!res.ok) {
        commit("IS_LOADING", false);
        throw data["error"];
      }

      commit("SET_RECIPES", data["recipes"]);
      commit("IS_LOADING", false);
    },
  },
  modules: {},
});
