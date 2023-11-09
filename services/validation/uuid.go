package validation

import (
	"github.com/google/uuid"
	"reflect"
)

func validateUUID(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(uuid.UUID); ok {
		return valuer.String()
	}

	return nil
}
