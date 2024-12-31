package image_product

import (
	"time"
)

type ImageProduct struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	URL       string
	IsPrimary int
	Status    string
	IDProduct uint
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt gorm.DeletedAt `gorm:"index"`
}
