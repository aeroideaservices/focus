package postgres

import (
	"context"
	"reflect"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/aeroideaservices/focus/models/plugin/actions"
	"github.com/aeroideaservices/focus/models/plugin/focus"
	"github.com/aeroideaservices/focus/models/postgres/utils"
	focusClause "github.com/aeroideaservices/focus/services/db/clause"
	"github.com/aeroideaservices/focus/services/errors"
)

// elementsRepository репозиторий элементов модели
type elementsRepository struct {
	db    *gorm.DB
	model *focus.Model
}

// newElementsRepository конструктор
func newElementsRepository(db *gorm.DB, model *focus.Model) *elementsRepository {
	return &elementsRepository{
		db:    db,
		model: model,
	}
}

// Has проверяет, существует ли элемент модели с таким первичным ключом
func (r elementsRepository) Has(ctx context.Context, pk any) (bool, error) {
	elem, _ := r.model.NewElement(nil, nil)
	err := r.db.WithContext(ctx).Table(r.model.TableName).
		Scopes(getSelectScopes(r.model.TableName, []*focus.Field{r.model.PrimaryKey}, true)).
		Where(clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: r.model.PrimaryKey.Column}, Value: pk}).
		First(elem).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, errors.NoType.Wrap(err, "error checking model element exists")
	}

	return true, nil
}

// Get получение элемента модели по первичному ключу
func (r elementsRepository) Get(ctx context.Context, pk any) (elem any, err error) {
	elem, err = r.model.NewElement(map[string]any{r.model.PrimaryKey.Code: pk}, nil)
	if err != nil {
		return nil, err
	}

	db := r.db.WithContext(ctx).Table(r.model.TableName)
	// подтягиваем элементы ассоциированных моделей/медиа
	for _, field := range r.model.Fields {
		// если тип ассоциации - many2many, при этом в смежной таблице связи сортируются
		if field.Association != nil && field.Association.Type == focus.ManyToMany && field.Association.JoinSort != "" {
			db = db.Preload(field.Name(), r.preload(field.Association, pk))
		} else if field.Association != nil || field.IsMedia { // если ассоциация без сортировки/привязка к медиа
			db = db.Preload(field.Name())
		}
	}
	err = db.First(elem).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.Wrapf(err, "model element with id %s not found", pk)
	}
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting model element")
	}

	return elem, nil
}

// List получение списка элементов модели по фильтру
func (r elementsRepository) List(ctx context.Context, filter actions.ListModelElementsQuery) (elems []any, err error) {
	es, _ := r.model.ElementsSlice(nil)

	var selectFields []*focus.Field
	for _, fieldCode := range filter.SelectFields {
		var field *focus.Field
		if r.model.PrimaryKey.Code == fieldCode {
			field = r.model.PrimaryKey
		} else {
			field = r.model.Fields.GetByCode(fieldCode)
		}
		if field != nil {
			selectFields = append(selectFields, field)
		}
	}

	model, _ := r.model.NewElement(nil, nil)
	err = r.db.WithContext(ctx).Table(r.model.TableName).Model(model).
		Scopes(
			r.getFilterScopes(clause.CurrentTable, filter.Filter.FieldsFilter),
			r.getIlikeScopes(clause.CurrentTable, filter.Filter.QueryFilter),
			getSelectScopes(clause.CurrentTable, selectFields, true),
			getPaginationScopes(filter.Pagination),
			getOrderByScopes(clause.CurrentTable, filter.OrderBy),
		).
		Find(&es).Error
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting model elements list")
	}

	ves := reflect.ValueOf(es)
	for i := 0; i < ves.Len(); i++ {
		elems = append(elems, ves.Index(i).Interface())
	}

	return elems, nil
}

// Create создание нового элемент модели
func (r elementsRepository) Create(ctx context.Context, elem any) (id any, err error) {
	err = r.db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(elem).Error
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error creating model element")
	}

	pks, _ := focus.GetPKs(elem, r.model.PrimaryKey.Name())
	return pks[0], nil
}

// Update обновление элемента модели
func (r elementsRepository) Update(ctx context.Context, elem any) error {
	db := r.db.WithContext(ctx).Table(r.model.TableName)

	// сохраняем элемент модели
	err := db.Save(elem).Error
	if err != nil {
		return errors.NoType.Wrap(err, "error updating model element")
	}

	// отдельно сохраняем все связи
	for _, field := range r.model.Fields {
		// если поле не является ассоциацией или медиа - пропускаем
		if field.Association == nil && !field.IsMedia {
			continue
		}

		// получаем значение поля
		fv := reflect.ValueOf(elem).Elem().FieldByName(field.Name())
		fieldValue := fv.Interface()

		// костыль, иначе падает с паникой reflect.Value.Addr of unaddressable value
		if fv.Kind() == reflect.Struct {
			fieldValue = fv.Addr().Interface()
		}

		// подменяем ассоциацию
		err := r.db.Model(elem).Association(field.Name()).Replace(fieldValue)
		if err != nil {
			return errors.NoType.Wrap(err, "error updating model element associations")
		}
	}

	return nil
}

// Delete удаление элемента модели по первичному ключу
func (r elementsRepository) Delete(ctx context.Context, pks ...any) error {
	model, _ := r.model.NewElement(nil, nil) // пустой элемент модели
	err := r.db.WithContext(ctx).Table(r.model.TableName).Delete(model, pks).Error
	if err != nil {
		return errors.NoType.Wrap(err, "error deleting model element")
	}

	return nil
}

// Count получение количества элементов модели по фильтру
func (r elementsRepository) Count(ctx context.Context, filter actions.ModelElementsFilter) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table(r.model.TableName).
		Scopes(
			r.getFilterScopes(clause.CurrentTable, filter.FieldsFilter),
			r.getIlikeScopes(clause.CurrentTable, filter.QueryFilter),
		).
		Clauses(clause.Select{Distinct: true, Columns: []clause.Column{{Table: clause.CurrentTable, Name: r.model.PrimaryKey.Column}}}).
		Count(&count).Error
	if err != nil {
		return 0, errors.NoType.Wrap(err, "error counting mode elements")
	}

	return count, nil
}

// GetFieldValues
func (r elementsRepository) ListFieldValues(ctx context.Context, filter actions.ListFieldValues) (any, error) {
	var field *focus.Field
	if filter.FieldCode == r.model.PrimaryKey.Code {
		field = r.model.PrimaryKey
	} else {
		field = r.model.Fields.GetByCode(filter.FieldCode)
	}
	if field == nil {
		return nil, errors.BadRequest.Newf("wrong field code given")
	}

	values := reflect.MakeSlice(reflect.SliceOf(field.RawType()), 0, 0).Interface()
	err := r.db.WithContext(ctx).Table(r.model.TableName).
		Scopes(
			r.getFieldValueFilterScopes(clause.CurrentTable, field, filter.Query),
			getSelectScopes(clause.CurrentTable, []*focus.Field{field}, true),
			getPaginationScopes(filter.Pagination),
			getOrderByScopes(clause.CurrentTable, actions.OrderBy{Sort: field.Column, Order: "asc"}),
		).
		Scan(&values).Error

	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting field values")
	}

	return values, nil
}

// GetFieldValuesCount
func (r elementsRepository) CountFieldValues(ctx context.Context, fieldCode string, query string) (int64, error) {
	field := r.model.Fields.GetByCode(fieldCode)
	if field == nil {
		return 0, errors.NotFound.New("field not found")
	}

	var count int64
	err := r.db.WithContext(ctx).Table(r.model.TableName).Distinct(field.Column).
		Scopes(r.getFieldValueFilterScopes(clause.CurrentTable, field, query)).
		Clauses(clause.Select{Distinct: true, Columns: []clause.Column{{Table: clause.CurrentTable, Name: field.Column}}}).
		Count(&count).Error

	if err != nil {
		return 0, errors.NoType.Wrap(err, "error counting field values")
	}

	return count, nil
}

// getPaginationScopes получает scopes для LIMIT OFFSET
func getPaginationScopes(dto actions.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(dto.Offset).Limit(dto.Limit)
	}
}

// getOrderByScopes получает scopes для ORDER BY
func getOrderByScopes(table string, dto actions.OrderBy) func(db *gorm.DB) *gorm.DB {
	if dto.Sort == "" {
		return func(db *gorm.DB) *gorm.DB { return db }
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(
			clause.OrderBy{Columns: []clause.OrderByColumn{{Column: clause.Column{Table: table, Name: dto.Sort}, Desc: dto.Order == "desc"}}},
		)
	}
}

// getIlikeScopes получает scopes для ILIKE
func (r elementsRepository) getIlikeScopes(table string, dto actions.ModelElementsQueryFilter) func(db *gorm.DB) *gorm.DB {
	if len(dto.FieldsCodes) == 0 || dto.Query == "" {
		return func(db *gorm.DB) *gorm.DB { return db }
	}

	query := utils.PrepareLikeQuery(dto.Query)
	var exprs []clause.Expression
	for _, fieldCode := range dto.FieldsCodes {
		field := r.model.Fields.GetByCode(fieldCode)
		exprs = append(exprs, focusClause.Ilike{Column: clause.Column{Table: table, Name: field.Column + "::varchar", Raw: true}, Value: query})
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(clause.OrConditions{Exprs: exprs})
	}
}

// getSelectScopes получает scopes для SELECT
func getSelectScopes(table string, fields []*focus.Field, distinct bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		columns := make([]clause.Column, 0, len(fields))
		for _, field := range fields {
			if field.Association != nil || field.IsMedia {
				continue
			}
			columns = append(columns, clause.Column{Table: table, Name: field.Column})
		}
		return db.Clauses(clause.Select{Distinct: distinct, Columns: columns})
	}
}

// getFilterScopes получает scopes для WHERE
func (r elementsRepository) getFilterScopes(table string, filter actions.FieldsFilter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for fieldCode, values := range filter {
			if len(values) == 0 {
				continue
			}
			field := r.model.Fields.GetByCode(fieldCode)

			if field.IsTime {
				if len(values) != 2 {
					db.Error = errors.BadRequest.Newf("wrong length of time filters for field %s", fieldCode)
				}
				from, to := values[0], values[1]

				// если переданы пустые значения
				if (from == nil || from == time.Time{}) && (to == nil || to == time.Time{}) {
					continue
				}

				var expr clause.Expression
				if from != nil && (from != time.Time{}) {
					expr = clause.Gte{Column: clause.Column{Table: table, Name: field.Column}, Value: from}
				}
				if to != nil && (to != time.Time{}) {
					if expr != nil {
						expr = clause.And(expr, clause.Lte{Column: clause.Column{Table: table, Name: field.Column}, Value: to})
					}
				}
				db.Where(clause.Or(clause.Eq{Column: clause.Column{Table: table, Name: field.Column}, Value: nil}, expr))
			} else {
				db.Where(clause.IN{Column: clause.Column{Table: table, Name: field.Column}, Values: values})
			}
		}

		return db
	}
}

// getFieldValueFilterScopes получение scope для полей модели для WHERE
func (r elementsRepository) getFieldValueFilterScopes(table string, field *focus.Field, query string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if field.IsTime {
			_ = db.AddError(errors.NoType.New("cannot use ILIKE query with datetime field"))
			return db
		}
		if field.IsMedia {
			_ = db.AddError(errors.NoType.New("cannot use ILIKE query with media field"))
			return db
		}
		if field.Association != nil {
			_ = db.AddError(errors.NoType.New("cannot use ILIKE query with associated field"))
			return db
		}

		if query != "" {
			db = db.Where(focusClause.Ilike{
				Column: focusClause.Cast{Column: clause.Column{Table: table, Name: field.Column}, Type: "TEXT"},
				Value:  utils.PrepareLikeQuery(query),
			})
		}

		return db.Where(clause.Neq{Column: clause.Column{Table: clause.CurrentTable, Name: field.Column}, Value: nil})
	}
}

// Правила preload ассоциаций для связи many2many с сортировкой
func (r elementsRepository) preload(assoc *focus.Association, pk any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// если не many2many или не сортируемое -> ничего не делаем
		if assoc.Type != focus.ManyToMany || assoc.JoinSort == "" {
			return db
		}

		// получаем все поля модели + сортировку // todo ограничиться pkey и sort
		selectScope := clause.Select{
			Distinct: true,
			Columns:  []clause.Column{{Name: "*", Raw: true}, {Table: assoc.Many2Many, Name: assoc.JoinSort + " AS HIDDEN", Raw: true}},
		}
		fromScope := clause.From{
			Tables: []clause.Table{{Name: assoc.Model.Code}}, // название таблицы ассоциируемой сущности
			Joins: []clause.Join{{
				Type:  clause.InnerJoin,
				Table: clause.Table{Name: assoc.Many2Many}, // название связующей таблицы
				ON: clause.Where{Exprs: []clause.Expression{
					clause.Eq{
						Column: clause.Column{Table: assoc.Many2Many, Name: assoc.JoinReferences},
						Value:  clause.Column{Table: clause.CurrentTable, Name: assoc.Model.PrimaryKey.Column},
					},
					clause.Eq{
						Column: clause.Column{Table: assoc.Many2Many, Name: assoc.JoinForeignKey},
						Value:  pk,
					},
				}},
			}},
		}
		// сортировка по полю sort из связующей таблицы
		orderByScope := clause.OrderBy{
			Columns: []clause.OrderByColumn{{
				Column: clause.Column{Table: assoc.Many2Many, Name: "sort"},
			}},
		}
		return db.Clauses(selectScope, fromScope, orderByScope)
	}
}
