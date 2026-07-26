package container

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
)

type components struct {
	Database *gorm.DB
	Logger   *slog.Logger
}

type services struct {
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
