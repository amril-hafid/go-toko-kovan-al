package material

import (
	"time"

	"gorm.io/gorm"
)

type MaterialType struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type MaterialProduct struct {
	ID         uint `gorm:"primaryKey"`
	IDMaterial uint
	IDProduct  uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
