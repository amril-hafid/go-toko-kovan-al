package category

import (
	"context"
	"errors"
	"fmt"
	"go-toko-kovan-al/helper"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Service interface {
	GetAllCategory(ctx context.Context) ([]CategoryView, error)
	GetCategoryByID(ctx context.Context, id uint) (Category, error)
	GetAllCategoryDeleted(ctx context.Context) ([]CategoryView, error)
	GetCategoryByIDDeleted(ctx context.Context, id uint) (Category, error)
	CreateCategory(ctx context.Context, input InputCategory) (Category, error)
	UpdateCategory(ctx context.Context, input UpdateCategory) (Category, error)
	DeleteCategorySoft(ctx context.Context, id uint) (Category, error)
	DeleteCategory(ctx context.Context, id uint) (Category, error)
	RestoreCategory(ctx context.Context, id uint) (Category, error)
}

type service struct {
	repository Repository
	validate   *validator.Validate
	db         *gorm.DB
}

func NewService(repository Repository, validate *validator.Validate, db *gorm.DB) *service {
	return &service{repository, validate, db}
}

func (s *service) GetAllCategory(ctx context.Context) ([]CategoryView, error) {
	var categorys []Category
	var categoryViews []CategoryView
	var categoryNil []CategoryView

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return categoryNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return categoryNil, ctx.Err()
	default:
		err := s.repository.FindAllCategory(&categorys, tx)
		if err != nil {
			tx.Rollback()
			return categoryNil, err
		}

		for i, category := range categorys {
			var categoryView CategoryView
			categoryView.ID = category.ID
			categoryView.Index = i + 1
			categoryView.Name = category.Name
			categoryView.Status = category.Status

			// Format waktu create dan update
			timeCategoryCreate, _ := helper.DatetimeToFormatIndo(category.CreatedAt)
			categoryView.CreatedAt = timeCategoryCreate
			timeCategoryUpdate, _ := helper.DatetimeToFormatIndo(category.UpdatedAt)
			categoryView.UpdatedAt = timeCategoryUpdate

			categoryViews = append(categoryViews, categoryView)
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return categoryNil, err
		}
	}

	return categoryViews, nil
}

func (s *service) GetCategoryByID(ctx context.Context, id uint) (Category, error) {
	var category Category
	var categoryNil Category

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return categoryNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return categoryNil, ctx.Err()
	default:

		err := s.repository.FindByIDCategory(&category, id, tx)
		if err != nil {
			tx.Rollback()
			return categoryNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return categoryNil, err
		}
	}

	return category, nil
}

func (s *service) GetAllCategoryDeleted(ctx context.Context) ([]CategoryView, error) {
	var categorys []Category
	var categoryViews []CategoryView
	var categoryNil []CategoryView

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return categoryNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return categoryNil, ctx.Err()
	default:
		err := s.repository.FindAllDeletedAtCategory(&categorys, tx)
		if err != nil {
			tx.Rollback()
			return categoryNil, err
		}

		for i, category := range categorys {
			var categoryView CategoryView
			categoryView.ID = category.ID
			categoryView.Index = i + 1
			categoryView.Name = category.Name
			categoryView.Status = category.Status

			// Format waktu create dan update
			timeCategoryCreate, _ := helper.DatetimeToFormatIndo(category.CreatedAt)
			categoryView.CreatedAt = timeCategoryCreate
			timeCategoryUpdate, _ := helper.DatetimeToFormatIndo(category.UpdatedAt)
			categoryView.UpdatedAt = timeCategoryUpdate

			categoryViews = append(categoryViews, categoryView)
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return categoryNil, err
		}
	}

	return categoryViews, nil
}

func (s *service) GetCategoryByIDDeleted(ctx context.Context, id uint) (Category, error) {
	var category Category
	var categoryNil Category

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return categoryNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return categoryNil, ctx.Err()
	default:

		err := s.repository.FindByIDDeletedAtCategory(&category, id, tx)
		if err != nil {
			tx.Rollback()
			return categoryNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return categoryNil, err
		}
	}

	return category, nil
}

func (s *service) CreateCategory(ctx context.Context, input InputCategory) (Category, error) {

	var category Category
	var categoryNil Category

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return categoryNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return categoryNil, ctx.Err()
	default:

		err := s.validate.Struct(input)
		if err != nil {
			return categoryNil, errors.New("isi form dengan benar!")
		}

		category.Name = input.Name
		category.Status = input.Status

		err = s.repository.SaveCategory(&category, tx)
		if err != nil {
			tx.Rollback()
			return categoryNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return categoryNil, err
		}
	}

	return category, nil
}

func (s *service) UpdateCategory(ctx context.Context, input UpdateCategory) (Category, error) {
	var category Category
	var categoryRow Category
	var categoryNil Category

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return categoryNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return categoryNil, ctx.Err()
	default:
		err := s.validate.Struct(input)
		if err != nil {
			return categoryNil, errors.New("isi form dengan benar")
		}

		err = s.repository.FindByIDCategory(&categoryRow, input.ID, tx)
		if err != nil {
			tx.Rollback()
			return categoryNil, err
		}

		category.ID = categoryRow.ID
		category.Name = input.Name
		category.Status = input.Status

		err = s.repository.UpdateCategory(&category, tx)
		if err != nil {
			tx.Rollback()
			return categoryNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return categoryNil, err
		}
	}

	return category, nil
}

func (s *service) DeleteCategorySoft(ctx context.Context, id uint) (Category, error) {
	var category Category
	var categoryNil Category

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return categoryNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return categoryNil, ctx.Err()
	default:
		err := s.repository.DeleteSoftCategory(&category, id, tx)
		if err != nil {
			tx.Rollback()
			return categoryNil, err
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return categoryNil, err
		}
	}
	fmt.Println("delete service :", category)

	return category, nil
}

func (s *service) DeleteCategory(ctx context.Context, id uint) (Category, error) {
	var category Category
	var categoryrNil Category

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return categoryrNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return categoryrNil, ctx.Err()
	default:

		err := s.repository.Delete(&category, id, tx)
		if err != nil {
			tx.Rollback()
			return categoryrNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return categoryrNil, err
		}
	}

	fmt.Println("delete service :", category)

	return category, nil
}

func (s *service) RestoreCategory(ctx context.Context, id uint) (Category, error) {
	var category Category
	var categoryRow Category
	var categoryNil Category

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return categoryNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return categoryNil, ctx.Err()
	default:

		err := s.repository.FindByIDDeletedAtCategory(&category, id, tx)
		if err != nil {
			tx.Rollback()
			return categoryNil, err
		}

		err = s.repository.UpdateDeletedAtCategory(&categoryRow, category.ID, tx)
		if err != nil {
			tx.Rollback()
			return categoryNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return categoryNil, err
		}
	}

	return categoryRow, nil
}
