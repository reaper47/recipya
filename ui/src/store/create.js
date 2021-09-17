import { showSnackbar, SNACKBAR_TYPE } from "../eventbus/action";
import router from "../router";

export default {
  namespaced: true,
  state: {
    isImporting: false,
    isPosting: false,
    isWebsitesLoading: false,
    websites: [],
  },
  actions: {
    async fetchWebsites({ commit, getters, rootGetters }) {
      const websites = getters.websites;
      if (websites.length > 0) {
        return websites;
      }

      commit("IS_WEBSITES_LOADING", true);

      try {
        const url = rootGetters.apiUrl("import/websites");
        const res = await fetch(url);
        const data = await res.json();

        commit("SET_WEBSITES", data["websites"]);
      } catch (error) {
        const title = `${error.status} (${error.code})`;
        showSnackbar(SNACKBAR_TYPE.ERROR, title, error.message);
      } finally {
        commit("IS_WEBSITES_LOADING", false);
      }
    },
    async importRecipe({ commit, dispatch, rootGetters }, url) {
      commit("IS_IMPORTING", true);

      try {
        const response = await fetch(rootGetters.apiUrl("import/url"), {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ url }),
        });

        const data = await response.json();
        if ("error" in data) {
          throw data["error"];
        }

        dispatch("browse/addRecipe", data, { root: true });
        dispatch("setStore", { store: "browse" }, { root: true });
        dispatch("browse/getCategories", {}, { root: true });
        router.push({ name: "Recipe Page", params: { id: data["id"] } });
      } catch (error) {
        const title = `${error.status} (${error.code})`;
        showSnackbar(SNACKBAR_TYPE.ERROR, title, error.message);
      } finally {
        commit("IS_IMPORTING", false);
      }
    },
    async postRecipe({ commit, dispatch, rootGetters }, recipe) {
      commit("IS_POSTING", true);

      try {
        const response = await fetch(rootGetters.apiUrl("recipes"), {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(recipe),
        });

        recipe = await response.json();
        if ("error" in recipe) {
          throw recipe["error"];
        }

        dispatch("browse/addRecipe", recipe, { root: true });
        dispatch("setStore", { store: "browse" }, { root: true });
        dispatch("browse/getCategories", {}, { root: true });
        router.push({ name: "Recipe Page", params: { id: recipe.id } });
      } catch (error) {
        const title = `${error.status} (${error.code})`;
        showSnackbar(SNACKBAR_TYPE.ERROR, title, error.message);
      } finally {
        commit("IS_POSTING", false);
      }
    },
    async sendRequest(_obj, url, payload) {
      const response = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ payload }),
      });

      const data = await response.json();
      if ("error" in data) {
        throw data["error"];
      }
      return data;
    },
  },
  mutations: {
    IS_IMPORTING: (state, value) => (state.isImporting = value),
    IS_POSTING: (state, value) => (state.isPosting = value),
    IS_WEBSITES_FETCHING: (state, value) => (state.isWebsitesLoading = value),
    SET_WEBSITES: (state, websites) => (state.websites = websites),
  },
  getters: {
    isImporting: (state) => state.isImporting,
    isPosting: (state) => state.isPosting,
    isWebsitesLoading: (state) => state.isWebsitesLoading,
    websites: (state) => state.websites,
  },
};
