package category

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CategoryView struct {
	Index     int
	ID        uint
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
