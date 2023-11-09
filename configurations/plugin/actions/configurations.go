package actions

import (
	"context"
	"github.com/aeroideaservices/focus/configurations/plugin/entity"
	"github.com/aeroideaservices/focus/services/callbacks"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/google/uuid"
)

// Configurations сервис работы с конфигурациями
type Configurations struct {
	callbacks.Callbacks
	confRepository ConfigurationsRepository
}

// NewConfigurations конструктор
func NewConfigurations(repository ConfigurationsRepository,
	callbacks callbacks.Callbacks,
) *Configurations {
	return &Configurations{
		confRepository: repository,
		Callbacks:      callbacks,
	}
}

// Create создание новой конфигурации
func (c Configurations) Create(ctx context.Context, action CreateConfiguration) (uuid.UUID, error) {
	if c.confRepository.HasByCode(ctx, action.Code) {
		return uuid.Nil, ErrConfAlreadyExists
	}

	newId := uuid.New()
	configuration := entity.Configuration{Id: newId, Code: action.Code, Name: action.Name}
	err := c.confRepository.Create(ctx, configuration)
	if err != nil {
		return uuid.Nil, err
	}

	c.GoAfterCreate(configuration.Id)

	return newId, nil
}

// Get получение конфигурации по id
func (c Configurations) Get(ctx context.Context, action GetConfiguration) (*entity.Configuration, error) {
	conf, err := c.confRepository.Get(ctx, action.Id)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting configuration")
	}

	return conf, nil
}

// Update обновление конфигурации
func (c Configurations) Update(ctx context.Context, action UpdateConfiguration) error {
	configuration, err := c.confRepository.Get(ctx, action.Id)
	if err != nil {
		return err
	}

	if configuration.Code != action.Code {
		return ErrFieldNotUpdatable.T("field-not-updatable", "Символьный код")
	}

	configuration.Name = action.Name
	if err := c.confRepository.Update(ctx, configuration); err != nil {
		return errors.NoType.Wrap(err, "error updating configuration")
	}

	c.GoAfterUpdate(configuration.Id)

	return nil
}

// Delete удаление конфигурации
func (c Configurations) Delete(ctx context.Context, action GetConfiguration) error {
	if hasConf := c.confRepository.Has(ctx, action.Id); !hasConf {
		return ErrConfNotFound
	}

	if err := c.confRepository.Delete(ctx, action.Id); err != nil {
		return errors.NoType.Wrap(err, "error deleting configuration")
	}

	c.GoAfterDelete(action.Id)

	return nil
}

// List получение списка конфигураций
func (c Configurations) List(ctx context.Context, dto ListConfigurations) (*entity.ConfigurationsList, error) {
	filter := ConfigurationsListFilter{
		Offset: dto.Offset,
		Limit:  dto.Limit,
		Sort:   dto.Sort,
		Order:  dto.Order,
		Filter: ConfigurationsFilter{Query: dto.Query},
	}
	items, err := c.confRepository.List(ctx, filter)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error listing configurations")
	}

	count, err := c.confRepository.Count(ctx, filter.Filter)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error counting configurations")
	}

	confList := &entity.ConfigurationsList{
		Total: count,
		Items: items,
	}
	return confList, nil
}
