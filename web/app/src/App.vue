<template>
  <v-app>
    <app-bar :actions="actions"></app-bar>
    <snack-bar></snack-bar>
    <v-main>
      <router-view />
    </v-main>
    <nice-footer></nice-footer>
  </v-app>
</template>

<script>
import AppBar from "@/components/AppBar.vue";
import NiceFooter from "@/components/NiceFooter.vue";
import SnackBar from "@/components/basic/SnackBar.vue";

export default {
  name: "App",
  components: {
    AppBar,
    NiceFooter,
    SnackBar,
  },
  data: () => ({
    actions: [
      {
        to: "/",
        name: "Home",
      },
      {
        to: "/search",
        name: "Search",
      },
    ],
  }),
};
</script>
