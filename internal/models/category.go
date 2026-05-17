package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name         string `gorm:"type:varchar(50);not null" json:"name"`
	RestaurantID uint   `gorm:"not null" json:"restaurant_id"`
	Items        []Item `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}