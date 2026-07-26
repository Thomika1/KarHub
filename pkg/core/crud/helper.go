package crud

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Filter struct {
	field    string
	Value    any
	Operator string // "=", ">", "<", ">=", "<=", "LIKE", "IN"
}

func (r *Repository) ApplyPreloads(query *gorm.DB, filter Query) *gorm.DB {
	return query
}

func (r *Repository) ApplyFilters(query *gorm.DB, filters []Filter) *gorm.DB {
	for _, f := range filters {
		if f.field == "" || f.Operator == "" {
			continue
		}

		op := strings.ToUpper(f.Operator)
		switch op {
		case "LIKE":
			query = query.Where(fmt.Sprintf("%s LIKE ?", f.field), fmt.Sprintf("%%%s%%", f.Value))
		case "IN":
			query = query.Where(f.field+" IN (?)", f.Value)
		default:
			query = query.Where(fmt.Sprintf("%s %s ?", f.field, op), f.Value)
		}
	}
	return query
}

func (r *Repository) ApplyOrganizationFilter(ctx context.Context, entities any, query *gorm.DB) *gorm.DB {
	return query
}

func (r *Repository) ApplyPaginationRules(query *gorm.DB, filter Query) *gorm.DB {
	if filter.PageSize > 0 {
		query = query.Limit(filter.PageSize)
	}
	if filter.Page > 0 && filter.PageSize > 0 {
		query = query.Offset(filter.Page * filter.PageSize)
	}
	return query
}

func (r *Repository) ApplyOrderBy(query *gorm.DB, filter Query) *gorm.DB {
	if filter.OrderBy != "" {
		dir := "asc"
		if strings.ToLower(filter.OrderDirection) == "desc" {
			dir = "desc"
		}
		query = query.Order(fmt.Sprintf("%s %s", filter.OrderBy, dir))
	}
	return query
}
