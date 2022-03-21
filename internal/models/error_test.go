package models

import "testing"

func TestModelError(t *testing.T) {
	t.Run("NewErrorJSON creates a proper JSON object for the error", func(t *testing.T) {
		actual, err := NewErrorJSON(404, "An error message")
		if err != nil {
			t.Fatalf("error creating a a JSON object: %s", err)
		}

		expected := `{"error":{"code":404,"message":"An error message","status":"Not Found"}}`
		if string(actual) != expected {
			t.Fatalf("wanted %s but got %s", expected, actual)
		}
	})
}
