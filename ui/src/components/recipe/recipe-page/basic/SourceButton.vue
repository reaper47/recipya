<template>
  <v-btn
    v-if="url === ''"
    :href="moshPotatoesUrl"
    target="_blank"
    class="float-right"
    :small="isSmall"
  >
    Source
  </v-btn>
  <v-btn v-else :href="url" target="_blank" class="float-right"> Source </v-btn>
</template>

<script>
export default {
  name: "SourceButton",
  props: {
    isSmall: {
      type: Boolean,
      required: false,
      default: false,
    },
    url: {
      type: String,
      required: true,
    },
  },
  data: () => ({
    moshPotatoesUrl:
      "https://www.amazon.com/Mosh-Potatoes-Recipes-Anecdotes-Heavyweights/dp/1439181322",
  }),
};
</script>
