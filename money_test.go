package money_test

import (
	"encoding/json"
	"testing"

	"github.com/mvndaai/money"
	"github.com/stretchr/testify/assert"
)

func TestParseInt(t *testing.T) {
	const in = 10
	expected := "10.0000"

	mI := money.ParseInt(in)
	mI64 := money.ParseInt64(in)
	assert.True(t, mI.Equal(mI64))
	assert.Equal(t, expected, mI.String())
}

func TestParseFloat64(t *testing.T) {
	tests := []struct {
		name     string
		in       float64
		expected float64
	}{
		{name: "normal", in: 1.25, expected: 1.2500},
		{name: "round down", in: 1.00004, expected: 1.0000},
		{name: "round up", in: 1.000051, expected: 1.0001},
		{name: "cents only", in: .000051, expected: .0001},
		{name: "neagive", in: -.25, expected: -.25},
		{name: "neagive round ", in: -.00006, expected: -.0001},
		{name: "large number", in: 123456789.123456, expected: 123456789.1235},
		{name: "roundToEven 0.5", in: .0000_5, expected: .0000},
		{name: "roundToEven 1.5", in: .0001_5, expected: .0002},
		{name: "roundToEven 2.5", in: .0002_5, expected: .0002},
		{name: "roundToEven 3.5", in: .0003_5, expected: .0004},
		{name: "roundToEven 4.5", in: .0004_5, expected: .0004},
		{name: "roundToEven 5.5", in: .0005_5, expected: .0006},
		{name: "roundToEven 6.5", in: .0006_5, expected: .0006},
		{name: "roundToEven 7.5", in: .0007_5, expected: .0008},
		{name: "roundToEven 8.5", in: .0008_5, expected: .0008},
		{name: "roundToEven 9.5", in: .0009_5, expected: .0010},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := money.ParseFloat64(tt.in)
			actual := m.Float64()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		name     string
		in       []float64
		expected float64
	}{
		{name: "simple", in: []float64{5.0005, 5.0005}, expected: 25.0050},
		{name: "negatives", in: []float64{5, -5}, expected: -25},
		{name: "double negative", in: []float64{-5, -5}, expected: 25},
		{name: "empty", in: []float64{}, expected: 0},
		{name: "nil", in: nil, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ms []money.Money
			for _, v := range tt.in {
				ms = append(ms, money.ParseFloat64(v))
			}
			actual := money.Mul(ms...).Float64()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestQuo(t *testing.T) {
	tests := []struct {
		name        string
		in          []float64
		expected    float64
		expectError bool
	}{
		{name: "ints", in: []float64{25, 5}, expected: 5},
		{name: "simple", in: []float64{25.0050, 5.0005}, expected: 5.0005},
		{name: "negative", in: []float64{25, -5}, expected: -5},
		{name: "empty", in: []float64{}, expected: 0, expectError: true},
		{name: "nil", in: nil, expected: 0, expectError: true},
		{name: "divide by 0", in: []float64{10, 0}, expected: 0, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ms []money.Money
			for _, v := range tt.in {
				ms = append(ms, money.ParseFloat64(v))
			}
			actual, err := money.Quo(ms...)
			assert.Equal(t, tt.expectError, err != nil)
			assert.Equal(t, tt.expected, actual.Float64())
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		in       []float64
		expected float64
	}{
		{name: "ints", in: []float64{25, 5}, expected: 30},
		{name: "simple", in: []float64{5.0005, 5.0005}, expected: 10.0010},
		{name: "negative", in: []float64{10, -5}, expected: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ms []money.Money
			for _, v := range tt.in {
				ms = append(ms, money.ParseFloat64(v))
			}
			actual := money.Add(ms...).Float64()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		name     string
		in       []float64
		expected float64
	}{
		{name: "ints", in: []float64{25, 5}, expected: 20},
		{name: "simple", in: []float64{25.0050, 5.0005}, expected: 20.0045},
		{name: "negative", in: []float64{25, -5}, expected: 30},
		{name: "double negative", in: []float64{-25, -5}, expected: -20},
		{name: "empty", in: []float64{}, expected: 0},
		{name: "nil", in: nil, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ms []money.Money
			for _, v := range tt.in {
				ms = append(ms, money.ParseFloat64(v))
			}
			actual := money.Sub(ms...).Float64()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected float64
		hasError bool
	}{
		{name: "number", in: "1.25", expected: 1.2500, hasError: false},
		{name: "string", in: "$$$$", expected: 0, hasError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := money.ParseString(tt.in)
			assert.Equal(t, tt.hasError, err != nil)

			actual := m.Float64()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCurrencyRound(t *testing.T) {
	tests := []struct {
		name         string
		value        float64
		currencyCode string
		expected     float64
		expectError  bool
	}{
		{name: "2 decmails", value: 1.2666, currencyCode: "USD", expected: 1.27},
		{name: "3 decmails", value: 1.2666, currencyCode: "BHD", expected: 1.267},
		{name: "4 decmails", value: 1.2666, currencyCode: "CLF", expected: 1.2666},
		{name: "4 decmails rounded", value: 1.26666, currencyCode: "CLF", expected: 1.2667},
		{name: "unknown code", value: 1, currencyCode: "ABC", expected: 0, expectError: true},
		{name: "2 decmails down", value: 1.2222, currencyCode: "USD", expected: 1.22},
		{name: "3 decmails down", value: 1.2222, currencyCode: "BHD", expected: 1.222},
		{name: "4 decmails down", value: 1.2222, currencyCode: "CLF", expected: 1.2222},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := money.ParseFloat64(tt.value)
			actual, err := m.CurrencyFloat64(tt.currencyCode)
			assert.Equal(t, tt.expectError, err != nil)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCurrencyString(t *testing.T) {
	tests := []struct {
		name         string
		value        float64
		currencyCode string
		expected     string
		expectError  bool
	}{
		{name: "2 decmails", value: 1.2555, currencyCode: "USD", expected: "1.26"},
		{name: "3 decmails", value: 1.2555, currencyCode: "BHD", expected: "1.256"},
		{name: "4 decmails", value: 1.2555, currencyCode: "CLF", expected: "1.2555"},
		{name: "4 decmails rounded input", value: 1.25555, currencyCode: "CLF", expected: "1.2556"},
		{name: "cents only", value: 0.0051, currencyCode: "USD", expected: "0.01"},
		{name: "cased insensitivity", value: 0.0060, currencyCode: "usD", expected: "0.01"},
		{name: "unknown code", value: 1, currencyCode: "ABC", expected: "", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := money.ParseFloat64(tt.value)
			actual, err := m.CurrencyString(tt.currencyCode)
			assert.Equal(t, tt.expectError, err != nil)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestJSON(t *testing.T) {
	tests := []struct {
		name              string
		inJSON            string
		expectedJSON      string
		currencyCode      string
		hasUnmarshalError bool
	}{
		{
			name:              "number",
			inJSON:            "1.00009",
			expectedJSON:      `"1.0001"`,
			hasUnmarshalError: false,
		},
		{
			name:              "str number",
			inJSON:            `"1.00009"`,
			expectedJSON:      `"1.0001"`,
			hasUnmarshalError: false,
		},
		{
			name:              "non number",
			inJSON:            "a",
			expectedJSON:      "",
			hasUnmarshalError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m money.Money
			err := json.Unmarshal([]byte(tt.inJSON), &m)
			assert.Equal(t, tt.hasUnmarshalError, err != nil)
			if err != nil {
				return
			}

			b, err := json.Marshal(m)
			assert.Nil(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, tt.expectedJSON, string(b))
		})
	}
}

func TestJSONStruct(t *testing.T) {
	tests := []struct {
		name              string
		inJSON            string
		expectedJSON      string
		hasUnmarshalError bool
	}{
		{
			name:              "number",
			inJSON:            `{"A":1.00009}`,
			expectedJSON:      `{"A":"1.0001"}`,
			hasUnmarshalError: false,
		},
		{
			name:              "str number",
			inJSON:            `{"A":"1.00009"}`,
			expectedJSON:      `{"A":"1.0001"}`,
			hasUnmarshalError: false,
		},
		{
			name:              "non number",
			inJSON:            `{"A":"a"}`,
			expectedJSON:      "",
			hasUnmarshalError: true,
		},
	}

	type B struct {
		A money.Money
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c B
			err := json.Unmarshal([]byte(tt.inJSON), &c)
			assert.Equal(t, tt.hasUnmarshalError, err != nil)
			if err != nil {
				return
			}

			b, err := json.Marshal(c)
			assert.Nil(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, tt.expectedJSON, string(b))
		})
	}
}
