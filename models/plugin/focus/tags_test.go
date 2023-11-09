package focus

import (
	"github.com/aeroideaservices/focus/models/plugin/form"
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func Test_blockFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "succeed",
			args: args{
				field: &Field{},
				value: "Block name",
			},
			wantField: &Field{Block: "Block name"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("blockFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_codeDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "succeed",
			args: args{
				field: &Field{name: "Field"},
			},
			wantField: &Field{Code: "field", name: "Field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codeDefault(tt.args.field)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("codeDefault() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_codeFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "",
			args: args{
				field: &Field{},
				value: "field code",
			},
			wantField: &Field{Code: "field code"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codeFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("codeFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_columnDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "",
			args: args{
				field: &Field{name: "Field"},
			},
			wantField: &Field{Column: "field", name: "Field"},
		},
		{
			name: "",
			args: args{
				field: &Field{Column: "field_col", name: "Field"},
			},
			wantField: &Field{Column: "field", name: "Field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columnDefault(tt.args.field)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("columnDefault() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_columnFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "",
			args: args{
				field: &Field{},
				value: "field_col",
			},
			wantField: &Field{Column: "field_col"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columnFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("columnFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_disabledDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "not primary key",
			args: args{
				field: &Field{},
			},
			wantField: &Field{},
		},
		{
			name: "primary key",
			args: args{
				field: &Field{primaryKey: true},
			},
			wantField: &Field{Disabled: []view{CreateView, UpdateView}, primaryKey: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			disabledDefault(tt.args.field)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("disabledDefault() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_disabledFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
		wantPanic bool
	}{
		{
			name: "empty",
			args: args{
				field: &Field{},
				value: "",
			},
			wantField: &Field{Disabled: []view{CreateView, UpdateView}},
		},
		{
			name: "true",
			args: args{
				field: &Field{},
				value: "true",
			},
			wantField: &Field{Disabled: []view{CreateView, UpdateView}},
		},
		{
			name: "create",
			args: args{
				field: &Field{},
				value: "create",
			},
			wantField: &Field{Disabled: []view{CreateView}},
		},
		{
			name: "update",
			args: args{
				field: &Field{},
				value: "update",
			},
			wantField: &Field{Disabled: []view{UpdateView}},
		},
		{
			name: "create,update",
			args: args{
				field: &Field{},
				value: "create,update",
			},
			wantField: &Field{Disabled: []view{CreateView, UpdateView}},
		},
		{
			name: "false",
			args: args{
				field: &Field{},
				value: "false",
			},
			wantField: &Field{},
		},
		{
			name: "panic",
			args: args{
				field: &Field{},
				value: "creat",
			},
			wantField: nil,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("disabledFill() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			disabledFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("disabledFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_filterFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "empty",
			args: args{
				field: &Field{},
				value: "",
			},
			wantField: &Field{Filterable: true},
		},
		{
			name: "true",
			args: args{
				field: &Field{},
				value: "true",
			},
			wantField: &Field{Filterable: true},
		},
		{
			name: "any",
			args: args{
				field: &Field{},
				value: "sdfsdf",
			},
			wantField: &Field{Filterable: false},
		},
		{
			name: "is media",
			args: args{
				field: &Field{IsMedia: true},
				value: "",
			},
			wantField: &Field{IsMedia: true, Filterable: false},
		},
		{
			name: "is association",
			args: args{
				field: &Field{Association: &Association{}},
				value: "",
			},
			wantField: &Field{Association: &Association{}, Filterable: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filterFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("filterFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_firstToLower(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "starts with upper",
			args: args{
				s: "Upper",
			},
			want: "upper",
		},
		{
			name: "starts with lower",
			args: args{
				s: "lower",
			},
			want: "lower",
		},
		{
			name: "starts with digit",
			args: args{
				s: "1digit",
			},
			want: "1digit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := firstToLower(tt.args.s); got != tt.want {
				t.Errorf("firstToLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hiddenDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "not a primary key",
			args: args{
				field: &Field{},
			},
			wantField: &Field{},
		},
		{
			name: "a primary key",
			args: args{
				field: &Field{primaryKey: true},
			},
			wantField: &Field{Hidden: []view{CreateView, UpdateView}, primaryKey: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hiddenDefault(tt.args.field)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("hiddenDefault() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_hiddenFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
		wantPanic bool
	}{
		{
			name: "empty",
			args: args{
				field: &Field{},
				value: "",
			},
			wantField: &Field{Hidden: []view{CreateView, UpdateView, ListView}},
		},
		{
			name: "true",
			args: args{
				field: &Field{},
				value: "true",
			},
			wantField: &Field{Hidden: []view{CreateView, UpdateView, ListView}},
		},
		{
			name: "create",
			args: args{
				field: &Field{},
				value: "create",
			},
			wantField: &Field{Hidden: []view{CreateView}},
		},
		{
			name: "update",
			args: args{
				field: &Field{},
				value: "update",
			},
			wantField: &Field{Hidden: []view{UpdateView}},
		},
		{
			name: "list",
			args: args{
				field: &Field{},
				value: "list",
			},
			wantField: &Field{Hidden: []view{ListView}},
		},
		{
			name: "create,update",
			args: args{
				field: &Field{},
				value: "create,update",
			},
			wantField: &Field{Hidden: []view{CreateView, UpdateView}},
		},
		{
			name: "create,update,list",
			args: args{
				field: &Field{},
				value: "create,update,list",
			},
			wantField: &Field{Hidden: []view{CreateView, UpdateView, ListView}},
		},
		{
			name: "false",
			args: args{
				field: &Field{},
				value: "false",
			},
			wantField: &Field{},
		},
		{
			name: "any",
			args: args{
				field: &Field{},
				value: "asdfasd",
			},
			wantField: &Field{},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("hiddenFill() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			hiddenFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("hiddenFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_mediaFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "empty",
			args: args{
				field: &Field{},
				value: "",
			},
			wantField: &Field{IsMedia: true},
		},
		{
			name: "true",
			args: args{
				field: &Field{},
				value: "true",
			},
			wantField: &Field{IsMedia: true},
		},
		{
			name: "any",
			args: args{
				field: &Field{},
				value: "sdfsdf",
			},
			wantField: &Field{IsMedia: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mediaFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("mediaFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_multipleDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "a slice",
			args: args{
				field: &Field{t: reflect.TypeOf("")},
			},
			wantField: &Field{t: reflect.TypeOf("")},
		},
		{
			name: "not a slice",
			args: args{
				field: &Field{t: reflect.TypeOf([]string{})},
			},
			wantField: &Field{Multiple: true, t: reflect.TypeOf([]string{})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			multipleDefault(tt.args.field)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("multipleDefault() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_multipleFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "empty",
			args: args{
				field: &Field{},
				value: "",
			},
			wantField: &Field{Multiple: true},
		},
		{
			name: "true",
			args: args{
				field: &Field{},
				value: "true",
			},
			wantField: &Field{Multiple: true},
		},
		{
			name: "any",
			args: args{
				field: &Field{},
				value: "alkdhfasdf",
			},
			wantField: &Field{Multiple: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			multipleFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("multipleFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_positionFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "1",
			args: args{
				field: &Field{},
				value: "1",
			},
			wantField: &Field{Position: 1},
		},
		{
			name: "-31",
			args: args{
				field: &Field{},
				value: "-31",
			},
			wantField: &Field{Position: -31},
		},
		{
			name: "any",
			args: args{
				field: &Field{},
				value: "asdf",
			},
			wantField: &Field{Position: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positionFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("positionFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_primaryKeyFill(t *testing.T) {
	model := &Model{}
	pk := &Field{Model: model, primaryKey: true}
	model.PrimaryKey = pk
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "empty",
			args: args{
				field: &Field{Model: &Model{}},
				value: "",
			},
			wantField: pk,
		},
		{
			name: "true",
			args: args{
				field: &Field{Model: &Model{}},
				value: "true",
			},
			wantField: pk,
		},
		{
			name: "any",
			args: args{
				field: &Field{Model: &Model{}},
				value: "lddjfasdf",
			},
			wantField: &Field{Model: &Model{}, primaryKey: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primaryKeyFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("primaryKeyFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_sortFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "empty",
			args: args{
				field: &Field{},
				value: "",
			},
			wantField: &Field{Sortable: true},
		},
		{
			name: "true",
			args: args{
				field: &Field{},
				value: "true",
			},
			wantField: &Field{Sortable: true},
		},
		{
			name: "any",
			args: args{
				field: &Field{},
				value: "fgdfg",
			},
			wantField: &Field{Sortable: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("sortFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_timeDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "is a time",
			args: args{
				field: &Field{t: reflect.TypeOf(&time.Time{})},
			},
			wantField: &Field{IsTime: true, t: reflect.TypeOf(&time.Time{})},
		},
		{
			name: "isn't a time",
			args: args{
				field: &Field{t: reflect.TypeOf("")},
			},
			wantField: &Field{IsTime: false, t: reflect.TypeOf("")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeDefault(tt.args.field)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("timeDefault() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_timeFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "empty",
			args: args{
				field: &Field{},
				value: "",
			},
			wantField: &Field{IsTime: true},
		},
		{
			name: "true",
			args: args{
				field: &Field{},
				value: "true",
			},
			wantField: &Field{IsTime: true},
		},
		{
			name: "any",
			args: args{
				field: &Field{},
				value: "sdlkfsdf",
			},
			wantField: &Field{IsTime: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("timeFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_titleDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "",
			args: args{
				field: &Field{name: "Field"},
			},
			wantField: &Field{Title: "Field", name: "Field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			titleDefault(tt.args.field)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("titleDefault() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_titleFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "non-empty",
			args: args{
				field: &Field{},
				value: "Field title",
			},
			wantField: &Field{Title: "Field title"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			titleFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("titleFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_uniqueFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "empty",
			args: args{
				field: &Field{},
				value: "",
			},
			wantField: &Field{IsUnique: true},
		},
		{
			name: "true",
			args: args{
				field: &Field{},
				value: "true",
			},
			wantField: &Field{IsUnique: true},
		},
		{
			name: "any",
			args: args{
				field: &Field{},
				value: "laskjdfadf",
			},
			wantField: &Field{IsUnique: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uniqueFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("uniqueFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_viewDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
		wantPanic bool
	}{
		{
			name: "is a media",
			args: args{
				field: &Field{IsMedia: true},
			},
			wantField: &Field{IsMedia: true, View: form.Media},
		},
		{
			name: "is a time",
			args: args{
				field: &Field{IsTime: true},
			},
			wantField: &Field{IsTime: true, View: form.DateTimePicker},
		},
		{
			name: "is an association",
			args: args{
				field: &Field{Association: &Association{}},
			},
			wantField: &Field{Association: &Association{}, View: form.Select},
		},
		{
			name: "uuid",
			args: args{
				field: &Field{t: reflect.TypeOf(uuid.UUID{})},
			},
			wantField: &Field{t: reflect.TypeOf(uuid.UUID{}), View: form.TextInput},
		},
		{
			name: "bool",
			args: args{
				field: &Field{t: reflect.TypeOf(false)},
			},
			wantField: &Field{t: reflect.TypeOf(false), View: form.Checkbox},
		},
		{
			name: "int",
			args: args{
				field: &Field{t: reflect.TypeOf(0)},
			},
			wantField: &Field{t: reflect.TypeOf(0), View: form.IntInput},
		},
		{
			name: "uint",
			args: args{
				field: &Field{t: reflect.TypeOf(uint(0))},
			},
			wantField: &Field{t: reflect.TypeOf(uint(0)), View: form.UintInput},
		},
		{
			name: "float64",
			args: args{
				field: &Field{t: reflect.TypeOf(0.0)},
			},
			wantField: &Field{t: reflect.TypeOf(0.0), View: form.FloatInput},
		},
		{
			name: "string",
			args: args{
				field: &Field{t: reflect.TypeOf("")},
			},
			wantField: &Field{t: reflect.TypeOf(""), View: form.TextInput},
		},
		{
			name: "struct",
			args: args{
				field: &Field{t: reflect.TypeOf(Example{})},
			},
			wantField: &Field{t: reflect.TypeOf(Example{}), View: form.Select},
		},
		{
			name: "panic",
			args: args{
				field: &Field{t: reflect.TypeOf(new(any))},
			},
			wantField: &Field{t: reflect.TypeOf(new(any))},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("viewDefault() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			viewDefault(tt.args.field)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("viewDefault() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_viewExtraFill(t *testing.T) {
	form.RegisterViewsExtras(form.ViewsExtras{
		"first":  form.ViewExtras{"a": 1, "b": 2},
		"second": form.ViewExtras{"display": "field"},
	})
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
		wantPanic bool
	}{
		{
			name: "first",
			args: args{
				field: &Field{},
				value: "first",
			},
			wantField: &Field{ViewExtra: form.ViewExtras{"a": 1, "b": 2}},
		},
		{
			name: "first",
			args: args{
				field: &Field{},
				value: "second",
			},
			wantField: &Field{ViewExtra: form.ViewExtras{"display": "field"}},
		},
		{
			name: "not registered",
			args: args{
				field: &Field{},
				value: "alsfkjsdf",
			},
			wantField: &Field{},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("viewExtraFill() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			viewExtraFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("viewExtraFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_viewFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
		wantPanic bool
	}{
		{
			name: "none",
			args: args{
				field: &Field{},
				value: "none",
			},
			wantField: &Field{View: form.None},
		},
		{
			name: "checkbox",
			args: args{
				field: &Field{},
				value: "checkbox",
			},
			wantField: &Field{View: form.Checkbox},
		},
		{
			name: "intInput",
			args: args{
				field: &Field{},
				value: "intInput",
			},
			wantField: &Field{View: form.IntInput},
		},
		{
			name: "uintInput",
			args: args{
				field: &Field{},
				value: "uintInput",
			},
			wantField: &Field{View: form.UintInput},
		},
		{
			name: "floatInput",
			args: args{
				field: &Field{},
				value: "floatInput",
			},
			wantField: &Field{View: form.FloatInput},
		},
		{
			name: "rating",
			args: args{
				field: &Field{},
				value: "rating",
			},
			wantField: &Field{View: form.Rating},
		},
		{
			name: "select",
			args: args{
				field: &Field{},
				value: "select",
			},
			wantField: &Field{View: form.Select},
		},
		{
			name: "datePickerInput",
			args: args{
				field: &Field{},
				value: "datePickerInput",
			},
			wantField: &Field{View: form.DatePickerInput},
		},
		{
			name: "dateTimePicker",
			args: args{
				field: &Field{},
				value: "dateTimePicker",
			},
			wantField: &Field{View: form.DateTimePicker},
		},
		{
			name: "textarea",
			args: args{
				field: &Field{},
				value: "textarea",
			},
			wantField: &Field{View: form.Textarea},
		},
		{
			name: "textInput",
			args: args{
				field: &Field{},
				value: "textInput",
			},
			wantField: &Field{View: form.TextInput},
		},
		{
			name: "wysiwyg",
			args: args{
				field: &Field{},
				value: "wysiwyg",
			},
			wantField: &Field{View: form.Wysiwyg},
		},
		{
			name: "editorJs",
			args: args{
				field: &Field{},
				value: "editorJs",
			},
			wantField: &Field{View: form.EditorJs},
		},
		{
			name: "media",
			args: args{
				field: &Field{},
				value: "media",
			},
			wantField: &Field{View: form.Media},
		},
		{
			name: "phoneInput",
			args: args{
				field: &Field{},
				value: "phoneInput",
			},
			wantField: &Field{View: form.PhoneInput},
		},
		{
			name: "emailInput",
			args: args{
				field: &Field{},
				value: "emailInput",
			},
			wantField: &Field{View: form.EmailInput},
		},
		{
			name: "not registered",
			args: args{
				field: &Field{},
				value: "ljdfsdfg",
			},
			wantField: &Field{},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("viewFill() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			viewFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("viewFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_associationDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "belongs to",
			args: args{
				field: &Field{t: reflect.TypeOf(Example{}), Association: &Association{}},
			},
			wantField: &Field{t: reflect.TypeOf(Example{}), Association: &Association{Type: BelongsTo}},
		},
		{
			name: "has many",
			args: args{
				field: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{}},
			},
			wantField: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{Type: HasMany}},
		},
		{
			name: "with association type",
			args: args{
				field: &Field{t: reflect.TypeOf(Example{}), Association: &Association{Type: HasOne}},
			},
			wantField: &Field{t: reflect.TypeOf(Example{}), Association: &Association{Type: HasOne}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			associationDefault(tt.args.field)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("associationDefault() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_associationFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "already set",
			args: args{
				field: &Field{t: reflect.TypeOf(Example{}), Association: &Association{Type: ManyToMany}},
				value: "belongsTo",
			},
			wantField: &Field{t: reflect.TypeOf(Example{}), Association: &Association{Type: ManyToMany}},
		},
		{
			name: "belongs to",
			args: args{
				field: &Field{t: reflect.TypeOf(Example{}), Association: &Association{}},
				value: "belongsTo",
			},
			wantField: &Field{t: reflect.TypeOf(Example{}), Association: &Association{Type: BelongsTo}},
		},
		{
			name: "has one",
			args: args{
				field: &Field{t: reflect.TypeOf(Example{}), Association: &Association{}},
				value: "hasOne",
			},
			wantField: &Field{t: reflect.TypeOf(Example{}), Association: &Association{Type: HasOne}},
		},
		{
			name: "has many",
			args: args{
				field: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{}},
				value: "hasMany",
			},
			wantField: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{Type: HasMany}},
		},
		{
			name: "many to many",
			args: args{
				field: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{}},
				value: "manyToMany",
			},
			wantField: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{Type: ManyToMany}},
		},
		{
			name: "many2many",
			args: args{
				field: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{}},
				value: "many2many",
			},
			wantField: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{Type: ManyToMany}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			associationFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("associationFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_foreignKeyDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "already set",
			args: args{
				field: &Field{Association: &Association{ForeignKey: "sdkjfsd"}},
			},
			wantField: &Field{Association: &Association{ForeignKey: "sdkjfsd"}},
		},
		{
			name: "no association",
			args: args{
				field: &Field{Association: &Association{}},
			},
			wantField: &Field{Association: &Association{}},
		},
		{
			name: "belongs to",
			args: args{
				field: &Field{
					Column:      "example",
					Association: &Association{Type: BelongsTo, Model: &Model{PrimaryKey: &Field{Column: "id"}}},
				},
			},
			wantField: &Field{
				Column:      "example",
				Association: &Association{Type: BelongsTo, ForeignKey: "example_id", Model: &Model{PrimaryKey: &Field{Column: "id"}}},
			},
		},
		{
			name: "has one",
			args: args{
				field: &Field{
					Column:      "element",
					Model:       &Model{TableName: "examples", PrimaryKey: &Field{Column: "id"}},
					Association: &Association{Type: HasOne},
				},
			},
			wantField: &Field{
				Column:      "element",
				Model:       &Model{TableName: "examples", PrimaryKey: &Field{Column: "id"}},
				Association: &Association{Type: HasOne, ForeignKey: "example_id"},
			},
		},
		{
			name: "has many",
			args: args{
				field: &Field{
					Column:      "elements",
					Model:       &Model{TableName: "examples", PrimaryKey: &Field{Column: "id"}},
					Association: &Association{Type: HasMany},
				},
			},
			wantField: &Field{
				Column:      "elements",
				Model:       &Model{TableName: "examples", PrimaryKey: &Field{Column: "id"}},
				Association: &Association{Type: HasMany, ForeignKey: "example_id"},
			},
		},
		{
			name: "many to many",
			args: args{
				field: &Field{
					Model:       &Model{TableName: "examples", PrimaryKey: &Field{Column: "id"}},
					Association: &Association{Type: ManyToMany},
				},
			},
			wantField: &Field{
				Model:       &Model{PrimaryKey: &Field{Column: "id"}},
				Association: &Association{Type: ManyToMany, ForeignKey: "id"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foreignKeyDefault(tt.args.field)
		})
	}
}

func Test_foreignKeyFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "succeed",
			args: args{
				field: &Field{t: reflect.TypeOf(Example{}), Association: &Association{}},
				value: "field_column",
			},
			wantField: &Field{t: reflect.TypeOf(Example{}), Association: &Association{ForeignKey: "field_column"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foreignKeyFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("foreignKeyFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_joinForeignKeyDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "many to many",
			args: args{
				field: &Field{
					Model:       &Model{TableName: "examples"},
					Association: &Association{Type: ManyToMany, ForeignKey: "id"},
				},
			},
			wantField: &Field{
				Model:       &Model{TableName: "examples"},
				Association: &Association{Type: ManyToMany, ForeignKey: "id", JoinForeignKey: "example_id"},
			},
		},
		{
			name: "not many to many",
			args: args{
				field: &Field{Association: &Association{Type: BelongsTo}},
			},
			wantField: &Field{Association: &Association{Type: BelongsTo}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			joinForeignKeyDefault(tt.args.field)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("joinForeignKeyDefault() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_joinForeignKeyFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "succeed",
			args: args{
				field: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{}},
				value: "field_column",
			},
			wantField: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{JoinForeignKey: "field_column"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			joinForeignKeyFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("joinForeignKeyFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_joinReferencesDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "already written",
			args: args{
				field: &Field{Association: &Association{JoinReferences: "qweqwe"}},
			},
			wantField: &Field{Association: &Association{JoinReferences: "qweqwe"}},
		},
		{
			name: "not many to many",
			args: args{
				field: &Field{Association: &Association{Type: BelongsTo}},
			},
			wantField: &Field{Association: &Association{Type: BelongsTo}},
		},
		{
			name: "many to many",
			args: args{
				field: &Field{Association: &Association{
					Type:       ManyToMany,
					Model:      &Model{TableName: "examples"},
					References: "id",
				}},
			},
			wantField: &Field{Association: &Association{
				Type:           ManyToMany,
				Model:          &Model{TableName: "examples"},
				References:     "id",
				JoinReferences: "example_id",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			joinReferencesDefault(tt.args.field)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("joinReferencesDefault() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_joinReferencesFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "succeed",
			args: args{
				field: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{}},
				value: "field_col",
			},
			wantField: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{JoinReferences: "field_col"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			joinReferencesFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("joinReferencesFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_many2manyDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "without association type",
			args: args{
				field: &Field{Association: &Association{}},
			},
			wantField: &Field{Association: &Association{}},
		},
		{
			name: "succeed",
			args: args{
				field: &Field{Association: &Association{Many2Many: "lkskdfg", Type: ManyToMany}},
			},
			wantField: &Field{Association: &Association{Many2Many: "lkskdfg", Type: ManyToMany}},
		},
		{
			name: "succeed",
			args: args{
				field: &Field{Association: &Association{Many2Many: "lkskdfg", Type: BelongsTo}},
			},
			wantField: &Field{Association: &Association{Many2Many: "lkskdfg", Type: BelongsTo}},
		},
		{
			name: "succeed",
			args: args{
				field: &Field{
					Model:       &Model{TableName: "first"},
					Association: &Association{Type: ManyToMany, Model: &Model{TableName: "second"}},
				},
			},
			wantField: &Field{
				Model:       &Model{TableName: "first"},
				Association: &Association{Type: ManyToMany, Many2Many: "first_second", Model: &Model{TableName: "second"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			many2manyDefault(tt.args.field)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("many2manyDefault() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_many2manyFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "succeed",
			args: args{
				field: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{}},
				value: "example_example",
			},
			wantField: &Field{t: reflect.TypeOf([]Example{}), Association: &Association{Many2Many: "example_example", Type: ManyToMany}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			many2manyFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("many2manyFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_referencesDefault(t *testing.T) {
	type args struct {
		field *Field
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "without association",
			args: args{
				field: &Field{Association: &Association{}},
			},
			wantField: &Field{Association: &Association{}},
		},
		{
			name: "with association",
			args: args{
				field: &Field{Association: &Association{References: "id"}},
			},
			wantField: &Field{Association: &Association{References: "id"}},
		},
		{
			name: "many to many",
			args: args{
				field: &Field{Association: &Association{
					Model: &Model{PrimaryKey: &Field{Column: "pk_col"}},
					Type:  ManyToMany,
				}},
			},
			wantField: &Field{Association: &Association{
				References: "pk_col",
				Model:      &Model{PrimaryKey: &Field{Column: "pk_col"}},
				Type:       ManyToMany,
			}},
		},
		{
			name: "belongs to",
			args: args{
				field: &Field{Association: &Association{
					Model: &Model{PrimaryKey: &Field{Column: "pk_col"}},
					Type:  BelongsTo,
				}},
			},
			wantField: &Field{Association: &Association{
				References: "pk_col",
				Model:      &Model{PrimaryKey: &Field{Column: "pk_col"}},
				Type:       BelongsTo,
			}},
		},
		{
			name: "has one",
			args: args{
				field: &Field{
					Model:       &Model{PrimaryKey: &Field{Column: "pk_col"}},
					Association: &Association{Type: HasOne},
				},
			},
			wantField: &Field{
				Model:       &Model{PrimaryKey: &Field{Column: "pk_col"}},
				Association: &Association{Type: HasOne, References: "pk_col"},
			},
		},
		{
			name: "has one",
			args: args{
				field: &Field{
					Model:       &Model{PrimaryKey: &Field{Column: "pk_col"}},
					Association: &Association{Type: HasMany},
				},
			},
			wantField: &Field{
				Model:       &Model{PrimaryKey: &Field{Column: "pk_col"}},
				Association: &Association{Type: HasMany, References: "pk_col"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			referencesDefault(tt.args.field)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("referencesDefault() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_referencesFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "",
			args: args{
				field: &Field{t: reflect.TypeOf(Example{}), Association: &Association{}},
				value: "ref_col",
			},
			wantField: &Field{t: reflect.TypeOf(Example{}), Association: &Association{References: "ref_col"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			referencesFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("referencesFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}

func Test_joinSortFill(t *testing.T) {
	type args struct {
		field *Field
		value string
	}
	tests := []struct {
		name      string
		args      args
		wantField *Field
	}{
		{
			name: "succeed",
			args: args{
				field: &Field{Association: &Association{Type: ManyToMany}},
				value: "sort_col",
			},
			wantField: &Field{Association: &Association{Type: ManyToMany, JoinSort: "sort_col"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			joinSortFill(tt.args.field, tt.args.value)
			if !reflect.DeepEqual(tt.args.field, tt.wantField) {
				t.Errorf("joinSortFill() gotField = %v, wantField %v", tt.args.field, tt.wantField)
			}
		})
	}
}
