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
      const prep = this.extractTimesHelper(this.prep);
      const cook = this.extractTimesHelper(this.cook);

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
    extractTimesHelper(str) {
      const time = { hours: 0, minutes: 0 };
      if (str !== "") {
        if (str.startsWith("PT")) {
          const parts = str.replace("PT", "").split("H");
          time.hours = parseInt(parts[0]);
          time.minutes = parseInt(parts[1].replace("M", ""));
        } else {
          let parts = str.replace("P", "").split("DT");
          const days = parseInt(parts[0]);
          parts = parts[1].split("H");
          time.hours = 24 * days + parseInt(parts[0]);
          time.minutes = parseInt(parts[1].replace("M", ""));
        }
      }
      return time;
    },
    formatTime: (time) => `${time.hours}h${time.minutes}m`,
  },
};
</script>
