package common

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// ToString converts any value to string (optimized with strconv)
func ToString(v interface{}) string {
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

func CaptureTypeReflect(v interface{}) {
	t := reflect.TypeOf(v) // returns the type
	k := t.Kind()          // returns the kind (basic category)
	fmt.Printf("Type: %s, Kind: %s\n", t.String(), k.String())
}
