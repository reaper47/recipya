<template>
  <v-container>
    <v-layout row wrap justify-center align-center>
      <v-card elevation="4" max-width="15rem">
        <v-card-title>{{ title }}</v-card-title>
        <v-card-text>{{ text }}</v-card-text>
        <v-divider></v-divider>
        <v-card-actions>
          <v-btn text block @click="$router.go(-1)">
            <v-icon left dark> mdi-arrow-left </v-icon>
            Back
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-layout>
  </v-container>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      required: true,
    },
    text: {
      type: String,
      required: true,
    },
  },
};
</script>
