<template>
  <v-card elevation="2" class="mx-auto" max-width="50ch" min-width="20ch">
    <v-card-title>Scan</v-card-title>
    <v-card-text class="text--primary">
      <div>
        Upload an image of the recipe you want to add or take a picture using
        your device's camera.
      </div>
    </v-card-text>
    <v-divider> </v-divider>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn text>Camera</v-btn>
      <v-btn text>Upload</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
export default {
  name: "CreateScan",
};
</script>
