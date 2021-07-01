<template>
  <v-col v-if="!isMobile">
    <source-button :url="url"></source-button>
  </v-col>
  <v-col v-else class="pa-0">
    <v-container class="pa-0" fill-height>
      <v-row justify="center" align="center">
        <source-button :url="url" :isSmall="true"></source-button>
      </v-row>
    </v-container>
  </v-col>
</template>

<script>
import SourceButton from "./basic/SourceButton.vue";

export default {
  name: "RecipeSource",
  components: {
    SourceButton,
  },
  props: {
    isMobile: {
      type: Boolean,
      required: false,
      default: false,
    },
    url: {
      type: String,
      required: true,
    },
  },
};
</script>
