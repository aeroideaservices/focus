package array

import (
	"database/sql/driver"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type UuidArray []uuid.UUID

func (a *UuidArray) Scan(src interface{}) error {
	if src == nil {
		*a = nil
		return nil
	}

	if src == "{}" {
		*a = make([]uuid.UUID, 0)
		return nil
	}

	sa := pq.StringArray{}
	err := sa.Scan(src)
	if err != nil {
		return err
	}

	*a = make([]uuid.UUID, len(sa))
	for i := range sa {
		(*a)[i], err = uuid.Parse(sa[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (a UuidArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	sa := make(pq.StringArray, len(a))
	for i := range a {
		sa[i] = a[i].String()
	}

	return sa.Value()
}

func (a UuidArray) GormDataType() string {
	return "uuid[]"
}

func (a UuidArray) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "uuid[]"
	default:
		return "text"
	}
}
