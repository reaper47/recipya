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

// NewErrorAuthToast creates an action-less error Toast with the title of "Auth Error".
func NewErrorAuthToast(message string) Toast {
	return NewErrorToast("Auth Error", message, "")
}

// NewErrorDBToast creates an action-less error Toast with the title of "Database Error".
func NewErrorDBToast(message string) Toast {
	return NewErrorToast("Database Error", message, "")
}

// NewErrorFilesToast creates an action-less error Toast with the title of "Files Error".
func NewErrorFilesToast(message string) Toast {
	return NewErrorToast("Files Error", message, "")
}

// NewErrorFormToast creates an action-less error Toast with the title of "Form Error".
func NewErrorFormToast(message string) Toast {
	return NewErrorToast("Form Error", message, "")
}

// NewErrorGeneralToast creates an action-less error Toast with the title of "General Error".
func NewErrorGeneralToast(message string) Toast {
	return NewErrorToast("General Error", message, "")
}

// NewErrorReqToast creates an action-less error Toast with the title of "Request Error".
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

// NewWarningWSToast creates an action-less warning Toast with the title of "Websocket".
func NewWarningWSToast(message string) Toast {
	return NewWarningToast("Websocket", message, "")
}
