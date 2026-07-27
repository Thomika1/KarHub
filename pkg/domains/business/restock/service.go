package restock

import (
	"context"

	"github.com/Thomika1/KarHub/pkg/core/crud"
	"github.com/Thomika1/KarHub/pkg/domains/shared/models/business"
)

func (s *Service) Priorities(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
	var products []business.Product

	products, err := s.domainRepository.Priorities(ctx, parameters)
	if err != nil {
		return nil, err
	}

	return products, nil
}
