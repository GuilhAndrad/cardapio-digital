package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name        string  `gorm:"type:varchar(100);not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"price"`
	ImageURL    string          `gorm:"type:text" json:"image_url"`
	CategoryID  uint    `gorm:"not null" json:"category_id"`
}