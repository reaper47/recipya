import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../views/Home.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
    meta: {
      title: "Home | Recipe Hunter",
      metaTags: [
        {
          name: "description",
          content: "The home page of Recipe Hunter.",
        },
        {
          name: "og:description",
          content: "The home page of Recipe Hunter.",
        },
      ],
    },
  },
  {
    path: "/search",
    name: "Search",
    component: () =>
      import(/* webpackChunkName: "search" */ "../views/Search.vue"),
    meta: {
      title: "Search | Recipe Hunter",
      metaTags: [
        {
          name: "description",
          content: "Search for recipes based on what is in your fridge.",
        },
        {
          name: "og:description",
          content: "Search for recipes based on what is in your fridge.",
        },
      ],
    },
  },
  {
    path: "/results",
    name: "Results",
    component: () =>
      import(/* webpackChunkName: "results" */ "../views/Results.vue"),
    meta: {
      title: "Results | Recipe Hunter",
      metaTags: [
        {
          name: "description",
          content: "Recipes from the search are show on this page.",
        },
        {
          name: "og:description",
          content: "Recipes from the search are show on this page.",
        },
      ],
    },
  },
  {
    path: "/:pathMatch(.*)*",
    component: import(/* webpackChunkName: "404" */ "../views/404.vue"),
    name: "NotFound",
    meta: {
      title: "Page Not Found | Recipe Hunter",
    },
  },
];

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
