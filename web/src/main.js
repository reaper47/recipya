import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import vuetify from "./plugins/vuetify";

import "@/assets/global.css";

Vue.config.productionTip = false;

const StringsPlugin = {
  install(Vue) {
    Vue.prototype.$toTitleCase = (str) =>
      str
        .toLowerCase()
        .replace(/\.\s*([a-z])|^[a-z]/gm, (s) => s.toUpperCase());
  },
};

Vue.use(StringsPlugin);

new Vue({
  router,
  store,
  vuetify,

  render: function (h) {
    return h(App);
  },
}).$mount("#app");
