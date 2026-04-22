package domain

import "github.com/shopspring/decimal"

type Product struct {
	BaseModel
	SKU         string          `gorm:"uniqueIndex;size:64;not null" json:"sku"`
	Name        string          `gorm:"size:255;not null;index" json:"name"`
	Slug        string          `gorm:"uniqueIndex;size:255;not null" json:"slug"`
	Description string          `gorm:"type:text" json:"description"`
	Price       decimal.Decimal `gorm:"type:numeric(12,2);not null" json:"price"`
	Stock       int             `gorm:"not null;default:0" json:"stock"`
	IsActive    bool            `gorm:"not null;default:true;index" json:"is_active"`
}

func (Product) TableName() string { return "products" }
