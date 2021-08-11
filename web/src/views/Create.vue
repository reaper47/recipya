<template>
  <v-container fill-height align-start>
    <v-row wrap>
      <v-col>
        <create-manual></create-manual>
      </v-col>
      <v-col>
        <create-import></create-import>
      </v-col>
      <v-col>
        <create-scan></create-scan>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import CreateManual from "@/components/create/manual.vue";
import CreateImport from "@/components/create/import/import.vue";
import CreateScan from "@/components/create/scan.vue";

export default {
  name: "Create",
  components: {
    CreateManual,
    CreateImport,
    CreateScan,
  },
};
</script>
