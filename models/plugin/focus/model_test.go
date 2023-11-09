package focus

import (
	"github.com/aeroideaservices/focus/models/plugin/form"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"reflect"
	"testing"
	"time"
)

func ptrTo[T any](val T) *T {
	return &val
}

type Example struct {
	Id           uuid.UUID  `focus:"title:ID;primaryKey" validate:"required"`
	String       string     `focus:"title:Строка" validate:"required,min=3,max=50"`
	Int          int        `focus:"title:Целое число" validate:"omitempty,max=100"`
	Uint         uint       `focus:"title:Целое неотрицательное"`
	Float        float64    `focus:"" validate:"-"`
	Time         time.Time  `focus:"title:Время" validate:""`
	TimePtr      *time.Time `focus:"" validate:"required"`
	Examples     []Example  `focus:"title:Примеры;many2many:examples_examples;joinForeignKey:next_example_id" validate:"required,dive,notBlank,unique"`
	ExamplePtr   *Example   `focus:"title:Пример;association" validate:"omitempty,notBlank"`
	ExamplePtrId *uuid.UUID `focus:"-" validate:"omitempty"`
}

func (e Example) TableName() string {
	return "examples"
}

func (e Example) ModelTitle() string {
	return "Примеры"
}

func TestModel_NewElement(t *testing.T) {
	var model Model
	pk := &Field{
		primaryKey: true,
		t:          reflect.TypeOf(uuid.UUID{}),
		Code:       "id",
		name:       "Id",
	}
	model = Model{
		Code:       "",
		PrimaryKey: pk,
		Fields: Fields{
			pk,
			{
				t:    reflect.TypeOf(""),
				Code: "string",
				name: "String",
			},
			{
				t:    reflect.TypeOf(0),
				Code: "int",
				name: "Int",
			},
			{
				t:    reflect.TypeOf(uint(0)),
				Code: "uint",
				name: "Uint",
			},
			{
				t:    reflect.TypeOf(0.0),
				Code: "float",
				name: "Float",
			},
			{
				t:    reflect.TypeOf(time.Time{}),
				Code: "time",
				name: "Time",
			},
			{
				t:    reflect.TypeOf(ptrTo(time.Time{})),
				Code: "timePtr",
				name: "TimePtr",
			},
			{
				Association: &Association{
					Type:  ManyToMany,
					Model: &model,
				},
				Multiple: true,
				t:        reflect.TypeOf([]Example{}),
				Code:     "examples",
				name:     "Examples",
			},
			{
				Association: &Association{
					Type:  BelongsTo,
					Model: &model,
				},
				t:    reflect.TypeOf(ptrTo(Example{})),
				Code: "examplePtr",
				name: "ExamplePtr",
			},
		},
		t: reflect.TypeOf(Example{}),
	}
	type args struct {
		fieldsMap map[string]any
	}
	tests := []struct {
		name      string
		args      args
		wantModel any
		wantErr   bool
		model     Model
	}{
		{
			name:  "nil",
			model: model,
			args: args{
				fieldsMap: nil,
			},
			wantModel: &Example{},
			wantErr:   false,
		},
		{
			name:  "empty map",
			model: model,
			args: args{
				fieldsMap: make(map[string]any),
			},
			wantModel: &Example{},
			wantErr:   false,
		},
		{
			name:  "succeed",
			model: model,
			args: args{
				fieldsMap: map[string]any{
					"id":         "11111111-1111-1111-1111-111111111111",
					"string":     "str",
					"int":        -30,
					"uint":       25,
					"float":      123.456,
					"time":       "2023-02-12T13:33:09Z",
					"timePtr":    "2023-02-25T13:33:09Z",
					"examples":   []string{"11111111-1111-1111-1111-111111111112", "11111111-1111-1111-1111-111111111113"},
					"examplePtr": "11111111-1111-1111-1111-111111111114",
				},
			},
			wantModel: &Example{
				Id:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				String:       "str",
				Int:          -30,
				Uint:         25,
				Time:         func() time.Time { t, _ := time.Parse(time.RFC3339, "2023-02-12T13:33:09Z"); return t }(),
				TimePtr:      func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-02-25T13:33:09Z"); return &t }(),
				Float:        123.456,
				Examples:     []Example{{Id: uuid.MustParse("11111111-1111-1111-1111-111111111112")}, {Id: uuid.MustParse("11111111-1111-1111-1111-111111111113")}},
				ExamplePtr:   &Example{Id: uuid.MustParse("11111111-1111-1111-1111-111111111114")},
				ExamplePtrId: nil,
			},
			wantErr: false,
		},
		{
			name:  "wrong value type",
			model: model,
			args: args{
				fieldsMap: map[string]any{
					"id":         "11111111-1111-1111-1111-111111111111",
					"string":     123,
					"int":        -30,
					"uint":       25,
					"time":       "2023-02-12T13:33:09Z",
					"timePtr":    "2023-02-25T13:33:09Z",
					"float":      123.456,
					"examples":   []string{"11111111-1111-1111-1111-111111111112", "11111111-1111-1111-1111-111111111113"},
					"examplePtr": "11111111-1111-1111-1111-111111111114",
				},
			},
			wantModel: nil,
			wantErr:   true,
		},
		{
			name:  "wrong value type",
			model: model,
			args: args{
				fieldsMap: map[string]any{
					"id":         "11111111-1111-1111-1111-11111111111",
					"string":     "123",
					"int":        -30,
					"uint":       25,
					"time":       "2023-02-12T13:33:09Z",
					"timePtr":    "2023-02-25T13:33:09Z",
					"float":      123.456,
					"examples":   []string{"11111111-1111-1111-1111-111111111112", "11111111-1111-1111-1111-111111111113"},
					"examplePtr": "11111111-1111-1111-1111-111111111114",
				},
			},
			wantModel: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model
			gotModel, err := m.NewElement(tt.args.fieldsMap, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotModel, tt.wantModel) {
				t.Errorf("NewElement() gotModel = %v, want %v", gotModel, tt.wantModel)
			}
		})
	}
}

func TestModel_ElementToMap(t *testing.T) {
	var model Model
	model = Model{
		Code: "",
		PrimaryKey: &Field{
			primaryKey: true,
			t:          reflect.TypeOf(uuid.UUID{}),
			Code:       "id",
		},
		Fields: Fields{
			{
				primaryKey: true,
				t:          reflect.TypeOf(uuid.UUID{}),
				Code:       "id",
			},
			{
				t:    reflect.TypeOf(""),
				Code: "string",
			},
			{
				t:    reflect.TypeOf(0),
				Code: "int",
			},
			{
				t:    reflect.TypeOf(uint(0)),
				Code: "uint",
			},
			{
				t:    reflect.TypeOf(0.0),
				Code: "float",
			},
			{
				t:    reflect.TypeOf(time.Time{}),
				Code: "time",
			},
			{
				t:    reflect.TypeOf(ptrTo(time.Time{})),
				Code: "timePtr",
			},
			{
				Association: &Association{
					Type:  ManyToMany,
					Model: &model,
				},
				t:    reflect.TypeOf([]Example{}),
				Code: "examples",
			},
			{
				Association: &Association{
					Type:  BelongsTo,
					Model: &model,
				},
				t:    reflect.TypeOf(Example{}),
				Code: "examplePtr",
			},
		},
		t: reflect.TypeOf(Example{}),
	}
	type args struct {
		modelElement any
		filter       func(field Field) bool
	}
	tests := []struct {
		name      string
		args      args
		wantModel map[string]any
		wantErr   bool
		model     Model
	}{
		{
			name:  "empty",
			model: model,
			args: args{
				modelElement: Example{},
				filter:       nil,
			},
			wantModel: map[string]any{
				"id":         "00000000-0000-0000-0000-000000000000",
				"string":     "",
				"int":        float64(0),
				"uint":       float64(0),
				"time":       "0001-01-01T00:00:00Z",
				"timePtr":    nil,
				"float":      float64(0),
				"examples":   nil,
				"examplePtr": nil,
			},
			wantErr: false,
		},
		{
			name:  "non-empty",
			model: model,
			args: args{
				modelElement: Example{
					Id:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
					String:       "str",
					Int:          -30,
					Uint:         25,
					Time:         func() time.Time { t, _ := time.Parse(time.RFC3339, "2023-02-12T13:33:09Z"); return t }(),
					TimePtr:      func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-02-25T13:33:09Z"); return &t }(),
					Float:        123.456,
					Examples:     []Example{{Id: uuid.MustParse("11111111-1111-1111-1111-111111111112")}, {Id: uuid.MustParse("11111111-1111-1111-1111-111111111113")}},
					ExamplePtr:   &Example{Id: uuid.MustParse("11111111-1111-1111-1111-111111111114")},
					ExamplePtrId: ptrTo(uuid.MustParse("11111111-1111-1111-1111-111111111114")),
				},
				filter: nil,
			},
			wantModel: map[string]any{
				"id":         "11111111-1111-1111-1111-111111111111",
				"string":     "str",
				"int":        float64(-30),
				"uint":       float64(25),
				"time":       "2023-02-12T13:33:09Z",
				"timePtr":    "2023-02-25T13:33:09Z",
				"float":      123.456,
				"examples":   []any{"11111111-1111-1111-1111-111111111112", "11111111-1111-1111-1111-111111111113"},
				"examplePtr": "11111111-1111-1111-1111-111111111114",
			},
			wantErr: false,
		},
		{
			name:  "with filter",
			model: model,
			args: args{
				modelElement: Example{
					Id:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
					String:       "str",
					Int:          -30,
					Uint:         25,
					Time:         func() time.Time { t, _ := time.Parse(time.RFC3339, "2023-02-12T13:33:09Z"); return t }(),
					TimePtr:      func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-02-25T13:33:09Z"); return &t }(),
					Float:        123.456,
					Examples:     []Example{{Id: uuid.MustParse("11111111-1111-1111-1111-111111111112")}, {Id: uuid.MustParse("11111111-1111-1111-1111-111111111113")}},
					ExamplePtr:   &Example{Id: uuid.MustParse("11111111-1111-1111-1111-111111111114")},
					ExamplePtrId: ptrTo(uuid.MustParse("11111111-1111-1111-1111-111111111114")),
				},
				filter: func(field Field) bool {
					return slices.Contains([]string{"id", "string"}, field.Code)
				},
			},
			wantModel: map[string]any{
				"id":     "11111111-1111-1111-1111-111111111111",
				"string": "str",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model
			gotModel, err := m.ElementToMap(tt.args.modelElement, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ModelElementToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotModel, tt.wantModel) {
				t.Errorf("ModelElementToMap() gotModel = %v, want %v", gotModel, tt.wantModel)
			}
		})
	}
}

func TestGetPKs(t *testing.T) {
	type args struct {
		obj    any
		pkName string
	}
	tests := []struct {
		name    string
		args    args
		want    []any
		wantErr bool
	}{
		{
			name: "struct",
			args: args{
				obj:    struct{ Id int }{Id: 123},
				pkName: "Id",
			},
			want:    []any{123},
			wantErr: false,
		},
		{
			name: "[]struct",
			args: args{
				obj:    []struct{ Id int }{{Id: 123}, {Id: 456}},
				pkName: "Id",
			},
			want:    []any{123, 456},
			wantErr: false,
		},
		{
			name: "*struct",
			args: args{
				obj:    &struct{ Id int }{Id: 123},
				pkName: "Id",
			},
			want:    []any{123},
			wantErr: false,
		},
		{
			name: "*[]*struct",
			args: args{
				obj:    &[]*struct{ Id int }{{Id: 123}, {Id: 456}, {Id: 789}},
				pkName: "Id",
			},
			want:    []any{123, 456, 789},
			wantErr: false,
		},
		{
			name: "wrong type",
			args: args{
				obj:    "123456",
				pkName: "Id",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong key",
			args: args{
				obj:    struct{ Id int }{Id: 123},
				pkName: "id",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong type",
			args: args{
				obj:    []string{"123"},
				pkName: "id",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPKs(tt.args.obj, tt.args.pkName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPKs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPKs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_ElementsSlice(t *testing.T) {
	model := NewModel(RawType(reflect.TypeOf(Example{})))

	type args struct {
		modelElementsMaps []map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "nil",
			args: args{
				modelElementsMaps: nil,
			},
			want:    []Example{},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				modelElementsMaps: []map[string]any{},
			},
			want:    []Example{},
			wantErr: false,
		},
		{
			name: "not empty",
			args: args{
				modelElementsMaps: []map[string]any{
					{"id": "11111111-1111-1111-1111-111111111111"},
					{"id": "11111111-2222-1111-1111-111111111111"},
				},
			},
			want: []Example{
				{Id: uuid.MustParse("11111111-1111-1111-1111-111111111111")},
				{Id: uuid.MustParse("11111111-2222-1111-1111-111111111111")},
			},
			wantErr: false,
		},
		{
			name: "err",
			args: args{
				modelElementsMaps: []map[string]any{
					{"id": 123},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := model.ElementsSlice(tt.args.modelElementsMaps)
			if (err != nil) != tt.wantErr {
				t.Errorf("ElementsSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ElementsSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_UpdateElement(t *testing.T) {
	type fields struct {
		TableName  string
		Code       string
		Title      string
		PrimaryKey *Field
		Fields     Fields
		name       string
		t          reflect.Type
	}
	type args struct {
		old    any
		new    any
		filter func(field *Field) bool
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantModel *Example
	}{
		{
			name:   "without filter",
			fields: fields{},
			args: args{
				old: &Example{
					Id:     uuid.New(),
					String: "abc",
				},
				new: &Example{
					Id:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
					String: "123",
					Int:    123,
					Examples: []Example{
						{Id: uuid.MustParse("11111111-2222-1111-1111-111111111111")},
						{Id: uuid.MustParse("11111111-3333-1111-1111-111111111111")},
					},
					ExamplePtr: &Example{Id: uuid.MustParse("11111111-4444-1111-1111-111111111111")},
				},
				filter: nil,
			},
			wantErr: false,
			wantModel: &Example{
				Id:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				String: "123",
				Int:    123,
				Examples: []Example{
					{Id: uuid.MustParse("11111111-2222-1111-1111-111111111111")},
					{Id: uuid.MustParse("11111111-3333-1111-1111-111111111111")},
				},
				ExamplePtr: &Example{Id: uuid.MustParse("11111111-4444-1111-1111-111111111111")},
			},
		},
		{
			name:   "with filter",
			fields: fields{},
			args: args{
				old: &Example{
					Id:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
					String: "abc",
				},
				new: &Example{
					Id:     uuid.New(),
					String: "123",
					Int:    123,
					Examples: []Example{
						{Id: uuid.MustParse("11111111-2222-1111-1111-111111111111")},
						{Id: uuid.MustParse("11111111-3333-1111-1111-111111111111")},
					},
					ExamplePtr: &Example{Id: uuid.MustParse("11111111-4444-1111-1111-111111111111")},
				},
				filter: func(field *Field) bool {
					return slices.Contains([]string{"string", "int", "examplePtr"}, field.Code)
				},
			},
			wantErr: false,
			wantModel: &Example{
				Id:         uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				String:     "123",
				Int:        123,
				ExamplePtr: &Example{Id: uuid.MustParse("11111111-4444-1111-1111-111111111111")},
			},
		},
		{
			name:   "different types",
			fields: fields{},
			args: args{
				old: &Example{},
				new: &Example2{},
			},
			wantErr:   true,
			wantModel: &Example{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewModelsRegistry(false)
			registry.Register(Example{})
			m := registry.GetModel("examples")
			if err := m.UpdateElement(tt.args.old, tt.args.new, tt.args.filter); (err != nil) != tt.wantErr {
				t.Errorf("UpdateElement() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.old, tt.wantModel) {
				t.Errorf("UpdateElement() got = %v, want %v", tt.args.old, tt.wantModel)
			}
		})
	}
}

func TestNewModel(t *testing.T) {
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
					Many2Many:      "examples_examples",
					JoinForeignKey: "next_example_id",
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
					modelCode: "examples",
				},
				name: "ExamplePtr",
				t:    reflect.TypeOf(ptrTo(Example{})),
			},
		},
		name: "Example",
		t:    reflect.TypeOf(Example{}),
	}

	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name      string
		args      args
		want      *Model
		wantPanic bool
	}{
		{
			name: "succeed",
			args: args{
				t: RawType(reflect.TypeOf(Example{})),
			},
			want:      &model,
			wantPanic: false,
		},
		{
			name: "panic",
			args: args{
				t: RawType(reflect.TypeOf("")),
			},
			want:      nil,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("NewModel() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			if got := NewModel(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_modelCode(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantPanic bool
	}{
		{
			name: "success",
			args: args{
				t: reflect.TypeOf(Example{}),
			},
			want:      "examples",
			wantPanic: false,
		},
		{
			name: "success",
			args: args{
				t: reflect.TypeOf(struct{}{}),
			},
			want:      "",
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("modelCode() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			if got := modelCode(tt.args.t); got != tt.want {
				t.Errorf("modelCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
