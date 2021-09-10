<template>
  <v-col class="pa-0">
    <basic-list
      title="Instructions"
      :items="items"
      :fontSize="fontSize"
    ></basic-list>
  </v-col>
</template>

<script>
import BasicList from "../../../components/basic/BasicList.vue";

export default {
  name: "RecipeInstructions",
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
