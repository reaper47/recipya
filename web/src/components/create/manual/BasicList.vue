<template>
  <v-container>
    <v-subheader class="text-subtitle-1 text-weight-medium">
      {{ title }}
    </v-subheader>

    <v-row v-for="(item, i) in items" :key="key(i)" dense>
      <v-col cols="1" align-self="center">{{ i + 1 }}.</v-col>
      <v-col>
        <v-text-field
          full-width
          v-model="item.value"
          :rules="rule"
          required
          :placeholder="label(i + 1)"
          hide-details
          dense
        ></v-text-field>
      </v-col>
      <v-col cols="3" class="d-flex" align-self="center">
        <v-btn
          v-if="i === items.length - 1"
          small
          color="primary"
          @click="addRow"
          :disabled="item.value === ''"
          class="elevation-0"
        >
          +
        </v-btn>
        <v-spacer v-if="i === items.length - 1"></v-spacer>
        <v-btn
          v-if="items.length > 1"
          small
          color="error"
          @click="removeRow(i)"
          class="elevation-0"
        >
          -
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: "ManualBasicList",
  props: {
    title: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      items: [{ value: "", key: this.key(0) }],
      rule: [(v) => !!v || `${this.title} is required`],
    };
  },
  methods: {
    key(i) {
      return `${this.title}-${i}`;
    },
    label(i) {
      return `${this.title.substr(0, this.title.length - 1)} #${i}`;
    },
    addRow() {
      this.items.push({ value: "", key: this.key(this.items.length) });
    },
    removeRow(i) {
      this.items.splice(i, 1);
      this.items.forEach((item, index) => (item.key = index));
    },
    getItems() {
      return this.items.map((item) => item.value);
    },
  },
};
</script>
