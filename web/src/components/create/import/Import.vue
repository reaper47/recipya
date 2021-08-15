<template>
  <v-dialog v-model="dialog" scrollable max-width="60ch">
    <template v-slot:activator="{ on, attrs }">
      <v-card elevation="2" class="mx-auto" max-width="50ch" min-width="20ch">
        <v-card-title>Import</v-card-title>
        <v-card-text class="text--primary pb-0">
          <p>
            Import a recipe from one of the supported
            <span class="blue--text" v-bind="attrs" v-on="on">websites</span>.
          </p>
          <v-text-field v-model="url" label="Paste a URL"></v-text-field>
        </v-card-text>
        <v-divider> </v-divider>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            text
            @click="importRecipe"
            :disabled="!isUrlValid"
            :loading="isImporting"
          >
            Import
          </v-btn>
        </v-card-actions>
      </v-card>
    </template>

    <supported-websites @closeDialog="dialog = false"></supported-websites>
  </v-dialog>
</template>

<script>
import SupportedWebsites from "./SupportWebsites.vue";

export default {
  name: "CreateImport",
  components: {
    SupportedWebsites,
  },
  data: () => ({
    dialog: false,
    url: "",
  }),
  computed: {
    isImporting() {
      return this.$store.getters["create/isImporting"];
    },
    isUrlValid() {
      const a = document.createElement("a");
      a.href = this.url;
      return !!a.host && a.host !== window.location.host;
    },
  },
  methods: {
    importRecipe() {
      this.$store.dispatch("create/importRecipe", this.url);
    },
  },
};
</script>
