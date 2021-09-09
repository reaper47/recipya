<template>
  <loading-fullscreen v-if="$store.getters.isLoading"></loading-fullscreen>
  <v-container
    v-else-if="!$vuetify.breakpoint.mdAndDown"
    fill-height
    fluid
    class="pl-0 pt-0"
  >
    <v-treeview
      v-if="!$vuetify.breakpoint.mdAndDown"
      class="capitalize"
      style="align-self: start"
      v-model="tree"
      activatable
      :open="initiallyOpen"
      :items="items"
      item-key="name"
      open-on-click
      :active="[selectedNode]"
      @update:active="changeNode"
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
    <v-container fill-height align-start style="width: 80%">
      <v-row wrap>
        <v-col
          v-for="(recipe, i) in recipes"
          :key="recipe.name"
          style="min-width: 50ch"
        >
          <loading-fullscreen
            v-if="$store.getters['browse/isLoading']"
          ></loading-fullscreen>
          <recipe-card v-else :index="i + 1" :recipe="recipe"></recipe-card>
        </v-col>
      </v-row>
      <v-container class="pa-0" style="align-self: end">
        <v-row fill-height>
          <v-col>
            <browse-pagination :selectedNode="selectedNode"></browse-pagination>
          </v-col>
        </v-row>
      </v-container>
    </v-container>
  </v-container>
  <v-container v-else fluid class="pt-0 align-start">
    <v-row>
      <v-col>
        <v-select
          v-model="nodeMobile"
          :items="itemsMobile"
          item-text="name"
          item-value="id"
          label="Filter"
          @change="changeFilter"
        ></v-select>
      </v-col>
    </v-row>
    <loading-fullscreen
      v-if="$store.getters['browse/isLoading']"
    ></loading-fullscreen>
    <v-row v-else>
      <v-col>
        <v-expansion-panels>
          <v-expansion-panel
            v-for="(recipe, i) in recipes"
            :key="`${i}-mobile`"
          >
            <v-expansion-panel-header class="capitalize">
              {{ recipe.name }}
            </v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-card flat>
                <v-card-text>
                  {{ recipe.description }}
                </v-card-text>
                <v-card-actions class="pt-0">
                  <v-spacer></v-spacer>
                  <v-btn block @click="openRecipe(recipe.id)"> Open </v-btn>
                </v-card-actions>
              </v-card>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-col>
    </v-row>
    <v-container class="pa-0" style="align-self: end">
      <v-row fill-height>
        <v-col>
          <browse-pagination :selectedNode="selectedNode"></browse-pagination>
        </v-col>
      </v-row>
    </v-container>
  </v-container>
</template>
<script>
import BrowsePagination from "@/components/browse/Pagination.vue";
import LoadingFullscreen from "@/components/basic/LoadingFullscreen.vue";
import RecipeCard from "@/components/recipe/RecipeCard.vue";

export default {
  name: "Browse",
  components: {
    BrowsePagination,
    LoadingFullscreen,
    RecipeCard,
  },
  data: () => ({
    nodeMobile: null,
    tree: [],
    items: [],
    initiallyOpen: ["categories"],
    icons: {
      txt: "mdi-file-document-outline",
    },
  }),
  beforeCreate() {
    this.$store.dispatch("browse/getPaginationLengths");
  },
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
    this.updateNodeMobile(this.selectedNode);
  },
  computed: {
    categories() {
      return this.$store.getters["browse/categories"];
    },
    itemsMobile() {
      return [
        { id: "all", name: "All" },
        ...this.categories.map((name) => ({
          id: name,
          name: `[Category] ${this.titleCase(name)}`,
        })),
      ];
    },
    recipes() {
      return this.$store.getters["browse/recipes"];
    },
    selectedNode() {
      return this.$store.getters["browse/selectedNode"];
    },
  },
  methods: {
    changeFilter(filter) {
      this.changeNode([filter.replace("[Category]", "").trim().toLowerCase()]);
    },
    async changeNode(nodes) {
      const node = nodes[0];

      let category;
      if (node === "all") {
        category = null;
      } else if (this.categories.includes(node)) {
        category = node;
      }

      this.$store.dispatch("browse/setPage", 1);
      this.$store.dispatch("browse/setSelectedNode", node);
      await this.$store.dispatch("browse/getRecipes", { category });
      this.updateNodeMobile(node);
    },
    openRecipe(id) {
      this.$store.dispatch("setStore", { store: "browse" });
      this.$router.push({
        name: "Recipe Page",
        params: { id },
      });
    },
    updateNodeMobile(node) {
      this.nodeMobile = {
        id: node,
        name: this.titleCase(node),
      };
    },
    titleCase(str) {
      return str.toLowerCase().replace(/\b(\w)/g, (s) => s.toUpperCase());
    },
  },
};
</script>

<style scoped>
::v-deep .v-treeview-node__root.v-treeview-node--active {
  pointer-events: none;
}

::v-deep .v-treeview-node__toggle {
  pointer-events: all;
}
</style>
