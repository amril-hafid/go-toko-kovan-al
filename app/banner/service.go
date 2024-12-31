package banner

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Service interface {
	CreateBanner(ctx context.Context, input InputBanner) (Banner, error)
	GetAllBannder(ctx context.Context) ([]BannerViews, error)
	GetBannerByLimit(ctx context.Context, limit int) ([]Banner, error)
	GetBannerByID(ctx context.Context, id uint) (Banner, error)
	UpdateBanner(ctx context.Context, input UpdateBanner) (Banner, error)
	DeleteBanner(ctx context.Context, id uint) (Banner, error)
}

type service struct {
	repository Repository
	validate   *validator.Validate
	db         *gorm.DB
}

func NewService(repository Repository, validate *validator.Validate, db *gorm.DB) *service {
	return &service{repository, validate, db}
}

func (s *service) CreateBanner(ctx context.Context, input InputBanner) (Banner, error) {
	var banner Banner
	// var resultBanner Banner
	var BannerNil Banner

	// var resultBannerUpdateIsPrimary Banner

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return BannerNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return BannerNil, ctx.Err()
	default:

		err := s.validate.Struct(input)
		if err != nil {
			return BannerNil, errors.New("Isi form dengan benar!")
		}

		banner.Name = input.Name
		banner.NameImage = input.NameImage
		banner.ImageURL = input.ImageURL
		banner.URL = input.URL
		banner.Description = input.Description
		banner.IsPrimary = input.IsPrimary
		banner.Status = input.Status

		err = s.repository.Save(&banner, tx)

		if err != nil {
			tx.Rollback()
			return BannerNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return BannerNil, err
		}
	}

	return banner, nil
}

func (s *service) GetAllBannder(ctx context.Context) ([]BannerViews, error) {
	var banners []Banner
	var bannerViews []BannerViews
	var bannerNil []BannerViews

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return bannerNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return bannerNil, ctx.Err()
	default:

		err := s.repository.FindAll(&banners, tx)
		if err != nil {
			tx.Rollback()
			return bannerNil, err
		}

		for i, banner := range banners {
			var bannerRowViews BannerViews
			bannerRowViews.Index = i + 1
			bannerRowViews.ID = banner.ID
			bannerRowViews.Name = banner.Name
			bannerRowViews.ImageURL = banner.ImageURL
			bannerRowViews.URL = banner.URL
			bannerRowViews.Description = banner.Description
			bannerRowViews.IsPrimary = banner.IsPrimary
			bannerRowViews.Status = banner.Status
			bannerRowViews.CreatedAt = banner.CreatedAt
			bannerRowViews.UpdatedAt = banner.UpdatedAt
			bannerViews = append(bannerViews, bannerRowViews)
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return bannerNil, err
		}
	}

	return bannerViews, nil

}

func (s *service) GetBannerByID(ctx context.Context, id uint) (Banner, error) {
	var banner Banner
	var bannerNil Banner

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return bannerNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return bannerNil, ctx.Err()
	default:

		err := s.repository.FindByID(&banner, id, tx)
		if err != nil {
			tx.Rollback()
			return bannerNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return bannerNil, err
		}
	}

	return banner, nil
}

func (s *service) GetBannerByLimit(ctx context.Context, limit int) ([]Banner, error) {
	var banner []Banner
	var bannerNil []Banner

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return bannerNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return bannerNil, ctx.Err()
	default:

		err := s.repository.FindbyLimit(&banner, limit, tx)
		if err != nil {
			tx.Rollback()
			return bannerNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return bannerNil, err
		}
	}

	return banner, nil
}

func (s *service) UpdateBanner(ctx context.Context, input UpdateBanner) (Banner, error) {
	var banner Banner
	var bannerResult Banner
	var bannerNil Banner

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return bannerNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return bannerNil, ctx.Err()
	default:

		err := s.repository.FindByID(&banner, input.ID, tx)
		if err != nil {
			tx.Rollback()
			return bannerNil, err
		}

		bannerResult.ID = banner.ID
		bannerResult.Name = input.Name
		bannerResult.ImageURL = banner.ImageURL
		bannerResult.URL = input.URL
		bannerResult.Description = input.Description
		bannerResult.IsPrimary = input.IsPrimary
		bannerResult.Status = input.Status

		err = s.repository.Update(&bannerResult, banner.ID, tx)
		if err != nil {
			tx.Rollback()
			return bannerNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return bannerNil, err
		}
	}

	return bannerResult, nil
}

func (s *service) DeleteBanner(ctx context.Context, id uint) (Banner, error) {
	var banner Banner
	var bannerResult Banner
	var bannerNil Banner

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return bannerNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return bannerNil, ctx.Err()
	default:

		err := s.repository.FindByID(&banner, id, tx)
		if err != nil {
			tx.Rollback()
			return bannerNil, err
		}

		err = s.repository.Delete(&bannerResult, banner.ID, tx)
		if err != nil {
			tx.Rollback()
			return bannerNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return bannerNil, err
		}
	}

	return bannerResult, nil
}
