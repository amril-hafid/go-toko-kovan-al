package product

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	FindAllWeb(products *[]Product, paging *Paging, tx *gorm.DB) error
	FindWebHome(product *[]Product, limit int, tx *gorm.DB) error
	FindByIDWeb(product *Product, id uint, tx *gorm.DB) error
	FindAllProduct(products *[]Product, tx *gorm.DB) error
	FindByIDProduct(product *Product, id uint, tx *gorm.DB) error
	SaveProduct(product *Product, tx *gorm.DB) error
	UpdateProduct(product *Product, tx *gorm.DB) error
	DeleteSoftProduct(product *Product, id uint, tx *gorm.DB) error
	DeleteProduct(product *Product, id uint, tx *gorm.DB) error
	FindAllDeletedAtProducts(products *[]Product, tx *gorm.DB) error
	FindByIDDeletedAtProduct(product *Product, id uint, tx *gorm.DB) error
	UpdateDeletedAtProduct(product *Product, id uint, tx *gorm.DB) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAllWeb(products *[]Product, paging *Paging, tx *gorm.DB) error {

	offset := (paging.Page - 1) * paging.Limit

	query := tx.Model(&Product{})

	if paging.Search != "" {
		query = query.Where("Name LIKE ?", "%"+paging.Search+"%")
	}

	if paging.Category != "" {
		query = query.Where("category = ?", paging.Category)
	}

	query.Count(&paging.Total)
	query.Offset(offset).Limit(paging.Limit).Find(&products)

	return query.Error
}

func (r *repository) FindWebHome(product *[]Product, limit int, tx *gorm.DB) error {
	return tx.Limit(limit).Find(product).Error

}

func (r *repository) FindByIDWeb(product *Product, id uint, tx *gorm.DB) error {
	return tx.Where("id = ?", id).First(product).Error
}

func (r *repository) FindAllProduct(products *[]Product, tx *gorm.DB) error {
	return tx.Find(&products).Error
}

func (r *repository) FindByIDProduct(product *Product, id uint, tx *gorm.DB) error {
	return tx.Where("id = ?", id).First(product).Error
}

func (r *repository) SaveProduct(product *Product, tx *gorm.DB) error {
	return tx.Create(&product).Error
}

func (r *repository) UpdateProduct(product *Product, tx *gorm.DB) error {
	return tx.Save(&product).Error
}

func (r *repository) DeleteSoftProduct(product *Product, id uint, tx *gorm.DB) error {
	return tx.Where("id = ?", id).Delete(&product).Error
}

func (r *repository) DeleteProduct(product *Product, id uint, tx *gorm.DB) error {
	return tx.Unscoped().Where("id = ?", id).Delete(&product).Error
}

func (r *repository) FindAllDeletedAtProducts(products *[]Product, tx *gorm.DB) error {
	return tx.Unscoped().Where("deleted_at > 0").Find(&products).Error
}

func (r *repository) FindByIDDeletedAtProduct(product *Product, id uint, tx *gorm.DB) error {
	err := tx.Unscoped().Where("id = ?", id).First(product).Error

	if product.ID == 0 {
		return errors.New("Data tidak di temukan.")
	}

	return err
}

func (r *repository) UpdateDeletedAtProduct(product *Product, id uint, tx *gorm.DB) error {
	return tx.Unscoped().Model(&product).Where("id", id).Update("deleted_at", nil).Error
}
