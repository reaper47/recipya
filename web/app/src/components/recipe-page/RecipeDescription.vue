<template>
  <v-col>
    <p class="px-6 mb-0" :class="fontSize">{{ text }}</p>
  </v-col>
</template>

<script>
export default {
  name: "RecipeDescription",
  props: {
    fontSize: {
      type: String,
      required: false,
      default: "text-body-1",
    },
    text: {
      type: String,
      required: true,
    },
  },
};
</script>
