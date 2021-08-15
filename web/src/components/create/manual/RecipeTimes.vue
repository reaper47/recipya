<template>
  <v-data-table
    :headers="headers"
    :items="items"
    hide-default-footer
    class="elevation-1"
    dense
  >
    <template #item.hours="{ item }">
      <v-text-field
        v-model.number="item.hours"
        :rules="numberRule"
        required
        placeholder="Number of hours"
        hide-details
        dense
        type="number"
        min="0"
        solo
      ></v-text-field>
    </template>

    <template #item.minutes="{ item }">
      <v-text-field
        v-model.number="item.minutes"
        :rules="numberRule"
        required
        placeholder="Number of minutes"
        hide-details
        dense
        type="number"
        min="0"
        solo
      ></v-text-field>
    </template>
  </v-data-table>
</template>

<script>
export default {
  name: "ManualRecipeTimes",
  props: {
    numberRule: {
      type: Array,
      required: true,
    },
  },
  data: () => ({
    headers: [
      {
        text: "Time",
        sortable: false,
        value: "time",
      },
      { text: "Hours", value: "hours", sortable: false },
      { text: "Minutes", value: "minutes", sortable: false },
    ],
    items: [
      { time: "Prep.", hours: 0, minutes: 10 },
      { time: "Cook", hours: 1, minutes: 30 },
    ],
  }),
  methods: {
    getTimes() {
      const prep = this.items.find((item) => item.time === "Prep.");
      const cook = this.items.find((item) => item.time === "Cook");

      const totalHours = (prep.hours | 0) + (cook.hours | 0);
      const totalMinutes = (prep.minutes | 0) + (cook.minutes | 0);

      return {
        prep: `PT${prep.hours | 0}H${prep.minutes | 0}M`,
        cook: `PT${cook.hours | 0}H${cook.minutes | 0}M`,
        total: `PT${totalHours}H${totalMinutes}M`,
      };
    },
  },
};
</script>
