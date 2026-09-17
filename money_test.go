package money_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/mvndaai/money"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const delta = .0001

func TestParseInt(t *testing.T) {
	const in = 10
	expected := "10.0000"

	mI := money.Parse(in)
	mI64 := money.Parse(10.0)
	assert.True(t, mI.IsEqual(mI64))
	assert.Equal(t, expected, mI.String())
}

func TestParseBigFloat(t *testing.T) {
	in := big.NewFloat(12.34567)
	expected := "12.3457"

	mI := money.ParseBigFloat(*in)
	assert.Equal(t, expected, mI.StringDecimals(4))
}

func TestParseFloat64(t *testing.T) {
	tests := []struct {
		name     string
		in       float64
		expected float64
	}{
		{name: "normal", in: 1.25, expected: 1.2500},
		{name: "round down", in: 1.0000_4, expected: 1.0000},
		{name: "round up", in: 1.0000_51, expected: 1.0001},
		{name: "cents only", in: .0000_51, expected: .0001},
		{name: "neagive", in: -.25, expected: -.25},
		{name: "neagive round ", in: -.00006, expected: -.0001},
		{name: "large number", in: 123456789.1234_56, expected: 123456789.1235},
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
			m := money.Parse(tt.in)
			actual := m.Float64()
			assert.InDelta(t, tt.expected, actual, delta)
		})
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		name     string
		in1      float64
		in2      []float64
		expected float64
	}{
		{name: "simple", in1: 5.0005, in2: []float64{5.0005}, expected: 25.0050},
		{name: "negatives", in1: 5, in2: []float64{-5}, expected: -25},
		{name: "double negative", in1: -5, in2: []float64{-5}, expected: 25.0000},
		{name: "empty", in1: 1, in2: []float64{}, expected: 1},
		{name: "nil", in1: 1, in2: nil, expected: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := money.Mul(tt.in1, tt.in2...).Float64()
			assert.InDelta(t, tt.expected, actual, delta)
		})
	}
}

func TestQuo(t *testing.T) {
	tests := []struct {
		name        string
		in1         float64
		in2         []float64
		expected    float64
		expectError bool
	}{
		{name: "ints", in1: 25, in2: []float64{5, 1}, expected: 5},
		{name: "simple", in1: 25.0050, in2: []float64{5.0005}, expected: 5.0005},
		{name: "negative", in1: 25, in2: []float64{-5}, expected: -5},
		{name: "empty", in1: 0, in2: []float64{}, expected: 0, expectError: false},
		{name: "nil", in1: 0, in2: nil, expected: 0, expectError: false},
		{name: "divide by 0", in1: 10, in2: []float64{0}, expected: 0, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := money.Quo(tt.in1, tt.in2...)
			require.Equal(t, tt.expectError, err != nil)
			assert.InDelta(t, tt.expected, actual.Float64(), delta)
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
			actual := money.Add(tt.in[0], tt.in[1:]...).Float64()
			assert.InDelta(t, tt.expected, actual, delta)
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
		{name: "empty", in: []float64{0}, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := money.Sub(tt.in[0], tt.in[1:]...).Float64()
			assert.InDelta(t, tt.expected, actual, delta)
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
			m, err := money.ParseString(context.Background(), tt.in)
			assert.Equal(t, tt.hasError, err != nil)

			actual := m.Float64()
			assert.InDelta(t, tt.expected, actual, delta)
		})
	}
}

func TestCurrencyRound(t *testing.T) {
	tests := []struct {
		name         string
		value        float64
		currencyCode string
		expected     float64
	}{
		{name: "2 decimals", value: 1.2666, currencyCode: "USD", expected: 1.27},
		{name: "3 decimals", value: 1.2666, currencyCode: "BHD", expected: 1.267},
		{name: "4 decimals", value: 1.2666, currencyCode: "CLF", expected: 1.2666},
		{name: "4 decimals rounded", value: 1.26666, currencyCode: "CLF", expected: 1.2667},
		{name: "unknown code", value: 1, currencyCode: "ABC", expected: 1},
		{name: "2 decimals down", value: 1.2222, currencyCode: "USD", expected: 1.22},
		{name: "3 decimals down", value: 1.2222, currencyCode: "BHD", expected: 1.222},
		{name: "4 decimals down", value: 1.2222, currencyCode: "CLF", expected: 1.2222},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := money.Parse(tt.value)
			actual := m.Float64CurrencyCode(tt.currencyCode)
			assert.InDelta(t, tt.expected, actual, delta)
		})
	}
}

func TestCurrencyString(t *testing.T) {
	tests := []struct {
		name         string
		value        float64
		currencyCode string
		expected     string
	}{
		{name: "2 decimals", value: 1.2555, currencyCode: "USD", expected: "1.26"},
		{name: "3 decimals", value: 1.2555, currencyCode: "BHD", expected: "1.256"},
		{name: "4 decimals", value: 1.2555, currencyCode: "CLF", expected: "1.2555"},
		{name: "4 decimals rounded input", value: 1.25555, currencyCode: "CLF", expected: "1.2556"},
		{name: "cents only", value: 0.0051, currencyCode: "USD", expected: "0.01"},
		{name: "cased insensitivity", value: 0.0060, currencyCode: "usD", expected: "0.01"},
		{name: "unknown code", value: 1, currencyCode: "ABC", expected: "1.0000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := money.Parse(tt.value)
			actual := m.StringCurrencyCode(tt.currencyCode)
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
			inJSON:            "1.009",
			expectedJSON:      `"1.0090"`,
			hasUnmarshalError: false,
		},
		{
			name:              "str number",
			inJSON:            `"1.009"`,
			expectedJSON:      `"1.0090"`,
			hasUnmarshalError: false,
		},
		{
			name:              "non number",
			inJSON:            "a",
			expectedJSON:      "",
			hasUnmarshalError: true,
		},
		{
			name:              "huge number",
			inJSON:            "1e9999",
			expectedJSON:      "",
			hasUnmarshalError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m money.Money
			err := json.Unmarshal([]byte(tt.inJSON), &m)
			assert.Equal(t, tt.hasUnmarshalError, err != nil, err)
			if err != nil {
				return
			}

			b, err := json.Marshal(m)
			require.NoError(t, err)
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
			inJSON:            `{"A":1.009}`,
			expectedJSON:      `{"A":"1.0090"}`,
			hasUnmarshalError: false,
		},
		{
			name:              "str number",
			inJSON:            `{"A":"1.009"}`,
			expectedJSON:      `{"A":"1.0090"}`,
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
			require.NoError(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, tt.expectedJSON, string(b))
		})
	}
}

func TestComparors(t *testing.T) {
	tests := []struct {
		name             string
		in               float64
		expectedZero     bool
		expectedPositive bool
		expectedNegative bool
	}{
		{name: "zero", in: 0, expectedZero: true, expectedPositive: false, expectedNegative: false},
		{name: "pos", in: 1, expectedZero: false, expectedPositive: true, expectedNegative: false},
		{name: "neg", in: -.0001, expectedZero: false, expectedPositive: false, expectedNegative: true},
		{name: "pos small", in: .0001, expectedZero: false, expectedPositive: true, expectedNegative: false},
		{name: "neg small", in: -.0001, expectedZero: false, expectedPositive: false, expectedNegative: true},
		{name: "pos rounded out", in: .00001, expectedZero: true, expectedPositive: false, expectedNegative: false},
		{name: "neg rounded out", in: -.00001, expectedZero: true, expectedPositive: false, expectedNegative: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := money.Parse(tt.in)
			assert.Equal(t, tt.expectedZero, m.IsZero())
			assert.Equal(t, tt.expectedPositive, m.IsPositive())
			assert.Equal(t, tt.expectedNegative, m.IsNegative())
		})
	}
}

func TestNegate(t *testing.T) {
	tests := []struct {
		name     string
		in       float64
		expected float64
	}{
		{name: "zero", in: 0, expected: 0},
		{name: "pos", in: 1, expected: -1},
		{name: "neg", in: -1, expected: 1},
		{name: "pos small", in: .0001, expected: -.0001},
		{name: "neg small", in: -.0001, expected: .0001},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := money.Parse(tt.in)
			actual := money.Negate(m).Float64()
			assert.InDelta(t, tt.expected, actual, delta)
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		in       float64
		expected float64
	}{
		{name: "zero", in: 0, expected: 0},
		{name: "pos", in: 1, expected: 1},
		{name: "neg", in: -1, expected: 1},
		{name: "pos small", in: .0001, expected: .0001},
		{name: "neg small", in: -.0001, expected: .0001},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := money.Parse(tt.in)
			actual := money.Abs(m).Float64()
			assert.InDelta(t, tt.expected, actual, delta)
		})
	}
}

func TestRoundToDecimals(t *testing.T) {
	tests := []struct {
		name     string
		in       float64
		dec      int
		expected float64
	}{
		{name: "zero", in: 0, dec: 0, expected: 0},
		{name: "pos", in: 1.1234, dec: 2, expected: 1.12},
		{name: "neg", in: -1.1234, dec: 2, expected: -1.12},
		{name: "pos-4", in: 1.1234, dec: 4, expected: 1.1234},
		{name: "neg-4", in: -1.1234, dec: 4, expected: -1.1234},
		{name: "pos-over4", in: 1.1234, dec: 10, expected: 1.1234},
		{name: "neg-over4", in: -1.1234, dec: 10, expected: -1.1234},
		{name: "pos-dec0", in: 1.1234, dec: 0, expected: 1},
		{name: "neg-dec0", in: -1.1234, dec: 0, expected: -1},
		{name: "pos-dec-neg", in: 1.1234, dec: -2, expected: 1.1234},
		{name: "neg-dec-neg", in: -1.1234, dec: -2, expected: -1.1234},
		{name: "+d0r", in: 1.9111, dec: 0, expected: 2.0000},
		{name: "+d1r", in: 1.1911, dec: 1, expected: 1.2000},
		{name: "+d2r", in: 1.1191, dec: 2, expected: 1.1200},
		{name: "+d3r", in: 1.1119, dec: 3, expected: 1.1120},
		{name: "+d4r", in: 1.11119, dec: 4, expected: 1.1112},
		{name: "+d5r", in: 1.111119, dec: 9, expected: 1.1111},
		{name: "-d0r", in: -1.9111, dec: 0, expected: -2.0000},
		{name: "-d1r", in: -1.1911, dec: 1, expected: -1.2000},
		{name: "-d2r", in: -1.1191, dec: 2, expected: -1.1200},
		{name: "-d3r", in: -1.1119, dec: 3, expected: -1.1120},
		{name: "-d4r", in: -1.11119, dec: 4, expected: -1.1112},
		{name: "-d5r", in: -1.111119, dec: 5, expected: -1.1111},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := money.Parse(tt.in)
			m = money.RoundToDecimals(m, tt.dec)
			actual := m.Float64()
			assert.InDelta(t, tt.expected, actual, delta)
		})
	}
}

func TestCurrencyStringWithFallback(t *testing.T) {
	tests := []struct {
		name          string
		value         float64
		currencyCodes []string
		expected      string
	}{
		{name: "USD", value: 1.5678, currencyCodes: []string{"USD"}, expected: "1.57"},
		{name: "UNKNOWN", value: 1.5678, currencyCodes: []string{"UNKNOWN"}, expected: "1.5678"},
		{name: "fallback", value: 1.5678, currencyCodes: []string{"UNKNOWN", "BHD", "USD"}, expected: "1.568"},
		{name: "empty", value: 1.5678, currencyCodes: []string{""}, expected: "1.5678"}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := money.Parse(tt.value)
			actual := m.StringCurrencyCode(tt.currencyCodes...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCmp(t *testing.T) {
	tests := []struct {
		name     string
		in1      float64
		in2      float64
		expected int
	}{
		{name: "equal", in1: 1.25, in2: 1.25, expected: 0},
		{name: "less", in1: 1.24, in2: 1.25, expected: -1},
		{name: "greater", in1: 1.26, in2: 1.25, expected: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m1 := money.Parse(tt.in1)
			m2 := money.Parse(tt.in2)
			actual := m1.Cmp(m2)
			assert.Equal(t, tt.expected, actual)

			// Ensure this matches big.Float
			bfAcutal := big.NewFloat(tt.in1).Cmp(big.NewFloat(tt.in2))
			assert.Equal(t, bfAcutal, actual)
		})
	}
}

func TestCents(t *testing.T) {
	tests := []struct {
		name     string
		in       int
		expected float64
	}{
		{name: "zero", in: 0, expected: 0},
		{name: "pos one", in: 1, expected: .01},
		{name: "neg one", in: -1, expected: -.01},
		{name: "pos ten", in: 10, expected: .10},
		{name: "neg ten", in: -10, expected: -.10},
		{name: "pos hund", in: 100, expected: 1},
		{name: "neg hund", in: -100, expected: -1},
		{name: "pos thou", in: 1000, expected: 10},
		{name: "neg thou", in: -1000, expected: -10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := money.ParseCents(tt.in)
			assert.InDelta(t, tt.expected, m.Float64(), delta)
			assert.Equal(t, tt.in, m.Cents())
		})
	}
}

func TestScan(t *testing.T) {
	tests := []struct {
		name        string
		in          any
		expected    money.Money
		expectedErr bool
	}{
		{name: "number", in: []byte("1.25"), expected: money.Parse(1.25)},
		{name: "float64", in: 1.25, expected: money.Parse(1.25)},
		{name: "nil", in: nil, expected: money.Parse(0)},
		{name: "non number", in: "a", expected: money.Parse(0), expectedErr: true},
		{name: "unsupported type", in: true, expected: money.Parse(0), expectedErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m money.Money
			err := m.Scan(tt.in)
			assert.Equal(t, tt.expectedErr, err != nil)
			assert.Equal(t, tt.expected, m)
		})
	}
}

func TestValue(t *testing.T) {
	m := money.Parse(1.25)
	v, err := m.Value()
	require.NoError(t, err)
	assert.InEpsilon(t, 1.25, v, delta)
}

func TestPercentage(t *testing.T) {
	tests := []struct {
		amount     float64
		percentage float64
		expected   money.Money
	}{
		{amount: 100, percentage: 10, expected: money.Parse(10)},
		{amount: 1, percentage: 25, expected: money.Parse(0.25)},
		{amount: 0, percentage: 50, expected: money.Parse(0)},
		{amount: 100, percentage: 0, expected: money.Parse(0)},
		{amount: 100, percentage: 100, expected: money.Parse(100)},
		{amount: 100, percentage: 150, expected: money.Parse(150)},
		{amount: -100, percentage: 10, expected: money.Parse(-10)},
		{amount: 100, percentage: -10, expected: money.Parse(-10)},
		// rounds up: 0.335 * 33% = 0.11055 -> ties to even at the 4th decimal (5 -> 6)
		{amount: 0.335, percentage: 33, expected: money.Parse(0.1106)},
		// rounds down: 0.665 * 33% = 0.21945 -> ties to even at the 4th decimal (4 stays)
		{amount: 0.665, percentage: 33, expected: money.Parse(0.2194)},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v-%v%s", tt.amount, tt.percentage, "%")
		t.Run(name, func(t *testing.T) {
			actual := money.Percentage(tt.amount, tt.percentage)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestGeneric(t *testing.T) {
	assert.Equal(t, "1.0100", money.Add(1, .01).String())
	assert.Equal(t, "10.0000", money.Sub(big.NewInt(20), 10).String())
	assert.Equal(t, "10.0000", money.Mul(big.NewFloat(1), 10).String())
	quo, err := money.Quo(*big.NewFloat(100), 10)
	require.NoError(t, err)
	assert.Equal(t, "10.0000", quo.String())
	assert.True(t, money.Parse(10).IsEqual(*big.NewInt(10)))

	assert.Equal(t, "1.0000", money.Parse(float32(1)).String())
	assert.Equal(t, "1.0000", money.Parse(int8(1)).String())
	assert.Equal(t, "1.0000", money.Parse(int16(1)).String())
	assert.Equal(t, "1.0000", money.Parse(int32(1)).String())
	assert.Equal(t, "1.0000", money.Parse(int64(1)).String())
	assert.Equal(t, "1.0000", money.Parse(uint(1)).String())
	assert.Equal(t, "1.0000", money.Parse(uint8(1)).String())
	assert.Equal(t, "1.0000", money.Parse(uint16(1)).String())
	assert.Equal(t, "1.0000", money.Parse(uint32(1)).String())
	assert.Equal(t, "1.0000", money.Parse(uint64(1)).String())
	assert.Equal(t, "1.0000", money.Parse(uintptr(1)).String())
	m := money.Parse(1)
	m = money.Add(0, &m)
	assert.Equal(t, "1.0000", money.Parse(&m).String())
}

func TestCurrencyCodeOffset(t *testing.T) {
	t.Parallel()
	tests := []struct {
		amount       float64
		currencyCode string
		expected     string
	}{
		{amount: 10.0099, currencyCode: "USD", expected: "10.01"},
		{amount: 10.0099, currencyCode: "", expected: "10.0099"},
		{amount: 10.0099, currencyCode: "BIF", expected: "10"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%v-%s", tt.amount, tt.currencyCode)
		t.Run(name, func(t *testing.T) {
			m := money.Parse(tt.amount, "", tt.currencyCode)
			assert.Equal(t, tt.expected, m.String())
		})
	}

	centTests := []struct {
		amount       int
		currencyCode string
		expected     string
	}{
		{amount: 99, currencyCode: "USD", expected: "0.99"},
		{amount: 99, currencyCode: "", expected: "0.9900"},
		{amount: 99, currencyCode: "BIF", expected: "1"},
	}

	for _, tt := range centTests {
		name := fmt.Sprintf("cnts: %v-%s", tt.amount, tt.currencyCode)
		t.Run(name, func(t *testing.T) {
			m := money.ParseCents(tt.amount, "", tt.currencyCode)
			assert.Equal(t, tt.expected, m.String())
		})
	}

	t.Run("bigInt", func(t *testing.T) {
		m := money.ParseBigFloat(*big.NewFloat(.99), "BIF")
		assert.Equal(t, "1", m.String())
	})
	t.Run("string", func(t *testing.T) {
		m, err := money.ParseString(context.Background(), ".99", "BIF")
		require.NoError(t, err)
		assert.Equal(t, "1", m.String())
	})

	t.Run("floatRounding", func(t *testing.T) {
		m := money.Parse(.99, "BIF")
		assert.InEpsilon(t, 1.0, m.Float64(), delta)
	})
}

func TestCurrencyCodeOffsetSavedThroughActions(t *testing.T) {
	t.Parallel()
	m := money.Abs(money.Parse(10, "USD"))
	assert.Equal(t, "10.00", m.String())

	m = money.Abs(money.Parse(-10, "USD"))
	assert.Equal(t, "10.00", m.String())

	m = money.Negate(money.Parse(-10, "USD"))
	assert.Equal(t, "10.00", m.String())

	m = money.Percentage(money.Parse(100, "USD"), 10)
	assert.Equal(t, "10.00", m.String())

	tests := []struct {
		first            string
		second           string
		expectedDecimals int
	}{
		{first: "USD", second: "", expectedDecimals: 2},
		{first: "", second: "USD", expectedDecimals: 2},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("multi-%s-%s", tt.first, tt.second), func(t *testing.T) {
			m := money.Add(money.Parse(5, tt.first), money.Parse(5, tt.second))
			assert.Equal(t, "10.00", m.String())

			m = money.Sub(money.Parse(20, tt.first), money.Parse(10, tt.second))
			assert.Equal(t, "10.00", m.String())

			m = money.Mul(money.Parse(5, tt.first), money.Parse(2, tt.second))
			assert.Equal(t, "10.00", m.String())

			m, err := money.Quo(money.Parse(100, tt.first), money.Parse(10, tt.second))
			require.NoError(t, err)
			assert.Equal(t, "10.00", m.String())
		})
	}
}

func TestCurrencySet(t *testing.T) {
	m := money.Parse(1.23456)
	assert.Equal(t, "1.2346", m.String())

	m = money.Parse(1.23456)
	m.CurrencySet("BIF")
	assert.Equal(t, "1", m.String())

	m = money.Parse(1.23456)
	m.CurrencySet("")
	assert.Equal(t, "1.2346", m.String())

	m = money.Parse(1.23456)
	m.CurrencySet("USD")
	assert.Equal(t, "1.23", m.String())

	m = money.Parse(1.23456)
	m.CurrencySet()
	assert.Equal(t, "1.2346", m.String())

	m = money.Parse(1.23456)
	m.CurrencySet("", "CLF")
	assert.Equal(t, "1.2346", m.String())
}

func TestCurrencyGet(t *testing.T) {
	m := money.Parse(1.23456)
	assert.Equal(t, "1.2346", m.String())

	m2 := money.Parse(0, "USD")
	m.CurrencySet(m2.CurrencyGet())
	assert.Equal(t, "1.23", m.String())
}

func TestMinMax(t *testing.T) {
	a := money.Parse(1, "USD")
	b := money.Parse(2, "BIF")

	ac := money.Parse(a, "BIF")
	bc := money.Parse(b, "USD")
	assert.Equal(t, a, money.Min(a, b))
	assert.Equal(t, bc, money.Max(a, b))

	assert.Equal(t, ac, money.Min(b, a))
	assert.Equal(t, b, money.Max(b, a))
}

func TestFinishCoverage(t *testing.T) {
	var m *money.Money
	m.CurrencySet("USD")
	assert.Nil(t, m)
}
