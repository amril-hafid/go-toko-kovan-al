package banner

import "time"

type Banner struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	NameImage   string
	ImageURL    string
	URL         string
	Description string
	IsPrimary   int
	Status      string
	IDProduct   uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type BannerViews struct {
	Index       int
	ID          uint
	Name        string
	ImageURL    string
	URL         string
	Description string
	IsPrimary   int
	Status      string
	IDProduct   uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
