package auth_test

import (
	"encoding/base64"
	"github.com/reaper47/recipya/internal/auth"
	"testing"
)

func TestGenerateSelectorAndValidator(t *testing.T) {
	t.Run("generated selector and validator have the correct lengths", func(t *testing.T) {
		selector, validator := auth.GenerateSelectorAndValidator()
		selectorBytes, err := base64.URLEncoding.DecodeString(selector)
		if err != nil {
			t.Errorf("failed to decode selector: %q", err)
		}
		if len(selectorBytes) != 12 {
			t.Errorf("invalid selector length, expected 12 but got %d", len(selectorBytes))
		}

		validatorBytes, err := base64.URLEncoding.DecodeString(validator)
		if err != nil {
			t.Errorf("failed to decode validator: %q", err)
		}
		if len(validatorBytes) != 32 {
			t.Errorf("invalid validator length, expected 32, got %d", len(validatorBytes))
		}
	})

	t.Run("generated selector and validator are unique across multiple invocations", func(t *testing.T) {
		selectorSet := make(map[string]bool)
		validatorSet := make(map[string]bool)
		for i := 0; i < 1000; i++ {
			selector, validator := auth.GenerateSelectorAndValidator()

			if _, exists := selectorSet[selector]; exists {
				t.Errorf("duplicate selector found: %s", selector)
			}
			selectorSet[selector] = true

			if _, exists := validatorSet[validator]; exists {
				t.Errorf("duplicate validator found: %s", validator)
			}
			validatorSet[validator] = true
		}
	})
}
