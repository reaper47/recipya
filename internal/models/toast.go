package models

import (
	"encoding/json"
)

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

// NewErrorToast creates a new error notification.
func NewErrorToast(title, message, action string) Toast {
	return Toast{
		Action:     action,
		Background: "alert-error",
		Message:    message,
		Title:      title,
	}
}

func NewErrorAuthToast(message string) Toast {
	return NewErrorToast("Auth Error", message, "")
}

func NewErrorDBToast(message string) Toast {
	return NewErrorToast("Database Error", message, "")
}

func NewErrorFilesToast(message string) Toast {
	return NewErrorToast("Files Error", message, "")
}

func NewErrorFormToast(message string) Toast {
	return NewErrorToast("Form Error", message, "")
}

func NewErrorGeneralToast(message string) Toast {
	return NewErrorToast("General Error", message, "")
}

func NewErrorReqToast(message string) Toast {
	return NewErrorToast("Request Error", message, "")
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

func NewWarningWSToast(message string) Toast {
	return NewWarningToast("Websocket", message, "")
}
