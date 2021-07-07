import Recipe from "@/models/recipe";

export default {
  namespaced: true,
  state: () => ({
    categories: [],
    recipes: [],
  }),
  actions: {
    async getCategories({ commit, rootGetters }) {
      commit("IS_LOADING", true, { root: true });

      const res = await fetch(rootGetters.apiUrl("categories"));
      const data = await res.json();
      if (res.ok) {
        commit("SET_CATEGORIES", data["categories"]);
      }

      commit("IS_LOADING", false, { root: false });
    },
    async getRecipes({ commit, rootGetters }, { category }) {
      commit("IS_LOADING", true, { root: true });

      let data = null;
      if (category === null) {
        const res = await fetch(rootGetters.apiUrl("recipes"));
        data = await res.json();
      } else {
        const res = await fetch(rootGetters.apiUrl(`recipes?c=${category}`));
        data = await res.json();
      }

      commit("SET_RECIPES", data["recipes"]);
      commit("IS_LOADING", false, { root: true });
    },
  },
  mutations: {
    SET_CATEGORIES: (state, categories) => (state.categories = categories),
    SET_RECIPES: (state, recipes) => {
      state.recipes.splice(
        0,
        state.recipes.length,
        ...recipes.map((item) => new Recipe(item))
      );
    },
  },
  getters: {
    categories: (state) => state.categories,
    recipe: (state) => (id) => state.recipes.find((recipe) => recipe.id === id),
    recipes: (state) => state.recipes,
  },
};
