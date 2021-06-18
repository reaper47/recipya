<template>
  <div id="search-results">
    <v-container style="width: 80%">
      <v-layout row wrap>
        <v-flex v-for="(recipe, index) in recipes" :key="recipe.name">
          <recipe-card :index="index + 1" :recipe="recipe"></recipe-card>
        </v-flex>
      </v-layout>
    </v-container>
  </div>
</template>

<script>
import RecipeCard from "@/components/RecipeCard.vue";

export default {
  components: { RecipeCard },
  computed: {
    recipes() {
      return this.$store.getters.recipes;
    },
  },
};
</script>
