package crud

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Query struct {
	Page           int
	PageSize       int
	OrderBy        string
	OrderDirection string
	Filters        []Filter
}

type Filter struct {
	Field    string `json:"field"`
	Value    any    `json:"value"`
	Operator string `json:"operator"` // "=", ">", "<", ">=", "<=", "LIKE", "IN"
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
		if f.Field == "" {
			return errors.New("filter column cannot be empty")
		}
	}
	return nil
}

func (r *Repository) ApplyFilters(query *gorm.DB, filters []Filter) *gorm.DB {
	for _, f := range filters {
		if f.Field == "" || f.Operator == "" {
			continue
		}

		op := strings.ToUpper(f.Operator)
		switch op {
		case "LIKE":
			query = query.Where(fmt.Sprintf("%s LIKE ?", f.Field), fmt.Sprintf("%%%s%%", f.Value))
		case "IN":
			query = query.Where(f.Field+" IN (?)", f.Value)
		default:
			query = query.Where(fmt.Sprintf("%s %s ?", f.Field, op), f.Value)
		}
	}
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
