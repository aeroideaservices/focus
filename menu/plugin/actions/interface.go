package actions

import (
	"context"
	"github.com/aeroideaservices/focus/menu/plugin/entity"
	"github.com/aeroideaservices/focus/services/db/db_types/json"
	"github.com/google/uuid"
)

// MenuRepository интерфейс репозитория работы с меню
type MenuRepository interface {
	Has(ctx context.Context, id uuid.UUID) (bool, error)
	HasByCode(ctx context.Context, code string) (bool, error)
	Create(ctx context.Context, menu entity.Menu) error
	Update(ctx context.Context, menu *entity.Menu) error
	GetByCode(ctx context.Context, code string) (*entity.Menu, error)
	Get(ctx context.Context, Id uuid.UUID) (*entity.Menu, error)
	Delete(ctx context.Context, Id uuid.UUID) error
	List(ctx context.Context, filter MenuFilter) ([]entity.Menu, error)
	Count(ctx context.Context) (int64, error)
}

type MenuFilter struct {
	Offset int
	Limit  int
	Sort   string
	Order  string
}

// MenuItemRepository интерфейс репозитория работы с элементами меню
type MenuItemRepository interface {
	Has(ctx context.Context, menuId, menuItemId uuid.UUID) (bool, error)
	Create(ctx context.Context, menuItem entity.MenuItem) error
	Update(ctx context.Context, menuItem entity.MenuItem) error
	Move(ctx context.Context, move MoveMenuItemQuery) error
	Delete(ctx context.Context, Id uuid.UUID) error
	Get(ctx context.Context, Id uuid.UUID) (*entity.MenuItem, error)
	List(ctx context.Context, filter MenuItemsListFilter) ([]entity.MenuItem, error)
	Count(ctx context.Context, filter MenuItemFilter) (int64, error)
	GetAsTree(ctx context.Context, menuCode string) ([]*entity.MenuItem, error)
	GetDepthLevel(ctx context.Context, id uuid.UUID) (int, error)
	GetMaxDepthLevel(ctx context.Context, parentId uuid.UUID) (int, error)
}
type MenuItemsListFilter struct {
	Sort   string
	Order  string
	Filter MenuItemFilter
}

type MenuItemFilter struct {
	MenuId   uuid.UUID  `json:"menuId" validate:"required,notBlank"`
	ParentId *uuid.UUID `json:"parentId,omitempty" validate:"omitempty,notBlank"`
}

type MenuItemPreview struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Url      string    `json:"url"`
	Position int64     `json:"position"`
}

type MenuItemTree struct {
	Id               uuid.UUID       `json:"id"`
	ParentMenuItemId *uuid.UUID      `json:"-"`
	Name             string          `json:"name"`
	Url              string          `json:"url"`
	Position         int64           `json:"position"`
	AdditionalFields json.JSONB      `json:"additionalFields"`
	MenuItems        []*MenuItemTree `json:"menuItems,omitempty" gorm:"-"`
}

type MoveMenuItemQuery struct {
	MenuItemId          uuid.UUID
	OldParentMenuItemId *uuid.UUID
	NewParentMenuItemId *uuid.UUID
	OldPosition         int64
	NewPosition         int64
}

type DomainsRepository interface {
	Has(ctx context.Context, domain string) (bool, error)
	Create(ctx context.Context, domain entity.Domain) error
	List(ctx context.Context, query DomainsListQuery) ([]entity.Domain, error)
	Count(ctx context.Context) (int64, error)
}

type DomainsListQuery struct {
	Pagination
}
