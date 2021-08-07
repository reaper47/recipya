<template>
  <v-card elevation="2" class="mx-auto" max-width="50ch" min-width="20ch">
    <v-card-title>Manual</v-card-title>
    <v-card-text class="text--primary">
      <div>Add a new recipe by filling out its content manually.</div>
    </v-card-text>
    <v-divider> </v-divider>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn text :to="{ name: 'CreateManual' }">Go</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
export default {
  name: "CreateManual",
};
</script>
