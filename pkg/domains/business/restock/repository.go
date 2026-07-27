package restock

import (
	"context"

	"github.com/Thomika1/KarHub/pkg/core/crud"
	"github.com/Thomika1/KarHub/pkg/domains/shared/models/business"
	"gorm.io/gorm"
)

type domainRepository struct {
	db *gorm.DB
}

func NewDomainRepository(db *gorm.DB) (*domainRepository, error) {
	return &domainRepository{db: db}, nil
}

func (r *domainRepository) Priorities(ctx context.Context, parameters crud.Query) ([]business.Product, error) {
	var products []business.Product

	err := r.db.WithContext(ctx).
		Select("*, (current_stock - (average_daily_sales * lead_time_days)) AS projected_stock, (minimum_stock - (current_stock - (average_daily_sales * lead_time_days))) * criticality_level AS urgency_score").
		Order("urgency_score DESC, criticality_level DESC, average_daily_sales DESC, name ASC").
		Limit(parameters.PageSize).
		Offset(parameters.Page * parameters.PageSize).
		Find(&products).Error

	return products, err
}
