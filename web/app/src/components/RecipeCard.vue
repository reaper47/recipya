<template>
  <v-card class="ma-3" elevation="2" outlined tile flat max-width="22rem">
    <v-card-title class="capitalize">
      {{ recipe.name }}
      <v-spacer></v-spacer>

      <v-tooltip bottom :disabled="!isBestMatch">
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
      this.$router.push({
        name: "Search Result Recipe Page",
        params: { id: this.recipe.id },
      });
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
</style>
