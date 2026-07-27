package restock

import (
	"context"
	"log/slog"

	"github.com/Thomika1/KarHub/pkg/core/crud"
	"github.com/Thomika1/KarHub/pkg/domains/shared/models/business"
)

type ServiceI interface {
	Priorities(context.Context, crud.Query) ([]business.Product, error)
}

type DomainRepository interface {
	Priorities(context.Context, crud.Query) ([]business.Product, error)
}

type Service struct {
	repository       crud.RepositoryI
	logger           *slog.Logger
	domainRepository DomainRepository
}

func NewService(repository crud.RepositoryI, logger *slog.Logger, domainRepository DomainRepository) (*Service, error) {
	return &Service{
		repository:       repository,
		logger:           logger,
		domainRepository: domainRepository,
	}, nil
}
