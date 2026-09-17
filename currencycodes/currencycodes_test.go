package currencycodes_test

import (
	"fmt"
	"testing"

	"github.com/mvndaai/money/currencycodes"
	"github.com/stretchr/testify/assert"
)

func TestBestBestDecimals(t *testing.T) {
	decimal0 := currencycodes.BIF
	decimal2 := currencycodes.USD
	decimal3 := currencycodes.BHD
	decimal4 := currencycodes.CLF

	tests := []struct {
		codes    []string
		expected int
	}{
		{codes: nil, expected: 4},
		{codes: []string{""}, expected: 4},
		{codes: []string{decimal0}, expected: 0},
		{codes: []string{decimal2}, expected: 2},
		{codes: []string{decimal3}, expected: 3},
		{codes: []string{decimal4}, expected: 4},
		{codes: []string{"", "unknown", decimal4}, expected: 4},
	}
	for _, tt := range tests {
		name := fmt.Sprint(tt.codes)
		t.Run(name, func(t *testing.T) {
			actual := currencycodes.BestDecimal(tt.codes...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestFirstValidCode(t *testing.T) {
	tests := []struct {
		name     string
		codes    []string
		expected string
	}{
		{name: "nil", codes: nil, expected: ""},
		{name: "empty", codes: []string{""}, expected: ""},
		{name: "unknown", codes: []string{"unknown"}, expected: ""},
		{name: "valid", codes: []string{currencycodes.USD}, expected: currencycodes.USD},
		{name: "lowercase", codes: []string{"usd"}, expected: currencycodes.USD},
		{name: "first valid wins", codes: []string{currencycodes.USD, currencycodes.BIF}, expected: currencycodes.USD},
		{name: "skips invalid", codes: []string{"", "unknown", currencycodes.BIF}, expected: currencycodes.BIF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := currencycodes.FirstValidCode(tt.codes...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
