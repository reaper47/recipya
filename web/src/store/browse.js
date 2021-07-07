import Recipe from "@/models/recipe";
import { showSnackbar, SNACKBAR_TYPE } from "@/eventbus/action";

export default {
  namespaced: true,
  state: () => ({
    categories: [],
    isLoading: false,
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

      commit("IS_LOADING", false, { root: true });
    },
    async getRecipes({ commit, rootGetters }, { category }) {
      try {
        commit("IS_LOADING", true);

        let data = null;
        if (category === null) {
          const res = await fetch(rootGetters.apiUrl("recipes"));
          data = await res.json();
        } else {
          const res = await fetch(rootGetters.apiUrl(`recipes?c=${category}`));
          data = await res.json();
        }
        commit("SET_RECIPES", data["recipes"]);
      } catch (error) {
        const title = `${error.status} (${error.code})`;
        showSnackbar(SNACKBAR_TYPE.ERROR, title, error.message);
      } finally {
        commit("IS_LOADING", false);
      }
    },
  },
  mutations: {
    SET_CATEGORIES: (state, categories) => (state.categories = categories),
    IS_LOADING: (state, value) => (state.isLoading = value),
    SET_RECIPES: (state, recipes) => {
      state.recipes.splice(
        0,
        state.recipes.length,
        ...recipes.map((item) => new Recipe(item))
      );
    },
  },
  getters: {
    isLoading: (state) => state.isLoading,
    categories: (state) => state.categories,
    recipe: (state) => (id) => state.recipes.find((recipe) => recipe.id === id),
    recipes: (state) => state.recipes,
  },
};
