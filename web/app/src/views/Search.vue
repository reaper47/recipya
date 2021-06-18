<template>
  <v-container id="search" @keydown.enter="handleKey">
    <v-layout justify-center>
      <search-form ref="form"></search-form>
    </v-layout>
  </v-container>
</template>

<script>
import SearchForm from "@/components/SearchForm.vue";

export default {
  name: "Search",
  components: {
    SearchForm,
  },
  methods: {
    handleKey(event) {
      event, this.$refs.form.submit(event);
    },
  },
};
</script>
