<template>
  <div id="search-results">
    <goback-card
      v-if="!hasRecipes"
      title="Uh Oh..."
      text="No recipes have been found for your given query."
    ></goback-card>
    <v-container v-else style="width: 80%">
      <v-layout row wrap>
        <v-flex v-for="(recipe, index) in recipes" :key="recipe.name">
          <recipe-card
            :index="index + 1"
            :recipe="recipe"
            isResult
          ></recipe-card>
        </v-flex>
      </v-layout>
    </v-container>
  </div>
</template>

<script>
import GobackCard from "@/components/basic/GobackCard.vue";
import RecipeCard from "@/components/RecipeCard.vue";

export default {
  name: "Results",
  components: {
    GobackCard,
    RecipeCard,
  },
  computed: {
    hasRecipes() {
      return this.recipes.length > 0;
    },
    recipes() {
      return this.$store.getters["search/recipes"];
    },
  },
};
</script>
