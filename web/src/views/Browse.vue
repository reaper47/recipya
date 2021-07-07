<template>
  <loading-fullscreen v-if="$store.getters.isLoading"></loading-fullscreen>
  <v-container v-else fill-height fluid class="pl-0 pt-0">
    <v-treeview
      v-model="tree"
      :open="initiallyOpen"
      activatable
      @update:active="changeNode"
      :items="items"
      item-key="name"
      open-on-click
      class="capitalize"
      style="align-self: start"
    >
      <template v-slot:prepend="{ item, open }">
        <v-icon v-if="!item.icon">
          {{ open ? "mdi-folder-open" : "mdi-folder" }}
        </v-icon>
        <v-icon v-else>
          {{ icons[item.icon] }}
        </v-icon>
      </template>
    </v-treeview>
    <v-divider vertical></v-divider>
    <loading-fullscreen
      v-if="$store.getters['browse/isLoading']"
    ></loading-fullscreen>
    <v-container v-else fill-height style="width: 80%">
      <v-layout row wrap>
        <v-flex v-for="(recipe, index) in recipes" :key="recipe.name">
          <recipe-card :index="index + 1" :recipe="recipe"></recipe-card>
        </v-flex>
      </v-layout>
    </v-container>
  </v-container>
</template>
<script>
import LoadingFullscreen from "@/components/basic/LoadingFullscreen.vue";
import RecipeCard from "@/components/RecipeCard.vue";

export default {
  name: "Browse",
  components: {
    LoadingFullscreen,
    RecipeCard,
  },
  data: () => ({
    active: null,
    tree: [],
    items: [],
    initiallyOpen: ["categories"],
    icons: {
      txt: "mdi-file-document-outline",
    },
  }),
  async created() {
    if (this.categories.length === 0) {
      await this.$store.dispatch("browse/getCategories");
    }
    this.items = [
      { name: "all", icon: "txt" },
      {
        name: "categories",
        children: this.categories.map((name) => {
          return { name, icon: "txt" };
        }),
      },
    ];
  },
  computed: {
    categories() {
      return this.$store.getters["browse/categories"];
    },
    recipes() {
      return this.$store.getters["browse/recipes"];
    },
  },
  methods: {
    async changeNode(nodes) {
      const node = nodes[0];
      if (node === "all") {
        await this.$store.dispatch("browse/getRecipes", { category: null });
      } else if (this.categories.includes(node)) {
        await this.$store.dispatch("browse/getRecipes", { category: node });
      }
    },
  },
};
</script>
