package focus

import (
	"github.com/aeroideaservices/focus/services/formatting/strings"
	"github.com/pkg/errors"
	"reflect"
)

type Model struct {
	TableName  string `json:"tableName"`  // Название таблицы
	Code       string `json:"code"`       // Код модели
	Title      string `json:"name"`       // Название модели для отображения
	PrimaryKey *Field `json:"primaryKey"` // Первичный ключ
	Fields     Fields `json:"fields"`     // Поля модели, кроме первичного ключа

	name string       // Системное название модели
	t    reflect.Type // t тип элементов модели
}

type Focusable interface {
	TableName() string
	ModelTitle() string
}

var focusableType = reflect.TypeOf((*Focusable)(nil)).Elem()

// modelCode получает код модели для использования в запросах
func modelCode(t reflect.Type) string {
	if ok := t.Implements(focusableType); !ok {
		panic("model must implement the Focusable interface")
	}

	elem := reflect.New(t).Elem().Interface().(Focusable)
	modelCode := strings.SnakeToDashedCase(elem.TableName())

	return modelCode
}

// NewModel конструктор модели
func NewModel(t reflect.Type) *Model {
	if !t.Implements(focusableType) {
		panic("model must implements Focusable")
	}
	elem := reflect.New(t).Elem().Interface().(Focusable)
	m := &Model{
		TableName: elem.TableName(),
		Code:      modelCode(t),
		Title:     elem.ModelTitle(),
		name:      t.Name(),
		t:         t,
	}
	for _, structField := range reflect.VisibleFields(m.t) {
		field := &Field{Model: m}
		field.Scan(structField)
	}

	// если первичный ключ не был определен, ищем поле с названием Id или ID
	if m.PrimaryKey == nil {
		panic("primary key must be specified")
	}

	return m
}

func (m Model) Type() reflect.Type {
	return m.t
}

// Model.NewElement возвращает указатель на элемент модели
func (m Model) NewElement(fieldsMap map[string]any, filter func(field *Field) bool) (modelElement any, err error) {
	elem := reflect.New(m.t)
	if fieldsMap == nil {
		return elem.Interface(), nil
	}

	for _, field := range m.Fields {
		if value, ok := fieldsMap[field.Code]; ok && (filter == nil || filter(field)) {
			fv, err := field.NewValue(value)
			if err != nil {
				return nil, err
			}
			elem.Elem().FieldByName(field.name).Set(reflect.ValueOf(fv))
		}
	}

	return elem.Interface(), nil
}

// Model.ElementsSlice получение среза элементов
func (m Model) ElementsSlice(modelElementsMaps []map[string]any) (any, error) {
	res := reflect.MakeSlice(reflect.SliceOf(m.t), 0, len(modelElementsMaps))
	for _, elementMap := range modelElementsMaps {
		elem, err := m.NewElement(elementMap, nil)
		if err != nil {
			return nil, err
		}
		res = reflect.Append(res, reflect.ValueOf(elem).Elem())
	}
	return res.Interface(), nil
}

// ElementToMap преобразование элемента модели в карту
func (m Model) ElementToMap(modelElement any, filter func(field Field) bool) (model map[string]any, err error) {
	data, err := json.Marshal(modelElement)
	if err != nil {
		return nil, err
	}

	fieldsMap := make(map[string]any)
	err = json.Unmarshal(data, &fieldsMap)
	if err != nil {
		return nil, err
	}

	for key, value := range fieldsMap {
		field := m.Fields.GetByCode(key)
		if field == nil || (filter != nil && !filter(*field)) {
			delete(fieldsMap, key)
			continue
		}

		if field.Association != nil || field.IsMedia {
			if value == nil {
				continue
			}
			var pkCode string
			if field.IsMedia {
				pkCode = "id"
			} else {
				pkCode = field.Association.Model.PrimaryKey.Code
			}
			if field.Multiple {
				val := reflect.ValueOf(value)
				var pks []any
				for i := 0; i < val.Len(); i++ {
					pk := val.Index(i).Elem().MapIndex(reflect.ValueOf(pkCode)).Interface()
					pks = append(pks, pk)
				}
				fieldsMap[key] = pks
			} else {
				pk := reflect.ValueOf(value).MapIndex(reflect.ValueOf(pkCode)).Interface()
				fieldsMap[key] = pk
			}
		}
	}

	return fieldsMap, nil
}

// Model.UpdateElement проставляет поля из source в destination
func (m Model) UpdateElement(dst any, src any, filter func(field *Field) bool) error {
	dstV := reflect.ValueOf(dst).Elem()
	srcV := reflect.ValueOf(src).Elem()
	if dstV.Type() != srcV.Type() {
		return errors.New("given elements of different types")
	}
	for _, field := range m.Fields {
		if filter != nil && !filter(field) {
			continue
		}
		dstV.FieldByName(field.name).Set(srcV.FieldByName(field.name))
	}

	return nil
}

// GetPKs получение первичных ключей
func GetPKs(obj any, pkName string) ([]any, error) {
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		if value.IsZero() {
			return nil, nil
		}
		return GetPKs(value.Elem().Interface(), pkName)
	case reflect.Slice:
		var res []any
		for i := 0; i < value.Len(); i++ {
			v, err := GetPKs(value.Index(i).Interface(), pkName)
			if err != nil {
				return nil, err
			}
			res = append(res, v...)
		}
		return res, nil
	case reflect.Struct:
		field := value.FieldByName(pkName)
		if !field.IsValid() {
			return nil, errors.New("wrong pKey name given")
		}
		pk := value.FieldByName(pkName).Interface()
		return []any{pk}, nil
	default:
		return nil, errors.New("wrong object type given")
	}
}
