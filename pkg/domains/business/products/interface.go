package products

import (
	"context"
	"log/slog"

	"github.com/Thomika1/KarHub/pkg/core/crud"
	"github.com/Thomika1/KarHub/pkg/domains/shared/models/business"
)

type ServiceI interface {
	Create(context.Context, *business.Product) error
	Get(context.Context, string) (business.Product, error)
	List(context.Context, crud.Query) ([]business.Product, error)
	Update(context.Context, string, *business.ProductData) error
	Delete(context.Context, string) error
}

type Service struct {
	repository crud.RepositoryI
	logger     slog.Logger
}

func NewService(repository crud.RepositoryI, logger *slog.Logger) (*Service, error) {
	return &Service{
		repository: repository,
		logger:     *logger,
	}, nil
}
