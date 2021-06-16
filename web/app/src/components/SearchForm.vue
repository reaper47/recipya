<template>
  <v-form
    ref="form"
    v-model="valid"
    class="mx-4"
    @submit.prevent="submit"
    lazy-validation
  >
    <v-row>
      <v-col cols="6">
        <v-select
          v-model="mode"
          :items="modes"
          item-text="label"
          item-value="mode"
          :rules="[(v) => !!v || 'Item is required']"
          label="Search Mode"
          outlined
        ></v-select>
      </v-col>
      <v-col cols="6">
        <v-select
          v-model="limit"
          :items="numRecipes"
          :rules="[(v) => !!v || 'Item is required']"
          label="Number of Recipes"
          outlined
        ></v-select>
      </v-col>
    </v-row>
    <v-col>
      <v-text-field
        v-for="(ingredient, index) in ingredients"
        :key="`${ingredient}${index}`"
        v-model.lazy="ingredient.name"
        :label="`Ingredient ${index + 1}`"
        :rules="[rules.required]"
        required
        outlined
        :append-outer-icon="ingredient.icon"
        :autofocus="ingredient.focus"
        @click:append-outer="addOrDeleteItem(index)"
        @keydown.enter.prevent="addOrDeleteItem(index)"
      ></v-text-field>
    </v-col>
    <div class="text-center">
      <v-btn type="submit" color="primary" rounded class="pa-8 mb-8">
        Search
      </v-btn>
    </div>
  </v-form>
</template>

<script>
export default {
  data() {
    return {
      ingredients: [{ name: "", icon: "mdi-plus", focus: true }],

      limit: 10,
      numRecipes: [5, 10, 15, 20, 25, 30],

      mode: 1,
      modes: [
        { label: "Maximize items taken from the fridge", mode: 1 },
        { label: "Minimize number of items to buy", mode: 2 },
      ],

      valid: true,
      rules: {
        required: (v) => !!v || "This field is required.",
      },
    };
  },
  methods: {
    submit() {
      if (!this.$refs.form.validate()) {
        return;
      }

      this.$store.dispatch("search", {
        ingredientsList: this.ingredients
          .filter((item) => item.name !== "")
          .map((item) => item.name),
        limit: this.limit,
        mode: this.mode,
      });

      this.$router.push("results");
    },
    addOrDeleteItem(index) {
      const item = this.ingredients[index];
      if (item.icon === "mdi-plus") {
        if (item.name !== "") {
          item.icon = "mdi-minus";
          item.focus = false;
          this.ingredients.push({
            name: "",
            icon: "mdi-plus",
            focus: true,
          });
        }
      } else {
        this.ingredients.splice(index, 1);
      }
    },
  },
};
</script>
