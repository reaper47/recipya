import { showSnackbar, SNACKBAR_TYPE } from "@/eventbus/action";

export default {
  namespaced: true,
  state: {
    isImporting: false,
    isWebsitesLoading: false,
    websites: [
      "https://claudia.abril.com.br",
      "https://acouplecooks.com/",
      "http://www.afghankitchenrecipes.com/",
      "https://1claudia.abril.com.br",
      "https://2acouplecooks.com/",
      "http://3www.afghankitchenrecipes.com/",
      "https://4claudia.abril.com.br",
      "https://5acouplecooks.com/",
      "http://6www.afghankitchenrecipes.com/",
      "https://7claudia.abril.com.br",
      "https://8acouplecooks.com/",
      "http://www.9afghankitchenrecipes.com/",
    ],
  },
  actions: {
    async fetchWebsites({ commit, getters, rootGetters }) {
      const websites = getters.websites;
      if (websites.length > 0) {
        return websites;
      }

      commit("IS_WEBSITES_LOADING", true);

      try {
        const url = rootGetters.apiUrl("new/import/websites");

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
    async importRecipe({ commit, rootGetters }, url) {
      commit("IS_IMPORTING", true);

      try {
        const apiUrl = rootGetters.apiUrl("new?type=import");
        console.warn(apiUrl, { url });
      } catch (error) {
        const title = `${error.status} (${error.code})`;
        showSnackbar(SNACKBAR_TYPE.ERROR, title, error.message);
      } finally {
        commit("IS_IMPORTING", false);
      }
    },
  },
  mutations: {
    IS_IMPORTING: (state, value) => (state.isImporting = value),
    IS_WEBSITES_FETCHING: (state, value) => (state.isWebsitesLoading = value),
    SET_WEBSITES: (state, websites) => (state.websites = websites),
  },
  getters: {
    isImporting: (state) => state.isImporting,
    isWebsitesLoading: (state) => state.isWebsitesLoading,
    websites: (state) => state.websites,
  },
};
