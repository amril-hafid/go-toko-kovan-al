package banner

import "gorm.io/gorm"

type Repository interface {
	FindAll(banner *[]Banner, tx *gorm.DB) error
	FindByID(banner *Banner, id uint, tx *gorm.DB) error
	FindbyLimit(product *[]Banner, limit int, tx *gorm.DB) error
	Save(banner *Banner, tx *gorm.DB) error
	Update(banner *Banner, id uint, tx *gorm.DB) error
	Delete(banner *Banner, id uint, tx *gorm.DB) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll(banner *[]Banner, tx *gorm.DB) error {
	return tx.Find(&banner).Error
}
func (r *repository) FindByID(banner *Banner, id uint, tx *gorm.DB) error {
	return tx.Where("id = ?", id).First(banner).Error
}

func (r *repository) FindbyLimit(banner *[]Banner, limit int, tx *gorm.DB) error {
	return tx.Limit(limit).Find(banner).Error

}

func (r *repository) Save(banner *Banner, tx *gorm.DB) error {
	return tx.Create(&banner).Error
}

func (r *repository) Update(banner *Banner, id uint, tx *gorm.DB) error {
	return tx.Save(&banner).Where("id = ?", id).Error
}

func (r *repository) Delete(banner *Banner, id uint, tx *gorm.DB) error {
	return tx.Where("id = ?", id).Delete(&banner).Error
}
