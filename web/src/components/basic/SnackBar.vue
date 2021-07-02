<template>
  <v-snackbar
    v-model="snackbar.show"
    :color="snackbar.color"
    :multi-line="snackbar.mode === 'multi-line'"
    :timeout="snackbar.timeout"
    :top="snackbar.position === 'top'"
    width="25rem"
  >
    <v-layout align-center pr-4>
      <v-icon class="pr-3" dark large>{{ snackbar.icon }}</v-icon>
      <v-layout column>
        <div>
          <strong>{{ snackbar.title }}</strong>
        </div>
        <div>{{ snackbar.message }}</div>
      </v-layout>
    </v-layout>
    <template v-slot:action="{ attrs }">
      <v-btn
        fill-height
        color="#FFF"
        text
        v-bind="attrs"
        @click="snackbar.show = false"
      >
        Dismiss
      </v-btn>
    </template>
  </v-snackbar>
</template>

<script>
import EventBus from "@/eventbus";
import { ACTION } from "@/eventbus/action";

export default {
  data: () => ({
    snackbar: {
      show: false,
      title: "",
      message: "",
      color: "red",
      type: "multi-line",
      position: "top",
      timeout: 2500,
      icon: "mdi-heart",
    },
  }),
  mounted() {
    EventBus.$on(ACTION.SNACKBAR, (snackbarType, title, message) => {
      this.snackbar = { ...snackbarType };
      this.snackbar.title = title;
      this.snackbar.message = message;
      this.snackbar.show = true;
    });
  },
};
</script>
