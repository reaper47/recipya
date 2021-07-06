<template>
  <v-card class="ma-3" elevation="2" outlined tile flat max-width="22rem">
    <v-card-title class="capitalize">
      <div class="truncate" :class="{ shortWidth: isResult }">
        {{ recipe.name }}
      </div>
      <v-spacer></v-spacer>

      <v-tooltip v-if="isResult" bottom :disabled="!isBestMatch">
        <template #activator="{ on }">
          <v-chip v-on="on" color="secondary" text-color="white">
            <v-icon v-if="isBestMatch" left> mdi-trophy-outline </v-icon>
            {{ index }}
          </v-chip>
        </template>
        <span>Best Match</span>
      </v-tooltip>
    </v-card-title>

    <v-card-subtitle class="capitalize">{{ recipe.category }}</v-card-subtitle>

    <v-card-text class="text--primary">
      <p class="three-lines mt-5">{{ recipe.description }}</p>
    </v-card-text>

    <v-divider></v-divider>

    <v-card-actions>
      <v-btn text block @click="openRecipe">Open</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import Recipe from "@/models/recipe";

export default {
  name: "RecipeCard",
  props: {
    isResult: {
      type: Boolean,
      required: false,
      default: false,
    },
    index: {
      type: Number,
      required: true,
    },
    recipe: {
      type: Recipe,
      required: true,
    },
  },
  computed: {
    isBestMatch() {
      return this.index === 1;
    },
  },
  methods: {
    openRecipe() {
      const id = this.recipe.id;
      let name = "";
      let store = "";

      if (this.isResult) {
        name = "Search Result Recipe Page";
        store = "search";
      } else {
        name = "Recipe Page";
        store = "browse";
      }

      this.$router.push({ name, params: { id, store } });
    },
  },
};
</script>

<style scoped>
.three-lines {
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
  white-space: normal;
}
.shortWidth {
  width: 250px;
}
.truncate {
  white-space: nowrap;
  word-break: normal;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
