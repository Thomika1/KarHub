package crud

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryI interface {
	Create(ctx context.Context, entity any) error
	Get(ctx context.Context, entity any, id any, preload bool) error
	GetBy(ctx context.Context, entity any, column string, value any, preload bool) error
	Update(ctx context.Context, entity any, id any) error
	UpdateBy(ctx context.Context, entity any, column string, value any) error
	Delete(ctx context.Context, entity any, id any) error
	DeleteBy(ctx context.Context, entity any, column string, value any) error
	List(ctx context.Context, entities any, filter Query) error
}

type MigrationStrategy int

const (
	// Disabled is a migration strategy that disables the migration
	Disabled MigrationStrategy = iota
	// OnlyCreate is a migration strategy that only creates the table
	OnlyCreate
	// AutoMigrate is a migration strategy that creates the table if it doesn't exist
	AutoMigrate
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB, migration MigrationStrategy, args ...any) (*Repository, error) {
	if len(args) == 0 {
		return nil, errors.New("args cannot be empty")
	}

	if migration != Disabled {
		err := applyMigration(db, migration, args)
		if err != nil {
			return nil, err
		}
	}
	return &Repository{db: db}, nil
}

func applyMigration(db *gorm.DB, migration MigrationStrategy, models []any) error {
	dst := make([]any, 0)

	for _, model := range models {
		shouldCreate := migration == OnlyCreate && !db.Migrator().HasTable(model)
		if migration == AutoMigrate || shouldCreate {
			dst = append(dst, model)
		}
	}

	if len(dst) > 0 {
		if err := db.AutoMigrate(dst...); err != nil {
			return err
		}
	}

	return nil
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

	query = r.ApplyFilters(query, filter.Filters)

	query = r.ApplyPaginationRules(query, filter)

	query = r.ApplyOrderBy(query, filter)

	return query.Find(entities).Error
}
