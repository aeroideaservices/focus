package focus

import (
	"github.com/aeroideaservices/focus/services/validation"
	"github.com/google/uuid"
	"reflect"
	"strings"
)

const (
	arrayType   = "array"
	booleanType = "boolean"
	dateType    = "date"
	numberType  = "number"
	objectType  = "object"
	stringType  = "string"

	patternTag = "pattern"

	errMessageRule = "errMessages"
)

type rules map[string]rule
type rule map[string]any

func GetValidationRules(model Model, filter func(field *Field) bool) rules {
	res := make(rules)

	for _, field := range model.Fields {
		if filter != nil && !filter(field) {
			continue
		}
		structField, _ := model.t.FieldByName(field.name)
		if tag, _ := structField.Tag.Lookup("validate"); tag != "-" {
			res[field.Code] = readFieldValidation(*field, tag)
		}
	}

	return res
}

func readFieldValidation(field Field, tag string) rule {
	if tag == "" {
		tag = "omitempty"
	}
	tagParts := strings.Split(tag, ",")
	res := readTagParts(tagParts)

	res["title"] = field.Title
	typ := typeOfField(field)
	res["type"] = typ
	if _, ok := res["items"]; ok {
		if typ == arrayType {
			var typ string
			switch {
			case field.Association != nil:
				typ = typeOfField(*field.Association.Model.PrimaryKey)
			case field.IsMedia:
				typ = stringType
			default:
				typ = typeOf(field.RawType())
			}
			(res["items"]).(rule)["type"] = typ
		} else {
			delete(res, "items")
		}
	}

	return res
}

// readTagParts
func readTagParts(tagParts []string) rule {
	res := make(rule)
	for i, part := range tagParts {
		keyVal := strings.SplitN(part, "=", 2)
		if keyVal[0] == "dive" {
			res["items"] = readTagParts(tagParts[i+1:])
			return res
		}
		if len(keyVal) == 2 {
			if keyVal[0] == patternTag {
				p, ok := validation.Patterns[keyVal[1]]
				if ok {
					res[keyVal[0]] = p.RegexpString
					res[errMessageRule] = map[string]string{patternTag: p.ErrMessage}
				}
			} else {
				res[keyVal[0]] = keyVal[1]
			}
		} else {
			res[keyVal[0]] = true
		}
	}

	return res
}

func typeOfField(field Field) string {
	if field.Multiple {
		return arrayType
	}
	if field.IsTime {
		return dateType
	}
	if field.IsMedia {
		return stringType // Id медиа
	}
	if field.Association != nil {
		return typeOfField(*field.Association.Model.PrimaryKey)
	}

	t := field.t
	if t.Kind() == reflect.Ptr {
		t = field.t.Elem()
	}
	return typeOf(t)
}

func typeOf(typ reflect.Type) string {
	if typ == reflect.TypeOf(uuid.UUID{}) {
		return stringType
	}

	switch typ.Kind() {
	case reflect.String:
		return stringType
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float64, reflect.Float32:
		return numberType
	case reflect.Map:
		return objectType
	case reflect.Bool:
		return booleanType
	default:
		return ""
	}
}
