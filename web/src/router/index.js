import Vue from "vue";
import VueRouter from "vue-router";

import routes from "./routes";

Vue.use(VueRouter);

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

router.beforeEach((to, from, next) => {
  const nearestTitle = to.matched
    .slice()
    .reverse()
    .find((r) => r.meta && r.meta.title);

  const nearestMeta = to.matched
    .slice()
    .reverse()
    .find((r) => r.meta && r.meta.metaTags);

  const previousMeta = from.matched
    .slice()
    .reverse()
    .find((r) => r.meta && r.meta.metaTags);

  if (nearestTitle) {
    document.title = nearestTitle.meta.title;
  } else if (previousMeta) {
    document.title = previousMeta.meta.title;
  }

  const staleMeta = "data-vue-router-controlled";
  Array.from(document.querySelectorAll(`[${staleMeta}]`)).map((el) =>
    el.parentNode.removeChild(el)
  );

  if (!nearestMeta) {
    return next();
  }

  nearestMeta.meta.metaTags
    .map((def) => {
      const tag = document.createElement("meta");
      Object.keys(def).forEach((key) => tag.setAttribute(key, def[key]));
      tag.setAttribute(staleMeta, "");
      return tag;
    })
    .forEach((tag) => document.head.appendChild(tag));

  next();
});

export default router;
