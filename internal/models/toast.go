package models

import (
	"encoding/json"
)

// NewErrorToast creates a new error notification.
func NewErrorToast(title, message string) Toast {
	return Toast{
		Background: "alert-error",
		Message:    message,
		Title:      title,
	}
}

// NewInfoToast creates a new info notification.
func NewInfoToast(title, message string) Toast {
	return Toast{
		Background: "alert-info",
		Message:    message,
		Title:      title,
	}
}

// NewWarningToast creates a new warning notification.
func NewWarningToast(title, message string) Toast {
	return Toast{
		Background: "alert-warning",
		Message:    message,
		Title:      title,
	}
}

// Toast holds data related to a notification toast.
type Toast struct {
	Background string `json:"background"`
	Message    string `json:"message"`
	Title      string `json:"title"`
}

// Render returns the JSON encoding of the toast.
func (t Toast) Render() string {
	toast, _ := json.Marshal(t)
	xb, _ := json.Marshal(map[string]string{"showToast": string(toast)})
	return string(xb)
}
