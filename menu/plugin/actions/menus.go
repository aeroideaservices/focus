package actions

import (
	"context"
	"github.com/aeroideaservices/focus/menu/plugin/entity"
	"github.com/aeroideaservices/focus/services/callbacks"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/google/uuid"
)

// MenusCallbacks колбэки сервиса работы с меню
type MenusCallbacks struct {
	AfterCreate func(menu *entity.Menu)
	AfterUpdate func(oldMenu, newMenu *entity.Menu)
}

// Menus сервис работы с меню
type Menus struct {
	callbacks.Callbacks
	menuRepository MenuRepository
}

// NewMenus конструктор
func NewMenus(
	menuRepository MenuRepository,
	callbacks callbacks.Callbacks,
) *Menus {
	return &Menus{
		menuRepository: menuRepository,
		Callbacks:      callbacks,
	}
}

// List получение списка меню
func (m Menus) List(ctx context.Context, action ListMenus) (*MenusList, error) {
	items, err := m.menuRepository.List(ctx, MenuFilter{
		Offset: action.Offset,
		Limit:  action.Limit,
		Sort:   action.Sort,
		Order:  action.Order,
	})
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error listing menus")
	}

	total, err := m.menuRepository.Count(ctx)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error counting menus")
	}

	return &MenusList{
		Total: total,
		Items: items,
	}, nil
}

// Create создание меню
func (m Menus) Create(ctx context.Context, action CreateMenu) (*uuid.UUID, error) {
	hasMenu, err := m.menuRepository.HasByCode(ctx, action.Code)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error checking that menu exists")
	}
	if hasMenu {
		return nil, ErrMenuAlreadyExists
	}

	menu := entity.Menu{Id: uuid.New(), Code: action.Code, Name: action.Name}
	err = m.menuRepository.Create(ctx, menu)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error creating menu")
	}

	m.GoAfterCreate(menu.Id)

	return &menu.Id, nil
}

// Get получение меня по id
func (m Menus) Get(ctx context.Context, action GetMenu) (*entity.Menu, error) {
	menu, err := m.menuRepository.Get(ctx, action.MenuId)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting menu")
	}

	return menu, nil
}

// GetByCode получение меню по коду
func (m Menus) GetByCode(ctx context.Context, action GetMenuByCode) (*entity.Menu, error) {
	menu, err := m.menuRepository.GetByCode(ctx, action.MenuCode)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting menu by code")
	}

	return menu, nil
}

// Update обновление меню
func (m Menus) Update(ctx context.Context, action UpdateMenu) error {
	menu, err := m.menuRepository.Get(ctx, action.Id)
	if err != nil {
		return errors.NoType.Wrap(err, "error getting menu")
	}
	if menu.Code != action.Code {
		return errors.BadRequest.Wrap(ErrFieldNotUpdatable, "field \"code\" is not updatable").T("field-not-updatable", "Код")
	}

	menu.Code = action.Code
	menu.Name = action.Name
	err = m.menuRepository.Update(ctx, menu)
	if err != nil {
		return errors.NoType.Wrap(err, "error updating menu")
	}

	m.GoAfterUpdate(menu.Id)

	return nil
}

// Delete удаление меню
func (m Menus) Delete(ctx context.Context, action GetMenu) error {
	hasMenu, err := m.menuRepository.Has(ctx, action.MenuId)
	if err != nil {
		return errors.NoType.Wrap(err, "error checking that menu exists")
	}
	if !hasMenu {
		return ErrMenuNotFound
	}

	err = m.menuRepository.Delete(ctx, action.MenuId)
	if err != nil {
		return errors.NoType.Wrap(err, "error deleting menu")
	}

	m.GoAfterDelete(action.MenuId)

	return nil
}
