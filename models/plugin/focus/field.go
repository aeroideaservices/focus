package focus

import (
	"errors"
	"github.com/aeroideaservices/focus/models/plugin/form"
	"golang.org/x/exp/slices"
	"reflect"
	"strings"
)

type FieldType string

type view int

const (
	CreateView view = iota + 1
	UpdateView
	ListView
)

type Fields []*Field

func (f Fields) GetByCode(code string) *Field {
	idx := slices.IndexFunc(f, func(field *Field) bool {
		return field.Code == code
	})

	if idx < 0 {
		return nil
	}
	return f[idx]
}

type Field struct {
	Title            string         // Title Название поля (как оно будет отображаться)
	Column           string         // Column Название колонки в таблице (в базе)
	Code             string         // Code код поля
	Sortable         bool           // Sortable Участвует в сортировке или нет (в списке элементов)
	Filterable       bool           // Filterable Участвует в фильтрации или нет (в списке элементов)
	Position         int            // Position Порядок отображения в списке и карточке
	Block            string         // Block Блок, в котором располагается поле
	IsUnique         bool           // IsUnique Уникальность
	IsMedia          bool           // IsMedia поле является ассоциацией к медиа
	IsTime           bool           // IsTime поле является датой/временем
	Hidden           []view         // Hidden описывает, где поле не отображается
	Disabled         []view         // Disabled описывает, где поле не доступно для редактирования
	View             form.FieldType // View тип отображения
	ViewExtra        map[string]any // ViewExtra доп параметры отображения
	Multiple         bool           // Multiple множественное
	Model            *Model         // Model Описание модели поля
	*FloatProperties                // FloatProperties настройки для типа float

	Association *Association // Association Описание ассоциации

	primaryKey bool         // primaryKey является первичным ключом
	name       string       // name Название поля в модели
	t          reflect.Type // t Структура поля
}

type FloatProperties struct {
	Step      float64
	Precision int
}

// Name получение наименования поля в объекте модели
func (f Field) Name() string {
	return f.name
}

// Scan заполнение поля значениями на основе структуры поля
func (f *Field) Scan(structField reflect.StructField) {
	// получаем тэг в виде "name:users"
	fieldTag, hasTags := structField.Tag.Lookup("focus")
	// игнорируем поле, если в тэг проставлен "-"
	if fieldTag == `-` {
		return
	}

	f.Model.Fields = append(f.Model.Fields, f)
	f.name = structField.Name
	f.t = structField.Type
	// Заполняем поле модели
	if hasTags {
		f.ScanTag(fieldTag)
	}
	// если тип поля - ассоциация
	if f.Association != nil {
		f.Association.modelCode = modelCode(f.RawType()) // получаем код модели
	}
}

// ScanTag заполнение поля значениями на основе тегов
func (f *Field) ScanTag(tagString string) {
	tags := make(map[string]string)
	// Теги приходят в формате "тэг:значение тега;другойТег:другое значение;тегБезЗначения;..."
	// Разбиваем теги на ключ-значение
	tagStrings := strings.Split(tagString, ";")
	for _, tagString := range tagStrings {
		separated := strings.SplitN(tagString, ":", 2)
		key := separated[0]
		value := ""
		if len(separated) > 1 {
			value = separated[1]
		}
		tags[key] = value
	}

	// Для каждого существующего тега заполняем поле модели определенными значениями
	for _, tag := range commonTags {
		if value, ok := tags[tag.Code]; ok {
			tag.Fill(f, value)
		} else if tag.Default != nil {
			tag.Default(f)
		}
	}

	// если тип поля - модель, заполняем также информацию об ассоциации
	if f.RawType().Implements(focusableType) {
		if f.Association == nil {
			f.Association = &Association{}
		}
		for _, tag := range associationTags {
			if value, ok := tags[tag.Code]; ok {
				tag.Fill(f, value)
			}
		}
	}
}

// setAssociationsDefaults проставляет ассоциированную модель и параметры ассоциации по-умолчанию
func (f *Field) setAssociationsDefaults(r map[string]*Model) {
	if f.Association != nil { // если поле - это ассоциация на другую модель
		assocModel := r[f.Association.modelCode]
		if assocModel == nil { // если модель ассоциации не зарегистрирована: panic
			panic("association type specified incorrectly")
		}
		f.Association.Model = assocModel

		for _, tag := range associationTags {
			if tag.Default != nil {
				tag.Default(f)
			}
		}
	}
}

// NewValue приводит переданное значение к типу поля
func (f Field) NewValue(value any) (any, error) {
	// если передан nil
	if value == nil {
		return reflect.Zero(f.t).Interface(), nil
	}

	// если передано значение того же типа, что и тип поля
	if reflect.TypeOf(value) == f.t {
		return value, nil
	}

	// если передан указатель
	if v := reflect.ValueOf(value); v.Kind() == reflect.Ptr {
		return f.NewValue(v.Elem().Interface())
	}

	if f.Association != nil || f.IsMedia {
		var pkCode string
		if f.Association != nil {
			pkCode = f.Association.Model.PrimaryKey.Code
		} else {
			pkCode = "id"
		}
		if f.Multiple {
			if reflect.TypeOf(value).Kind() != reflect.Slice {
				return nil, errors.New("value type must be a slice")
			}
			var values []map[string]any
			val := reflect.ValueOf(value)
			for i := 0; i < val.Len(); i++ {
				v := val.Index(i).Interface()
				values = append(values, map[string]any{pkCode: v})
			}
			value = values
		} else {
			value = map[string]any{pkCode: value}
		}
	}

	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	field := reflect.New(f.t).Interface()
	err = json.Unmarshal(data, field)
	if err != nil {
		return nil, err
	}

	return reflect.ValueOf(field).Elem().Interface(), nil
}

// Slice приводит переданный срез значений к срезу типа поля
func (f Field) Slice(values []any) ([]any, error) {
	// если поле множественное - возвращаем элементы, преобразованные в значения поля
	if f.Multiple {
		fieldValue, err := f.NewValue(values)
		if err != nil {
			return nil, err
		}
		fvv := reflect.ValueOf(fieldValue)
		if fvv.Kind() != reflect.Slice {
			return nil, errors.New("123") // todo
		}
		var fieldValues []any
		for i := 0; i < fvv.Len(); i++ {
			fieldValues = append(fieldValues, fvv.Index(i).Interface())
		}
		return fieldValues, nil
	}

	var res []any
	for _, value := range values {
		fieldValue, err := f.NewValue(value)
		if err != nil {
			return nil, err
		}
		res = append(res, fieldValue)
	}

	return res, nil
}

// Field.RawKind извлекает вид элемента
func (f Field) RawKind() reflect.Kind {
	t := RawType(f.t)
	return t.Kind()
}

// Field.RawType извлекает вид элемента
func (f Field) RawType() reflect.Type {
	return RawType(f.t)
}

// RawType извлекает тип элемента из указателей и срезов
func RawType(elem reflect.Type) reflect.Type {
	switch elem.Kind() {
	case reflect.Slice, reflect.Array, reflect.Pointer:
		return RawType(elem.Elem())
	default:
		return elem
	}
}
