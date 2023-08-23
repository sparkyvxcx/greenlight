package data

import (
	"testing"

	"greenlight.sparkyvxcx.co/internal/assert"
	"greenlight.sparkyvxcx.co/internal/validator"
)

func TestValidateFilters(t *testing.T) {
	t.Run("Reject page value smaller than 0", func(t *testing.T) {
		v := validator.New()
		filters := Filters{Page: -1, PageSize: 1, Sort: "drama", SortWhitelist: []string{"drama"}}

		ValidateFilters(v, filters)

		assert.Equal(t, v.Valid(), false)
		assert.Equal(t, v.Errors["page"], "must be greater than zero")
	})

	t.Run("Reject page value 0", func(t *testing.T) {
		v := validator.New()
		filters := Filters{Page: 0, PageSize: 1, Sort: "drama", SortWhitelist: []string{"drama"}}

		ValidateFilters(v, filters)

		assert.Equal(t, v.Valid(), false)
		assert.Equal(t, v.Errors["page"], "must be greater than zero")
	})

	t.Run("Reject page value larger than 10,000,000", func(t *testing.T) {
		v := validator.New()
		filters := Filters{Page: 10_000_001, PageSize: 1, Sort: "drama", SortWhitelist: []string{"drama"}}

		ValidateFilters(v, filters)

		assert.Equal(t, v.Valid(), false)
		assert.Equal(t, v.Errors["page"], "must be a maximum of 10 million")
	})

	t.Run("Page value should within range 1 to 10 million", func(t *testing.T) {
		v := validator.New()
		filters := Filters{Page: 9_999_999, PageSize: 1, Sort: "drama", SortWhitelist: []string{"drama"}}

		ValidateFilters(v, filters)

		assert.Equal(t, v.Valid(), true)
		assert.Equal(t, v.Errors["page"], "")
	})

	t.Run("PageSize value must be greater than zero", func(t *testing.T) {
		v := validator.New()
		filters := Filters{Page: 10_000, PageSize: 0, Sort: "drama", SortWhitelist: []string{"drama"}}

		ValidateFilters(v, filters)

		assert.Equal(t, v.Valid(), false)
		assert.Equal(t, v.Errors["page_size"], "must be greater than zero")
	})

	t.Run("PageSize value must be not greater than 100", func(t *testing.T) {
		v := validator.New()
		filters := Filters{Page: 10000, PageSize: 101, Sort: "drama", SortWhitelist: []string{"drama"}}

		ValidateFilters(v, filters)

		assert.Equal(t, v.Valid(), false)
		assert.Equal(t, v.Errors["page_size"], "must be a maximum of 100")
	})

	t.Run("PageSize value should within greater than 100", func(t *testing.T) {
		v := validator.New()
		filters := Filters{Page: 10000, PageSize: 99, Sort: "drama", SortWhitelist: []string{"drama"}}

		ValidateFilters(v, filters)

		assert.Equal(t, v.Valid(), true)
		assert.Equal(t, v.Errors["page_size"], "")
	})

	t.Run("Sort value should in whitelist", func(t *testing.T) {
		v := validator.New()
		filters := Filters{Page: 10000, PageSize: 100, Sort: "id", SortWhitelist: []string{"drama"}}

		ValidateFilters(v, filters)

		assert.Equal(t, v.Valid(), false)
		assert.Equal(t, v.Errors["sort"], "invalid sort value")
	})
}
