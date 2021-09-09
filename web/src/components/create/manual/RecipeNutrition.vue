<template>
  <v-data-table
    :headers="headers"
    :items="items"
    hide-default-footer
    class="elevation-1"
    dense
  >
    <template #item.amount="{ item }">
      <v-text-field
        v-model.number="item.amount"
        :placeholder="item.placeholder"
        single-line
        dense
        hide-details
        class="pb-2"
      ></v-text-field> </template
  ></v-data-table>
</template>

<script>
export default {
  name: "RecipeNutrition",
  data: () => ({
    headers: [
      { text: "Nutrition (per 100g)", value: "nutrition" },
      { text: "Amount (not required)", value: "amount" },
    ],
    items: [
      {
        nutrition: "Calories",
        amount: "",
        placeholder: "368 calories",
        field: "calories",
      },
      {
        nutrition: "Carbohydrate",
        amount: "",
        placeholder: "35g",
        field: "carbohydrateContent",
      },
      { nutrition: "Fat", amount: "", placeholder: "15g", field: "fatContent" },
      {
        nutrition: "Saturated Fat",
        amount: "",
        placeholder: "1.8g",
        field: "saturatedFatContent",
      },
      {
        nutrition: "Cholesterol",
        amount: "",
        placeholder: "1.1mg",
        field: "cholesterolContent",
      },
      {
        nutrition: "Protein",
        amount: "",
        placeholder: "21g",
        field: "proteinContent",
      },
      {
        nutrition: "Sodium",
        amount: "",
        placeholder: "100mg",
        field: "sodiumContent",
      },
      {
        nutrition: "Fiber",
        amount: "",
        placeholder: "8g",
        field: "fiberContent",
      },
      {
        nutrition: "Sugar",
        amount: "",
        placeholder: "3g",
        field: "sugarContent",
      },
    ],
  }),
  methods: {
    getData() {
      const data = {};
      this.items.forEach((item) => (data[item.field] = item.amount.toString()));
      return data;
    },
  },
};
</script>
