import Recipe from "../models/recipe";

export default {
  namespaced: true,
  state: {
    recipes: [],
  },
  actions: {
    async search({ commit, rootGetters }, { ingredientsList, limit, mode }) {
      commit("IS_LOADING", true, { root: true });

      const ingredients = Array.from(new Set(ingredientsList)).join(",");
      const url = rootGetters.apiUrl(
        `search?ingredients=${ingredients}&mode=${mode}&n=${limit}`
      );

      try {
        const res = await fetch(url);
        const data = await res.json();
        if (!res.ok) {
          throw data["error"];
        }

        commit("SET_RECIPES", data["recipes"]);
      } finally {
        commit("IS_LOADING", false, { root: true });
      }
    },
  },
  getters: {
    recipe: (state) => (id) => state.recipes.find((recipe) => recipe.id === id),
    recipes: (state) => state.recipes,
  },
  mutations: {
    SET_RECIPES(state, recipes) {
      state.recipes.splice(
        0,
        state.recipes.length,
        ...recipes.map((item) => new Recipe(item))
      );
    },
  },
};
