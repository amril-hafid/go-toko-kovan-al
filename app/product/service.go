package product

import (
	"context"
	"errors"
	"fmt"
	"go-toko-kovan-al/app/category"
	"go-toko-kovan-al/app/image_product"
	"go-toko-kovan-al/config"
	"go-toko-kovan-al/helper"
	"html/template"
	"math"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Service interface {
	GetAllWebProduct(ctx context.Context, paging *Paging) (ProductListWebViews, error)
	GetbyIDWebProduct(ctx context.Context, id uint) (ProductDetail, error)
	GetAllWebProductHOme(ctx context.Context, limit int) (ProductListWebViews, error)
	GetProductWebByID(ctx context.Context, id uint) (ProductEndImageDetail, error)
	GetAllProduct(ctx context.Context) ([]ProductView, error)
	GetProductByID(ctx context.Context, id uint) (ProductEndImageDetail, error)

	GetAllProductDeleted(ctx context.Context) ([]ProductView, error)
	GetProductByIDDeleted(ctx context.Context, id uint) (Product, error)

	CreateProduct(ctx context.Context, input InputProduct, userID uint) (Product, error)
	UpdateProduct(ctx context.Context, input UpdateProduct) (Product, error)
	DeleteProductSoft(ctx context.Context, id uint) (Product, error)
	DeleteProduct(ctx context.Context, id uint) (Product, error)
	RestoreProduct(ctx context.Context, id uint) (Product, error)
}

type service struct {
	repository         Repository
	repositoryCategory category.Repository
	imageRepository    image_product.Repository
	validate           *validator.Validate
	db                 *gorm.DB
	conf               config.Config
}

func NewService(repository Repository, repositoryCategory category.Repository, imageRepository image_product.Repository, validate *validator.Validate, db *gorm.DB, conf config.Config) *service {
	return &service{repository, repositoryCategory, imageRepository, validate, db, conf}
}

func (s *service) GetAllWebProduct(ctx context.Context, paging *Paging) (ProductListWebViews, error) {
	var productResult []ProductWebList
	var products []Product
	var productViews ProductListWebViews
	var productViewNil ProductListWebViews

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productViewNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productViewNil, ctx.Err()
	default:
		err := s.repository.FindAllWeb(&products, paging, tx)
		if err != nil {
			tx.Rollback()
			return productViewNil, err
		}
		for i, product := range products {
			var productRow ProductWebList
			productRow.Index = i + 1
			productRow.ID = product.ID
			productRow.SKU = product.SKU
			productRow.Name = product.Name
			productRow.Price = helper.FormatRupiah(product.Price)
			productRow.Stock = product.Stock

			var image image_product.ImageProduct
			s.imageRepository.FindAllByIDProductEndPrimary(&image, product.ID, tx)
			// if err != nil {
			// 	tx.Rollback()
			// 	return productViewNil, err
			// }

			// Pastikan URL tidak kosong
			if image.URL != "" {
				baseURL := fmt.Sprintf("http://%s:%s/image/%s", s.conf.Srv.Host, s.conf.Srv.Port, image.URL)
				productRow.ImagePrimary = template.HTML(baseURL)
			} else {
				productRow.ImagePrimary = template.HTML("") // Gambar placeholder
			}
			productResult = append(productResult, productRow)
		}

		lastPage := int(math.Ceil(float64(paging.Total) / float64(paging.Limit)))

		productViews.Product = productResult
		productViews.Total = int(paging.Total)
		productViews.Page = paging.Page
		productViews.LastPage = lastPage
		productViews.PrevPage = paging.Page - 1
		productViews.NextPage = paging.Page + 1
		productViews.Search = paging.Search
		productViews.Category = paging.Category

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productViewNil, err
		}
	}

	return productViews, nil
}

func (s *service) GetAllWebProductHOme(ctx context.Context, limit int) (ProductListWebViews, error) {
	var productResult []ProductWebList
	var products []Product
	var productViews ProductListWebViews
	var productViewNil ProductListWebViews

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productViewNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productViewNil, ctx.Err()
	default:
		err := s.repository.FindWebHome(&products, limit, tx)
		if err != nil {
			tx.Rollback()
			return productViewNil, err
		}
		for i, product := range products {
			var productRow ProductWebList
			productRow.Index = i + 1
			productRow.ID = product.ID
			productRow.SKU = product.SKU
			productRow.Name = product.Name
			productRow.Price = helper.FormatRupiah(product.Price)
			productRow.Stock = product.Stock

			var image image_product.ImageProduct
			s.imageRepository.FindAllByIDProductEndPrimary(&image, product.ID, tx)

			if image.URL != "" {
				baseURL := fmt.Sprintf("http://%s:%s/image/%s", s.conf.Srv.Host, s.conf.Srv.Port, image.URL)
				productRow.ImagePrimary = template.HTML(baseURL)
			} else {
				productRow.ImagePrimary = template.HTML("") // Gambar placeholder
			}
			productResult = append(productResult, productRow)
		}

		productViews.Product = productResult
		// productViews.Total = int(paging.Total)
		// productViews.Page = paging.Page
		// productViews.LastPage = lastPage
		// productViews.PrevPage = paging.Page - 1
		// productViews.NextPage = paging.Page + 1
		productViews.Search = ""
		// productViews.Category = paging.Category

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productViewNil, err
		}
	}

	return productViews, nil
}

func (s *service) GetbyIDWebProduct(ctx context.Context, id uint) (ProductDetail, error) {
	var product Product
	var imageProduct []image_product.ImageProduct
	var productDetail ProductDetail
	var productNil ProductDetail
	var productWebDetail ProductWebDetail

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productNil, ctx.Err()
	default:

		err := s.repository.FindByIDWeb(&product, id, tx)
		if err != nil {
			tx.Rollback()
			return productNil, err
		}

		productWebDetail.Index = product.Index
		productWebDetail.ID = product.ID
		productWebDetail.SKU = product.SKU
		productWebDetail.Name = product.Name
		productWebDetail.Price = helper.FormatRupiah(product.Price)
		productWebDetail.Stock = product.Stock
		productWebDetail.ShortDescription = product.ShortDescription
		productWebDetail.LongDescription = product.LongDescription
		productWebDetail.SizeType = product.SizeType
		productWebDetail.Long = product.Long
		productWebDetail.Wide = product.Wide
		productWebDetail.Tall = product.Tall
		productWebDetail.Diameter = product.Diameter
		productWebDetail.IDKategori = product.IDKategori
		productWebDetail.Status = product.Status
		productWebDetail.CreatedAt = product.CreatedAt
		productWebDetail.UpdatedAt = product.UpdatedAt
		productWebDetail.DeletedAt = product.DeletedAt

		err = s.imageRepository.FindAllByIDProduct(&imageProduct, id, tx)
		if err != nil {
			tx.Rollback()
			return productNil, err
		}
		productDetail.Product = productWebDetail
		productDetail.Image = imageProduct

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productNil, err
		}
	}

	return productDetail, nil
}

func (s *service) GetProductWebByID(ctx context.Context, id uint) (ProductEndImageDetail, error) {
	var product Product
	var image image_product.ImageProduct
	var images []image_product.ImageProduct
	var productResult ProductEndImageDetail
	var productNil ProductEndImageDetail

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productNil, ctx.Err()
	default:

		err := s.repository.FindByIDProduct(&product, id, tx)
		if err != nil {
			tx.Rollback()
			return productNil, err
		}

		s.imageRepository.FindAllByIDProductEndPrimary(&image, product.ID, tx)

		s.imageRepository.FindAllByIDProductEndNoPrimary(&images, product.ID, tx)

		productResult.Product = product
		productResult.Image = image_product.ImageProductFormatter(image, images)

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productNil, err
		}
	}

	return productResult, nil
}

func (s *service) GetAllProduct(ctx context.Context) ([]ProductView, error) {

	// var productResult []Product
	var products []Product
	var productViews []ProductView
	var productViewNil []ProductView

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productViewNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productViewNil, ctx.Err()
	default:
		err := s.repository.FindAllProduct(&products, tx)
		if err != nil {
			tx.Rollback()
			return productViewNil, err
		}

		for i, product := range products {
			var productResult ProductView
			productResult.ID = product.ID
			productResult.Index = i + 1
			productResult.Name = product.Name
			productResult.SKU = product.SKU
			productResult.Price = product.Price
			productResult.Stock = product.Stock
			productResult.ShortDescription = product.ShortDescription
			productResult.LongDescription = product.LongDescription
			productResult.SizeType = product.SizeType
			productResult.Long = product.Long
			productResult.Diameter = product.Diameter
			productResult.Tall = product.Tall
			productResult.Wide = product.Wide
			productResult.Status = product.Status

			var categoryProduct category.Category
			s.repositoryCategory.FindByIDCategory(&categoryProduct, product.IDKategori, tx)

			productResult.IDKategori = product.IDKategori
			productResult.CategoryName = categoryProduct.Name

			// Format waktu create dan update
			timeProductCreate, _ := helper.DatetimeToFormatIndo(product.CreatedAt)
			productResult.CreatedAt = timeProductCreate
			timeProductUpdate, _ := helper.DatetimeToFormatIndo(product.UpdatedAt)
			productResult.UpdatedAt = timeProductUpdate

			productViews = append(productViews, productResult)
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productViewNil, err
		}
	}

	return productViews, nil
}
func (s *service) GetProductByID(ctx context.Context, id uint) (ProductEndImageDetail, error) {
	var product Product
	var image image_product.ImageProduct
	var images []image_product.ImageProduct
	var productResult ProductEndImageDetail
	var productNil ProductEndImageDetail

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productNil, ctx.Err()
	default:

		err := s.repository.FindByIDProduct(&product, id, tx)
		if err != nil {
			tx.Rollback()
			return productNil, err
		}

		s.imageRepository.FindAllByIDProductEndPrimary(&image, product.ID, tx)

		s.imageRepository.FindAllByIDProductEndNoPrimary(&images, product.ID, tx)

		productResult.Product = product
		productResult.Image = image_product.ImageProductFormatter(image, images)

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productNil, err
		}
	}

	return productResult, nil
}
func (s *service) GetAllProductDeleted(ctx context.Context) ([]ProductView, error) {

	// var productResult []Product
	var products []Product
	var productViews []ProductView
	var productViewNil []ProductView

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productViewNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productViewNil, ctx.Err()
	default:
		err := s.repository.FindAllDeletedAtProducts(&products, tx)
		if err != nil {
			tx.Rollback()
			return productViewNil, err
		}

		for i, product := range products {
			var productResult ProductView
			productResult.ID = product.ID
			productResult.Index = i + 1
			productResult.Name = product.Name
			productResult.SKU = product.SKU
			productResult.Price = product.Price
			productResult.Stock = product.Stock
			productResult.ShortDescription = product.ShortDescription
			productResult.LongDescription = product.LongDescription
			productResult.SizeType = product.SizeType
			productResult.Long = product.Long
			productResult.Diameter = product.Diameter
			productResult.Tall = product.Tall
			productResult.Wide = product.Wide
			productResult.Status = product.Status

			var categoryProduct category.Category
			err := s.repositoryCategory.FindByIDCategory(&categoryProduct, product.IDKategori, tx)
			if err != nil {
				tx.Rollback()
				return productViewNil, err
			}

			productResult.IDKategori = product.IDKategori
			productResult.CategoryName = categoryProduct.Name

			// Format waktu create dan update
			timeProductCreate, _ := helper.DatetimeToFormatIndo(product.CreatedAt)
			productResult.CreatedAt = timeProductCreate
			timeProductUpdate, _ := helper.DatetimeToFormatIndo(product.UpdatedAt)
			productResult.UpdatedAt = timeProductUpdate

			productViews = append(productViews, productResult)
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productViewNil, err
		}
	}

	return productViews, nil
}
func (s *service) GetProductByIDDeleted(ctx context.Context, id uint) (Product, error) {
	var product Product
	var productNil Product

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productNil, ctx.Err()
	default:

		err := s.repository.FindByIDDeletedAtProduct(&product, id, tx)
		if err != nil {
			tx.Rollback()
			return productNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productNil, err
		}
	}

	return product, nil
}
func (s *service) CreateProduct(ctx context.Context, input InputProduct, userID uint) (Product, error) {
	var product Product
	var productNil Product

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productNil, ctx.Err()
	default:

		err := s.validate.Struct(input)
		if err != nil {
			return productNil, errors.New("isi form dengan benar!")
		}

		product.SKU = input.SKU
		product.Name = input.Name
		product.Price = input.Price
		product.ShortDescription = input.ShortDescription
		product.LongDescription = input.LongDescription
		product.Stock = input.Stock
		product.SizeType = input.SizeType
		product.Long = input.Long
		product.Wide = input.Wide
		product.Tall = input.Tall
		product.Diameter = input.Diameter
		product.IDKategori = uint(input.IDKategori)
		product.Status = input.Status

		err = s.repository.SaveProduct(&product, tx)
		if err != nil {
			tx.Rollback()
			return productNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productNil, err
		}
	}

	return product, nil
}
func (s *service) UpdateProduct(ctx context.Context, input UpdateProduct) (Product, error) {
	var product Product
	var productRow Product
	var productNil Product

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productNil, ctx.Err()
	default:
		err := s.validate.Struct(input)
		if err != nil {
			return productNil, errors.New("isi form dengan benar")
		}

		err = s.repository.FindByIDProduct(&productRow, input.ID, tx)
		if err != nil {
			tx.Rollback()
			return productNil, err
		}

		product.ID = productRow.ID
		product.SKU = input.SKU
		product.Name = input.Name
		product.Price = input.Price
		product.Stock = input.Stock
		product.ShortDescription = input.ShortDescription
		product.LongDescription = input.LongDescription
		product.SizeType = input.SizeType
		product.Long = input.Long
		product.Wide = input.Wide
		product.Tall = input.Tall
		product.Diameter = input.Diameter
		product.IDKategori = input.IDKategori
		product.Status = input.Status

		err = s.repository.UpdateProduct(&product, tx)
		if err != nil {
			tx.Rollback()
			return productNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productNil, err
		}
	}

	return product, nil
}
func (s *service) DeleteProductSoft(ctx context.Context, id uint) (Product, error) {
	var product Product
	var productNil Product

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productNil, ctx.Err()
	default:
		err := s.repository.DeleteSoftProduct(&product, id, tx)
		if err != nil {
			tx.Rollback()
			return productNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productNil, err
		}
	}

	return product, nil
}
func (s *service) DeleteProduct(ctx context.Context, id uint) (Product, error) {
	var product Product
	var productrNil Product

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productrNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productrNil, ctx.Err()
	default:

		err := s.repository.DeleteProduct(&product, id, tx)
		if err != nil {
			tx.Rollback()
			return productrNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productrNil, err
		}
	}

	return product, nil
}
func (s *service) RestoreProduct(ctx context.Context, id uint) (Product, error) {
	var product Product
	var productRow Product
	var productNil Product

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productNil, ctx.Err()
	default:

		err := s.repository.FindByIDDeletedAtProduct(&product, id, tx)
		if err != nil {
			tx.Rollback()
			return productNil, err
		}

		err = s.repository.UpdateDeletedAtProduct(&productRow, product.ID, tx)
		if err != nil {
			tx.Rollback()
			return productNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productNil, err
		}
	}

	return productRow, nil
}
