package focus

import (
	"reflect"
	"testing"
)

func Test_readTagParts(t *testing.T) {
	type args struct {
		tagParts []string
	}
	tests := []struct {
		name string
		args args
		want rule
	}{
		{
			name: "without dive",
			args: args{
				tagParts: []string{"required", "min=2", "unique=id"},
			},
			want: rule{
				"required": true,
				"min":      "2",
				"unique":   "id",
			},
		},
		{
			name: "with dive",
			args: args{
				tagParts: []string{"required", "min=2", "dive", "unique=id"},
			},
			want: rule{
				"required": true,
				"min":      "2",
				"items": rule{
					"unique": "id",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readTagParts(tt.args.tagParts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getValidationRules(t *testing.T) {
	type args struct {
		model Model
	}
	tests := []struct {
		name string
		args args
		want rules
	}{
		{
			args: args{
				model: *NewModel(reflect.TypeOf(Example{})),
			},
			want: rules{
				"id": rule{
					"required": true,
					"title":    "ID",
					"type":     stringType,
				},
				"string": rule{
					"required": true,
					"min":      "3",
					"max":      "50",
					"title":    "Строка",
					"type":     stringType,
				},
				"int": rule{
					"omitempty": true,
					"max":       "100",
					"title":     "Целое число",
					"type":      numberType,
				},
				"uint": rule{
					"omitempty": true,
					"title":     "Целое неотрицательное",
					"type":      numberType,
				},
				"time": rule{
					"omitempty": true,
					"title":     "Время",
					"type":      dateType,
				},
				"timePtr": rule{
					"required": true,
					"title":    "TimePtr",
					"type":     dateType,
				},
				"examples": rule{
					"required": true,
					"items": rule{
						"notBlank": true,
						"unique":   true,
						"type":     stringType,
					},
					"title": "Примеры",
					"type":  arrayType,
				},
				"examplePtr": rule{
					"omitempty": true,
					"notBlank":  true,
					"title":     "Пример",
					"type":      stringType,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewModelsRegistry(false)
			registry.Register(Example{})
			model := *registry.GetModel("examples")
			if got := GetValidationRules(model, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValidationRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
