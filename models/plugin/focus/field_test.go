package focus

import (
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func TestField_Name(t *testing.T) {
	type fields struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "",
			fields: fields{
				name: "fieldName",
			},
			want: "fieldName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Field{
				name: tt.fields.name,
			}
			if got := f.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_New(t *testing.T) {
	type fields struct {
		Multiple    bool
		Association *Association
		t           reflect.Type
	}
	type args struct {
		value any
	}
	tTime, _ := time.Parse("2006-01-02T15:04:05", "2006-01-02T15:04:05")
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "string->string",
			fields: fields{
				t: reflect.TypeOf(""),
			},
			args: args{
				value: "string",
			},
			want:    "string",
			wantErr: false,
		},
		{
			name: "nil->string",
			fields: fields{
				t: reflect.TypeOf(""),
			},
			args: args{
				value: "string",
			},
			want:    "string",
			wantErr: false,
		},
		{
			name: "string->ptr to string",
			fields: fields{
				t: reflect.TypeOf(ptrTo("")),
			},
			args: args{
				value: "string",
			},
			want:    ptrTo("string"),
			wantErr: false,
		},
		{
			name: "nil->ptr to string",
			fields: fields{
				t: reflect.TypeOf(ptrTo("")),
			},
			args: args{
				value: nil,
			},
			want:    (*string)(nil),
			wantErr: false,
		},
		{
			name: "ptr to string->string",
			fields: fields{
				t: reflect.TypeOf(""),
			},
			args: args{
				value: ptrTo("string"),
			},
			want:    "string",
			wantErr: false,
		},
		{
			name: "int->string",
			fields: fields{
				t: reflect.TypeOf(ptrTo("")),
			},
			args: args{
				value: 123,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "int->int",
			fields: fields{
				t: reflect.TypeOf(0),
			},
			args: args{
				value: 123,
			},
			want:    123,
			wantErr: false,
		},
		{
			name: "uint->int",
			fields: fields{
				t: reflect.TypeOf(0),
			},
			args: args{
				value: uint(123),
			},
			want:    123,
			wantErr: false,
		},
		{
			name: "float->int",
			fields: fields{
				t: reflect.TypeOf(0),
			},
			args: args{
				value: uint(123),
			},
			want:    123,
			wantErr: false,
		},
		{
			name: "string->int",
			fields: fields{
				t: reflect.TypeOf(0),
			},
			args: args{
				value: "123",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "uint->uint",
			fields: fields{
				t: reflect.TypeOf(uint(0)),
			},
			args: args{
				value: uint(123),
			},
			want:    uint(123),
			wantErr: false,
		},
		{
			name: "int->uint",
			fields: fields{
				t: reflect.TypeOf(uint(0)),
			},
			args: args{
				value: 123,
			},
			want:    uint(123),
			wantErr: false,
		},
		{
			name: "negative int->uint",
			fields: fields{
				t: reflect.TypeOf(uint(0)),
			},
			args: args{
				value: -123,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "float->uint",
			fields: fields{
				t: reflect.TypeOf(uint(0)),
			},
			args: args{
				value: 123.456,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "string->uint",
			fields: fields{
				t: reflect.TypeOf(uint(0)),
			},
			args: args{
				value: "123",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "float->float",
			fields: fields{
				t: reflect.TypeOf(0.0),
			},
			args: args{
				value: 123.456,
			},
			want:    123.456,
			wantErr: false,
		},
		{
			name: "int->float",
			fields: fields{
				t: reflect.TypeOf(0.0),
			},
			args: args{
				value: 123,
			},
			want:    123.0,
			wantErr: false,
		},
		{
			name: "uint->float",
			fields: fields{
				t: reflect.TypeOf(0.0),
			},
			args: args{
				value: uint(123),
			},
			want:    123.0,
			wantErr: false,
		},
		{
			name: "nil->float",
			fields: fields{
				t: reflect.TypeOf(0.0),
			},
			args: args{
				value: nil,
			},
			want:    0.0,
			wantErr: false,
		},
		{
			name: "string->time",
			fields: fields{
				t: reflect.TypeOf(time.Time{}),
			},
			args: args{
				value: "2006-01-02T15:04:05Z",
			},
			want:    tTime,
			wantErr: false,
		},
		{
			name: "string->ptr to time",
			fields: fields{
				t: reflect.TypeOf(&time.Time{}),
			},
			args: args{
				value: "2006-01-02T15:04:05Z",
			},
			want:    &tTime,
			wantErr: false,
		},
		{
			name: "nil->time",
			fields: fields{
				t: reflect.TypeOf(time.Time{}),
			},
			args: args{
				value: nil,
			},
			want:    time.Time{},
			wantErr: false,
		},
		{
			name: "wrong format->time",
			fields: fields{
				t: reflect.TypeOf(time.Time{}),
			},
			args: args{
				value: "123",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "int->time",
			fields: fields{
				t: reflect.TypeOf(time.Time{}),
			},
			args: args{
				value: 123,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "string->hasOne association",
			fields: fields{
				Association: &Association{
					Model: &Model{PrimaryKey: &Field{Code: "id"}},
					Type:  HasOne,
				},
				t: reflect.TypeOf(Example{}),
			},
			args: args{
				value: "11111111-1111-1111-1111-111111111111",
			},
			want:    Example{Id: uuid.MustParse("11111111-1111-1111-1111-111111111111")},
			wantErr: false,
		},
		{
			name: "string->ptr to hasOne association",
			fields: fields{
				Association: &Association{
					Model: &Model{PrimaryKey: &Field{Code: "id"}},
					Type:  HasOne,
				},
				t: reflect.TypeOf(&Example{}),
			},
			args: args{
				value: "11111111-1111-1111-1111-111111111111",
			},
			want:    &Example{Id: uuid.MustParse("11111111-1111-1111-1111-111111111111")},
			wantErr: false,
		},
		{
			name: "[]string->hasMany association",
			fields: fields{
				Multiple: true,
				Association: &Association{
					Model: &Model{PrimaryKey: &Field{Code: "id"}},
					Type:  HasMany,
				},
				t: reflect.TypeOf([]Example{}),
			},
			args: args{
				value: []string{"11111111-1111-1111-1111-111111111111", "22222222-1111-1111-1111-111111111111"},
			},
			want: []Example{
				{Id: uuid.MustParse("11111111-1111-1111-1111-111111111111")},
				{Id: uuid.MustParse("22222222-1111-1111-1111-111111111111")},
			},
			wantErr: false,
		},
		{
			name: "string->hasMany association",
			fields: fields{
				Association: &Association{
					Model: &Model{PrimaryKey: &Field{Code: "id"}},
					Type:  HasMany,
				},
				t: reflect.TypeOf([]Example{}),
			},
			args: args{
				value: "11111111-1111-1111-1111-111111111111",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Field{
				Multiple:    tt.fields.Multiple,
				Association: tt.fields.Association,
				t:           tt.fields.t,
			}
			got, err := f.NewValue(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewElement() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_RawKind(t *testing.T) {
	type fields struct {
		t reflect.Type
	}
	tests := []struct {
		name   string
		fields fields
		want   reflect.Kind
	}{
		{
			name: "int",
			fields: fields{
				t: reflect.TypeOf(0),
			},
			want: reflect.Int,
		},
		{
			name: "*float",
			fields: fields{
				t: reflect.TypeOf(ptrTo(0.0)),
			},
			want: reflect.Float64,
		},
		{
			name: "[]bool",
			fields: fields{
				t: reflect.TypeOf([]bool{}),
			},
			want: reflect.Bool,
		},
		{
			name: "[]string",
			fields: fields{
				t: reflect.TypeOf([]string{}),
			},
			want: reflect.String,
		},
		{
			name: "[]*string",
			fields: fields{
				t: reflect.TypeOf([]*string{}),
			},
			want: reflect.String,
		},
		{
			name: "[]*struct",
			fields: fields{
				t: reflect.TypeOf([]*struct{}{}),
			},
			want: reflect.Struct,
		},
		{
			name: "*[]*struct",
			fields: fields{
				t: reflect.TypeOf(&[]*struct{}{}),
			},
			want: reflect.Struct,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Field{
				t: tt.fields.t,
			}
			if got := f.RawKind(); got != tt.want {
				t.Errorf("RawKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_RawType(t *testing.T) {
	type fields struct {
		t reflect.Type
	}
	tests := []struct {
		name   string
		fields fields
		want   reflect.Type
	}{
		{
			name: "*[]*struct",
			fields: fields{
				t: reflect.TypeOf(&[]*struct{}{}),
			},
			want: reflect.TypeOf(struct{}{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Field{
				t: tt.fields.t,
			}
			if got := f.RawType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RawType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_Slice(t *testing.T) {
	type fields struct {
		Association *Association
		t           reflect.Type
		Multiple    bool
	}
	type args struct {
		values []any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []any
		wantErr bool
	}{
		{
			name: "[]string->[]string",
			fields: fields{
				t: reflect.TypeOf(""),
			},
			args: args{
				[]any{"123", "456"},
			},
			want:    []any{"123", "456"},
			wantErr: false,
		},
		{
			name: "[]any->[]string",
			fields: fields{
				t: reflect.TypeOf(""),
			},
			args: args{
				[]any{"123", 456},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "[]string->[]associations",
			fields: fields{
				t:           reflect.TypeOf([]Example{}),
				Association: &Association{Type: HasMany, Model: &Model{PrimaryKey: &Field{Code: "id"}}},
				Multiple:    true,
			},
			args: args{
				[]any{"11111111-1111-1111-1111-111111111111", "22222222-1111-1111-1111-111111111111"},
			},
			want: []any{
				Example{Id: uuid.MustParse("11111111-1111-1111-1111-111111111111")},
				Example{Id: uuid.MustParse("22222222-1111-1111-1111-111111111111")},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Field{
				Association: tt.fields.Association,
				t:           tt.fields.t,
				Multiple:    tt.fields.Multiple,
			}
			got, err := f.Slice(tt.args.values)
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

func TestField_setAssociationsDefaults(t *testing.T) {
	type fields struct {
		Title       string
		Column      string
		Code        string
		Model       *Model
		Association *Association
	}
	type args struct {
		r map[string]*Model
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantField *Field
		wantPanic bool
	}{
		{
			name: "panic",
			fields: fields{
				Association: &Association{Model: &Model{}},
			},
			args: args{
				r: map[string]*Model{},
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("setAssociationsDefaults() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			f := &Field{
				Title:       tt.fields.Title,
				Column:      tt.fields.Column,
				Code:        tt.fields.Code,
				Model:       tt.fields.Model,
				Association: tt.fields.Association,
			}
			f.setAssociationsDefaults(tt.args.r)

			if !reflect.DeepEqual(f, tt.wantField) {
				t.Errorf("many2manyFill() gotField = %v, wantField %v", f, tt.wantField)
			}
		})
	}
}

func TestRawType(t *testing.T) {
	type args struct {
		elem reflect.Type
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RawType(tt.args.elem); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RawType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFields_GetByCode(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name string
		f    Fields
		args args
		want *Field
	}{
		{
			name: "found",
			f:    Fields{{Code: "code"}},
			args: args{
				code: "code",
			},
			want: &Field{Code: "code"},
		},
		{
			name: "not found",
			f:    Fields{},
			args: args{
				code: "123",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.GetByCode(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
