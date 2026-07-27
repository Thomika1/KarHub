package restock

import (
	"context"
	"errors"

	"github.com/Thomika1/KarHub/pkg/core/crud"
	"github.com/Thomika1/KarHub/pkg/domains/shared/models/business"
	"github.com/go-kit/kit/endpoint"
)

func mapToPriorityItem(p business.Product) PriorityItem {
	item := PriorityItem{
		PartID:         *p.ID,
		CurrentStock:   p.CurrentStock,
		ProjectedStock: p.ProjectedStock,
		MinimumStock:   p.MinimumStock,
		UrgencyScore:   p.UrgencyScore,
	}
	if p.Name != nil {
		item.Name = *p.Name
	}
	return item
}

func Priorities(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		query, ok := request.(crud.Query)
		if !ok {
			return nil, errors.New("invalid request")
		}

		products, err := svc.Priorities(ctx, query)
		if err != nil {
			return nil, err
		}

		items := make([]PriorityItem, len(products))
		for i, p := range products {
			items[i] = mapToPriorityItem(p)
		}

		return PrioritiesResponse{Priorities: items}, nil
	}
}
