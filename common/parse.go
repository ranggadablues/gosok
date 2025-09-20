package common

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func Payload(output any, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(output); err != nil {
		return err
	}
	return nil
}

func ToJSON(v interface{}) string {
	json, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(json)
}

// ToString converts any value to string (optimized with strconv)
func ParseString(v interface{}) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	case bson.ObjectID:
		return val.Hex()
	case fmt.Stringer: // types implementing String() string
		return val.String()
	default:
		// For slices, maps, structs → JSON
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Struct || rv.Kind() == reflect.Map || rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			if b, err := json.Marshal(v); err == nil {
				return string(b)
			}
		}
		// Fallback
		return fmt.Sprintf("%v", v)
	}
}

func KindDataType(v interface{}) {
	t := reflect.TypeOf(v) // returns the type
	k := t.Kind()          // returns the kind (basic category)
	fmt.Printf("Type: %s, Kind: %s\n", t.String(), k.String())
}

// ParseInt converts a given interface{} value into an integer.
// It handles integer types, string representations of integers, and
// float types by truncating the decimal part.
// If the input type is not supported or the string cannot be parsed, it returns 0.
func ParseInt(i interface{}) int {
	switch v := i.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float32:
		// Truncate the float to an integer.
		return int(v)
	case float64:
		// Truncate the float to an integer.
		return int(v)
	case string:
		// Attempt to parse the string into an integer.
		parsedInt, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return parsedInt
	default:
		return 0
	}
}

// ParseFloat64 converts any data type to float64 without rounding, keeping value as is
func ParseFloat64(v interface{}) float64 {
	if v == nil {
		return 0
	}

	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case string:
		if parsed, err := strconv.ParseFloat(val, 64); err == nil {
			return parsed
		}
		return 0
	case bool:
		if val {
			return 1
		}
		return 0
	default:
		// Try to convert to string first, then parse
		str := ParseString(v)
		if parsed, err := strconv.ParseFloat(str, 64); err == nil {
			return parsed
		}
		return 0
	}
}

// ParseFloat64RoundUp rounds up to the specified number of decimal places
func ParseFloat64RoundUp(v interface{}, decimalPlaces int) float64 {
	value := ParseFloat64(v)
	if decimalPlaces < 0 {
		decimalPlaces = 0
	}

	multiplier := math.Pow(10, float64(decimalPlaces))
	return math.Ceil(value*multiplier) / multiplier
}

// ParseFloat64RoundDown rounds down to the specified number of decimal places
func ParseFloat64RoundDown(v interface{}, decimalPlaces int) float64 {
	value := ParseFloat64(v)
	if decimalPlaces < 0 {
		decimalPlaces = 0
	}

	multiplier := math.Pow(10, float64(decimalPlaces))
	return math.Floor(value*multiplier) / multiplier
}

// ParseFloat64RoundAuto rounds to the nearest value with the specified number of decimal places
func ParseFloat64RoundAuto(v interface{}, decimalPlaces int) float64 {
	value := ParseFloat64(v)
	if decimalPlaces < 0 {
		decimalPlaces = 0
	}

	multiplier := math.Pow(10, float64(decimalPlaces))
	return math.Round(value*multiplier) / multiplier
}

// RoundingMode defines the type of rounding to apply
type RoundingMode int

const (
	RoundNone RoundingMode = iota // No rounding, keep value as is
	RoundUp                       // Round up (ceiling)
	RoundDown                     // Round down (floor)
	RoundAuto                     // Round to nearest (automatic)
)

// ParseFloat64Round provides flexible rounding based on the specified mode and decimal places
// If mode is RoundNone, decimalPlaces is ignored and the value is returned as is
func ParseFloat64Round(v interface{}, mode RoundingMode, decimalPlaces int) float64 {
	value := ParseFloat64(v)

	switch mode {
	case RoundNone:
		return value
	case RoundUp:
		return ParseFloat64RoundUp(value, decimalPlaces)
	case RoundDown:
		return ParseFloat64RoundDown(value, decimalPlaces)
	case RoundAuto:
		return ParseFloat64RoundAuto(value, decimalPlaces)
	default:
		return value
	}
}
