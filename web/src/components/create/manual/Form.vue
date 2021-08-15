<template>
  <v-form v-model="valid" ref="form">
    <v-container v-if="!isXs" class="mb-8" style="width: 90vw">
      <v-row style="border: 1px solid #212121">
        <v-col cols="3"></v-col>
        <v-col>
          <v-text-field
            v-model="title"
            :rules="rules.title"
            label="Title"
            required
            placeholder="Grandma's Thai Curry Chicken"
            hide-details
            dense
            clearable
          ></v-text-field>
        </v-col>
        <v-col cols="3"></v-col>
      </v-row>
      <v-row style="border: 1px solid black; border-top: none">
        <v-col>
          <div style="width: 25ch">
            <v-combobox
              v-model="category"
              :rules="rules.category"
              :items="categories"
              label="Category"
              placeholder="Category"
              required
              filled
              single-line
              hide-details
              dense
              rounded
            ></v-combobox>
          </div>
        </v-col>
        <v-col>
          <div style="width: 25ch; float: right">
            <v-text-field
              v-model="source"
              :rules="rules.source"
              label="Source"
              placeholder="https://www.recipes.com/..."
              required
              filled
              hide-details
              dense
              clearable
            ></v-text-field>
          </div>
        </v-col>
      </v-row>
      <v-row style="border: 1px solid black; border-top: none">
        <v-col>
          <v-textarea
            v-model="description"
            :rules="rules.description"
            outlined
            name="description"
            label="Description..."
            placeholder="This Thai curry chicken will make you drool..."
            dense
            clearable
            rows="4"
            hide-details
          ></v-textarea>
        </v-col>
      </v-row>
      <v-row v-if="!isXs" style="border: 1px solid black; border-top: none">
        <v-col style="border-right: 1px solid black">
          <recipe-nutrition ref="nutrition"></recipe-nutrition>
        </v-col>
        <v-col style="border-right: 1px solid black" cols="3">
          <v-row>
            <v-col cols="4">
              <p float-left class="mt-2 ml-2 text-subtitle-2">Yields:</p>
            </v-col>
            <v-col>
              <v-text-field
                v-model.number="yields"
                :rules="rules.number"
                label="Number of servings"
                required
                placeholder="Number of servings"
                hide-details
                dense
                solo
                type="number"
                single-line
                min="1"
              ></v-text-field>
            </v-col>
          </v-row>
        </v-col>
        <v-col>
          <recipe-times
            ref="times"
            :numberRule="rules.numberTimes"
          ></recipe-times>
        </v-col>
      </v-row>
      <v-row style="border: 1px solid black; border-top: none">
        <v-col class="pa-0" style="border-right: 1px solid black">
          <basic-list title="Ingredients" ref="ingredients"></basic-list>
        </v-col>
        <v-col class="pa-0">
          <basic-list title="Instructions" ref="instructions"></basic-list>
        </v-col>
      </v-row>
      <v-row>
        <v-btn class="mr-4" block @click="submit" x-large :loading="isPosting">
          Create
        </v-btn>
      </v-row>
    </v-container>
  </v-form>
</template>

<script>
import BasicList from "./BasicList.vue";
import RecipeNutrition from "./RecipeNutrition.vue";
import RecipeTimes from "./RecipeTimes.vue";

export default {
  name: "CreateManualForm",
  components: {
    BasicList,
    RecipeNutrition,
    RecipeTimes,
  },
  data: () => ({
    valid: false,
    title: "",
    category: "",
    source: "",
    description: "",
    yields: 1,
    rules: {
      title: [(v) => !!v || "Title is required"],
      category: [(v) => !!v || "Category is required"],
      source: [(v) => !!v || "Source is required"],
      description: [(v) => !!v || "Description is required"],
      number: [
        (val) => {
          if (val <= 0) {
            return "Please enter a positive number >= 0";
          }
          return true;
        },
      ],
      numberTimes: [
        (val) => {
          if (val < 0) {
            return "Please enter a positive number > 0";
          }
          return true;
        },
      ],
    },
  }),
  computed: {
    categories() {
      return this.$store.getters["browse/categories"].map((c) =>
        this.$toTitleCase(c)
      );
    },
    isPosting() {
      return this.$store.getters["create/isPosting"];
    },
    isXs() {
      return this.$vuetify.breakpoint.xs;
    },
  },
  methods: {
    submit() {
      if (!this.$refs.form.validate()) {
        return;
      }

      const now = new Date();
      const nowISO = now.toISOString();
      const times = this.$refs.times.getTimes();

      const recipe = {
        name: this.title,
        description: this.description,
        url: this.source,
        image: "",
        prepTime: times.prep,
        cookTime: times.cook,
        totalTime: times.total,
        recipeCategory: this.category,
        keywords: "",
        recipeYield: this.yields,
        tool: [],
        recipeIngredient: this.$refs.ingredients.getItems(),
        recipeInstructions: this.$refs.instructions.getItems(),
        nutrition: this.$refs.nutrition.getData(),
        dateCreated: nowISO,
        dateModified: nowISO,
      };

      this.$store.dispatch("create/postRecipe", recipe);
    },
  },
};
</script>
