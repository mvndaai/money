package money

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/mvndaai/money/currencycodes"
)

// Money holds a an int64 representing a float to 4 decimal places of precision
// 4 decimals is needed for some currency codes like 'UYW'
// Floats have issues with rouunding so doing the math as an it makes it work better
// Note: currencyCode does not get saved in the dbm it is just for rounding numbers
type Money struct {
	value        int64 // saved in ten thousandths
	currencyCode string
}

/* Parsing */

func Parse[T number](t T, currencyCodes ...string) Money {
	if m, ok := any(t).(Money); ok {
		currencyCodes = append(currencyCodes, m.currencyCode)
	}
	f := shiftDecimals(t, savedDecimals)
	cc := currencycodes.FirstValidCode(currencyCodes...)
	f = math.RoundToEven(f)
	value := int64(f)
	if cc != "" {
		value = roundTenThousandths(value, currencycodes.CurrencyDecimals[cc])
	}
	return Money{value: value, currencyCode: cc}
}

func ParseCents[T number](t T, currencyCodes ...string) Money {
	f := shiftDecimals(t, -centsDecimals)
	return Parse(f, currencyCodes...)
}

// Deprecated: use Parse
func ParseBigFloat(bf big.Float, currencyCodes ...string) Money { return Parse(bf, currencyCodes...) }

func ParseString(ctx context.Context, s string, currencyCodes ...string) (Money, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Money{}, err
	}
	return Parse(f, currencyCodes...), nil
}

/* Manipulating
These are package level to show that we are not manipulating the internal number
*/

func Add[T, U number](t T, u ...U) Money {
	m := numberToMoney(t)
	cc := m.currencyCode
	total := m.value
	for _, v := range u {
		mv := numberToMoney(v)
		if cc == "" {
			cc = mv.currencyCode
		}
		total += mv.value
	}
	return Money{value: total, currencyCode: cc}
}

func Sub[T, U number](t T, u ...U) Money {
	if len(u) == 0 {
		return numberToMoney(t)
	}
	m := numberToMoney(t)
	cc := m.currencyCode
	total := m.value
	for _, v := range u {
		mv := numberToMoney(v)
		if cc == "" {
			cc = mv.currencyCode
		}
		total -= mv.value
	}
	return Money{value: total, currencyCode: cc}
}

func Mul[T, U number](t T, u ...U) Money {
	m := numberToMoney(t)
	if len(u) == 0 {
		return m
	}
	cc := m.currencyCode
	totalbf := big.NewFloat(m.Float64())
	for _, v := range u {
		mv := numberToMoney(v)
		if cc == "" {
			cc = mv.currencyCode
		}
		bf := big.NewFloat(mv.Float64())
		totalbf.Mul(totalbf, bf)
	}
	return Parse(*totalbf, cc)
}

func Quo[T, U number](t T, u ...U) (Money, error) {
	m := numberToMoney(t)
	if len(u) == 0 {
		return m, nil
	}
	cc := m.currencyCode

	totalbf := big.NewFloat(m.Float64())
	for _, v := range u {
		mv := numberToMoney(v)
		if cc == "" {
			cc = mv.currencyCode
		}
		if mv.value == 0 {
			return Money{}, errors.New("cannot divide by zero")
		}
		bf := big.NewFloat(mv.Float64())
		totalbf.Quo(totalbf, bf)
	}
	return Parse(totalbf, cc), nil
}

func Negate[T number](t T) Money {
	m := numberToMoney(t)
	m.value = -m.value
	return m
}

func Abs[T number](t T) Money {
	m := numberToMoney(t)
	if m.IsNegative() {
		return Negate(m)
	}
	return m
}

func Percentage[N, P number](num N, percentage P) Money {
	f := shiftDecimals(percentage, -2)
	return Mul(num, f)
}

func Max[T, U number](a T, b U) Money {
	am := numberToMoney(a)
	if am.Cmp(b) >= 0 {
		return am
	}
	bm := numberToMoney(b)
	if am.currencyCode != "" {
		bm.currencyCode = am.currencyCode
	}
	return bm
}

func Min[T, U number](a T, b U) Money {
	am := numberToMoney(a)
	if am.Cmp(b) <= 0 {
		return am
	}
	bm := numberToMoney(b)
	if am.currencyCode != "" {
		bm.currencyCode = am.currencyCode
	}
	return bm
}

func roundTenThousandths(tenThousandths int64, decimalPlaces int) int64 {
	if decimalPlaces < 0 {
		return tenThousandths
	}
	digits := savedDecimals - decimalPlaces
	if digits < 1 || digits > savedDecimals {
		return tenThousandths
	}
	i := shiftDecimals(tenThousandths, -digits)
	i = math.RoundToEven(i)
	return int64(shiftDecimals(i, digits))
}

func RoundToDecimals[T number](t T, decimalPlaces int) Money {
	m := numberToMoney(t)
	if m.currencyCode != "" && (currencycodes.CurrencyDecimals[m.currencyCode] <= decimalPlaces) {
		return m // already smallser than what we are rounding to
	}
	return Money{
		value:        roundTenThousandths(m.value, decimalPlaces),
		currencyCode: m.currencyCode,
	}
}

/* Comparing */

// Cmp compares returns:
//
//	-1 if m <  y
//	 0 if m == y
//	+1 if m >  y
func (m Money) Cmp[T number](y T) int {
	return cmp.Compare(m.value, numberToMoney(y).value)
}

func (m Money) IsEqual[T number](x T) bool {
	return m.value == numberToMoney(x).value
}

func (m Money) IsZero() bool {
	return m.value == 0
}

func (m Money) IsPositive() bool {
	return m.value > 0
}

func (m Money) IsNegative() bool {
	return m.value < 0
}

/* Exporting */

func (m Money) Float64() float64 {
	return shiftDecimals(m.value, -savedDecimals)
}

func (m Money) Cents() int {
	f := m.Float64()
	f = shiftDecimals(f, centsDecimals)
	f = math.RoundToEven(f)
	return int(f)
}

func (m Money) String() string {
	return m.StringDecimals(currencycodes.BestDecimal(m.currencyCode))
}

// StringDecimals rounds to the correct precision then returns the string limited to the number of decimals
func (m Money) StringDecimals(decimals int) string {
	v := RoundToDecimals(m, decimals).Float64()
	return fmt.Sprintf("%.*f", decimals, v)
}

/* Export with currency code */

func (m *Money) CurrencySet(currencyCodes ...string) {
	if m == nil {
		return
	}
	m.currencyCode = currencycodes.FirstValidCode(currencyCodes...)
	if m.currencyCode != "" {
		m.value = roundTenThousandths(m.value, currencycodes.CurrencyDecimals[m.currencyCode])
	}
}

func (m Money) CurrencyGet() string {
	return m.currencyCode
}

// StringCurrencyCode returns a string rounded and truncated to the decimals of the currency code indicated with fallbacks if incorrect
func (m Money) StringCurrencyCode(currencyCodes ...string) string {
	d := currencycodes.BestDecimal(currencyCodes...)
	return m.StringDecimals(d)
}

// Float64CurrencyCode rounds a decimal to the correct currency precision
func (m Money) Float64CurrencyCode(currencyCodes ...string) float64 {
	d := currencycodes.BestDecimal(currencyCodes...)
	return RoundToDecimals(m, d).Float64()
}
