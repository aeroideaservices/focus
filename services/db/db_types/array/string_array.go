package array

import (
	"database/sql/driver"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type StringArray pq.StringArray

func (a *StringArray) Scan(src interface{}) error {
	pqa := pq.StringArray(*a)
	err := pqa.Scan(src)
	*a = StringArray(pqa)
	return err
}

func (a StringArray) Value() (driver.Value, error) {
	return pq.StringArray(a).Value()
}

func (a StringArray) GormDataType() string {
	return "text[]"
}

func (a StringArray) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "text[]"
	default:
		return "text"
	}
}
