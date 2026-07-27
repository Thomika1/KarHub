package products

import (
	"context"
	"errors"

	"github.com/Thomika1/KarHub/pkg/core/crud"
	"github.com/Thomika1/KarHub/pkg/domains/shared/models/business"
	"github.com/go-kit/kit/endpoint"
)

type UpdateRequest struct {
	ID   string
	Data business.ProductData
}

func Create(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req, ok := request.(*business.Product)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := svc.Create(ctx, req); err != nil {
			return nil, err
		}

		return req, nil
	}
}

func Get(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		id, ok := request.(string)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return svc.Get(ctx, id)
	}
}

func List(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		query, ok := request.(crud.Query)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return svc.List(ctx, query)
	}
}

func Update(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req, ok := request.(UpdateRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := svc.Update(ctx, req.ID, &req.Data); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func Delete(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		id, ok := request.(string)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := svc.Delete(ctx, id); err != nil {
			return nil, err
		}

		return map[string]string{"id": id}, nil
	}
}

func Priorities(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		query, ok := request.(crud.Query)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return svc.Priorities(ctx, query)
	}
}
