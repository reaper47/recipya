<template>
  <v-col class="pa-0" style="border-right: 1px solid #212121">
    <basic-list
      title="Ingredients"
      :items="items"
      :fontSize="fontSize"
    ></basic-list>
  </v-col>
</template>

<script>
import BasicList from "@/components/basic/BasicList.vue";

export default {
  name: "RecipeIngredients",
  props: {
    fontSize: {
      type: String,
      required: false,
      default: "text-body-1",
    },
    items: {
      type: Array,
      required: true,
    },
  },
  components: {
    BasicList,
  },
};
</script>
