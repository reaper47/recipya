import EventBus from ".";

export const ACTION = {
  SNACKBAR: "snackbar",
};

export const SNACKBAR_TYPE = Object.freeze({
  ERROR: {
    show: false,
    color: "#D32F2F",
    type: "multi-line",
    position: "top",
    timeout: 5000,
    icon: "mdi-alert",
  },
  INFO: 2,
  SUCCESS: 3,
  WARNING: 4,
});

export const showSnackbar = (snackbarType, title, message) => {
  EventBus.$emit(ACTION.SNACKBAR, snackbarType, title, message);
};
