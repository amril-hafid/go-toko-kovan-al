package category

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	FindAllCategory(categorys *[]Category, tx *gorm.DB) error
	FindByIDCategory(category *Category, id uint, tx *gorm.DB) error
	SaveCategory(category *Category, tx *gorm.DB) error
	UpdateCategory(category *Category, tx *gorm.DB) error
	DeleteSoftCategory(category *Category, id uint, tx *gorm.DB) error
	Delete(category *Category, id uint, tx *gorm.DB) error
	FindAllDeletedAtCategory(categorys *[]Category, tx *gorm.DB) error
	FindByIDDeletedAtCategory(category *Category, id uint, tx *gorm.DB) error
	UpdateDeletedAtCategory(category *Category, id uint, tx *gorm.DB) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAllCategory(categorys *[]Category, tx *gorm.DB) error {
	return tx.Find(&categorys).Error
}

func (r *repository) FindByIDCategory(category *Category, id uint, tx *gorm.DB) error {
	return tx.Where("id = ?", id).First(category).Error
}

func (r *repository) SaveCategory(category *Category, tx *gorm.DB) error {
	return tx.Create(&category).Error
}

func (r *repository) UpdateCategory(category *Category, tx *gorm.DB) error {
	return tx.Save(&category).Error
}

func (r *repository) DeleteSoftCategory(category *Category, id uint, tx *gorm.DB) error {
	return tx.Where("id = ?", id).Delete(&category).Error
}

func (r *repository) Delete(category *Category, id uint, tx *gorm.DB) error {
	return tx.Unscoped().Where("id = ?", id).Delete(&category).Error
}

func (r *repository) FindAllDeletedAtCategory(categorys *[]Category, tx *gorm.DB) error {
	return tx.Unscoped().Where("deleted_at > 0").Find(&categorys).Error
}

func (r *repository) FindByIDDeletedAtCategory(category *Category, id uint, tx *gorm.DB) error {
	err := tx.Unscoped().Where("id = ?", id).First(category).Error

	if category.ID == 0 {
		return errors.New("Data tidak di temukan.")
	}

	return err
}

func (r *repository) UpdateDeletedAtCategory(category *Category, id uint, tx *gorm.DB) error {
	return tx.Unscoped().Model(&category).Where("id", id).Update("deleted_at", nil).Error
}
