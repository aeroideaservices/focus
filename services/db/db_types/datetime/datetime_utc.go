package datetime

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

const RFC3339UTC = "2006-01-02T15:04:05Z"

type DateTimeUTC time.Time

func (dt *DateTimeUTC) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*dt = DateTimeUTC(nullTime.Time)
	return
}

func (dt DateTimeUTC) Value() (driver.Value, error) {
	y, m, d := time.Time(dt).Date()
	h, i, s := time.Time(dt).Clock()
	return time.Date(y, m, d, h, i, s, 0, time.UTC), nil
}

// GormDataType gorm common data type
func (dt DateTimeUTC) GormDataType() string {
	return "timestamp"
}

func (dt DateTimeUTC) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "timestamp without time zone"
	default:
		return "timestamp"
	}
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in RFC 3339 format, with sub-second precision added if present.
func (dt DateTimeUTC) MarshalJSON() ([]byte, error) {
	if y := time.Time(dt).Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(RFC3339UTC)+2)
	b = append(b, '"')
	b = time.Time(dt).AppendFormat(b, RFC3339UTC)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (dt *DateTimeUTC) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	t, err := time.ParseInLocation(`"`+RFC3339UTC+`"`, string(data), time.UTC)
	*dt = DateTimeUTC(t)
	return err
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in RFC 3339 format, with sub-second precision added if present.
func (dt DateTimeUTC) MarshalText() ([]byte, error) {
	if y := time.Time(dt).Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalText: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(RFC3339UTC))
	return time.Time(dt).AppendFormat(b, RFC3339UTC), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in RFC 3339 format.
func (dt *DateTimeUTC) UnmarshalText(data []byte) error {
	// Fractional seconds are handled implicitly by Parse.
	t, err := time.ParseInLocation(RFC3339UTC, string(data), time.UTC)
	*dt = DateTimeUTC(t)
	return err
}
