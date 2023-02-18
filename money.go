package money

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Money holds a numbers with 4 decimal points of precision
type Money struct {
	value        int64
	currencyCode string
}

const (
	defaultDecimals = 4
	multiplier      = 10000
)

func (m *Money) SetCurrencyCode(currencyCode string) {
	m.currencyCode = strings.ToUpper(currencyCode)
}

func (m Money) Float64() float64 {
	return float64(m.value) / multiplier
}

// CurrencyRound rounds a decimal to the correct currency precision
func (m Money) CurrencyRound() float64 {
	cd := CurrencyDecimals(m.currencyCode)
	return m.roundToDecimals(cd).Float64()
}

func ParseInt(i int) Money {
	return ParseInt64(int64(i))
}
func ParseInt64(i int64) Money {
	return Money{value: i * multiplier}
}

func ParseFloat64(f float64) Money {
	return Money{value: int64(math.Round(f * multiplier))}
}

func ParseString(ctx context.Context, s string) (Money, error) {
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

	v := math.Round(float64(m.value)/rounding) * rounding
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
	rounded := int64(math.Round(totalf * multiplier))
	return Money{value: rounded}
}

func Quo(m ...Money) Money {
	if len(m) == 0 {
		return Money{}
	}
	totalbf := big.NewFloat(float64(m[0].value) / multiplier)
	for _, v := range m[1:] {
		totalbf.Quo(totalbf, big.NewFloat(float64(v.value)/multiplier))
	}
	totalf, _ := totalbf.Float64()
	rounded := int64(math.Round(totalf * multiplier))
	return Money{value: rounded}
}

func (m Money) String() string {
	d := defaultDecimals
	if m.currencyCode != "" {
		d = CurrencyDecimals(m.currencyCode)
	}

	v := float64(m.roundToDecimals(d).value) / multiplier
	format := "%." + strconv.Itoa(d) + "f"
	return fmt.Sprintf(format, v)
}

func (m Money) Equal(x Money) bool {
	return m.value == x.value &&
		m.currencyCode == x.currencyCode
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
