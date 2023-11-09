package datetime

import (
	"database/sql/driver"
	"gorm.io/datatypes"
	"time"
)

type Date datatypes.Date

func (date Date) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in RFC 3339 format, with sub-second precision added if present.
func (date Date) MarshalJSON() ([]byte, error) {
	return time.Time(date).MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (date *Date) UnmarshalJSON(data []byte) error {
	t := time.Time{}
	if err := t.UnmarshalJSON(data); err != nil {
		return err
	}

	*date = Date(t)
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in RFC 3339 format, with sub-second precision added if present.
func (date Date) MarshalText() ([]byte, error) {
	return time.Time(date).MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in RFC 3339 format.
func (date *Date) UnmarshalText(data []byte) error {
	t := time.Time{}
	if err := t.UnmarshalText(data); err != nil {
		return err
	}

	*date = Date(t)
	return nil
}
