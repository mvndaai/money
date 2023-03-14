// This file is using package "money" instead of "money_test" to be able to unit test private functions
package money

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShiftForSave(t *testing.T) {
	tests := []struct {
		name     string
		in       float64
		digits   int
		expected int64
	}{
		{name: "base", in: 10.01, digits: 4, expected: 100100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := shiftForSave(tt.in, tt.digits)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestShiftForUse(t *testing.T) {
	tests := []struct {
		name     string
		in       int64
		digits   int
		expected float64
	}{
		{name: "base", in: 1, digits: 4, expected: 0.0001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := shiftForUse(tt.in, tt.digits)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
