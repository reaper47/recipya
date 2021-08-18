<template>
  <div class="text-center">
    <v-pagination v-model="page" :length="pages" @input="next"></v-pagination>
  </div>
</template>

<script>
export default {
  name: "BrowsePagination",
  props: {
    selectedNode: {
      type: String,
      required: true,
    },
  },
  data: () => ({
    currentPage: null,
  }),
  mounted() {
    this.currentPage = this.page;
  },
  computed: {
    page: {
      get() {
        return this.$store.getters["browse/page"];
      },
      set(value) {
        this.$store.dispatch("browse/setPage", value);
      },
    },
    pages() {
      return this.$store.getters["browse/pages"][this.selectedNode];
    },
  },
  methods: {
    next() {
      if (this.page != this.currentPage) {
        this.$store
          .dispatch("browse/getRecipes", this.selectedNode)
          .then(() => (this.currentPage = this.page));
      }
    },
  },
};
</script>
