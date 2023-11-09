package et

import (
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"testing"
)

func TestTranslator_Translate(t *testing.T) {
	type fields struct {
		Translator   ut.Translator
		transTagFunc map[string]TranslationFunc
	}
	type args struct {
		err error
	}
	translator, _ := ut.New(ru.New()).GetTranslator("ru")
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "existed",
			fields: fields{
				Translator: translator,
				transTagFunc: map[string]TranslationFunc{
					"existed": func(ut ut.Translator, msg string, params ...string) (string, error) {
						return "translation for existed", nil
					},
				},
			},
			args: args{
				err: errors.NoType.New("message").T("existed"),
			},
			want:    "translation for existed",
			wantErr: false,
		},
		{
			name: "wrapped",
			fields: fields{
				Translator: translator,
				transTagFunc: map[string]TranslationFunc{
					"existed": func(ut ut.Translator, msg string, params ...string) (string, error) {
						return "translation for existed", nil
					},
				},
			},
			args: args{
				err: errors.BadRequest.Wrap(errors.NoType.New("message").T("existed"), "error"),
			},
			wantErr: true,
		},
		{
			name: "wrapped existed",
			fields: fields{
				Translator: translator,
				transTagFunc: map[string]TranslationFunc{
					"existed": func(ut ut.Translator, msg string, params ...string) (string, error) {
						return "translation for existed", nil
					},
				},
			},
			args: args{
				err: errors.BadRequest.Wrap(errors.NoType.New("message"), "error").T("existed"),
			},
			want: "translation for existed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translator := Translator{
				Translator:   tt.fields.Translator,
				transTagFunc: tt.fields.transTagFunc,
			}
			got, err := translator.Translate(tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Translate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Translate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
