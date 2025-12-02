package common

const (
	DefaultConnectionName = "default"
)

// Time format constants for ParseTime function
const (
	// Standard formats
	TimeFormatRFC3339     = "2006-01-02T15:04:05Z07:00" // ISO 8601 / RFC3339
	TimeFormatRFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	TimeFormatRFC822      = "02 Jan 06 15:04 MST"
	TimeFormatRFC822Z     = "02 Jan 06 15:04 -0700"
	TimeFormatRFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	TimeFormatRFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700"

	// Common date-time formats
	TimeFormatDateTime        = "2006-01-02 15:04:05"            // YYYY-MM-DD HH:MM:SS
	TimeFormatDateTimeWithTZ  = "2006-01-02 15:04:05 -0700"      // YYYY-MM-DD HH:MM:SS -0700
	TimeFormatDateTimeMilli   = "2006-01-02 15:04:05.000"        // YYYY-MM-DD HH:MM:SS.mmm
	TimeFormatDateTimeMicro   = "2006-01-02 15:04:05.000000"     // YYYY-MM-DD HH:MM:SS.μμμμμμ
	TimeFormatDateTimeNano    = "2006-01-02 15:04:05.000000000"  // YYYY-MM-DD HH:MM:SS.nnnnnnnnn
	TimeFormatDateTimeT       = "2006-01-02T15:04:05"            // YYYY-MM-DDTHH:MM:SS
	TimeFormatDateTimeTMilli  = "2006-01-02T15:04:05.000"        // YYYY-MM-DDTHH:MM:SS.mmm
	TimeFormatDateTimeTZ      = "2006-01-02T15:04:05Z"           // YYYY-MM-DDTHH:MM:SSZ
	TimeFormatDateTimeTMilliZ = "2006-01-02T15:04:05.000Z"       // YYYY-MM-DDTHH:MM:SS.mmmZ
	TimeFormatDateTimeTOffset = "2006-01-02T15:04:05-07:00"      // YYYY-MM-DDTHH:MM:SS-07:00
	TimeFormatDateTimeTMicroZ = "2006-01-02T15:04:05.000000Z"    // YYYY-MM-DDTHH:MM:SS.μμμμμμZ
	TimeFormatDateTimeTNanoZ  = "2006-01-02T15:04:05.000000000Z" // YYYY-MM-DDTHH:MM:SS.nnnnnnnnnZ

	// Date only formats
	TimeFormatDate           = "2006-01-02"      // YYYY-MM-DD
	TimeFormatDateSlash      = "2006/01/02"      // YYYY/MM/DD
	TimeFormatDateDot        = "2006.01.02"      // YYYY.MM.DD
	TimeFormatDateUS         = "01/02/2006"      // MM/DD/YYYY
	TimeFormatDateUSWithDash = "01-02-2006"      // MM-DD-YYYY
	TimeFormatDateEU         = "02/01/2006"      // DD/MM/YYYY
	TimeFormatDateEUWithDash = "02-01-2006"      // DD-MM-YYYY
	TimeFormatDateCompact    = "20060102"        // YYYYMMDD
	TimeFormatDateReadable   = "02 Jan 2006"     // DD Mon YYYY
	TimeFormatDateLong       = "January 2, 2006" // Month DD, YYYY

	// Time only formats
	TimeFormatTime       = "15:04:05"        // HH:MM:SS
	TimeFormatTimeMilli  = "15:04:05.000"    // HH:MM:SS.mmm
	TimeFormatTimeMicro  = "15:04:05.000000" // HH:MM:SS.μμμμμμ
	TimeFormatTimeShort  = "15:04"           // HH:MM
	TimeFormatTime12Hour = "03:04:05 PM"     // 12-hour format with AM/PM
	TimeFormatTime12     = "03:04 PM"        // 12-hour short with AM/PM

	// Unix timestamp (string representation)
	TimeFormatUnix      = "unix"       // Unix timestamp in seconds
	TimeFormatUnixMilli = "unix-milli" // Unix timestamp in milliseconds
	TimeFormatUnixMicro = "unix-micro" // Unix timestamp in microseconds
	TimeFormatUnixNano  = "unix-nano"  // Unix timestamp in nanoseconds
)
