package common

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

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

// ParseBool converts any value to boolean
// Recognizes: true/false, 1/0, "true"/"false", "yes"/"no", "on"/"off", "1"/"0", "t"/"f", "y"/"n"
func ParseBool(v interface{}) bool {
	if v == nil {
		return false
	}

	switch val := v.(type) {
	case bool:
		return val
	case int, int8, int16, int32, int64:
		// Any integer type - convert to int64 for comparison
		intVal := reflect.ValueOf(val).Int()
		return intVal != 0
	case uint, uint8, uint16, uint32, uint64:
		// Any unsigned integer type - convert to uint64 for comparison
		uintVal := reflect.ValueOf(val).Uint()
		return uintVal != 0
	case float32, float64:
		// Any float type - convert to float64 for comparison
		floatVal := reflect.ValueOf(val).Float()
		return floatVal != 0
	case string:
		// Normalize string: trim spaces and convert to lowercase
		str := strings.TrimSpace(strings.ToLower(val))
		// True values
		switch str {
		case "true", "t", "yes", "y", "on", "1":
			return true
		}
		// False values (explicit check for clarity)
		switch str {
		case "false", "f", "no", "n", "off", "0", "":
			return false
		}
		// Try parsing as number
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			return num != 0
		}
		return false
	default:
		// For other types, check if they're non-zero using reflection
		rv := reflect.ValueOf(v)
		if !rv.IsValid() {
			return false
		}
		// Check if it's a zero value
		return !rv.IsZero()
	}
}

// ParseTime converts any value to time.Time with flexible format support
// Supports:
// - time.Time: returns as is
// - string: tries multiple formats (use format constants from constants.go)
// - int/int64: treats as Unix timestamp (seconds)
// - float64: treats as Unix timestamp with fractional seconds
// Returns zero time (time.Time{}) if parsing fails
// ParseBool examples
// ParseBool(true)              // true
// ParseBool(1)                 // true
// ParseBool("yes")             // true
// ParseBool("on")              // true
// ParseBool("Y")               // true (case-insensitive)
// ParseBool(0)                 // false
// ParseBool("no")              // false
// ParseBool(3.14)
func ParseTime(v interface{}, formats ...string) time.Time {
	if v == nil {
		return time.Time{}
	}

	switch val := v.(type) {
	case time.Time:
		return val
	case *time.Time:
		if val != nil {
			return *val
		}
		return time.Time{}
	case string:
		return parseTimeFromString(val, formats...)
	case int:
		return time.Unix(int64(val), 0)
	case int32:
		return time.Unix(int64(val), 0)
	case int64:
		return time.Unix(val, 0)
	case uint:
		return time.Unix(int64(val), 0)
	case uint32:
		return time.Unix(int64(val), 0)
	case uint64:
		return time.Unix(int64(val), 0)
	case float32:
		// Treat as Unix timestamp with fractional seconds
		sec := int64(val)
		nsec := int64((val - float32(sec)) * 1e9)
		return time.Unix(sec, nsec)
	case float64:
		// Treat as Unix timestamp with fractional seconds
		sec := int64(val)
		nsec := int64((val - float64(sec)) * 1e9)
		return time.Unix(sec, nsec)
	default:
		// Try converting to string first
		str := ParseString(v)
		if str != "" {
			return parseTimeFromString(str, formats...)
		}
		return time.Time{}
	}
}

// parseTimeFromString attempts to parse a time string using provided formats or common formats
// ParseTime examples - auto-detect format
// ParseTime("2024-10-14T15:04:05Z")           // RFC3339
// ParseTime("2024-10-14 15:04:05")            // DateTime
// ParseTime("01/02/2006")                     // US date format
// ParseTime(1697297045)                       // Unix timestamp (seconds)
// ParseTime(1697297045000)                    // Unix timestamp (milliseconds) - auto-detected
// ParseTime(1697297045.123)                   // Unix timestamp with fractional seconds

// // ParseTime with specific format
// ParseTime("2024-10-14", TimeFormatDate)
// ParseTime("14/10/2024", TimeFormatDateEU)
// ParseTime("1697297045", TimeFormatUnix)
// ParseTime("1697297045000", TimeFormatUnixMilli)

// // Multiple formats priority (tries in order)
// common.ParseTime("10-14-2024", common.TimeFormatDateUSWithDash, common.TimeFormatDateEUWithDash)
func parseTimeFromString(str string, formats ...string) time.Time {
	str = strings.TrimSpace(str)
	if str == "" {
		return time.Time{}
	}

	// If custom formats are provided, try them first
	if len(formats) > 0 {
		return parseCustomFormats(str, formats...)
	}

	// Default: try common formats in order of likelihood
	commonFormats := []string{
		TimeFormatRFC3339,
		TimeFormatRFC3339Nano,
		TimeFormatDateTime,
		TimeFormatDateTimeT,
		TimeFormatDateTimeTZ,
		TimeFormatDateTimeTMilliZ,
		TimeFormatDateTimeTMicroZ,
		TimeFormatDateTimeTNanoZ,
		TimeFormatDateTimeTOffset,
		TimeFormatDateTimeMilli,
		TimeFormatDateTimeMicro,
		TimeFormatDateTimeNano,
		TimeFormatDate,
		TimeFormatDateSlash,
		TimeFormatDateUS,
		TimeFormatDateEU,
		TimeFormatDateCompact,
		TimeFormatDateReadable,
		TimeFormatDateLong,
		TimeFormatRFC1123,
		TimeFormatRFC1123Z,
		TimeFormatRFC822,
		TimeFormatRFC822Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
	}

	for _, format := range commonFormats {
		if t, err := time.Parse(format, str); err == nil {
			return t
		}
	}

	// Try parsing as Unix timestamp (string)
	if t := parseUnixTimestamp(str, ""); !t.IsZero() {
		return t
	}

	return time.Time{}
}

func parseCustomFormats(str string, formats ...string) time.Time {
	for _, format := range formats {
		// Handle special unix timestamp formats
		if strings.HasPrefix(format, "unix") {
			if t := parseUnixTimestamp(str, format); !t.IsZero() {
				return t
			}
			continue
		}
		// Try parsing with the provided format
		if t, err := time.Parse(format, str); err == nil {
			return t
		}
	}
	// If custom formats are provided but none worked, return zero time
	return time.Time{}
}

// parseUnixTimestamp attempts to parse a Unix timestamp from string
func parseUnixTimestamp(str string, format string) time.Time {
	// Parse as number
	if num, err := strconv.ParseInt(str, 10, 64); err == nil {
		// Determine the scale based on format or number size
		switch format {
		case TimeFormatUnix:
			return time.Unix(num, 0)
		case TimeFormatUnixMilli:
			return time.Unix(num/1000, (num%1000)*1e6)
		case TimeFormatUnixMicro:
			return time.Unix(num/1e6, (num%1e6)*1e3)
		case TimeFormatUnixNano:
			return time.Unix(0, num)
		default:
			// Auto-detect based on magnitude
			if num > 1e12 { // Likely milliseconds or higher
				return uniqueDefaultParseTime(num)
			}
			// Likely seconds
			return time.Unix(num, 0)
		}
	}

	// Try parsing as float for fractional seconds
	if num, err := strconv.ParseFloat(str, 64); err == nil {
		sec := int64(num)
		nsec := int64((num - float64(sec)) * 1e9)
		return time.Unix(sec, nsec)
	}

	return time.Time{}
}

func uniqueDefaultParseTime(num int64) time.Time {
	if num > 1e15 { // Likely microseconds or nanoseconds
		if num > 1e18 { // Likely nanoseconds
			return time.Unix(0, num)
		}
		// Likely microseconds
		return time.Unix(num/1e6, (num%1e6)*1e3)
	}
	// Likely milliseconds
	return time.Unix(num/1000, (num%1000)*1e6)
}

func ParseObjectID(v interface{}) bson.ObjectID {
	if v == nil {
		return bson.ObjectID{}
	}

	objectID, err := bson.ObjectIDFromHex(ParseString(v))
	if err != nil {
		return bson.ObjectID{}
	}
	return objectID
}

func MapToStruct(in interface{}, out interface{}) error {
	// Convert map to JSON
	bytes, err := json.Marshal(in)
	if err != nil {
		return err
	}

	// Unmarshal into struct
	return json.Unmarshal(bytes, out)
}

func StructToMap(in interface{}, out map[string]interface{}) error {
	bytes, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &out)
}
