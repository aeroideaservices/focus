package repositories

import (
	"context"
	"github.com/aeroideaservices/focus/menu/plugin/actions"
	"github.com/aeroideaservices/focus/menu/plugin/entity"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) actions.MenuRepository {
	return &menuRepository{db: db}
}

func (r menuRepository) HasByCode(ctx context.Context, code string) (bool, error) {
	err := r.db.WithContext(ctx).Where("code", code).First(&entity.Menu{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, errors.WithStack(err)
	}

	return true, nil
}

func (r menuRepository) Has(ctx context.Context, id uuid.UUID) (bool, error) {
	err := r.db.WithContext(ctx).Where("id", id).First(&entity.Menu{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, errors.WithStack(err)
	}

	return true, nil
}

func (r menuRepository) Create(ctx context.Context, menu entity.Menu) error {
	err := r.db.WithContext(ctx).Create(menu).Error

	return errors.WithStack(err)
}

func (r menuRepository) GetByCode(ctx context.Context, code string) (*entity.Menu, error) {
	menu := &entity.Menu{}
	err := r.db.WithContext(ctx).Where("code", code).First(menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, actions.ErrMenuNotFound
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return menu, nil
}

func (r menuRepository) Get(ctx context.Context, Id uuid.UUID) (*entity.Menu, error) {
	menu := &entity.Menu{}
	err := r.db.WithContext(ctx).First(menu, Id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, actions.ErrMenuNotFound
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return menu, nil
}

func (r menuRepository) Update(ctx context.Context, menu *entity.Menu) error {
	err := r.db.WithContext(ctx).
		Model(entity.Menu{}).
		Where("id", menu.Id).
		Updates(map[string]any{
			"code": menu.Code,
			"name": menu.Name,
		}).Error

	return errors.WithStack(err)
}

func (r menuRepository) Delete(ctx context.Context, Id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&entity.Menu{Id: Id}).Error

	return errors.WithStack(err)
}
func (r menuRepository) List(ctx context.Context, filter actions.MenuFilter) ([]entity.Menu, error) {
	var menus []entity.Menu
	sort := filter.Sort
	if sort == "" {
		sort = "id"
	}
	order := filter.Order

	err := r.db.WithContext(ctx).
		Limit(filter.Limit).Offset(filter.Offset).
		Order(sort + " " + order).
		Find(&menus).Error

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return menus, nil
}
func (r menuRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&entity.Menu{}).
		Distinct("id").
		Count(&count).Error
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}
