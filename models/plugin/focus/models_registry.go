package focus

import (
	"golang.org/x/exp/maps"
	"reflect"
	"sort"
)

type ModelsRegistry struct {
	registered   map[string]*Model
	supportMedia bool
}

// NewModelsRegistry конструктор
func NewModelsRegistry(supportMedia bool) *ModelsRegistry {
	return &ModelsRegistry{
		registered:   make(map[string]*Model),
		supportMedia: supportMedia,
	}
}

// GetModel получение модели по коду
func (r ModelsRegistry) GetModel(code string) *Model {
	model, ok := r.registered[code]
	if !ok {
		return nil
	}
	return model
}

// ListModels получение всех зарегистрированных моделей
func (r ModelsRegistry) ListModels() []*Model {
	return maps.Values(r.registered)
}

// Register регистрация моделей
func (r *ModelsRegistry) Register(items ...any) {
	for _, task := range items {
		t := RawType(reflect.TypeOf(task))
		r.NewModel(t)
	}

	for _, model := range r.registered {
		fields := model.Fields
		// заполняем поля модели значениями по-умолчанию
		for i := range model.Fields {
			fields[i].setAssociationsDefaults(r.registered)
		}
		// упорядочиваем поля модели
		sort.Slice(fields, func(i, j int) bool {
			if fields[i].Position == fields[j].Position {
				first, _ := model.t.FieldByName(fields[i].name)
				second, _ := model.t.FieldByName(fields[j].name)
				return first.Index[0] < second.Index[0]
			}
			if fields[i].Position == 0 {
				return false
			}
			if fields[j].Position == 0 {
				return true
			}
			return fields[i].Position < fields[j].Position
		})
	}
}

// ModelsRegistry.NewModel сканирование модели
func (r *ModelsRegistry) NewModel(t reflect.Type) {
	if t.Kind() != reflect.Struct {
		panic("wrong model specified, expected struct")
	}
	// если модель уже зарегистрирована - выход
	if _, ok := r.registered[modelCode(t)]; ok { // если модель уже зарегистрирована - пропускаем ее
		return
	}
	// создаем новую модель
	model := NewModel(t)
	r.registered[model.Code] = model

	// проверяем ассоциации модели, если не зарегистрированы - добавляем в очередь
	for _, field := range model.Fields {
		if field.Association != nil {
			if r.GetModel(field.Association.modelCode) == nil {
				r.NewModel(field.RawType())
			}
			assocModel := r.GetModel(field.Association.modelCode)
			field.Association.Model = assocModel
		}
	}

	return
}
