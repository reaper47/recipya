<template>
  <v-app-bar
    class="mb-8"
    color="primary"
    dark
    hide-on-scroll
    dense
    max-height="48px"
  >
    <div class="d-flex align-center">
      <v-img
        alt="Vuetify Logo"
        class="shrink mr-2"
        contain
        src="/icons/android-chrome-192x192.png"
        transition="scale-transition"
        width="35"
      />
      <h2 class="font-weight-light">Recipe Hunter</h2>
    </div>

    <v-spacer></v-spacer>

    <nav>
      <v-btn
        v-for="action in actions"
        :key="action.name"
        text
        :to="action.to"
        active-class="active"
      >
        {{ action.name }}
      </v-btn>
    </nav>
  </v-app-bar>
</template>

<script>
export default {
  props: {
    actions: {
      type: Array,
      required: true,
    },
  },
};
</script>

<style>
.v-btn--active {
  background-color: #e53935;
}
</style>
