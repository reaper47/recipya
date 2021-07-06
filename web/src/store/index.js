import Vue from "vue";
import Vuex from "vuex";
import createPersistedState from "vuex-persistedstate";

import browse from "./browse";
import search from "./search";

Vue.use(Vuex);

export default new Vuex.Store({
  plugins: [createPersistedState()],
  state: {
    baseApiUrl: `${document.location.origin}/api/v1/`,
    //baseApiUrl: "http://localhost:3001/api/v1",
    isLoading: false,
  },
  getters: {
    apiUrl: (state) => (endpoint) => `${state.baseApiUrl}/${endpoint}`,
    isLoading: (state) => state.isLoading,
  },
  mutations: {
    IS_LOADING: (state, value) => (state.isLoading = value),
  },
  modules: {
    browse,
    search,
  },
});
