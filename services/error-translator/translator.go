package et

import (
	"github.com/aeroideaservices/focus/services/errors"
	ut "github.com/go-playground/universal-translator"
)

type Translator struct {
	ut.Translator
	transTagFunc map[string]TranslationFunc
}

func New(translator ut.Translator) *Translator {
	return &Translator{
		Translator:   translator,
		transTagFunc: make(map[string]TranslationFunc),
	}
}

type TranslationFunc func(ut ut.Translator, msg string, params ...string) (string, error)

type Translation struct {
	Tag             string
	Translation     string
	Override        bool
	CustomRegisFunc func(ut ut.Translator) error
	CustomTransFunc func(ut ut.Translator, params ...string) string
}

func (translator Translator) AddTranslation(ts ...Translation) (err error) {
	for _, t := range ts {
		if t.CustomRegisFunc != nil {
			err = t.CustomRegisFunc(translator)
		} else {
			err = translator.Add(t.Tag, t.Translation, t.Override)
		}

		if err != nil {
			return
		}
	}

	return err
}

func (translator Translator) Translate(err error) (string, error) {
	e := err

	if fe, ok := e.(errors.FocusError); ok {
		if fe.Trans != nil {
			fn := translator.transTagFunc[fe.Trans.Msg]
			if fn == nil {
				fn = translationFn
			}
			return fn(translator, fe.Trans.Msg, fe.Trans.Params...)
		}
	}

	fn := translator.transTagFunc[err.Error()]
	if fn == nil {
		fn = translationFn
	}
	return fn(translator, err.Error())
}

func translationFn(ut ut.Translator, msg string, params ...string) (string, error) {
	t, err := ut.T(msg, params...)
	if err != nil {
		return "", err
	}

	return t, nil
}
