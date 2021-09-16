<template>
  <v-form ref="form" v-model="valid" @submit.prevent="submit" lazy-validation>
    <v-row>
      <v-col cols="8">
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
      <v-col cols="4">
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
        :append-outer-icon="ingredient.icon.outer"
        :append-icon="ingredient.icon.inner"
        :autofocus="ingredient.focus"
        @click:append-outer="ingredient.action.outer($event, index)"
        @click:append="ingredient.action.inner($event, index)"
        @keydown.enter.prevent="addIngredient($event)"
      >
      </v-text-field>
    </v-col>
    <div class="text-center">
      <v-tooltip bottom>
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            v-bind="attrs"
            v-on="on"
            type="submit"
            color="primary"
            rounded
            class="pa-8 mb-8"
          >
            Search
          </v-btn>
        </template>
        <span>
          You can also submit with Ctrl + Enter
          <v-icon color="white">mdi-emoticon-cool</v-icon>
        </span>
      </v-tooltip>
    </div>
  </v-form>
</template>

<script>
import { showSnackbar, SNACKBAR_TYPE } from "../../eventbus/action";

const addIcon = "mdi-plus";
const removeIcon = "mdi-minus";

export default {
  name: "SearchForm",
  data() {
    return {
      ingredients: [
        {
          name: "",
          icon: {
            inner: null,
            outer: addIcon,
          },
          action: {
            inner: this.removeIngredient,
            outer: this.addIngredient,
          },
          focus: true,
        },
      ],

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
  computed: {
    hasEmptyIngredient() {
      return this.ingredients.filter((el) => el.name === "").length > 0;
    },
  },
  methods: {
    async submit(event) {
      if (
        (event.type === "keydown" && !event.ctrlKey) ||
        !this.$refs.form.validate()
      ) {
        return;
      }

      try {
        await this.$store.dispatch("search/search", {
          ingredientsList: this.ingredients
            .filter((item) => item.name !== "")
            .map((item) => item.name),
          limit: this.limit,
          mode: this.mode,
        });

        this.$router.push({ name: "Search Results" });
      } catch (error) {
        const title = `${error.status} (${error.code})`;
        showSnackbar(SNACKBAR_TYPE.ERROR, title, error.message);
      }
    },
    // eslint-disable-next-line no-unused-vars
    addIngredient(event, _index) {
      if (event.ctrlKey || this.hasEmptyIngredient) {
        return;
      }

      this.ingredients.forEach((el) => {
        el.icon.inner = null;
        el.icon.outer = removeIcon;
        el.action.inner = null;
        el.action.outer = this.removeIngredient;
        el.focus = false;
      });

      this.ingredients.push({
        name: "",
        icon: {
          inner: removeIcon,
          outer: addIcon,
        },
        action: {
          inner: this.removeIngredient,
          outer: this.addIngredient,
        },
        focus: true,
      });
    },
    removeIngredient(_event, index) {
      if (index === this.ingredients.length - 1) {
        const entry = this.ingredients[index - 1];
        entry.icon.inner = removeIcon;
        entry.icon.outer = addIcon;
        entry.action.inner = this.removeIngredient;
        entry.action.outer = this.addIngredient;
      }

      this.ingredients.splice(index, 1);

      if (this.ingredients.length === 1) {
        const entry = this.ingredients[0];
        entry.icon.inner = null;
        entry.icon.outer = addIcon;
        entry.action.inner = null;
        entry.action.outer = this.addIngredient;
      }
    },
  },
};
</script>
