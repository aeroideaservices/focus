package focus

import (
	"github.com/aeroideaservices/focus/models/plugin/form"
	"github.com/google/uuid"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"reflect"
	"testing"
	"time"
)

func TestNewModelsRegistry(t *testing.T) {
	type args struct {
		supportMedia bool
	}
	tests := []struct {
		name string
		args args
		want *ModelsRegistry
	}{
		{
			name: "without supporting media",
			args: args{
				supportMedia: false,
			},
			want: &ModelsRegistry{
				registered:   make(map[string]*Model),
				supportMedia: false,
			},
		},
		{
			name: "with supporting media",
			args: args{
				supportMedia: true,
			},
			want: &ModelsRegistry{
				registered:   make(map[string]*Model),
				supportMedia: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewModelsRegistry(tt.args.supportMedia); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewModelsRegistry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelsRegistry_GetRegisteredModels(t *testing.T) {
	type fields struct {
		registered   map[string]*Model
		supportMedia bool
	}
	a := &Model{Title: "a", Code: "a"}
	b := &Model{Title: "b", Code: "b"}
	c := &Model{Title: "c", Code: "c"}
	tests := []struct {
		name   string
		fields fields
		want   []*Model
	}{
		{
			name: "",
			fields: fields{
				registered: map[string]*Model{
					"a": a,
					"b": b,
					"c": c,
				},
				supportMedia: false,
			},
			want: []*Model{a, b, c},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ModelsRegistry{
				registered:   tt.fields.registered,
				supportMedia: tt.fields.supportMedia,
			}
			got := r.ListModels()
			slices.SortFunc(got, func(a, b *Model) bool {
				return a.Code < b.Code
			})
			slices.SortFunc(tt.want, func(a, b *Model) bool {
				return a.Code < b.Code
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRegisteredModels() = %v, want %v", got, tt.want)
			}
		})
	}
}

type Example2 struct {
	Id uuid.UUID `focusJson:"id" focus:"title:ID"`
}

func (e Example2) TableName() string {
	return "example_two"
}

func (e Example2) ModelTitle() string {
	return "Пример №2"
}

type Example3 struct {
	Id      uuid.UUID `focusJson:"id" focus:"primaryKey"`
	Code    string    `focusJson:"code"`
	Example Example   `focusJson:"example" focus:""`
}

func (e Example3) TableName() string {
	return "example_three"
}

func (e Example3) ModelTitle() string {
	return "Пример №3"
}

func TestModelsRegistry_Register(t *testing.T) {
	var model Model
	pk := &Field{
		Title:      "ID",
		Column:     "id",
		Code:       "id",
		Hidden:     []view{CreateView, UpdateView},
		Disabled:   []view{CreateView, UpdateView},
		View:       form.TextInput,
		Model:      &model,
		primaryKey: true,
		name:       "Id",
		t:          reflect.TypeOf(uuid.UUID{}),
	}
	model = Model{
		TableName:  "examples",
		Code:       "examples",
		Title:      "Примеры",
		PrimaryKey: pk,
		Fields: Fields{
			pk,
			{
				Title:  "Строка",
				Column: "string",
				Code:   "string",
				Model:  &model,
				View:   form.TextInput,
				name:   "String",
				t:      reflect.TypeOf(""),
			},
			{
				Title:  "Целое число",
				Column: "int",
				Code:   "int",
				View:   form.IntInput,
				Model:  &model,
				name:   "Int",
				t:      reflect.TypeOf(0),
			},
			{
				Title:  "Целое неотрицательное",
				Column: "uint",
				Code:   "uint",
				View:   form.UintInput,
				Model:  &model,
				name:   "Uint",
				t:      reflect.TypeOf(uint(0)),
			},
			{
				Title:  "Float",
				Column: "float",
				Code:   "float",
				View:   form.FloatInput,
				Model:  &model,
				name:   "Float",
				t:      reflect.TypeOf(0.0),
			},
			{
				Title:  "Время",
				Column: "time",
				Code:   "time",
				IsTime: true,
				View:   form.DateTimePicker,
				Model:  &model,
				name:   "Time",
				t:      reflect.TypeOf(time.Time{}),
			},
			{
				Title:  "TimePtr",
				Column: "time_ptr",
				Code:   "timePtr",
				IsTime: true,
				View:   form.DateTimePicker,
				Model:  &model,
				name:   "TimePtr",
				t:      reflect.TypeOf(ptrTo(time.Time{})),
			},
			{
				Title:     "Примеры",
				Column:    "examples",
				Code:      "examples",
				View:      form.Select,
				ViewExtra: nil, // todo
				Multiple:  true,
				Model:     &model,
				Association: &Association{
					Type:           ManyToMany,
					Model:          &model,
					Many2Many:      "examples_examples",
					ForeignKey:     "id",
					References:     "id",
					JoinForeignKey: "next_example_id",
					JoinReferences: "example_id",
					modelCode:      "examples",
				},
				name: "Examples",
				t:    reflect.TypeOf([]Example{}),
			},
			{
				Title:     "Пример",
				Column:    "example_ptr",
				Code:      "examplePtr",
				View:      form.Select,
				ViewExtra: nil, // todo
				Model:     &model,
				Association: &Association{
					Type:       BelongsTo,
					Model:      &model,
					ForeignKey: "example_ptr_id",
					References: "id",
					modelCode:  "examples",
				},
				name: "ExamplePtr",
				t:    reflect.TypeOf(ptrTo(Example{})),
			},
		},
		name: "Example",
		t:    reflect.TypeOf(Example{}),
	}

	type fields struct {
		registered   map[string]*Model
		supportMedia bool
	}
	type args struct {
		items []any
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantRegistered []*Model
		wantPanic      bool
	}{
		{
			name: "succeed",
			fields: fields{
				registered:   make(map[string]*Model),
				supportMedia: false,
			},
			args: args{
				items: []any{
					Example{},
				},
			},
			wantRegistered: []*Model{&model},
			wantPanic:      false,
		},
		{
			name: "panic without pk",
			fields: fields{
				registered:   make(map[string]*Model),
				supportMedia: false,
			},
			args: args{
				items: []any{
					Example2{},
				},
			},
			wantRegistered: nil,
			wantPanic:      true,
		},
		{
			name: "ptr to obj",
			fields: fields{
				registered:   make(map[string]*Model),
				supportMedia: false,
			},
			args: args{
				items: []any{
					&Example{},
				},
			},
			wantRegistered: []*Model{&model},
			wantPanic:      false,
		},
		{
			name: "panic",
			fields: fields{
				registered:   make(map[string]*Model),
				supportMedia: false,
			},
			args: args{
				items: []any{
					123,
				},
			},
			wantRegistered: nil,
			wantPanic:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("ModelsRegistry.Register() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			r := &ModelsRegistry{
				registered:   tt.fields.registered,
				supportMedia: tt.fields.supportMedia,
			}

			r.Register(tt.args.items...)

			gotRegistered := r.ListModels()
			if !reflect.DeepEqual(gotRegistered, tt.wantRegistered) {
				t.Errorf("GetRegisteredModels() gotRegistered = %v, wantRegistered %v", gotRegistered, tt.wantRegistered)
			}
		})
	}
}

func TestModelsRegistry_GetModel(t *testing.T) {
	type fields struct {
		registered   map[string]*Model
		supportMedia bool
	}
	type args struct {
		code string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Model
	}{
		{
			name: "found",
			fields: fields{
				registered:   map[string]*Model{"foo": {Code: "foo"}, "bar": {Code: "bar"}},
				supportMedia: false,
			},
			args: args{
				code: "foo",
			},
			want: &Model{Code: "foo"},
		},
		{
			name: "found",
			fields: fields{
				registered:   map[string]*Model{"foo": {Code: "foo"}, "bar": {Code: "bar"}},
				supportMedia: false,
			},
			args: args{
				code: "bar",
			},
			want: &Model{Code: "bar"},
		},
		{
			name: "found",
			fields: fields{
				registered:   map[string]*Model{"foo": {Code: "foo"}, "bar": {Code: "bar"}},
				supportMedia: false,
			},
			args: args{
				code: "baz",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ModelsRegistry{
				registered:   tt.fields.registered,
				supportMedia: tt.fields.supportMedia,
			}
			if got := r.GetModel(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelsRegistry_NewModel(t *testing.T) {
	type fields struct {
		registered   map[string]*Model
		supportMedia bool
	}
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		wantRegisteredCodes []string
		wantPanic           bool
	}{
		{
			name: "not a struct",
			fields: fields{
				registered:   map[string]*Model{},
				supportMedia: false,
			},
			args: args{
				t: reflect.TypeOf(""),
			},
			wantPanic: true,
		},
		{
			name: "already registered",
			fields: fields{
				registered:   map[string]*Model{"examples": {}},
				supportMedia: false,
			},
			args: args{
				t: reflect.TypeOf(Example{}),
			},
			wantRegisteredCodes: []string{"examples"},
		},
		{
			name: "with association",
			fields: fields{
				registered: map[string]*Model{},
			},
			args: args{
				t: reflect.TypeOf(Example3{}),
			},
			wantRegisteredCodes: []string{"example-three", "examples"},
		},
		{
			name: "with registered association",
			fields: fields{
				registered: map[string]*Model{"examples": {}},
			},
			args: args{
				t: reflect.TypeOf(Example3{}),
			},
			wantRegisteredCodes: []string{"examples", "example-three"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ModelsRegistry{
				registered:   tt.fields.registered,
				supportMedia: tt.fields.supportMedia,
			}
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("NewModel() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			r.NewModel(tt.args.t)
			gotRegisteredCodes := maps.Keys(tt.fields.registered)
			slices.Sort(gotRegisteredCodes)
			slices.Sort(tt.wantRegisteredCodes)
			if !reflect.DeepEqual(gotRegisteredCodes, tt.wantRegisteredCodes) {
				t.Errorf("blockFill() gotRegisteredCodes = %v, wantRegisteredCodes %v", gotRegisteredCodes, tt.wantRegisteredCodes)
			}
		})
	}
}
