package products

import (
	"context"

	"github.com/Thomika1/KarHub/pkg/core/crud"
	"github.com/Thomika1/KarHub/pkg/domains/shared/errors"
	errUtils "github.com/Thomika1/KarHub/pkg/domains/shared/errors"
	"github.com/Thomika1/KarHub/pkg/domains/shared/models/business"
)

func (s *Service) Create(ctx context.Context, product *business.Product) error {
	if err := product.ProductData.Validate(); err != nil {
		return errors.NewValidationError("invalid product: %v", err)
	}

	var existing []business.Product
	filter := crud.Query{
		Filters: []crud.Filter{
			{Field: "name", Value: *product.Name, Operator: "="},
			{Field: "category", Value: product.Category, Operator: "="},
		},
		PageSize: 1,
	}
	if err := s.repository.List(ctx, &existing, filter); err != nil {
		return err
	}
	if len(existing) > 0 {
		return errUtils.ErrProductAlreadyExists
	}

	return s.repository.Create(ctx, product)
}

func (s *Service) Get(ctx context.Context, id string) (business.Product, error) {
	var product business.Product

	err := s.repository.Get(ctx, &product, id, false)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (s *Service) List(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
	var products []business.Product

	return products, s.repository.List(ctx, &products, parameters)
}

func (s *Service) Update(ctx context.Context, id string, data *business.ProductData) error {
	if data == nil {
		return nil
	}

	product := &business.Product{ProductData: *data}
	return s.repository.UpdateBy(ctx, product, "id", id)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, &business.Product{}, id)
}

func (s *Service) Priorities(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
	var products []business.Product

	s.domainRepository.Priorities(ctx, parameters)

	return products, nil
}
