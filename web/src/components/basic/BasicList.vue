<template>
  <v-list flat subheader>
    <v-subheader class="text-subtitle-1 text-weight-medium">
      {{ title }}
    </v-subheader>

    <v-list-item-group v-model="checked" multiple active-class="">
      <v-list-item v-for="item in items" :key="item">
        <template v-slot:default="{ active }">
          <v-list-item-action>
            <v-checkbox :input-value="active"></v-checkbox>
          </v-list-item-action>

          <v-list-item-content>
            <v-list-item-title class="text-wrap" :class="fontSize">
              {{ toTitleCase(item) }}
            </v-list-item-title>
          </v-list-item-content>
        </template>
      </v-list-item>
    </v-list-item-group>
  </v-list>
</template>

<script>
export default {
  name: "BasicList",
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
    title: {
      type: String,
      required: true,
    },
  },
  data: () => ({
    checked: [],
  }),
  methods: {
    toTitleCase: (str) =>
      str
        .toLowerCase()
        .replace(/\.\s*([a-z])|^[a-z]/gm, (s) => s.toUpperCase()),
  },
};
</script>
