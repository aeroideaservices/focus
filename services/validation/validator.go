package validation

import (
	"context"
	"golang.org/x/text/language"
	"reflect"
	"strings"

	"github.com/aeroideaservices/focus/services/errors"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	ruTranslations "github.com/go-playground/validator/v10/translations/ru"
	"github.com/google/uuid"

	enLocal "github.com/aeroideaservices/focus/services/validation/translations/en"
	ruLocal "github.com/aeroideaservices/focus/services/validation/translations/ru"
)

const (
	sluggableTag        = "sluggable"
	slashedSluggableTag = "slashedSluggable"
	notBlankTag         = "notBlank"
	patternTag          = "pattern"
)

type Validator struct {
	validate    *validator.Validate
	translation ut.Translator
	translator  *ut.UniversalTranslator
}

func NewValidator(translator *ut.UniversalTranslator) *Validator {
	v := validator.New()
	// register function to get tag name from json tags.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// must copy alias validators for separate validations to be used in each validator instance
	addPatternAliases(patternTag, bakedInAliases)
	for k, val := range bakedInAliases {
		v.RegisterAlias(k, val)
	}

	// must copy validators for separate validations to be used in each instance
	for k, val := range bakedInValidators {
		switch k {
		// these require that even if the value is nil that the validation should run, omitempty still overrides this behaviour
		case sluggableTag, slashedSluggableTag, patternTag, notBlankTag:
			_ = v.RegisterValidation(k, val, true, true)
		default:
			// no need to error check here, baked in will always be valid
			_ = v.RegisterValidation(k, val, true, false)
		}
	}

	// todo optimize
	ruTranslator, _ := translator.GetTranslator(language.Russian.String())
	enTranslator, _ := translator.GetTranslator(language.English.String())
	_ = ruTranslations.RegisterDefaultTranslations(v, ruTranslator)
	_ = enTranslations.RegisterDefaultTranslations(v, enTranslator)
	_ = ruLocal.RegisterTranslations(v, ruTranslator)
	_ = enLocal.RegisterTranslations(v, enTranslator)

	v.RegisterCustomTypeFunc(validateUUID, uuid.UUID{})

	return &Validator{
		validate:   v,
		translator: translator,
	}
}

func (v Validator) Validate(ctx context.Context, value any) error {
	err := v.validate.StructCtx(ctx, value)
	if err == nil {
		return nil
	}

	return errors.BadRequest.Wrapf(err, "validation error: %s")
}

func (v Validator) ValidatePartial(ctx context.Context, value any) error {
	err := v.validate.StructPartialCtx(ctx, value)
	if err == nil {
		return nil
	}

	return errors.BadRequest.Wrapf(err, "validation error: %s")
}

func (v Validator) validationErrorsToString(validationErrors validator.ValidationErrors) string {
	errs := make([]string, len(validationErrors))
	translator, ok := v.translator.GetTranslator("ru")
	if !ok {
		translator = v.translator.GetFallback()
	}
	for i, err := range validationErrors {
		errs[i] = err.Translate(translator)
	}

	return strings.Join(errs, "; ")
}
