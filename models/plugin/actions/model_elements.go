package actions

import (
	"context"
	"github.com/aeroideaservices/focus/services/callbacks"
	"reflect"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/aeroideaservices/focus/models/plugin/focus"
	"github.com/aeroideaservices/focus/services/errors"
)

type ModelElements struct {
	modelsRegistry     ModelsRegistry
	repositoryResolver RepositoryResolver
	mediaService       MediaService
	validator          Validator
	callbacks          map[string]callbacks.Callbacks
}

func NewModelElements(
	modelsRegistry ModelsRegistry,
	repositoryResolver RepositoryResolver,
	mediaService MediaService,
	validator Validator,
	callbacks map[string]callbacks.Callbacks,
) *ModelElements {
	return &ModelElements{
		modelsRegistry:     modelsRegistry,
		repositoryResolver: repositoryResolver,
		mediaService:       mediaService,
		validator:          validator,
		callbacks:          callbacks,
	}
}

// List получение списка элементов модели.
func (s ModelElements) List(ctx context.Context, action ListModelElements) (*List, error) {
	model := s.modelsRegistry.GetModel(action.ModelCode)
	if model == nil {
		return nil, errors.NotFound.Newf("model with code \"%s\" not found", action.ModelCode)
	}

	// проверяем, что передано верное поле для сортировки
	if action.Sort != "" {
		if !slices.ContainsFunc(model.Fields, func(field *focus.Field) bool { return field.Code == action.Sort && field.Sortable }) {
			return nil, errors.BadRequest.New("got wrong field code for sort")
		}
	}

	selectFields := action.FieldsCodes
	if len(action.FieldsCodes) == 0 {
		// получаем коды полей, которые нужно вернуть в запросе
		for _, field := range model.Fields {
			if !(slices.Contains(field.Hidden, focus.ListView) || field.IsMedia || field.Association != nil) { // игнорируем скрытые поля, медиа и ассоциации
				selectFields = append(selectFields, field.Code)
			}
		}
	} else {
		// проверяем, все ли переданные поля присутствуют в модели
		for _, fieldCode := range action.FieldsCodes {
			if model.Fields.GetByCode(fieldCode) == nil {
				return nil, errors.BadRequest.New("got wrong field code")
			}
		}
	}

	// всегда возвращаем PK
	if !slices.Contains(selectFields, model.PrimaryKey.Code) {
		selectFields = append(selectFields, model.PrimaryKey.Code)
	}

	// убираем значения полей, которые не могут участвовать в фильтрации
	maps.DeleteFunc(action.Filter, func(key string, values []any) bool {
		return key != model.PrimaryKey.Code && !slices.ContainsFunc(model.Fields, func(field *focus.Field) bool {
			return field.Code == key && field.Filterable
		})
	})

	filter := ListModelElementsQuery{
		ModelCode: action.ModelCode,
		Filter: ModelElementsFilter{
			QueryFilter:  action.ModelElementsQueryFilter,
			FieldsFilter: make(FieldsFilter),
		},
		SelectFields: selectFields,
		Pagination:   action.Pagination,
		OrderBy:      action.OrderBy,
	}

	// преобразуем значения фильтров к нужным типам
	for fieldCode, values := range action.Filter {
		field := model.Fields.GetByCode(fieldCode)
		if field == nil || (!field.Filterable && field.Code != model.PrimaryKey.Code) {
			return nil, errors.BadRequest.New("got wrong field code for filter")
		}

		if field.IsTime && len(values) != 2 {
			return nil, errors.BadRequest.Newf("field \"%s\" must contains two values", fieldCode)
		}

		vs, err := field.Slice(values)
		if err != nil {
			return nil, errors.NoType.Wrap(err, "error converting filter values")
		}

		filter.Filter.FieldsFilter[fieldCode] = vs
	}

	repository := s.repositoryResolver.Resolve(model.Code)
	if repository == nil {
		return nil, errors.NoType.New("cannot resolve repository")
	}

	// получаем элементы модели
	elems, err := repository.List(ctx, filter)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error listing model elements")
	}
	// получаем общее количество элементов
	total, err := repository.Count(ctx, filter.Filter)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error counting model elements")
	}

	// преобразуем элементы модели
	items := make([]map[string]any, len(elems))
	for i := range elems {
		item, err := model.ElementToMap(elems[i], func(field focus.Field) bool {
			return slices.Contains(selectFields, field.Code)
		})
		if err != nil {
			return nil, errors.NoType.Wrap(err, "error converting model element to map")
		}
		items[i] = item
	}

	return &List{
		Items: items,
		Total: total,
	}, nil
}

// Get получение элемента модели по первичному ключу.
func (s ModelElements) Get(ctx context.Context, action GetModelElement) (map[string]any, error) {
	model := s.modelsRegistry.GetModel(action.ModelCode)
	if model == nil {
		return nil, errors.NotFound.Newf("model with code \"%s\" not found", action.ModelCode)
	}

	repository := s.repositoryResolver.Resolve(model.Code)
	if repository == nil {
		return nil, errors.NoType.New("cannot resolve repository")
	}

	pKey, err := model.PrimaryKey.NewValue(action.PKey)
	if err != nil {
		return nil, errors.BadRequest.Wrap(err, "error converting pKey").T("model-element.field.wrong", model.PrimaryKey.Title)
	}

	elem, err := repository.Get(ctx, pKey)
	if errors.GetType(err) == errors.NotFound {
		return nil, errors.NotFound.New("element not found")
	}
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting model element")
	}

	res, err := model.ElementToMap(elem, nil)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error encoding element")
	}

	return res, nil
}

// Create создание нового элемента модели.
func (s ModelElements) Create(ctx context.Context, action CreateModelElement) (pkey any, err error) {
	model := s.modelsRegistry.GetModel(action.ModelCode)
	if model == nil {
		return nil, errors.NotFound.Newf("model with code \"%s\" not found", action.ModelCode)
	}

	// Генерируем новый первичный ключ
	// "На текущий момент это (первичный ключ всех моделей - uuid) требование для всех гошных сервисов". (c) Лид проекта
	action.ModelElement[model.PrimaryKey.Code] = uuid.New().String()

	elem, err := model.NewElement(action.ModelElement, func(field *focus.Field) bool {
		return !slices.Contains(field.Disabled, focus.CreateView) || field == field.Model.PrimaryKey
	})
	if err != nil {
		return nil, errors.BadRequest.Wrap(err, "error parsing model element")
	}

	err = s.validator.Validate(ctx, elem)
	if err != nil {
		return nil, err
	}

	err = s.validate(ctx, model, elem, nil)
	if err != nil {
		return nil, err
	}

	repository := s.repositoryResolver.Resolve(model.Code)
	if repository == nil {
		return nil, errors.NoType.New("cannot resolve repository")
	}

	pkey, err = repository.Create(ctx, elem)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error creating model element")
	}

	s.afterCreate(model.Code, pkey)

	return pkey, nil
}

// Update обновление элемента модели.
func (s ModelElements) Update(ctx context.Context, action UpdateModelElement) error {
	model := s.modelsRegistry.GetModel(action.ModelCode)
	if model == nil {
		return errors.NotFound.Newf("model with code \"%s\" not found", action.ModelCode)
	}

	repository := s.repositoryResolver.Resolve(model.Code)
	if repository == nil {
		return errors.NoType.New("cannot resolve repository")
	}

	pKey, err := model.PrimaryKey.NewValue(action.PKey)
	if err != nil {
		return errors.BadRequest.Wrap(err, "error converting pKey").T("model-element.field.wrong", model.PrimaryKey.Title)
	}

	oldElem, err := repository.Get(ctx, pKey)
	if err != nil {
		return errors.NoType.Wrap(err, "error getting model element")
	}

	elem, err := model.NewElement(action.ModelElement, func(field *focus.Field) bool { return !slices.Contains(field.Disabled, focus.UpdateView) })
	if err != nil {
		return errors.BadRequest.Wrap(err, "error parsing model element")
	}

	fieldsFilter := func(field *focus.Field, fieldValue any) bool {
		if slices.Contains(field.Disabled, focus.UpdateView) {
			return false
		}
		old := reflect.ValueOf(oldElem).Elem().FieldByName(field.Name()).Interface()
		if reflect.DeepEqual(old, fieldValue) {
			return false
		}
		return true
	}
	err = s.validate(ctx, model, elem, fieldsFilter)
	if err != nil {
		return errors.BadRequest.Wrap(err, "validation error")
	}

	// наполняем полученную модель данными из запроса (заполняем только теми полями, которые доступны для обновления)
	err = model.UpdateElement(elem, oldElem, func(field *focus.Field) bool { return slices.Contains(field.Disabled, focus.UpdateView) })
	if err != nil {
		return errors.BadRequest.Wrap(err, "error filling struct")
	}

	err = s.validator.Validate(ctx, elem)
	if err != nil {
		return err
	}

	err = repository.Update(ctx, elem)
	if err != nil {
		return errors.NoType.Wrap(err, "an error occurred while updating model")
	}

	s.afterUpdate(model.Code, pKey)

	return nil
}

// DeleteList удаление нескольких элементов модели.
func (s ModelElements) DeleteList(ctx context.Context, action DeleteModelElements) error {
	model := s.modelsRegistry.GetModel(action.ModelCode)
	if model == nil {
		return errors.NotFound.Newf("model with code \"%s\" not found", action.ModelCode)
	}

	repository := s.repositoryResolver.Resolve(model.Code)
	if repository == nil {
		return errors.NoType.New("cannot resolve repository")
	}

	pKeys, err := model.PrimaryKey.Slice(action.PKeys)
	if err != nil {
		return errors.BadRequest.Wrap(err, "error converting primary keys")
	}

	count, err := repository.Count(ctx, ModelElementsFilter{FieldsFilter: map[string][]any{model.PrimaryKey.Code: pKeys}})
	if err != nil {
		return errors.NoType.Wrap(err, "error while counting model elements")
	}
	if count < int64(len(action.PKeys)) {
		return errors.NotFound.New("model element does not exist")
	}

	err = repository.Delete(ctx, pKeys...)
	if err != nil {
		return errors.NoType.Wrap(err, "error deleting model elements")
	}

	s.afterDelete(model.Code, pKeys...)

	return nil
}

// Delete удаление одного элемента модели.
func (s ModelElements) Delete(ctx context.Context, action DeleteModelElement) error {
	model := s.modelsRegistry.GetModel(action.ModelCode)
	if model == nil {
		return errors.NotFound.Newf("model with code \"%s\" not found", action.ModelCode)
	}

	repository := s.repositoryResolver.Resolve(model.Code)
	if repository == nil {
		return errors.NoType.New("cannot resolve repository")
	}

	pKey, err := model.PrimaryKey.NewValue(action.PKey)
	if err != nil {
		return errors.BadRequest.Wrap(err, "error converting pKey").T("model-element.field.wrong", model.PrimaryKey.Title)
	}

	has, err := repository.Has(ctx, pKey)
	if err != nil {
		return errors.NoType.Wrap(err, "error checking that model element exists")
	}
	if !has {
		return errors.NotFound.New("model element not found")
	}

	err = repository.Delete(ctx, pKey)
	if err != nil {
		return errors.NoType.Wrap(err, "error deleting model elements")
	}

	s.afterDelete(model.Code, pKey)

	return nil
}

// checkUnique проверяет, существует ли другой элемент модели с таким же значением поля.
// Если найдено хотя бы одно совпадение, возвращается ошибка.
func (s ModelElements) checkUnique(ctx context.Context, field *focus.Field, fieldValue any) error {
	repository := s.repositoryResolver.Resolve(field.Model.Code)
	if repository == nil {
		return errors.NoType.New("cannot resolve repository")
	}

	fieldsFilter := map[string][]any{field.Code: {fieldValue}}
	modelEntitiesCount, err := repository.Count(ctx, ModelElementsFilter{FieldsFilter: fieldsFilter})
	if err != nil {
		return errors.NoType.Wrap(err, "error counting model elements")
	}

	// если уже есть другие элементы с этим полем - ошибка
	if modelEntitiesCount > 0 {
		return errModelElementConflict.T("model-element.field.unique", field.Title)
	}

	return nil
}

// checkAssociations проверяет ассоциацию элемента модели на существование.
// Если ассоциируемый элемент модели не найден, возвращается ошибка.
func (s ModelElements) checkAssociations(ctx context.Context, assoc *focus.Association, fv any) error {
	pkName := assoc.Model.PrimaryKey.Name()

	// получаем первичные ключи элемента(-ов) ассоциированной модели
	pk, err := focus.GetPKs(fv, pkName)
	if err != nil {
		return err
	}
	if len(pk) == 0 {
		return nil
	}

	// получаем репозиторий ассоциированной модели
	repo := s.repositoryResolver.Resolve(assoc.Model.Code)
	if repo == nil {
		return errors.NoType.New("wrong repository given")
	}
	count, err := repo.Count(ctx, ModelElementsFilter{FieldsFilter: map[string][]any{assoc.Model.PrimaryKey.Code: pk}})
	if err != nil {
		return err
	}

	// count должно быть равно количеству элементов
	if count < int64(len(pk)) {
		return errors.NotFound.New("model element not found")
	}

	return nil
}

// checkMedias проверяет поле типа медиа.
// Если медиа не найдено, возвращается ошибка.
func (s ModelElements) checkMedias(ctx context.Context, fv any) error {
	// если плагин focus.media не подключен - ошибка
	if s.mediaService == nil {
		return errMediaPluginIsNotImported
	}
	pk, err := focus.GetPKs(fv, "Id")
	if err != nil {
		return err
	}

	var ids []uuid.UUID
	for _, id := range pk {
		ids = append(ids, id.(uuid.UUID))
	}
	if len(ids) == 0 {
		return nil
	}

	err = s.mediaService.CheckIds(ctx, ids...)
	if err != nil {
		return errors.BadRequest.Wrap(err, "validation error")
	}
	return nil
}

// validate проверяет поля элемента модели.
func (s ModelElements) validate(ctx context.Context, model *focus.Model, elem any, filter func(field *focus.Field, fieldValue any) bool) error {
	value := reflect.ValueOf(elem).Elem()
	for _, field := range model.Fields {
		fv := value.FieldByName(field.Name()).Interface()
		if filter != nil && !filter(field, fv) {
			continue
		}
		if field.IsUnique {
			err := s.checkUnique(ctx, field, fv)
			if err != nil {
				return err
			}
		}
		if field.Association != nil {
			err := s.checkAssociations(ctx, field.Association, fv)
			if err != nil {
				return err
			}
		}
		if field.IsMedia {
			// если плагин focus.media не подключен - ошибка
			if s.mediaService == nil {
				return errMediaPluginIsNotImported
			}
			err := s.checkMedias(ctx, fv)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s ModelElements) afterCreate(modelCode string, id any) {
	if callback, ok := s.callbacks[modelCode]; ok {
		go callback.GoAfterCreate(id.(uuid.UUID))
	}
}

func (s ModelElements) afterUpdate(modelCode string, id any) {
	if callback, ok := s.callbacks[modelCode]; ok {
		callback.GoAfterUpdate(id.(uuid.UUID))
	}
}

func (s ModelElements) afterDelete(modelCode string, ids ...any) {
	if callback, ok := s.callbacks[modelCode]; ok {
		var uuids []uuid.UUID
		for _, id := range ids {
			uuids = append(uuids, id.(uuid.UUID))
		}
		callback.GoAfterDelete(uuids...)
	}
}
