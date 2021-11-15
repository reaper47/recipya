import Home from "../views/Home.vue";

import browseRoutes from "./routes/browse";
import createRoutes from "./routes/create";
import searchRoutes from "./routes/search";

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
  ...browseRoutes,
  ...createRoutes,
  ...searchRoutes,
  {
    path: "/:pathMatch(.*)*",
    component: () => import(/* webpackChunkName: "404" */ "../views/404.vue"),
    name: "NotFound",
    meta: {
      title: "Page Not Found | Recipya",
    },
  },
];
