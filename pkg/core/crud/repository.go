package crud

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Query struct {
	Page           int
	PageSize       int
	OrderBy        string
	OrderDirection string
	Filters        []Filter
}

func (q Query) IsValid() error {
	if q.Page < 0 {
		return errors.New("page must be non-negative")
	}
	if q.PageSize < 0 || q.PageSize > 100 {
		return errors.New("page_size must be between 0 and 100")
	}
	if q.OrderDirection != "" && q.OrderDirection != "asc" && q.OrderDirection != "desc" {
		return errors.New("order_direction must be 'asc' or 'desc'")
	}
	for _, f := range q.Filters {
		if f.field == "" {
			return errors.New("filter column cannot be empty")
		}
	}
	return nil
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new record in the database
func (r *Repository) Create(ctx context.Context, entity any) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return err
	}

	return nil
}

// Get retrieves an entity with the given id
func (r *Repository) Get(ctx context.Context, entity any, id any, preload bool) error {
	return r.GetBy(ctx, entity, "id", id, preload)
}

// GetBy retrieves an entity based on the given field and value
func (r *Repository) GetBy(ctx context.Context, entity any, column string, value any, preload bool) error {
	query := r.db.WithContext(ctx)

	if preload {
		query = query.Preload(clause.Associations)
	}

	query = query.Where(column+" = ?", value)

	return query.Take(entity).Error
}

// Update updates an entity with the given id
func (r *Repository) Update(ctx context.Context, entity any, id any) error {
	return r.UpdateBy(ctx, entity, "id", id)
}

// Update updates an entity with the given column
func (r *Repository) UpdateBy(ctx context.Context, entity any, column string, value any) error {
	query := r.db.WithContext(ctx)

	return query.Where(column+" = ?", value).Updates(entity).Error
}

// Delete deletes an entity with the given id
func (r *Repository) Delete(ctx context.Context, entity any, id any) error {
	return r.DeleteBy(ctx, entity, "id", id)
}

func (r *Repository) DeleteBy(ctx context.Context, entity any, column string, value any) error {
	query := r.db.WithContext(ctx)

	query = query.Where(column+" = ?", value)

	return query.Delete(entity).Error
}

func (r *Repository) List(ctx context.Context, entities any, filter Query) error {
	err := filter.IsValid()
	if err != nil {
		return err
	}

	query := r.db.WithContext(ctx)

	query = r.ApplyPreloads(query, filter)

	query = r.ApplyFilters(query, filter.Filters)

	query = r.ApplyPaginationRules(query, filter)

	query = r.ApplyOrderBy(query, filter)

	return query.Find(entities).Error
}
