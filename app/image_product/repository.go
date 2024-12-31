package image_product

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(image *[]ImageProduct, tx *gorm.DB) error
	FindByID(image *ImageProduct, id uint, tx *gorm.DB) error

	FindAllByIDProduct(image *[]ImageProduct, idProduct uint, tx *gorm.DB) error

	FindAllByIDProductEndPrimary(image *ImageProduct, idProduct uint, tx *gorm.DB) error
	FindAllByIDProductEndNoPrimary(image *[]ImageProduct, idProduct uint, tx *gorm.DB) error

	Save(image *ImageProduct, tx *gorm.DB) error
	UpdateImage(image *ImageProduct, id uint, tx *gorm.DB) error

	DeleteSoftImageProduct(image *ImageProduct, id uint, tx *gorm.DB) error
	DeleteImageProduct(image *ImageProduct, id uint, tx *gorm.DB) error

	FindAllDeletedAtImageProducts(image *[]ImageProduct, tx *gorm.DB) error
	FindByIDDeletedAtImageProduct(image *ImageProduct, id uint, tx *gorm.DB) error
	UpdateDeletedAtImageProduct(image *ImageProduct, id uint, tx *gorm.DB) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}
func (r *repository) FindAll(image *[]ImageProduct, tx *gorm.DB) error {
	return tx.Find(&image).Error
}
func (r *repository) FindByID(image *ImageProduct, id uint, tx *gorm.DB) error {
	return tx.Where("id = ?", id).First(image).Error
}
func (r *repository) FindAllByIDProduct(image *[]ImageProduct, idProduct uint, tx *gorm.DB) error {
	return tx.Where("id_product = ?", idProduct).Find(image).Error
}
func (r *repository) FindAllByIDProductEndPrimary(image *ImageProduct, idProduct uint, tx *gorm.DB) error {
	return tx.Where("id_product = ?", idProduct).Where("is_primary = 1").First(image).Error
}

func (r *repository) FindAllByIDProductEndNoPrimary(image *[]ImageProduct, idProduct uint, tx *gorm.DB) error {
	return tx.Where("id_product = ?", idProduct).Where("is_primary = 0").Find(image).Error
}

func (r *repository) Save(image *ImageProduct, tx *gorm.DB) error {
	return tx.Create(&image).Error
}
func (r *repository) UpdateImage(image *ImageProduct, id uint, tx *gorm.DB) error {
	return tx.Save(&image).Where("id = ?", id).Error
}

func (r *repository) DeleteSoftImageProduct(image *ImageProduct, id uint, tx *gorm.DB) error {
	return tx.Where("id = ?", id).Delete(&image).Error
}
func (r *repository) DeleteImageProduct(image *ImageProduct, id uint, tx *gorm.DB) error {
	return tx.Unscoped().Where("id = ?", id).Delete(&image).Error
}
func (r *repository) FindAllDeletedAtImageProducts(image *[]ImageProduct, tx *gorm.DB) error {
	return tx.Unscoped().Where("deleted_at > 0").Find(&image).Error
}
func (r *repository) FindByIDDeletedAtImageProduct(image *ImageProduct, id uint, tx *gorm.DB) error {
	err := tx.Unscoped().Where("id = ?", id).Find(&image).Error

	if image.ID == 0 {
		return errors.New("Data tidak di temukan.")
	}

	return err
}
func (r *repository) UpdateDeletedAtImageProduct(image *ImageProduct, id uint, tx *gorm.DB) error {
	return tx.Unscoped().Model(&image).Where("id", id).Update("deleted_at", nil).Error
}
