import Vue from "vue";
import Vuex from "vuex";

import Recipe from "@/models/recipe";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    recipes: [],
  },
  mutations: {
    replaceRecipes(state, recipes) {
      state.recipes.splice(
        0,
        state.recipes.length,
        ...recipes.map((item) => new Recipe(item))
      );
    },
  },
  actions: {
    search({ commit }, { ingredientsList, limit, mode }) {
      const ingredients = ingredientsList.join(",");
      //const url = `${document.location.origin}/api/v1/search?ingredients=${ingredients}&mode=${this.mode}&n=${this.limit}`;
      const url = `http://localhost:3001/api/v1/search?ingredients=${ingredients}&mode=${mode}&n=${limit}`;

      fetch(url)
        .then((response) => response.json())
        .then((data) => commit("replaceRecipes", data["recipes"]))
        .catch((error) => console.error(error));
    },
  },
  modules: {},
});
