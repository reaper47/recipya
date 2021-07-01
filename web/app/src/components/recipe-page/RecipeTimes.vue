<template>
  <v-col>
    <v-data-table
      dense
      :headers="headers"
      :items="times"
      item-key="name"
      class="elevation-1 capitalize"
      hide-default-footer
      mobile-breakpoint="0"
    >
      <template slot="body.append">
        <tr>
          <th>Total</th>
          <th>{{ total }}</th>
        </tr>
      </template>
    </v-data-table>
  </v-col>
</template>

<script>
export default {
  name: "RecipeTimes",
  props: {
    prep: {
      type: String,
      required: true,
    },
    cook: {
      type: String,
      required: true,
    },
  },
  data: () => ({
    headers: [
      { text: "Time", value: "time", sortable: false },
      { text: "Length", value: "length", sortable: false },
    ],
    times: null,
    total: null,
  }),
  created() {
    const times = this.extractTimes();
    this.times = [
      { time: "Preparation", length: this.formatTime(times.prep) },
      { time: "Cooking", length: this.formatTime(times.cook) },
    ];
    this.total = this.formatTime(times.total);
  },
  methods: {
    extractTimes() {
      const prep = { hours: 0, minutes: 0 };
      if (this.prep !== "") {
        const parts = this.prep.replace("PT", "").split("H");
        prep.hours = parseInt(parts[0]);
        prep.minutes = parseInt(parts[1].replace("M", ""));
      }

      const cook = { hours: 0, minutes: 0 };
      if (this.cook !== "") {
        const parts = this.cook.replace("PT", "").split("H");
        cook.hours = parseInt(parts[0]);
        cook.minutes = parseInt(parts[1].replace("M", ""));
      }

      const totalMinutes = prep.minutes + cook.minutes;
      const additionalHours = (totalMinutes / 60) >> 0;
      return {
        prep,
        cook,
        total: {
          hours: prep.hours + cook.hours + additionalHours,
          minutes: (totalMinutes - 60 * additionalHours) % 60,
        },
      };
    },
    formatTime: (time) => `${time.hours}h${time.minutes}m`,
  },
};
</script>
