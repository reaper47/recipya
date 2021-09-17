<template>
  <v-col style="border-right: 1px solid black">
    <v-data-table
      dense
      :headers="headers"
      :items="items"
      item-key="name"
      class="elevation-1 capitalize"
      hide-default-footer
    >
      <template #item.amount="{ item }">
        <div class="lowercase">
          {{ item.amount.toLowerCase() }}
        </div>
      </template>
    </v-data-table>
  </v-col>
</template>

<script>
export default {
  name: "RecipeNutrition",
  props: {
    nutrition: {
      type: Object,
      required: true,
    },
  },
  data: () => ({
    headers: [
      {
        text: "Nutrition (per 100g)",
        align: "start",
        sortable: true,
        value: "name",
      },
      { text: "Amount", value: "amount" },
    ],
    items: [],
  }),
  mounted() {
    this.items = Object.entries(this.nutrition)
      .filter(([, amount]) => amount !== "")
      .map(([name, amount]) => {
        return {
          amount,
          name: name.replace("Content", "").replace("Fat", " Fat"),
        };
      });
  },
};
</script>
