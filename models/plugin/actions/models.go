package actions

import (
	"context"
	"sort"
	"strings"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"

	"github.com/aeroideaservices/focus/models/plugin/focus"
	"github.com/aeroideaservices/focus/models/plugin/form"
	"github.com/aeroideaservices/focus/services/errors"
)

// Models сервис работы с моделями
type Models struct {
	modelsRegistry     ModelsRegistry
	repositoryResolver RepositoryResolver
	selectRequest      form.Request
}

// NewModels конструктор
func NewModels(
	modelsRegistry ModelsRegistry,
	repositoryResolver RepositoryResolver,
	selectRequest form.Request,
) *Models {
	return &Models{
		modelsRegistry:     modelsRegistry,
		repositoryResolver: repositoryResolver,
		selectRequest:      selectRequest,
	}
}

// List получение списка моделей
func (s Models) List(action ListModels) ModelsList {
	models := s.modelsRegistry.ListModels()

	asc := action.Order != "desc"
	sort.Slice(models, func(i, j int) bool {
		if action.Sort == "code" {
			return (models[i].Code < models[j].Code) == asc
		}
		return (models[i].Title < models[j].Title) == asc
	})

	if action.Offset >= len(models) {
		return ModelsList{
			Items: []ModelShort{},
			Total: len(models),
		}
	}

	res := make([]ModelShort, len(models))
	for i, model := range models {
		res[i] = ModelShort{
			Code:  model.Code,
			Title: model.Title,
		}
	}

	return ModelsList{
		Items: res[action.Offset:min(action.Limit+action.Offset, len(res))],
		Total: len(res),
	}
}

// Get получение описания модели
func (s Models) Get(_ context.Context, action GetModel) (*ModelDescription, error) {
	model := s.modelsRegistry.GetModel(action.ModelCode)
	if model == nil {
		return nil, errors.NotFound.Newf("model with code \"%s\" not found", action.ModelCode)
	}

	return &ModelDescription{
		Code:           model.Code,
		Title:          model.Title,
		IdentifierCode: model.PrimaryKey.Code,
		Views: ModelViews{
			Create: s.createView(model),
			Update: s.updateView(model),
			Filter: s.filterView(model),
			List:   s.listView(model),
		},
	}, nil
}

// ListFieldValues получение значений поля модели
func (s Models) ListFieldValues(ctx context.Context, action ListFieldValues) (*FieldValuesList, error) {
	model := s.modelsRegistry.GetModel(action.ModelCode)
	if model == nil {
		return nil, errors.NotFound.Newf("model with code \"%s\" not found", action.ModelCode)
	}

	// ищем поле по всем обычным полям + первичный ключ
	field := model.Fields.GetByCode(action.FieldCode)
	if field == nil {
		return nil, errors.NotFound.Newf("field \"%s\" not found", action.FieldCode)
	}

	// для поля типа время получение значений недоступно
	if field.IsTime {
		return nil, errors.BadRequest.Newf("field \"%s\" is time type, cannot get field values", field.Code)
	}

	repository := s.repositoryResolver.Resolve(model.Code)
	if repository == nil {
		return nil, errors.NoType.New("cannot resolve repository")
	}
	items, err := repository.ListFieldValues(ctx, action)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error listing field values")
	}
	count, err := repository.CountFieldValues(ctx, action.FieldCode, action.Query)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error counting field values")
	}

	return &FieldValuesList{
		Total: count,
		Items: items,
	}, nil
}

// createView получение описания отображения формы создания элемента модели
func (s Models) createView(model *focus.Model) EditView {
	var formFields []FormField
	for _, field := range model.Fields {
		formField := FormField{
			Code:     field.Code,
			Title:    field.Title,
			Type:     field.View,
			Multiple: field.Multiple,
			Sortable: field.Association != nil && field.Association.JoinSort != "",
			Block:    field.Block,
			Extra:    field.ViewExtra,
			Hidden:   slices.Contains(field.Hidden, focus.CreateView),
			Disabled: slices.Contains(field.Disabled, focus.CreateView),
		}

		if field.FloatProperties != nil {
			formField.Step = field.Step
			formField.Precision = field.Precision
		}

		formFields = append(formFields, formField)
	}

	validation := focus.GetValidationRules(*model,
		func(field *focus.Field) bool { return !slices.Contains(field.Disabled, focus.CreateView) })

	return EditView{
		FormFields: formFields,
		Validation: validation,
	}
}

// updateView получение описания отображения формы обновления элемента модели
func (s Models) updateView(model *focus.Model) EditView {
	var formFields []FormField
	for _, field := range model.Fields {
		formField := FormField{
			Code:     field.Code,
			Title:    field.Title,
			Type:     field.View,
			Multiple: field.Multiple,
			Sortable: field.Association != nil && field.Association.JoinSort != "",
			Block:    field.Block,
			Extra:    field.ViewExtra,
			Hidden:   slices.Contains(field.Hidden, focus.UpdateView),
			Disabled: slices.Contains(field.Disabled, focus.UpdateView),
		}

		if field.FloatProperties != nil {
			formField.Step = field.Step
			formField.Precision = field.Precision
		}

		formFields = append(formFields, formField)
	}

	validation := focus.GetValidationRules(*model,
		func(field *focus.Field) bool { return !slices.Contains(field.Disabled, focus.UpdateView) })

	return EditView{
		FormFields: formFields,
		Validation: validation,
	}
}

// filterView получение описания отображения формы создания элемента модели
func (s Models) filterView(model *focus.Model) FormView {
	fv := FormView{}
	for _, field := range model.Fields {
		// пропускаем поля, которые скрыты в списке элементов модели
		if !field.Filterable {
			continue
		}

		formField := FormField{
			Code:     field.Code,
			Title:    field.Title,
			Multiple: true,
			Extra:    field.ViewExtra,
		}
		switch field.View {
		case form.None, form.Textarea, form.Media, form.Wysiwyg, form.EditorJs:
			continue
		case form.DateTimePicker, form.DatePickerInput, form.TimePicker:
			formField.Type = field.View
			formField.Extra = form.ViewExtras{
				"range": true,
			}
		case form.Checkbox:
			formField.Type = form.Select
			formField.Extra = form.ViewExtras{
				"selectData": form.SelectData{
					{Label: "Да", Value: true},
					{Label: "Нет", Value: false},
				},
			}
		case form.IntInput, form.UintInput, form.FloatInput, form.Rating, form.TextInput, form.PhoneInput, form.EmailInput:
			formField.Type = form.Select

			request := s.selectRequest
			request.URI = strings.ReplaceAll(request.URI, "{model-code}", field.Model.Code)
			request.URI = strings.ReplaceAll(request.URI, "{field-code}", field.Code)
			formField.Extra = form.ViewExtras{"request": request}
		default:
			if field.ViewExtra != nil {
				formField.Type = form.Select
			}
		}

		fv.FormFields = append(fv.FormFields, formField)
	}

	return fv
}

// listView получение описания отображения списка элементов модели
func (s Models) listView(model *focus.Model) ListView {
	lv := ListView{}
	for _, field := range model.Fields {
		// пропускаем поля, которые скрыты в списке элементов модели, а так же медиа и ассоциации
		if slices.Contains(field.Hidden, focus.ListView) || field.IsMedia || field.Association != nil {
			continue
		}

		listField := ListField{
			Code:     field.Code,
			Title:    field.Title,
			Sortable: field.Sortable,
			IsTime:   field.IsTime,
		}
		lv.Fields = append(lv.Fields, listField)
	}

	return lv
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
