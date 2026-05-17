package models

import (
	"fmt"
	"strings"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Restaurant struct {
	gorm.Model
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Slug        string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	LogoURL     string     `gorm:"type:text" json:"logo_url"`
	Categories  []Category `gorm:"foreignKey:RestaurantID;constraint:OnDelete:CASCADE" json:"categories,omitempty"`
}

// BeforeCreate é o Hook do GORM para gerar o slug antes de salvar no banco
func (r *Restaurant) BeforeCreate(tx *gorm.DB) (err error) {
	if r.Slug == "" {
		if strings.TrimSpace(r.Name) == "" {
			return fmt.Errorf("restaurant name is required when slug is empty")
		}
		r.Slug = slug.Make(r.Name)
	} else {
		r.Slug = slug.Make(r.Slug)
	}
	return nil
}
