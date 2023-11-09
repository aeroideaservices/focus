package actions

import (
	"context"
	"github.com/aeroideaservices/focus/configurations/plugin/entity"
	"github.com/aeroideaservices/focus/services/callbacks"
	"strconv"
	"time"

	"github.com/aeroideaservices/focus/services/errors"
	"github.com/google/uuid"
)

const IgnoreLimitValue = -1

// Options сервис работы с настройками
type Options struct {
	callbacks.Callbacks
	confRepository ConfigurationsRepository
	optRepository  OptionsRepository
	mediaService   MediaService
}

// NewOptions конструктор
func NewOptions(
	confRepository ConfigurationsRepository,
	optRepository OptionsRepository,
	mediaService MediaService,
	callbacks callbacks.Callbacks,
) *Options {
	return &Options{
		confRepository: confRepository,
		optRepository:  optRepository,
		mediaService:   mediaService,
		Callbacks:      callbacks,
	}
}

// Get получение настройки по id
func (o Options) Get(ctx context.Context, action GetOption) (*entity.Option, error) {
	if hasConf := o.confRepository.Has(ctx, action.ConfId); !hasConf {
		return nil, ErrConfNotFound
	}

	opt, err := o.optRepository.Get(ctx, action.Id)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting option by id")
	}
	if opt.ConfId != action.ConfId {
		return nil, ErrOptLinkedToAnotherConf
	}

	return opt, nil
}

// Create создание настройки
func (o Options) Create(ctx context.Context, action CreateOption) (*uuid.UUID, error) {
	if hasConf := o.confRepository.Has(ctx, action.ConfId); !hasConf {
		return nil, ErrConfNotFound
	}

	hasOpt := o.optRepository.HasByCode(ctx, action.ConfId, action.Code)
	if hasOpt {
		return nil, ErrOptAlreadyExists
	}

	newId := uuid.New()
	option := entity.Option{
		Id:     newId,
		ConfId: action.ConfId,
		Code:   action.Code,
		Name:   action.Name,
		Type:   action.Type,
	}
	if err := o.optRepository.Create(ctx, option); err != nil {
		return nil, errors.NoType.Wrap(err, "error creating option")
	}

	o.GoAfterCreate(option.Id)

	return &newId, nil
}

// Update обновление настройки
func (o Options) Update(ctx context.Context, action UpdateOption) error {
	if hasConf := o.confRepository.Has(ctx, action.ConfId); !hasConf {
		return ErrConfNotFound
	}

	option, err := o.optRepository.Get(ctx, action.Id)
	if err != nil {
		return errors.NoType.Wrap(err, "error getting option by id")
	}
	if option.ConfId != action.ConfId {
		return ErrOptLinkedToAnotherConf
	}

	if action.Type != option.Type {
		return ErrFieldNotUpdatable.T("field-not-updatable", "Тип")
	}
	if action.Code != option.Code {
		return ErrFieldNotUpdatable.T("field-not-updatable", "Символьный код")
	}

	updated := &entity.Option{}
	*updated = *option
	updated.Name = action.Name

	if err := o.optRepository.Update(ctx, *updated); err != nil {
		return errors.NoType.Wrap(err, "error updating option")
	}

	o.GoAfterUpdate(option.Id)

	return nil
}

// Delete удаление настройки
func (o Options) Delete(ctx context.Context, action GetOption) error {
	if hasConf := o.confRepository.Has(ctx, action.ConfId); !hasConf {
		return ErrConfNotFound
	}

	if hasOpt := o.optRepository.Has(ctx, action.ConfId, action.Id); !hasOpt {
		return ErrOptNotFound
	}

	if err := o.optRepository.Delete(ctx, action.Id); err != nil {
		return errors.NoType.Wrap(err, "error deleting option")
	}

	o.GoAfterDelete(action.Id)

	return nil
}

// List получение списка настроек
func (o Options) List(ctx context.Context, action ListOptions) (*entity.OptionsList, error) {
	if hasConf := o.confRepository.Has(ctx, action.ConfId); !hasConf {
		return nil, ErrConfNotFound
	}

	filter := OptionsListFilter{
		Offset: action.Offset,
		Limit:  action.Limit,
		Sort:   action.Sort,
		Order:  action.Order,
		Filter: OptionsFilter{ConfId: action.ConfId},
	}

	total, err := o.optRepository.Count(ctx, filter.Filter)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error counting options")
	}

	items, err := o.optRepository.List(ctx, filter)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error listing options")
	}

	optList := &entity.OptionsList{
		Total: total,
		Items: items,
	}

	return optList, nil
}

// ListPreviews получение списка превью настроек
func (o Options) ListPreviews(ctx context.Context, action ListOptionsPreviews) ([]entity.OptionShort, error) {
	if hasConf := o.confRepository.HasByCode(ctx, action.ConfCode); !hasConf {
		return nil, ErrConfNotFound
	}

	filter := OptionsListFilter{
		Limit: IgnoreLimitValue,
		Filter: OptionsFilter{
			ConfCode: action.ConfCode,
			OptCodes: action.OptCodes,
		},
	}
	optShorts, err := o.optRepository.ListPreviews(ctx, filter)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error listing options previews")
	}

	return optShorts, nil
}

// UpdateList обновление списка настроек
func (o Options) UpdateList(ctx context.Context, action UpdateOptionsList) error {
	optCodes := make([]string, len(action.Items))
	optShortsMap := make(map[string]string)
	for i, optShort := range action.Items {
		optCodes[i] = optShort.Code
		optShortsMap[optShort.Code] = optShort.Value
	}
	opts, err := o.optRepository.List(ctx, OptionsListFilter{
		Limit: IgnoreLimitValue,
		Filter: OptionsFilter{
			ConfId:   action.ConfId,
			OptCodes: optCodes,
		},
	})
	if err != nil {
		return errors.NoType.Wrap(err, "error listing options")
	}
	if len(opts) != len(action.Items) {
		return ErrOptNotFound
	}

	for _, opt := range opts {
		opt.Value = optShortsMap[opt.Code]

		err = o.checkValueType(ctx, opt.Type, opt.Value)
		if err != nil {
			return errors.BadRequest.Wrap(err, "error checking opt value type")
		}

		err = o.optRepository.Update(ctx, opt)
		if err != nil {
			return errors.NoType.Wrap(err, "error updating option")
		}
		o.GoAfterUpdate(opt.Id)
	}

	return nil
}

func (o Options) checkValueType(ctx context.Context, optType string, optValue string) error {
	switch optType {
	case "string", "text":
		return nil
	case "integer":
		_, err := strconv.Atoi(optValue)
		return err
	case "checkbox":
		_, err := strconv.ParseBool(optValue)
		return err
	case "file", "image":
		if o.mediaService == nil {
			return ErrMediaPluginIsNotImported
		}

		if optValue != "" {
			mediaId, err := uuid.Parse(optValue)
			if err != nil {
				return errors.BadRequest.Wrap(err, "value must be a valid uuid or nil")
			}
			err = o.mediaService.CheckIds(ctx, mediaId)
			if err != nil {
				return errors.BadRequest.Wrap(err, "error checking media ids")
			}
		}
	case "datetime", "date":
		if optValue == "" {
			return nil
		}
		_, err := time.Parse("2006-01-02T15:04:05Z", optValue)
		return err
	}

	return nil
}
