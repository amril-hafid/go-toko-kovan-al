package news

import (
	"time"

	"gorm.io/gorm"
)

type News struct {
	ID        uint `gorm:"primaryKey"`
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
