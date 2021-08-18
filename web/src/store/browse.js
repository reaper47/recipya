import Recipe from "@/models/recipe";
import { showSnackbar, SNACKBAR_TYPE } from "@/eventbus/action";

export default {
  namespaced: true,
  state: () => ({
    categories: [],
    isLoading: false,
    pagination: {
      page: 1,
      itemsPerPage: 12,
      lengths: {},
    },
    recipes: [],
    selectedNode: "all",
  }),
  actions: {
    addRecipe({ commit }, recipe) {
      commit("ADD_RECIPE", recipe);
    },
    async getCategories({ commit, rootGetters }) {
      commit("IS_LOADING", true, { root: true });

      const res = await fetch(rootGetters.apiUrl("categories"));
      const data = await res.json();
      if (res.ok) {
        commit("SET_CATEGORIES", data["categories"]);
      }

      commit("IS_LOADING", false, { root: true });
    },
    async getPaginationLengths({ commit, rootGetters }) {
      const res = await fetch(rootGetters.apiUrl("recipes/info"));
      const data = await res.json();
      if (res.ok) {
        commit("SET_PAGINATION_LENGTHS", data["info"]);
      }
    },
    async getRecipes({ commit, getters, rootGetters }, { category }) {
      try {
        commit("IS_LOADING", true);

        let data = null;
        const params = `page=${getters.page}&limit=${getters.itemsPerPage}`;
        if (!category) {
          const res = await fetch(rootGetters.apiUrl(`recipes?${params}`));
          data = await res.json();
        } else {
          const res = await fetch(
            rootGetters.apiUrl(`recipes?c=${category}&${params}`)
          );
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
    async setSelectedNode({ commit }, node) {
      commit("SET_SELECTED_NODE", node);
    },
    setPage({ commit }, page) {
      commit("SET_PAGE", page);
    },
  },
  mutations: {
    ADD_RECIPE: (state, recipe) => state.recipes.push(new Recipe(recipe)),
    SET_CATEGORIES: (state, categories) => (state.categories = categories),
    IS_LOADING: (state, value) => (state.isLoading = value),
    SET_PAGE: (state, page) => (state.pagination.page = page),
    SET_PAGINATION_LENGTHS: (state, info) => {
      const lengths = state.pagination.lengths;
      const itemsPerPage = state.pagination.itemsPerPage;

      lengths.all = Math.ceil(info.total / itemsPerPage);
      for (const [category, total] of Object.entries(info.totalPerCategory)) {
        lengths[category] = Math.ceil(total / itemsPerPage);
      }
    },
    SET_RECIPES: (state, recipes) => {
      state.recipes.splice(
        0,
        state.recipes.length,
        ...recipes.map((item) => new Recipe(item))
      );
    },
    SET_SELECTED_NODE: (state, node) => (state.selectedNode = node),
  },
  getters: {
    isLoading: (state) => state.isLoading,
    categories: (state) => state.categories,
    itemsPerPage: (state) => state.pagination.itemsPerPage,
    page: (state) => state.pagination.page,
    pages: (state) => state.pagination.lengths,
    recipe: (state) => (id) => state.recipes.find((recipe) => recipe.id === id),
    recipes: (state) => state.recipes,
    selectedNode: (state) => state.selectedNode,
  },
};
