package money

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

const (
	savedDecimals = 4
	centsDecimals = 2
)

type builtInNumber interface {
	~float32 | ~float64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type typedNumber interface {
	big.Int | *big.Int | big.Float | *big.Float | Money | *Money
}

type number interface {
	builtInNumber | typedNumber
}

type floatFunc interface {
	Float64() float64
}
type floatFuncWithOther[T any] interface {
	Float64() (float64, T)
}

//nolint:cyclop // exhaustive mapping table, complexity is inherent
func numberToFloat[T number](t T) float64 {
	switch v := any(t).(type) {
	case floatFunc:
		return v.Float64()
	case floatFuncWithOther[big.Accuracy]:
		f, _ := v.Float64()
		return f
	case big.Float:
		f, _ := v.Float64()
		return f
	case big.Int:
		f, _ := v.Float64()
		return f
	case float32:
		return float64(v)
	case float64:
		return v
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case uintptr:
		return float64(v)
	}
	// This should not be possible unless someone changes the 'number' interface
	panic(fmt.Sprintf("could not convert (%T) to float64", t))
}

func numberToMoney[T number](t T) (m Money) {
	if v, ok := any(t).(Money); ok {
		return v
	}
	if v, ok := any(t).(*Money); ok {
		return *v
	}
	return Parse(numberToFloat(t))
}

// shiftDecimals will takes a float and moves the decimal point +/- the number of digits
// This converts to string and moves rather than */÷ by 10s to prevent rounding of floats
func shiftDecimals[T number](t T, digits int) float64 {
	f := numberToFloat(t)
	var neg string
	if f < 0 {
		neg = "-"
		f = math.Abs(f)
	}

	var s string
	if digits < 0 {
		s = fmt.Sprintf("%s%0*d%f", neg, -digits, 0, f)
	} else {
		s = fmt.Sprintf("%s%f", neg, f)
	}

	dLoc := strings.Index(s, ".") + digits
	s = strings.ReplaceAll(s, ".", "")
	s = s[0:dLoc] + "." + s[dLoc:]
	f, _ = strconv.ParseFloat(s, 64)
	return f
}

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m *Money) UnmarshalJSON(b []byte) error {
	ctx := context.Background()
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	parsed, err := ParseString(ctx, n.String())
	if err != nil {
		return err
	}
	m.value = parsed.value
	return nil
}

// Scan implements the sql.Scanner interface
func (m *Money) Scan(src any) error {
	if src == nil {
		*m = Parse(0)
		return nil
	}

	switch v := src.(type) {
	case []byte:
		pm, err := ParseString(context.Background(), string(v))
		*m = pm
		return err
	case string:
		pm, err := ParseString(context.Background(), v)
		*m = pm
		return err
	case float64:
		*m = Parse(v)
		return nil
	default:
		return fmt.Errorf("failed to scan type '%T' as Money", src)
	}
}

// Value is used for sql saves
func (m *Money) Value() (driver.Value, error) {
	return m.Float64(), nil
}
