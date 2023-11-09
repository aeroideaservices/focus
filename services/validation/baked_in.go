package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"reflect"
)

var (
	// bakedInAliases is a default mapping of a single validation tag that
	// defines a common or complex set of validation(s) to simplify
	// adding validation to structs.
	bakedInAliases = map[string]string{
		"phone": "e164",
	}
	// bakedInValidators is the default map of ValidationFunc
	// you can add, remove or even replace items to suite your needs,
	// or even disregard and use your own map if so desired.
	bakedInValidators = map[string]validator.Func{
		"notBlank":         validators.NotBlank,
		"sluggable":        isSluggable,
		"slashedSluggable": isSlashedSluggable,
		"pattern":          matchPattern,
	}
)

func isSluggable(fl validator.FieldLevel) bool {
	return sluggableRegex.MatchString(fl.Field().String())
}

func isSlashedSluggable(fl validator.FieldLevel) bool {
	return slashedSluggableRegex.MatchString(fl.Field().String())
}

func matchPattern(fl validator.FieldLevel) bool {
	param := fl.Param()
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		p, ok := Patterns[param]

		if !ok {
			panic(fmt.Sprintf("Bad alias for pattern %s", param))
		}

		return p.GetRegexp().MatchString(field.String())
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func addPatternAliases(tag string, m map[string]string) {
	const comma = "|"
	var aliases string

	for _, p := range Patterns {
		aliases = p.Alias + comma
	}
	aliases = aliases[:len(aliases)-1]

	m[tag] = aliases
}
