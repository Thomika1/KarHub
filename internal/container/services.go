package container

import (
	"context"

	"github.com/Thomika1/KarHub/pkg/core/crud"
	"github.com/Thomika1/KarHub/pkg/domains/business/products"
	"github.com/Thomika1/KarHub/pkg/domains/business/restock"
	"github.com/Thomika1/KarHub/pkg/domains/shared/models/business"
)

func setupServices(ctx context.Context, cmp *components) (*services, error) {
	svc := &services{}

	if err := setupAPIServices(ctx, svc, cmp); err != nil {
		return nil, err
	}

	return svc, nil
}

func setupAPIServices(ctx context.Context, svc *services, cmp *components) error {
	if err := setupBusinessServices(ctx, svc, cmp); err != nil {
		return err
	}

	return nil
}

func setupBusinessServices(ctx context.Context, svc *services, cmp *components) error {

	service, err := setupProductService(ctx, svc, cmp)
	if err != nil {
		return err
	}
	svc.ProductService = service

	serviceRestock, err := setupRestockService(ctx, svc, cmp)
	if err != nil {
		return err
	}
	svc.RestockService = serviceRestock

	return nil
}

func setupProductService(ctx context.Context, svc *services, cmp *components) (*products.Service, error) {
	productRepository, err := crud.NewRepository(cmp.Database, crud.OnlyCreate, business.Product{})
	if err != nil {
		return nil, err
	}

	service, err := products.NewService(productRepository, cmp.Logger)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func setupRestockService(ctx context.Context, svc *services, cmp *components) (*restock.Service, error) {
	restockRepository, err := crud.NewRepository(cmp.Database, crud.OnlyCreate, business.Product{})
	if err != nil {
		return nil, err
	}

	restockDomainRepository, err := restock.NewDomainRepository(cmp.Database)
	if err != nil {
		return nil, err
	}

	restockService, err := restock.NewService(restockRepository, cmp.Logger, restockDomainRepository)
	if err != nil {
		return nil, err
	}

	return restockService, nil
}
