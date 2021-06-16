<template>
  <v-app>
    <app-bar :actions="actions"></app-bar>
    <v-main>
      <router-view />
    </v-main>
    <nice-footer></nice-footer>
  </v-app>
</template>

<script>
import AppBar from "@/components/basic/AppBar.vue";
import NiceFooter from "@/components/basic/NiceFooter.vue";

export default {
  name: "App",
  components: {
    AppBar,
    NiceFooter,
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
