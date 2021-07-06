import Recipe from "@/models/recipe";

export default {
  namespaced: true,
  state: () => ({
    isLoading: false,
    recipes: [],
  }),
  actions: {
    async getCategories({ rootGetters }) {
      const res = await fetch(rootGetters.apiUrl("categories"));
      const data = await res.json();
      if (!res.ok) {
        return [];
      }
      return data["categories"];
    },
    async getRecipes({ commit, rootGetters }, { category }) {
      let data = null;

      if (category === null) {
        const res = await fetch(rootGetters.apiUrl("recipes"));
        data = await res.json();
      } else {
        const res = await fetch(rootGetters.apiUrl(`recipes?c=${category}`));
        data = await res.json();
      }

      commit("SET_RECIPES", data["recipes"]);
    },
  },
  mutations: {
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
    recipe: (state) => (id) => state.recipes.find((recipe) => recipe.id === id),
    recipes: (state) => state.recipes,
  },
};
