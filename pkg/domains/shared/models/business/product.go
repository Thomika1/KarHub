package business

import (
	"errors"

	"github.com/Thomika1/KarHub/pkg/core/crud"
	productEnums "github.com/Thomika1/KarHub/pkg/domains/shared/enums/product"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Product struct {
	crud.BaseModel
	ProductData
}

type ProductData struct {
	Name              *string               `json:"name,omitempty" gorm:"type:varchar(255);uniqueIndex:idx_product_name_category" validate:"required"`
	Category          productEnums.Category `json:"category,omitempty" gorm:"type:varchar(255);uniqueIndex:idx_product_name_category" validate:"required"`
	CurrentStock      *int                  `json:"currentStock,omitempty" gorm:"type:int" validate:"gte=0"`
	MinimumStock      *int                  `json:"minimumStock,omitempty" gorm:"type:int" validate:"gte=0"`
	AverageDailySales *int                  `json:"averageDailySales,omitempty" gorm:"type:int" validate:"gte=0"`
	LeadTimeDays      *int                  `json:"leadTimeDays,omitempty" gorm:"type:int" validate:"gte=0"`
	UnitCost          *int                  `json:"unitCost,omitempty" gorm:"type:int" validate:"gte=0"`
	CriticalityLevel  *int                  `json:"criticalityLevel,omitempty" gorm:"type:int" validate:"gte=1,lte=5"`
}

func (p ProductData) ToModel() *Product {
	return &Product{
		ProductData: p,
	}
}

func (Product) TableName() string {
	return "products"
}

func (p ProductData) Validate() error {
	return validate.Struct(p)
}

func (p ProductData) ValidatePartial() error {
	if p.Name != nil && (*p.Name == "") {
		return errors.New("name cannot be empty")
	}
	if p.CurrentStock != nil && *p.CurrentStock < 0 {
		return errors.New("currentStock cannot be negative")
	}
	if p.MinimumStock != nil && *p.MinimumStock < 0 {
		return errors.New("minimumStock cannot be negative")
	}
	if p.AverageDailySales != nil && *p.AverageDailySales < 0 {
		return errors.New("averageDailySales cannot be negative")
	}
	if p.LeadTimeDays != nil && *p.LeadTimeDays < 0 {
		return errors.New("leadTimeDays cannot be negative")
	}
	if p.UnitCost != nil && *p.UnitCost < 0 {
		return errors.New("unitCost cannot be negative")
	}
	if p.CriticalityLevel != nil {
		if *p.CriticalityLevel < 1 || *p.CriticalityLevel > 5 {
			return errors.New("criticalityLevel must be between 1 and 5")
		}
	}
	return nil
}

func (p ProductData) ToUpdateMap() map[string]interface{} {
	m := make(map[string]interface{})
	if p.Name != nil {
		m["name"] = *p.Name
	}
	if p.Category != "" {
		m["category"] = p.Category
	}
	if p.CurrentStock != nil {
		m["current_stock"] = *p.CurrentStock
	}
	if p.MinimumStock != nil {
		m["minimum_stock"] = *p.MinimumStock
	}
	if p.AverageDailySales != nil {
		m["average_daily_sales"] = *p.AverageDailySales
	}
	if p.LeadTimeDays != nil {
		m["lead_time_days"] = *p.LeadTimeDays
	}
	if p.UnitCost != nil {
		m["unit_cost"] = *p.UnitCost
	}
	if p.CriticalityLevel != nil {
		m["criticality_level"] = *p.CriticalityLevel
	}
	return m
}
