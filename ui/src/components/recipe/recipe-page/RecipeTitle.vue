<template>
  <v-col>
    <h1 class="capitalize font-weight-light text-center">
      {{ title }}
    </h1>
  </v-col>
</template>

<script>
export default {
  name: "RecipeTitle",
  props: {
    title: {
      type: String,
      required: true,
    },
  },
};
</script>
