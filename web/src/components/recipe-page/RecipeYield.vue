<template>
  <v-col class="pb-0" style="border-right: 1px solid black">
    <v-container v-if="!isMobile">
      <v-row>
        <p float-left class="mt-2 ml-2 text-subtitle-2">Yields:</p>
      </v-row>
      <v-row class="mt-0" justify="center">
        <p class="text-h3 font-weight-light">
          {{ recipeYield }}
        </p>
      </v-row>
    </v-container>
    <v-row v-else class="mt-0" justify="center">
      <p class="text-subtitle-1">
        Yields: <b>{{ recipeYield }}</b>
      </p>
    </v-row>
  </v-col>
</template>

<script>
export default {
  name: "RecipeYield",
  props: {
    fontSize: {
      type: String,
      required: false,
      default: "text-h2",
    },
    isMobile: {
      type: Boolean,
      required: false,
      default: false,
    },
    recipeYield: {
      type: Number,
      required: true,
    },
  },
};
</script>
