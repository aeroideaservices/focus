package repositories

import (
	"context"
	"github.com/WinterYukky/gorm-extra-clause-plugin/exclause"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/aeroideaservices/focus/menu/plugin/actions"
	"github.com/aeroideaservices/focus/menu/plugin/entity"
	focusClause "github.com/aeroideaservices/focus/services/db/clause"
	"github.com/aeroideaservices/focus/services/errors"
)

type menuItemRepository struct {
	db *gorm.DB
}

func NewMenuItemRepository(db *gorm.DB) actions.MenuItemRepository {
	return &menuItemRepository{db: db}
}

func (r menuItemRepository) Has(ctx context.Context, menuId, menuItemId uuid.UUID) (bool, error) {
	err := r.db.WithContext(ctx).
		Where("id", menuItemId).
		Where("menu_id", menuId).
		First(&entity.MenuItem{}).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, errors.NoType.Wrap(err, "error checking if menu item exists")
	}

	return true, nil
}

func (r menuItemRepository) Create(ctx context.Context, menuItem entity.MenuItem) error {
	err := r.db.WithContext(ctx).
		Create(&menuItem).Error
	if err != nil {
		return errors.NoType.Wrap(err, "error creating menu item")
	}

	return nil
}

func (r menuItemRepository) List(ctx context.Context, filter actions.MenuItemsListFilter) ([]entity.MenuItem, error) {
	items := make([]entity.MenuItem, 0)
	sort := filter.Sort
	order := filter.Order
	if sort == "" {
		sort = "position"
	}

	recursiveTable := "fm"
	err := r.db.Table(recursiveTable).WithContext(ctx).Where("fm.menu_id", filter.Filter.MenuId).
		Scopes(withRecursive(recursiveTable), orderBy(sort, order)).
		Where("fm.parent_menu_item_id", filter.Filter.ParentId).
		Preload("Domain").
		Find(&items).Error

	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting menu items list")
	}

	return items, nil
}

func (r menuItemRepository) GetAsTree(ctx context.Context, menuCode string) ([]*entity.MenuItem, error) {
	items := make([]*entity.MenuItem, 0)

	tableName := "menu_items_urls"
	err := r.db.WithContext(ctx).
		Table(tableName).
		Scopes(withRecursive(tableName), orderBy("position", "asc")).
		Joins("INNER JOIN menus ON menus.id = "+tableName+".menu_id").
		Where("menus.code", menuCode).
		Preload("Domain").
		Find(&items).Error
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting menu items as tree")
	}

	return items, nil
}

func withRecursive(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(exclause.With{
			Recursive: true,
			CTEs: []exclause.CTE{{
				Name: tableName,
				Subquery: clause.Expr{
					SQL: `
					SELECT id, 
					       parent_menu_item_id, 
					       name, 
					       menu_id, 
					       position, 
					       additional_fields, 
					       domain_id,
					       url
                    FROM menu_items
                    WHERE parent_menu_item_id IS NULL 
                    UNION ALL
                    SELECT child.id,
                           child.parent_menu_item_id,
                           child.name,
                           child.menu_id,
                           child.position,
                           child.additional_fields,
                           child.domain_id,
                           child.url
                    FROM menu_items child
                             INNER JOIN ? parent ON parent.id = child.parent_menu_item_id
				 	`,
					Vars: []any{clause.Table{Name: tableName}},
				},
			}},
		})
	}
}

func (r menuItemRepository) GetDepthLevel(ctx context.Context, id uuid.UUID) (int, error) {
	tableName := "menu_items_depth_levels"
	var depthLevel int

	err := r.db.WithContext(ctx).
		Table(tableName).
		Scopes(withRecursiveDepthLevel(tableName, nil)).
		Select("level").
		Clauses(clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: "id"}, Value: id}).
		Order("level").Limit(1).
		Find(&depthLevel).Error
	if err != nil {
		return 0, errors.NoType.Wrap(err, "error getting menu item depth level")
	}
	if depthLevel == 0 {
		return 0, actions.ErrMenuItemNotFound
	}

	return depthLevel, nil
}

func (r menuItemRepository) GetMaxDepthLevel(ctx context.Context, parentId uuid.UUID) (int, error) {
	tableName := "menu_items_depth_levels"
	var depthLevel int

	err := r.db.WithContext(ctx).
		Table(tableName).
		Scopes(withRecursiveDepthLevel(tableName, &parentId)).
		Select("coalesce(MAX(level), 0)").
		Find(&depthLevel).Error
	if err != nil {
		return 0, errors.NoType.Wrap(err, "error getting max depth level of child menu items")
	}

	return depthLevel, nil
}

func withRecursiveDepthLevel(tableName string, parentMenuItemId *uuid.UUID) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(exclause.With{
			Recursive: true,
			CTEs: []exclause.CTE{{
				Name: tableName,
				Subquery: clause.Expr{
					SQL: `
					SELECT id, 
					       parent_menu_item_id, 
					       1 AS level
                    FROM menu_items
                    WHERE ?
                    UNION ALL
                    SELECT child.id,
                           child.parent_menu_item_id,
                           parent.level + 1 AS level
                    FROM menu_items child
                             INNER JOIN ? parent ON parent.id = child.parent_menu_item_id
				 	`,
					Vars: []any{
						clause.Eq{Column: clause.Column{Name: "parent_menu_item_id"}, Value: parentMenuItemId},
						clause.Table{Name: tableName},
					},
				},
			}},
		})
	}
}

func orderBy(sort string, order string) func(db *gorm.DB) *gorm.DB {
	if sort == "" {
		return func(db *gorm.DB) *gorm.DB { return db }
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{Column: clause.Column{Name: sort}, Desc: order == "desc"},
			},
		})
	}
}

func (r menuItemRepository) Update(ctx context.Context, menuItem entity.MenuItem) error {
	err := r.db.WithContext(ctx).
		Model(entity.MenuItem{}).
		Where("id", menuItem.Id).
		Updates(map[string]any{
			"name":              menuItem.Name,
			"domain_id":         menuItem.DomainId,
			"url":               menuItem.Url,
			"additional_fields": menuItem.AdditionalFields,
		}).Error
	if err != nil {
		return errors.NoType.Wrap(err, "error updating menu item")
	}

	return nil
}

func (r menuItemRepository) Move(ctx context.Context, move actions.MoveMenuItemQuery) error {
	var increaseCondition clause.Expression
	var decreaseCondition clause.Expression
	if move.OldParentMenuItemId != move.NewParentMenuItemId {
		increaseCondition = clause.And(
			clause.Eq{Column: "parent_menu_item_id", Value: move.NewParentMenuItemId},
			clause.Gte{Column: "position", Value: move.NewPosition},
		)
		decreaseCondition = clause.And(
			clause.Eq{Column: "parent_menu_item_id", Value: move.OldParentMenuItemId},
			clause.Gt{Column: "position", Value: move.OldPosition},
		)
	} else {
		increaseCondition = clause.And(
			clause.Gte{Column: "position", Value: move.NewPosition},
			clause.Lt{Column: "position", Value: move.OldPosition},
		)
		decreaseCondition = clause.And(
			clause.Gt{Column: "position", Value: move.OldPosition},
			clause.Lte{Column: "position", Value: move.NewPosition},
		)
	}

	posExpr := focusClause.Case{
		WhenThen: []focusClause.WhenThen{
			{When: clause.Eq{Column: "id", Value: move.MenuItemId}, Then: move.NewPosition},
			{When: increaseCondition, Then: clause.Expr{SQL: "position + 1"}},
			{When: decreaseCondition, Then: clause.Expr{SQL: "position - 1"}},
		},
		Else: clause.Column{Name: "position"},
	}
	parentMenuItemIdExpr := focusClause.Case{
		WhenThen: []focusClause.WhenThen{
			{When: clause.Eq{Column: "id", Value: move.MenuItemId}, Then: move.NewParentMenuItemId},
		},
		Else: clause.Column{Name: "parent_menu_item_id"},
	}

	err := r.db.WithContext(ctx).Model(entity.MenuItem{}).
		Where(clause.Or(
			clause.Eq{Column: "parent_menu_item_id", Value: move.OldParentMenuItemId},
			clause.Eq{Column: "parent_menu_item_id", Value: move.NewParentMenuItemId},
		)).
		Updates(
			map[string]any{
				"position":            posExpr,
				"parent_menu_item_id": parentMenuItemIdExpr,
			},
		).Error
	if err != nil {
		return errors.NoType.Wrap(err, "error moving menu item")
	}

	return nil
}

func (r menuItemRepository) Delete(ctx context.Context, Id uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("id", Id).Delete(&entity.MenuItem{}).Error
	if err != nil {
		return errors.NoType.Wrap(err, "error deleting menu item")
	}

	return nil
}

func (r menuItemRepository) Get(ctx context.Context, Id uuid.UUID) (*entity.MenuItem, error) {
	db := r.db
	menuItem := &entity.MenuItem{}
	err := db.WithContext(ctx).First(menuItem, Id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, actions.ErrMenuItemNotFound
	}
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting menu item")
	}

	return menuItem, nil
}

func (r menuItemRepository) Count(ctx context.Context, filter actions.MenuItemFilter) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&entity.MenuItem{})
	db = filterMenuItems(db, filter)
	err := db.Distinct("id").Count(&count).Error
	if err != nil {
		return 0, errors.NoType.Wrap(err, "error counting menu items")
	}

	return count, nil
}

func filterMenuItems(db *gorm.DB, filter actions.MenuItemFilter) *gorm.DB {
	db = db.Where("menu_id", filter.MenuId)

	if filter.ParentId == nil {
		db = db.Where("parent_menu_item_id IS NULL")
	} else {
		db = db.Where("parent_menu_item_id", filter.ParentId)
	}

	return db
}
