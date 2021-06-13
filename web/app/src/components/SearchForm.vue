<template>
  <div></div>
</template>
<script>
export default {
  data() {
    return {
      valid: null,
      validators: {
        required: (value) => !!value || "This field is required",
      },
      mode: 1,
      numRecipes: 10,
      ingredients: [],
      modes: [
        { label: "Maximize items taken from the fridge", value: 1 },
        { label: "Minimize number of items to buy", value: 2 },
      ],
    };
  },
  methods: {
    search() {
      return 1;
    },
  },
};
</script>
