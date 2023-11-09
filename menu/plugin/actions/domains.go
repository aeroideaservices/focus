package actions

import (
	"context"
	"github.com/aeroideaservices/focus/menu/plugin/entity"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/google/uuid"
)

// Domains сервис работы с доменами
type Domains struct {
	domainRepository DomainsRepository
}

// NewDomains конструктор
func NewDomains(domainRepository DomainsRepository) *Domains {
	return &Domains{domainRepository: domainRepository}
}

// List получение списка доменов
func (d Domains) List(ctx context.Context, action ListDomains) (*DomainsList, error) {
	items, err := d.domainRepository.List(ctx, DomainsListQuery{
		Pagination: action.Pagination,
	})
	if err != nil {
		return nil, err
	}
	total, err := d.domainRepository.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &DomainsList{
		Total: total,
		Items: items,
	}, nil
}

// Create создание нового домена
func (d Domains) Create(ctx context.Context, action CreateDomain) (*uuid.UUID, error) {
	has, err := d.domainRepository.Has(ctx, action.Domain)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error checking that domain exists")
	}
	if has {
		return nil, errors.BadRequest.New("domain already exists").T("domain.conflict")
	}

	domain := entity.Domain{
		Id:     uuid.New(),
		Domain: action.Domain,
	}
	if err := d.domainRepository.Create(ctx, domain); err != nil {
		return nil, errors.NoType.Wrap(err, "error creating domain")
	}

	return &domain.Id, nil
}
