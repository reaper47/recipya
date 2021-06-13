<template>
  <div class="search">
    <h1>This is the search page</h1>
    <search-form></search-form>
  </div>
</template>

<script>
import SearchForm from "@/components/SearchForm.vue";

export default {
  name: "Home",
  components: {
    SearchForm,
  },
};
</script>
