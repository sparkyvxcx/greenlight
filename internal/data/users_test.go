package data

import (
	"testing"

	"greenlight.sparkyvxcx.co/internal/assert"
	"greenlight.sparkyvxcx.co/internal/validator"
)

func TestValidateEmail(t *testing.T) {
	t.Run("Valid Email Should Pass", func(t *testing.T) {
		v := validator.New()
		email := "alice@example.com"

		ValidateEmail(v, email)

		assert.Equal(t, v.Valid(), true)
	})

	t.Run("Invalid Email Should Fail", func(t *testing.T) {
		v := validator.New()
		email := "alice@"

		ValidateEmail(v, email)

		assert.Equal(t, v.Valid(), false)
		assert.Equal(t, v.Errors["email"], "must be a valid email address")
	})
}
