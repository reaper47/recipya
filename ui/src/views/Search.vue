<template>
  <loading-fullscreen v-if="$store.getters.isLoading"></loading-fullscreen>
  <v-container v-else id="search" @keydown.enter="handleKey" fill-height>
    <v-layout justify-center>
      <search-form ref="form"></search-form>
    </v-layout>
  </v-container>
</template>

<script>
import LoadingFullscreen from "../components/basic/LoadingFullscreen.vue";
import SearchForm from "../components/search/SearchForm.vue";

export default {
  name: "Search",
  components: {
    LoadingFullscreen,
    SearchForm,
  },
  methods: {
    handleKey(event) {
      event, this.$refs.form.submit(event);
    },
  },
};
</script>
