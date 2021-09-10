<template>
  <v-app>
    <app-bar :actions="actions" @toggle-drawer="drawer = !drawer"></app-bar>

    <v-navigation-drawer v-model="drawer" absolute temporary right>
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="text-h6">Recipya</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-divider></v-divider>

      <v-list nav dense>
        <v-list-item-group v-model="page" color="primary">
          <v-list-item
            v-for="action in actions"
            :key="action.name"
            :to="action.to"
            exact
          >
            <v-list-item-icon>
              <v-icon>{{ action.icon }}</v-icon>
            </v-list-item-icon>
            <v-list-item-title>{{ action.name }}</v-list-item-title>
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-navigation-drawer>

    <snack-bar></snack-bar>
    <v-main>
      <router-view></router-view>
    </v-main>
    <nice-footer></nice-footer>
  </v-app>
</template>

<script>
import AppBar from "./components/nav/AppBar.vue";
import NiceFooter from "./components/nav/NiceFooter.vue";
import SnackBar from "./components/basic/SnackBar.vue";

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
        to: "/browse",
        name: "Browse",
        icon: "mdi-folder",
      },
      {
        to: "/search",
        name: "Search",
        icon: "mdi-briefcase-search",
      },
      {
        to: "/new",
        name: "New",
        icon: "mdi-plus-box",
      },
    ],
    drawer: false,
    page: null,
  }),
};
</script>
