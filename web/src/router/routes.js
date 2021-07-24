import Home from "../views/Home.vue";

export default [
  {
    path: "/",
    name: "Home",
    component: Home,
    meta: {
      title: "Home | Recipya",
      metaTags: [
        {
          name: "description",
          content: "The home page of Recipya.",
        },
        {
          name: "og:description",
          content: "The home page of Recipya.",
        },
      ],
    },
  },
  {
    path: "/browse",
    name: "Browse",
    component: () =>
      import(/* webpackChunkName: "browse" */ "../views/Browse.vue"),
    meta: {
      title: "Browse | Recipya",
      metaTags: [
        {
          name: "description",
          content: "Browse and filter your recipes.",
        },
        {
          name: "og:description",
          content: "Browse and filter your recipes.",
        },
      ],
    },
  },
  {
    path: "/browse/:id",
    name: "Recipe Page",
    component: () =>
      import(/* webpackChunkName: "recipe" */ "../views/Recipe.vue"),
    props: (route) => {
      const props = { ...route.params };
      props.id = +props.id;
      return props;
    },
  },
  {
    path: "/search",
    name: "Search",
    component: () =>
      import(/* webpackChunkName: "search" */ "../views/Search.vue"),
    meta: {
      title: "Search | Recipya",
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
    path: "/search/results",
    name: "Search Results",
    component: () =>
      import(/* webpackChunkName: "results" */ "../views/Results.vue"),
  },
  {
    path: "/search/results/:id",
    name: "Search Result Recipe Page",
    component: () =>
      import(/* webpackChunkName: "recipe" */ "../views/Recipe.vue"),
    props: (route) => {
      const props = { ...route.params };
      props.id = +props.id;
      return props;
    },
  },
  {
    path: "/:pathMatch(.*)*",
    component: () => import(/* webpackChunkName: "404" */ "../views/404.vue"),
    name: "NotFound",
    meta: {
      title: "Page Not Found | Recipya",
    },
  },
];
