package validation

import "regexp"

const (
	sluggableAlias        = "sluggable"
	slashedSluggableAlias = "slashedSluggable"
	urlAlias              = "url"
	phoneAlias            = "phone"
	emailAlias            = "email"

	sluggableRegexString        = "^[a-z0-9]+(?:-[a-z0-9]+)*$"
	slashedSluggableRegexString = "^[a-z0-9]+(?:[-/][a-z0-9]+)*$"
	urlRegexString              = "^(\\/)|((?:http(s)?:\\/\\/))"
	phoneRegexString            = "^((8|\\+7)[\\- ]?)?(\\(?\\d{3}\\)?[\\- ]?)?[\\d\\- ]{7,10}$"
	emailRegexString            = "^[a-zA-Z0-9_!#$%&'*+/=?`{|}~^.-]+@[a-zA-Z0-9.-]+$"

	sluggableErrMessage        = "Использовать латиницу в нижнем регистре, цифры и дефис. Дефис не должен быть первым или последним символом"
	slashedSluggableErrMessage = "Использовать латиницу в нижнем регистре, цифры, дефис и слэш"
	urlErrMessage              = "Ссылка невалидная"
	phoneErrMessage            = "Неверный формат телефона"
	emailErrMessage            = "Неверный формат e-mail"
)

var (
	sluggableRegex        = regexp.MustCompile(sluggableRegexString)
	slashedSluggableRegex = regexp.MustCompile(slashedSluggableRegexString)
)

var Patterns = map[string]Pattern{
	sluggableAlias: {
		Alias:        sluggableAlias,
		RegexpString: sluggableRegexString,
		ErrMessage:   sluggableErrMessage},
	slashedSluggableAlias: {
		Alias:        slashedSluggableAlias,
		RegexpString: slashedSluggableRegexString,
		ErrMessage:   slashedSluggableErrMessage},
	urlAlias: {
		Alias:        urlAlias,
		RegexpString: urlRegexString,
		ErrMessage:   urlErrMessage,
	},
	phoneAlias: {
		Alias:        phoneAlias,
		RegexpString: phoneRegexString,
		ErrMessage:   phoneErrMessage,
	},
	emailAlias: {
		Alias:        emailAlias,
		RegexpString: emailRegexString,
		ErrMessage:   emailErrMessage,
	},
}

type Pattern struct {
	Alias        string
	RegexpString string
	ErrMessage   string
	re           *regexp.Regexp
}

func (p Pattern) GetRegexp() *regexp.Regexp {
	if p.re == nil {
		p.re = regexp.MustCompile(p.RegexpString)
	}

	return p.re
}
