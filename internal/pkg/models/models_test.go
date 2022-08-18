package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat64(t *testing.T) {
	t.Run("return float64 from currency", func(t *testing.T) {
		c := Currency(10535)
		expected := 105.35
		actual := c.Float64()
		assert.Equal(t, expected, actual)
	})
}

func TestToCurrency(t *testing.T) {
	t.Run("return currency from float64", func(t *testing.T) {
		value := 105.35
		expected := Currency(10535)
		actual := ToCurrency(value)
		assert.Equal(t, expected, actual)
	})
}
