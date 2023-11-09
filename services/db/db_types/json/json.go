package json

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type JSON struct {
	v any
}

func (j *JSON) Scan(src any) error {
	switch value := src.(type) {
	case []byte:
		return json.Unmarshal(value, &j.v)
	case string:
		return json.Unmarshal([]byte(value), &j.v)
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (j JSON) Value() (driver.Value, error) {
	return json.Marshal(j.v)
}

func (j JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.v)
}

func (j *JSON) UnmarshalJSON(b []byte) error {
	if b == nil {
		*j = JSON{v: nil}
		return nil
	}

	err := json.Unmarshal(b, &j.v)
	return err
}

func (j JSON) MarshalText() ([]byte, error) {
	return j.MarshalJSON()
}

func (j *JSON) UnmarshalText(b []byte) error {
	return j.UnmarshalJSON(b)
}

func (JSON) GormDataType() string {
	return "json"
}

func (JSON) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	default:
		return "json"
	}
}
