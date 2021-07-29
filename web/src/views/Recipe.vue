<template>
  <v-container v-if="!isXs" class="mb-8" style="width: 90vw">
    <v-row style="border: 1px solid #212121">
      <r-title :title="recipe.name"></r-title>
    </v-row>
    <v-row style="border: 1px solid black; border-top: none">
      <r-category
        :category="recipe.category"
        :align="isXs ? 'center' : 'start'"
      ></r-category>
      <r-yield
        v-if="isXs"
        :isMobile="true"
        :recipeYield="recipe.yield"
        style="border-left: 1px solid black"
      ></r-yield>
      <r-source :url="recipe.url" :is-mobile="isXs"></r-source>
    </v-row>
    <v-row style="border: 1px solid black; border-top: none">
      <r-description
        :text="recipe.description"
        :fontSize="isXs ? 'text-body-2' : 'text-body-1'"
      ></r-description>
    </v-row>
    <v-row v-if="!isXs" style="border: 1px solid black; border-top: none">
      <r-nutrition :nutrition="recipe.nutrition"></r-nutrition>
      <r-yield :recipeYield="recipe.yield"></r-yield>
      <r-times :prep="recipe.prepTime" :cook="recipe.cookTime"></r-times>
    </v-row>
    <v-row
      v-if="isXs && recipe.nutrition.length > 0"
      style="border: 1px solid black; border-top: none; border-right: none"
    >
      <r-nutrition :nutrition="recipe.nutrition"></r-nutrition>
    </v-row>
    <v-row v-if="isXs" style="border: 1px solid black; border-top: none">
      <r-times :prep="recipe.prepTime" :cook="recipe.cookTime"></r-times>
    </v-row>
    <v-row style="border: 1px solid black; border-top: none">
      <r-ingredients :items="recipe.ingredients"></r-ingredients>
      <v-divider vertical dark class="mx-1"></v-divider>
      <r-instructions :items="recipe.instructions"></r-instructions>
    </v-row>
  </v-container>
  <recipe-mobile v-else :recipe="recipe"></recipe-mobile>
</template>

<script>
import Components from "@/components/recipe-page";
import RecipeMobile from "./mobile/RecipeMobile.vue";

export default {
  name: "RecipePage",
  props: {
    id: {
      type: Number,
      required: true,
    },
  },
  components: { ...Components, RecipeMobile },
  data: () => ({
    recipe: null,
  }),
  created() {
    this.recipe = this.$store.getters.recipe(this.id);
  },
  computed: {
    isXs() {
      return this.$vuetify.breakpoint.name === "xs";
    },
  },
};
</script>
