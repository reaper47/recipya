export default [
  {
    path: "/browse",
    name: "Browse",
    component: () =>
      import(/* webpackChunkName: "browse" */ "../../views/Browse.vue"),
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
      import(/* webpackChunkName: "recipe" */ "../../views/Recipe.vue"),
    props: (route) => {
      const props = { ...route.params };
      props.id = +props.id;
      return props;
    },
  },
];
