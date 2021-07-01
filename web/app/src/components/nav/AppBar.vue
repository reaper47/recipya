<template>
  <v-app-bar
    class="mb-8"
    color="primary"
    dark
    hide-on-scroll
    dense
    max-height="48px"
  >
    <v-btn
      class="d-flex align-center"
      text
      exact
      @click.stop="$router.push({ name: 'Home' }, () => {})"
    >
      <v-img
        alt="Vuetify Logo"
        class="shrink mr-2"
        contain
        src="/icons/android-chrome-192x192.png"
        transition="scale-transition"
        width="35"
      />
      <p class="mb-0 font-weight-light">Recipe Hunter</p>
    </v-btn>

    <v-spacer></v-spacer>

    <nav class="hidden-xs-only">
      <v-btn
        v-for="action in actions"
        :key="action.name"
        text
        :to="action.to"
        exact
        active-class="active"
      >
        {{ action.name }}
      </v-btn>
    </nav>

    <span class="hidden-sm-and-up">
      <v-app-bar-nav-icon @click="$emit('toggle-drawer')"></v-app-bar-nav-icon>
    </span>
  </v-app-bar>
</template>

<script>
export default {
  name: "AppBar",
  props: {
    actions: {
      type: Array,
      required: true,
    },
  },
  data: () => ({
    drawer: false,
  }),
};
</script>

<style>
.v-btn--active {
  background-color: #e53935;
}
</style>
