package actions

import (
	"context"
	"encoding/json"
	"github.com/aeroideaservices/focus/menu/plugin/entity"
	"github.com/aeroideaservices/focus/services/callbacks"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/google/uuid"
	"strings"
)

// MenuItems сервис работы с элементами меню
type MenuItems struct {
	callbacks.Callbacks
	menuRepository     MenuRepository
	menuItemRepository MenuItemRepository
	maxMenuItemsDepth  int
}

// NewMenuItems конструктор
func NewMenuItems(
	menuRepository MenuRepository,
	menuItemRepository MenuItemRepository,
	maxMenuItemsDepth int,
	callbacks callbacks.Callbacks,
) *MenuItems {
	return &MenuItems{
		menuRepository:     menuRepository,
		menuItemRepository: menuItemRepository,
		maxMenuItemsDepth:  maxMenuItemsDepth,
		Callbacks:          callbacks,
	}
}

// GetTree получение дерева пунктов меню
func (mi MenuItems) GetTree(ctx context.Context, action GetMenuByCode) ([]*MenuItemTree, error) {
	hasMenu, err := mi.menuRepository.HasByCode(ctx, action.MenuCode)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error checking that menu exists")
	}
	if !hasMenu {
		return nil, errors.NoType.Wrapf(ErrMenuNotFound, "menu with code \"%s\" not found", action.MenuCode)
	}

	menuItems, err := mi.menuItemRepository.GetAsTree(ctx, action.MenuCode)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting menu items at tree")
	}

	idItemMap := make(map[uuid.UUID]*MenuItemTree)
	res := make([]*MenuItemTree, 0)
	for _, menuItem := range menuItems {
		item := &MenuItemTree{
			Id:               menuItem.Id,
			ParentMenuItemId: menuItem.ParentMenuItemId,
			Name:             menuItem.Name,
			Url:              mi.joinURL(menuItem.Domain, menuItem.Url),
			Position:         menuItem.Position,
			AdditionalFields: menuItem.AdditionalFields,
		}
		idItemMap[item.Id] = item // формируем карту (id => пункт меню) для дальнейшего поиска по id
		// если у пункта меню нет родительского элемента
		if menuItem.ParentMenuItemId == nil {
			res = append(res, item) // формируем список пунктов меню первого уровня
		}
	}

	// для каждого пункта меню проставляем его дочерние элементы
	for _, item := range menuItems {
		if item.ParentMenuItemId == nil { // пропускаем пункты меню первого уровня
			continue
		}
		if parent, ok := idItemMap[*item.ParentMenuItemId]; ok {
			parent.MenuItems = append(parent.MenuItems, idItemMap[item.Id]) // добавляем пункты меню в родительские элементы
		} else {
			// в случае, если не удалось найти родительский элемент, выбрасываем ошибку
			return nil, errors.NoType.New("one of the menu menuItems has a non-existent parent")
		}
	}

	return res, nil
}

// List получение списка пунктов меню
func (mi MenuItems) List(ctx context.Context, action ListMenuItems) ([]MenuItemPreview, error) {
	hasMenu, err := mi.menuRepository.Has(ctx, action.Filter.MenuId)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error checking that menu exists")
	}
	if !hasMenu {
		return nil, ErrMenuNotFound
	}

	if action.Filter.ParentId != nil {
		hasMenuItem, err := mi.menuItemRepository.Has(ctx, action.Filter.MenuId, *action.Filter.ParentId)
		if err != nil {
			return nil, errors.NoType.Wrap(err, "error checking that menu item exists")
		}
		if !hasMenuItem {
			return nil, ErrMenuItemNotFound
		}
	}

	menuItems, err := mi.menuItemRepository.List(ctx, MenuItemsListFilter{
		Sort:   action.Sort,
		Order:  action.Order,
		Filter: action.Filter,
	})
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error listing menu items")
	}

	var items []MenuItemPreview
	for _, menuItem := range menuItems {
		items = append(items, MenuItemPreview{
			Id:       menuItem.Id,
			Name:     menuItem.Name,
			Url:      mi.joinURL(menuItem.Domain, menuItem.Url),
			Position: menuItem.Position,
		})
	}

	return items, nil
}

// Create создание пункта меню
func (mi MenuItems) Create(ctx context.Context, action CreateMenuItem) (*uuid.UUID, error) {
	hasMenu, err := mi.menuRepository.Has(ctx, action.MenuId)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error checking that menu exists")
	}
	if !hasMenu {
		return nil, ErrMenuNotFound
	}

	if action.ParentMenuItemId != nil {
		hasMenuItem, err := mi.menuItemRepository.Has(ctx, action.MenuId, *action.ParentMenuItemId)
		if err != nil {
			return nil, errors.NoType.Wrap(err, "error checking that menu item exists")
		}
		if !hasMenuItem {
			return nil, ErrMenuItemNotFound
		}

		if mi.maxMenuItemsDepth != 0 {
			depthLevel, err := mi.menuItemRepository.GetDepthLevel(ctx, *action.ParentMenuItemId)
			if err != nil {
				return nil, errors.NoType.Wrap(err, "error getting menu item depth level")
			}
			if depthLevel >= mi.maxMenuItemsDepth {
				return nil, ErrMaxDepthExceeded
			}
		}
	}

	menuItemsCount, err := mi.menuItemRepository.Count(ctx, MenuItemFilter{
		MenuId:   action.MenuId,
		ParentId: action.ParentMenuItemId,
	})
	menuItemId := uuid.New()
	menuItem := entity.MenuItem{
		Id:               menuItemId,
		Name:             action.Name,
		DomainId:         action.DomainId,
		Url:              action.Url,
		Position:         menuItemsCount + 1,
		ParentMenuItemId: action.ParentMenuItemId,
		MenuId:           action.MenuId,
	}
	jsonAF, err := json.Marshal(action.AdditionalFields)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error converting additional fields to json")
	}
	err = menuItem.AdditionalFields.Scan(jsonAF)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error scanning additional fields")
	}

	err = mi.menuItemRepository.Create(ctx, menuItem)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error creating menu item")
	}

	mi.GoAfterCreate(menuItem.Id)

	return &menuItemId, nil
}

// Get получение пункта меню по id
func (mi MenuItems) Get(ctx context.Context, action GetMenuItem) (*entity.MenuItem, error) {
	hasMenu, err := mi.menuRepository.Has(ctx, action.MenuId)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error checking that menu exists")
	}
	if !hasMenu {
		return nil, ErrMenuNotFound
	}

	hasMenuItem, err := mi.menuItemRepository.Has(ctx, action.MenuId, action.MenuItemId)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error checking that menu item exists")
	}
	if !hasMenuItem {
		return nil, ErrMenuItemNotFound
	}

	menuItem, err := mi.menuItemRepository.Get(ctx, action.MenuItemId)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting menu item")
	}

	return menuItem, nil
}

// Update обновление пункта меню
func (mi MenuItems) Update(ctx context.Context, action UpdateMenuItem) error {
	hasMenu, err := mi.menuRepository.Has(ctx, action.MenuId)
	if err != nil {
		return errors.NoType.Wrap(err, "error checking that menu exists")
	}
	if !hasMenu {
		return ErrMenuNotFound
	}

	hasMenuItem, err := mi.menuItemRepository.Has(ctx, action.MenuId, action.MenuItemId)
	if err != nil {
		return errors.NoType.Wrap(err, "error checking that menu item exists")
	}
	if !hasMenuItem {
		return ErrMenuItemNotFound
	}

	menuItem := entity.MenuItem{
		Id:       action.MenuItemId,
		Name:     action.Name,
		DomainId: action.DomainId,
		Url:      action.Url,
	}
	jsonAF, err := json.Marshal(action.AdditionalFields)
	if err != nil {
		return errors.NoType.Wrap(err, "error converting additional fields to json")
	}
	err = menuItem.AdditionalFields.Scan(jsonAF)
	if err != nil {
		return errors.NoType.Wrap(err, "error scanning additional fields")
	}
	err = mi.menuItemRepository.Update(ctx, menuItem)
	if err != nil {
		return errors.NoType.Wrap(err, "error updating menu item")
	}

	mi.GoAfterUpdate(menuItem.Id)

	return nil
}

// Delete удаление пункта меню
func (mi MenuItems) Delete(ctx context.Context, action GetMenuItem) error {
	hasMenu, err := mi.menuRepository.Has(ctx, action.MenuId)
	if err != nil {
		return errors.NoType.Wrap(err, "error checking that menu exists")
	}
	if !hasMenu {
		return ErrMenuNotFound
	}

	hasMenuItem, err := mi.menuItemRepository.Has(ctx, action.MenuId, action.MenuItemId)
	if err != nil {
		return errors.NoType.Wrap(err, "error checking that menu item exists")
	}
	if !hasMenuItem {
		return ErrMenuItemNotFound
	}

	err = mi.menuItemRepository.Delete(ctx, action.MenuItemId)
	if err != nil {
		return errors.NoType.Wrap(err, "error deleting menu item")
	}

	mi.GoAfterDelete(action.MenuItemId)

	return nil
}

// Move перемещение пункта меню
func (mi MenuItems) Move(ctx context.Context, action MoveMenuItem) error {
	hasMenu, err := mi.menuRepository.Has(ctx, action.MenuId)
	if err != nil {
		return errors.NoType.Wrap(err, "error checking that menu exists")
	}
	if !hasMenu {
		return ErrMenuNotFound
	}

	hasMenuItem, err := mi.menuItemRepository.Has(ctx, action.MenuId, action.MenuItemId)
	if err != nil {
		return errors.NoType.Wrap(err, "error checking that menu item exists")
	}
	if !hasMenuItem {
		return ErrMenuItemNotFound
	}

	menuItem, err := mi.menuItemRepository.Get(ctx, action.MenuItemId)
	if err != nil {
		return errors.NoType.Wrap(err, "error getting menu item")
	}

	if menuItem.ParentMenuItemId != action.ParentMenuItemId && action.ParentMenuItemId != nil {
		hasMenuItem, err := mi.menuItemRepository.Has(ctx, action.MenuId, *action.ParentMenuItemId)
		if err != nil {
			return errors.NoType.Wrap(err, "error checking that menu item exists")
		}
		if !hasMenuItem {
			return ErrMenuItemNotFound
		}

		if mi.maxMenuItemsDepth != 0 {
			depthLevel, err := mi.menuItemRepository.GetDepthLevel(ctx, *action.ParentMenuItemId)
			if err != nil {
				return errors.NoType.Wrap(err, "error updating menu item")
			}
			if depthLevel >= mi.maxMenuItemsDepth {
				return ErrMaxDepthExceeded
			}

			maxChildrenDepthLevel, err := mi.menuItemRepository.GetMaxDepthLevel(ctx, action.MenuItemId)
			if err != nil {
				return errors.NoType.Wrap(err, "error getting menu item max depth level")
			}
			// если вложенность родительского пункта меню + максимальная вложенность дочерних пунктов меню
			// больше либо равно максимальной вложенности -> ошибка
			if depthLevel+maxChildrenDepthLevel >= mi.maxMenuItemsDepth {
				return ErrMaxDepthExceeded
			}
		}
	}

	menuItemsCount, err := mi.menuItemRepository.Count(ctx, MenuItemFilter{
		MenuId:   action.MenuId,
		ParentId: action.ParentMenuItemId,
	})
	if err != nil {
		return errors.NoType.Wrap(err, "error counting menu items")
	}

	maxPosition := menuItemsCount
	if action.ParentMenuItemId != menuItem.ParentMenuItemId {
		maxPosition++
	}
	if action.Position > maxPosition {
		return ErrMaxPosition
	}

	err = mi.menuItemRepository.Move(ctx, MoveMenuItemQuery{
		MenuItemId:          action.MenuItemId,
		OldParentMenuItemId: menuItem.ParentMenuItemId,
		NewParentMenuItemId: action.ParentMenuItemId,
		OldPosition:         menuItem.Position,
		NewPosition:         action.Position,
	})
	if err != nil {
		return errors.NoType.Wrap(err, "error moving menu item")
	}

	return nil
}

func (mi MenuItems) joinURL(domain *entity.Domain, uri string) string {
	var domainStr string
	if domain != nil {
		domainStr = domain.Domain
	}
	if uri == "" {
		return domainStr
	}

	return strings.Join([]string{domainStr, uri}, "/")
}
