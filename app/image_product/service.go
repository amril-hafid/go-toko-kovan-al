package image_product

import (
	"context"
	"errors"
	"fmt"
	"go-toko-kovan-al/helper"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Service interface {
	CreateImageProduc(ctx context.Context, input InputImageProduct) (ImageProduct, error)
	DeleteImageProduct(ctx context.Context, id uint) (ImageProduct, error)
}

type service struct {
	repository Repository
	validate   *validator.Validate
	db         *gorm.DB
}

func NewService(repository Repository, validate *validator.Validate, db *gorm.DB) *service {
	return &service{repository, validate, db}
}

func (s *service) CreateImageProduc(ctx context.Context, input InputImageProduct) (ImageProduct, error) {
	var image ImageProduct
	var imageNil ImageProduct
	var resultImage ImageProduct
	var resultImageUpdateIsPrimary ImageProduct

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return imageNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return imageNil, ctx.Err()
	default:

		err := s.validate.Struct(input)
		fmt.Println("serive image product err 1 : ", err)
		if err != nil {
			return imageNil, errors.New("Isi form dengan benar!")
		}

		s.repository.FindAllByIDProductEndPrimary(&resultImage, input.IDProduct, tx)
		fmt.Println("serive image product 2 : ", resultImage)

		if resultImage.ID > 0 {

			resultImageUpdateIsPrimary.ID = resultImage.ID
			resultImageUpdateIsPrimary.Name = resultImage.Name
			resultImageUpdateIsPrimary.URL = resultImage.URL
			resultImageUpdateIsPrimary.Status = resultImage.Status
			resultImageUpdateIsPrimary.IDProduct = resultImage.IDProduct
			resultImageUpdateIsPrimary.IsPrimary = int(0)
			err = s.repository.UpdateImage(&resultImageUpdateIsPrimary, resultImage.ID, tx)
			fmt.Println("serive image product err 3 : ", err)

			if err != nil {
				tx.Rollback()
				return imageNil, err
			}
		}

		image.Name = input.NameFile
		image.URL = input.URL
		image.IsPrimary = input.IsPrimary
		image.Status = input.Status
		image.IDProduct = input.IDProduct

		err = s.repository.Save(&image, tx)
		fmt.Println("serive image product err 4 : ", err)

		if err != nil {
			tx.Rollback()
			return imageNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return imageNil, err
		}
	}

	return image, nil
}

func (s *service) DeleteImageSoft(ctx context.Context, id uint) (ImageProduct, error) {
	var resultImage ImageProduct
	var resultImageUpdate ImageProduct

	var image ImageProduct
	var imagetNil ImageProduct

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return imagetNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return imagetNil, ctx.Err()
	default:

		s.repository.FindByID(&resultImage, id, tx)

		resultImageUpdate.ID = resultImage.ID
		resultImageUpdate.Name = resultImage.Name
		resultImageUpdate.URL = resultImage.URL
		resultImageUpdate.Status = resultImage.Status
		resultImageUpdate.IDProduct = resultImage.IDProduct
		resultImageUpdate.IsPrimary = resultImage.IsPrimary
		err := s.repository.UpdateImage(&resultImageUpdate, resultImage.ID, tx)
		fmt.Println("serive image product err 3 : ", err)
		if err != nil {
			tx.Rollback()
			return imagetNil, err
		}

		err = s.repository.DeleteSoftImageProduct(&image, id, tx)
		if err != nil {
			tx.Rollback()
			return imagetNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return imagetNil, err
		}
	}

	return image, nil
}
func (s *service) DeleteImageProduct(ctx context.Context, id uint) (ImageProduct, error) {
	var product ImageProduct
	var productrNil ImageProduct

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return productrNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return productrNil, ctx.Err()
	default:

		err := s.repository.DeleteSoftImageProduct(&product, id, tx)
		if err != nil {
			tx.Rollback()
			return productrNil, err
		}

		err = helper.DeleteFile(product.URL)
		if err != nil {
			return productrNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return productrNil, err
		}
	}

	return product, nil
}
