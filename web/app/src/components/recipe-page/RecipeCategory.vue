<template>
  <v-col :align="align" justify="center" fill-height>
    <v-chip class="capitalize">
      {{ category }}
    </v-chip>
  </v-col>
</template>

<script>
export default {
  name: "RecipeCategory",
  props: {
    align: {
      type: String,
      required: false,
      default: "start",
    },
    category: {
      type: String,
      required: true,
    },
  },
};
</script>
