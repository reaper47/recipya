export default [
  {
    path: "/new",
    name: "Create",
    component: () =>
      import(/* webpackChunkName: "create" */ "../../views/Create.vue"),
    meta: {
      title: "New | Recipya",
      metaTags: [
        {
          name: "description",
          content: "Add a recipe to your collection.",
        },
        {
          name: "og:description",
          content: "Add a recipe to your collection.",
        },
      ],
    },
  },
  {
    path: "/new/manual",
    name: "CreateManual",
    component: () =>
      import(/* webpackChunkName: "create" */ "../../views/Create.vue"),
    meta: {
      title: "New | Recipya",
      metaTags: [
        {
          name: "description",
          content: "Add a recipe to your collection.",
        },
        {
          name: "og:description",
          content: "Add a recipe to your collection.",
        },
      ],
    },
  },
];
