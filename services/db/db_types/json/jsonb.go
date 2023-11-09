package json

import (
	"database/sql/driver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type JSONB JSON

func (j JSONB) MarshalText() (text []byte, err error) {
	return j.MarshalJSON()
}

func (j *JSONB) UnmarshalText(text []byte) error {
	return j.UnmarshalJSON(text)
}

func (j *JSONB) Scan(src any) error {
	return (*JSON)(j).Scan(src)
}

func (j JSONB) Value() (driver.Value, error) {
	return (JSON)(j).Value()
}

func (j JSONB) MarshalJSON() ([]byte, error) {
	return (JSON)(j).MarshalJSON()
}

func (j *JSONB) UnmarshalJSON(b []byte) error {
	return (*JSON)(j).UnmarshalJSON(b)
}

func (JSONB) GormDataType() string {
	return "jsonb"
}

func (JSONB) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "jsonb"
	default:
		return "json"
	}
}
