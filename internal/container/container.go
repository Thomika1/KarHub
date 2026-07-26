package container

import (
	"context"
	"log/slog"

	"github.com/Thomika1/KarHub/pkg/domains/business/products"
	"gorm.io/gorm"
)

type components struct {
	Database *gorm.DB
	Logger   *slog.Logger
}

type services struct {
	ProductService products.ServiceI
}

type Dependency struct {
	Components components
	Services   services
}

func New(ctx context.Context) (*Dependency, error) {
	cmp, err := setupComponents(ctx)
	if err != nil {
		return nil, err
	}

	svc, err := setupServices(ctx, cmp)
	if err != nil {
		return nil, err
	}

	dep := Dependency{
		Components: *cmp,
		Services:   *svc,
	}

	return &dep, err
}
