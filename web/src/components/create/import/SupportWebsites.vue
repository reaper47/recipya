<template>
  <v-card height="70vh">
    <v-card-title class="text-center"> Supported Websites </v-card-title>
    <v-divider class="mb-3"></v-divider>

    <v-card-subtitle class="pb-0">
      <v-text-field
        v-model="search"
        dense
        solo
        label="Search for a website..."
        prepend-inner-icon="mdi-magnify"
      ></v-text-field>
    </v-card-subtitle>

    <v-card-text>
      <loading-fullscreen
        v-if="$store.getters['create/isWebsitesLoading']"
      ></loading-fullscreen>
      <v-list v-else class="pt-0">
        <v-list-item
          v-for="website in filteredWebsites"
          :key="website"
          link
          :href="website"
          target="_blank"
        >
          <v-list-item-content>
            <v-list-item-title>
              {{ website }}
            </v-list-item-title>
          </v-list-item-content>

          <v-list-item-action>
            <v-icon> mdi-open-in-new </v-icon>
          </v-list-item-action>
        </v-list-item>
      </v-list>
    </v-card-text>

    <v-divider></v-divider>

    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn color="blue darken-1" text block @click="$emit('closeDialog')">
        Close
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import LoadingFullscreen from "../../../components/basic/LoadingFullscreen.vue";

export default {
  name: "SupportedWebsites",
  components: {
    LoadingFullscreen,
  },
  data: () => ({
    search: "",
  }),
  computed: {
    filteredWebsites() {
      return this.websites.filter((website) =>
        website.toLowerCase().match(this.search.toLowerCase())
      );
    },
    websites() {
      return this.$store.getters["create/websites"];
    },
  },
  created() {
    this.$store.dispatch("create/fetchWebsites");
  },
};
</script>
