export default [
  {
    path: "/search",
    name: "Search",
    component: () =>
      import(/* webpackChunkName: "search" */ "../../views/Search.vue"),
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
      import(/* webpackChunkName: "results" */ "../../views/Results.vue"),
  },
  {
    path: "/search/results/:id",
    name: "Search Result Recipe Page",
    component: () =>
      import(/* webpackChunkName: "recipe" */ "../../views/Recipe.vue"),
    props: (route) => {
      const props = { ...route.params };
      props.id = +props.id;
      return props;
    },
  },
];
