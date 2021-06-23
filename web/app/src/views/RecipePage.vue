<template>
  <v-container class="mb-8" style="width: 90vw">
    <v-row style="border: 1px solid #212121">
      <v-col>
        <h1 class="capitalize font-weight-light text-center">
          {{ recipe.name }}
        </h1>
      </v-col>
    </v-row>
    <v-row style="border: 1px solid black; border-top: none">
      <v-col>
        <v-chip class="capitalize">
          {{ recipe.category }}
        </v-chip>
      </v-col>
      <v-col>
        <v-btn
          v-if="recipe.url === ''"
          :href="moshPotatoesUrl"
          target="_blank"
          class="float-right"
        >
          Source
        </v-btn>
        <v-btn v-else :href="recipe.url" target="_blank" class="float-right">
          Source
        </v-btn>
      </v-col>
    </v-row>
    <v-row style="border: 1px solid black; border-top: none">
      <v-col>
        <p>{{ recipe.description }}</p>
      </v-col>
    </v-row>
    <v-row style="border: 1px solid black; border-top: none">
      <v-col style="border-right: 1px solid black">
        <v-data-table
          dense
          :headers="nutritionHeaders"
          :items="nutrition"
          item-key="name"
          class="elevation-1 capitalize"
          hide-default-footer
        ></v-data-table>
      </v-col>
      <v-col style="border-right: 1px solid black">
        <v-row>
          <p float-left class="mt-2 ml-1 text-subtitle-2">Yields:</p>
        </v-row>
        <v-row justify="center">
          <h1 class="text-h2">{{ recipe.yield }}</h1>
        </v-row>
      </v-col>
      <v-col>
        <v-data-table
          dense
          :headers="timeHeaders"
          :items="timesTable"
          item-key="name"
          class="elevation-1 capitalize"
          hide-default-footer
        >
          <template slot="body.append">
            <tr>
              <th>Total</th>
              <th>{{ formatTime(times.total) }}</th>
            </tr>
          </template>
        </v-data-table>
      </v-col>
    </v-row>
    <v-row style="border: 1px solid black; border-top: none">
      <v-col class="pa-0" style="border-right: 1px solid #212121">
        <v-list flat subheader>
          <v-subheader class="text-subtitle-1 text-weight-medium"
            >Ingredients
          </v-subheader>

          <v-list-item-group
            v-model="checkedIngredients"
            multiple
            active-class=""
          >
            <v-list-item
              v-for="ingredient in recipe.ingredients"
              :key="ingredient"
            >
              <template v-slot:default="{ active }">
                <v-list-item-action>
                  <v-checkbox :input-value="active"></v-checkbox>
                </v-list-item-action>

                <v-list-item-content>
                  <v-list-item-title class="text-wrap">
                    {{ toTitleCase(ingredient) }}
                  </v-list-item-title>
                </v-list-item-content>
              </template>
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </v-col>
      <v-divider vertical dark class="mx-1"></v-divider>
      <v-col class="pa-0">
        <v-list flat subheader>
          <v-subheader class="text-subtitle-1 text-weight-medium">
            Instructions
          </v-subheader>

          <v-list-item-group
            v-model="checkedInstructions"
            multiple
            active-class=""
          >
            <v-list-item
              v-for="instruction in recipe.instructions"
              :key="instruction"
            >
              <template v-slot:default="{ active }">
                <v-list-item-action>
                  <v-checkbox :input-value="active"></v-checkbox>
                </v-list-item-action>

                <v-list-item-content>
                  <v-list-item-title class="text-wrap">
                    {{ toTitleCase(instruction) }}
                  </v-list-item-title>
                </v-list-item-content>
              </template>
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: "Recipe",
  props: {
    id: {
      type: Number,
      required: true,
    },
  },
  data: () => ({
    recipe: null,
    moshPotatoesUrl:
      "https://www.amazon.com/Mosh-Potatoes-Recipes-Anecdotes-Heavyweights/dp/1439181322",

    nutrition: [],
    nutritionHeaders: [
      {
        text: "Nutrition",
        align: "start",
        sortable: true,
        value: "name",
      },
      { text: "Amount", value: "amount" },
    ],

    times: null,
    timesTable: [],
    timeHeaders: [
      { text: "Time", value: "time", sortable: false },
      { text: "Length", value: "length", sortable: false },
    ],

    checkedIngredients: [],
    checkedInstructions: [],
  }),
  methods: {
    extractTimes: (recipe) => {
      const prep = { hours: 0, minutes: 0 };
      if (recipe.prepTime !== "") {
        const parts = recipe.prepTime.replace("PT", "").split("H");
        prep.hours = parseInt(parts[0]);
        prep.minutes = parseInt(parts[1].replace("M", ""));
      }

      const cook = { hours: 0, minutes: 0 };
      if (recipe.cookTime !== "") {
        const parts = recipe.cookTime.replace("PT", "").split("H");
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
    toTitleCase: (str) =>
      str
        .toLowerCase()
        .replace(/\.\s*([a-z])|^[a-z]/gm, (s) => s.toUpperCase()),
  },
  created() {
    const recipe = this.$store.getters.recipe(this.id);
    if (recipe === undefined) {
      console.warn("not found");
      return;
    }

    this.nutrition = Object.entries(recipe.nutrition)
      .filter(([, amount]) => amount !== "")
      .map(([name, amount]) => {
        return {
          name: name.replace("Content", "").replace("Fat", " Fat"),
          amount,
        };
      });

    this.times = this.extractTimes(recipe);
    this.timesTable = [
      { time: "Preparation", length: this.formatTime(this.times.prep) },
      { time: "Cooking", length: this.formatTime(this.times.cook) },
    ];

    this.recipe = recipe;
  },
};
</script>
