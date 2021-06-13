<template>
  <div>
    <app-bar></app-bar>
    <main>
      <router-view />
    </main>
    <footer>Hello</footer>
  </div>
</template>

<script>
import AppBar from "@/components/AppBar.vue";

export default {
  name: "App",
  components: {
    AppBar,
  },
  data: () => ({
    //
  }),
};
</script>
