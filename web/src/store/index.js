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
    store: null,
  },
  actions: {
    setStore: ({ commit }, { store }) => commit("SET_STORE", store),
  },
  getters: {
    apiUrl: (state) => (endpoint) => `${state.baseApiUrl}/${endpoint}`,
    isLoading: (state) => state.isLoading,
    recipe: (state, getters) => (id) => getters[`${state.store}/recipe`](id),
    store: (state) => state.store,
  },
  mutations: {
    IS_LOADING: (state, value) => (state.isLoading = value),
    SET_STORE: (state, value) => (state.store = value),
  },
  modules: {
    browse,
    search,
  },
});
