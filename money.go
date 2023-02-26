package money

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

// Money holds a numbers to 4 decimal points of precision
type Money struct {
	value int64
}

const (
	defaultDecimals = 4
	multiplier      = 10000
)

func (m Money) Float64() float64 {
	return float64(m.value) / multiplier
}

func ParseInt(i int) Money {
	return ParseInt64(int64(i))
}
func ParseInt64(i int64) Money {
	return Money{value: i * multiplier}
}

func ParseFloat64(f float64) Money {
	return Money{value: int64(math.RoundToEven(f * multiplier))}
}

func ParseString(s string) (Money, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Money{}, err
	}
	return ParseFloat64(f), nil
}

func (m Money) roundToDecimals(decimalPlaces int) Money {
	rounding := math.Pow10(defaultDecimals - decimalPlaces)
	if rounding <= 1 {
		return m
	}

	v := math.RoundToEven(float64(m.value)/rounding) * rounding
	return Money{value: int64(v)}
}

func Add(m ...Money) Money {
	var total int64
	for _, v := range m {
		total += v.value
	}
	return Money{value: total}
}

func Sub(m ...Money) Money {
	if len(m) == 0 {
		return Money{}
	}
	total := m[0].value
	for _, v := range m[1:] {
		total -= v.value
	}
	return Money{value: total}
}

func Mul(m ...Money) Money {
	if len(m) == 0 {
		return Money{}
	}
	totalbf := big.NewFloat(1)
	for _, v := range m {
		totalbf.Mul(totalbf, big.NewFloat(float64(v.value)/multiplier))
	}
	totalf, _ := totalbf.Float64()
	rounded := int64(math.RoundToEven(totalf * multiplier))
	return Money{value: rounded}
}

func Quo(m ...Money) (Money, error) {
	if len(m) == 0 {
		return Money{}, fmt.Errorf("no params")
	}
	totalbf := big.NewFloat(float64(m[0].value) / multiplier)
	for _, v := range m[1:] {
		if v.value == 0 {
			return Money{}, fmt.Errorf("cannot divide by zero")
		}
		totalbf.Quo(totalbf, big.NewFloat(float64(v.value)/multiplier))
	}
	totalf, _ := totalbf.Float64()
	rounded := int64(math.RoundToEven(totalf * multiplier))
	return Money{value: rounded}, nil
}

func (m Money) Equal(x Money) bool {
	return m.value == x.value
}

func (m Money) String() string {
	s, _ := m.CurrencyString(CurrencyCodeCLF)
	return s
}

// CurrencyString returns a string rounded and formatted for a currency code
func (m Money) CurrencyString(currencyCode string) (string, error) {
	d, err := CurrencyDecimals(currencyCode)
	if err != nil {
		return "", err
	}
	v := float64(m.roundToDecimals(d).value) / multiplier
	format := "%." + strconv.Itoa(d) + "f"
	return fmt.Sprintf(format, v), nil
}

// CurrencyFloat64 rounds a decimal to the correct currency precision
func (m Money) CurrencyFloat64(currencyCode string) (float64, error) {
	cd, err := CurrencyDecimals(currencyCode)
	if err != nil {
		return 0, err
	}
	return m.roundToDecimals(cd).Float64(), nil
}

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m *Money) UnmarshalJSON(b []byte) error {
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	f, err := n.Float64()
	if err != nil {
		return err // unreachable
	}
	m.value = ParseFloat64(f).value
	return nil
}

// Scan implements the sql.Scanner interface
func (m *Money) Scan(src interface{}) error {
	if src == nil {
		*m = ParseInt(0)
		return nil
	}

	switch v := src.(type) {
	case []byte:
		pm, err := ParseString(string(v))
		*m = pm
		return err
	default:
		return fmt.Errorf("failed to scan type '%T' as Money", src)
	}
}
