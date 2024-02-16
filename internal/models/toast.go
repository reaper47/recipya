package models

import (
	"encoding/json"
)

// NewErrorToast creates a new error notification.
func NewErrorToast(title, message, action string) Toast {
	return Toast{
		Action:     action,
		Background: "alert-error",
		Message:    message,
		Title:      title,
	}
}

// NewInfoToast creates a new info notification.
func NewInfoToast(title, message, action string) Toast {
	return Toast{
		Action:     action,
		Background: "alert-info",
		Message:    message,
		Title:      title,
	}
}

// NewWarningToast creates a new warning notification.
func NewWarningToast(title, message, action string) Toast {
	return Toast{
		Action:     action,
		Background: "alert-warning",
		Message:    message,
		Title:      title,
	}
}

// Toast holds data related to a notification toast.
type Toast struct {
	Action     string `json:"action"`
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
