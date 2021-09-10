<template>
  <v-container class="mb-8" style="width: 90vw">
    <v-row style="border: 1px solid #212121">
      <r-title :title="recipe.name"></r-title>
    </v-row>
    <v-row style="border: 1px solid black; border-top: none">
      <r-category :category="recipe.category" align="center"></r-category>
      <r-yield
        :isMobile="true"
        :recipeYield="recipe.yield"
        style="border-left: 1px solid black"
      ></r-yield>
      <r-source :url="recipe.url" :is-mobile="true"></r-source>
    </v-row>
    <v-row style="border: 1px solid black; border-top: none">
      <r-description
        :text="recipe.description"
        fontSize="text-body-2"
      ></r-description>
    </v-row>
    <v-row
      v-if="hasNutrition"
      style="border: 1px solid black; border-top: none; border-right: none"
    >
      <r-nutrition :nutrition="recipe.nutrition"></r-nutrition>
    </v-row>
    <v-row style="border: 1px solid black; border-top: none">
      <r-times :prep="recipe.prepTime" :cook="recipe.cookTime"></r-times>
    </v-row>
    <v-row style="border: 1px solid black; border-top: none">
      <r-ingredients
        :items="recipe.ingredients"
        fontSize="text-body-2"
      ></r-ingredients>
    </v-row>
    <v-row style="border: 1px solid black; border-top: none">
      <r-instructions
        :items="recipe.instructions"
        fontSize="text-body-2"
      ></r-instructions>
    </v-row>
  </v-container>
</template>

<script>
import Components from "../../components/recipe/recipe-page";
import Recipe from "../../models/recipe";

export default {
  name: "RecipeMobile",
  props: {
    recipe: {
      type: Recipe,
      required: true,
    },
  },
  components: Components,
  computed: {
    hasNutrition() {
      return (
        Object.entries(this.recipe.nutrition).filter(
          ([, amount]) => amount !== ""
        ).length > 0
      );
    },
  },
};
</script>
