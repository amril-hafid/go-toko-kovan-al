package users

import (
	"context"
	"errors"
	"fmt"
	"go-toko-kovan-al/helper"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Service interface {
	GetAllUsers(ctx context.Context) ([]UserView, error)
	GetUserByID(ctx context.Context, id uint) (User, error)

	GetAllUsersDeleted(ctx context.Context) ([]UserView, error)
	GetUserByIDDeleted(ctx context.Context, id uint) (User, error)

	CreateImageProfile(ctx context.Context, file string, id uint) (User, error)
	RegisterUser(ctx context.Context, input RegisterUserInput) (User, error)
	Login(ctx context.Context, input LoginInput) (User, error)

	UpdateUser(ctx context.Context, input UpdateUserInput) (User, error)
	UpdatePassword(ctx context.Context, input UpdatePasswordInput) (User, error)
	DeleteUserSoft(ctx context.Context, id uint) (User, error)
	DeleteUser(ctx context.Context, id uint) (User, error)
	RestoreUser(ctx context.Context, id uint) (User, error)
}

type service struct {
	repository Repository
	validate   *validator.Validate
	db         *gorm.DB
}

func NewService(repository Repository, validate *validator.Validate, db *gorm.DB) *service {
	return &service{repository, validate, db}
}

func (s *service) GetAllUsers(ctx context.Context) ([]UserView, error) {
	var userResult []UserView
	var users []User
	var userNil []UserView

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return userNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return userNil, ctx.Err()
	default:
		err := s.repository.FindAll(&users, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		for i, user := range users {
			var userRow UserView
			userRow.ID = user.ID
			userRow.Index = i + 1
			userRow.Nama = user.Nama
			userRow.Email = user.Email
			userRow.NoHp = user.NoHp

			// Format tanggal lahir
			date, _ := helper.DateToFormatIndo(user.TanggalLahir)
			userRow.TanggalLahir = date
			userRow.Role = user.Role

			// Menentukan file profile picture
			if user.ProfileFile == "" {
				userRow.ProfileFile = "image-user-no-poto.png"
			} else {
				userRow.ProfileFile = user.ProfileFile
			}

			// Format waktu create dan update
			timeUserCreate, _ := helper.DatetimeToFormatIndo(user.CreatedAt)
			userRow.CreatedAt = timeUserCreate
			timeUserUpdate, _ := helper.DatetimeToFormatIndo(user.UpdatedAt)
			userRow.UpdatedAt = timeUserUpdate

			userResult = append(userResult, userRow)
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return userNil, err
		}
	}

	return userResult, nil
}
func (s *service) GetUserByID(ctx context.Context, id uint) (User, error) {
	var user User
	var userNil User

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return userNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return userNil, ctx.Err()
	default:

		err := s.repository.FindByID(&user, id, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		date, _ := helper.DateToFormatIndo(user.TanggalLahir)
		user.TanggalLahir = date

		if user.ProfileFile == "" {
			user.ProfileFile = "image-user-no-poto.png"
		}

		if user.ID == 0 {
			tx.Rollback()
			return userNil, errors.New("No user found no with that ID")
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return userNil, err
		}
	}
	return user, nil
}
func (s *service) GetAllUsersDeleted(ctx context.Context) ([]UserView, error) {
	var userResult []UserView
	var users []User
	var userNil []UserView

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return userNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return userNil, ctx.Err()
	default:
		err := s.repository.FindAllDeletedAt(&users, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		for i, user := range users {
			var userRow UserView
			userRow.ID = user.ID
			userRow.Index = i + 1
			userRow.Nama = user.Nama
			userRow.Email = user.Email
			userRow.NoHp = user.NoHp

			date, _ := helper.DateToFormatIndo(user.TanggalLahir)
			userRow.TanggalLahir = date

			userRow.Role = user.Role
			if user.ProfileFile == "" {
				userRow.ProfileFile = "image-user-no-poto.png"
			} else {
				userRow.ProfileFile = user.ProfileFile
			}

			timeUserCreate, _ := helper.DatetimeToFormatIndo(user.CreatedAt)
			userRow.CreatedAt = timeUserCreate
			timeUserDelete, _ := helper.StringToDateTimeIndoFormat(fmt.Sprint(user.DeletedAt))
			userRow.DeletedAt = timeUserDelete

			userResult = append(userResult, userRow)

		}

		if err := tx.Commit().Error; err != nil {
			return userNil, err
		}
	}

	return userResult, nil
}
func (s *service) GetUserByIDDeleted(ctx context.Context, id uint) (User, error) {
	var user User
	var userNil User

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return userNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return userNil, ctx.Err()
	default:

		err := s.repository.FindByIDDeletedAt(&user, id, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		date, _ := helper.DateToFormatIndo(user.TanggalLahir)
		user.TanggalLahir = date

		if user.ID == 0 {
			tx.Rollback()
			return userNil, errors.New("No user found no with that ID")
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return userNil, err
		}
	}

	return user, nil
}
func (s *service) CreateImageProfile(ctx context.Context, file string, id uint) (User, error) {
	var user User
	var userRow User
	var userNil User

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return userNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return userNil, ctx.Err()
	default:

		err := s.repository.FindByID(&userRow, id, tx)

		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		user.ID = userRow.ID
		user.Nama = userRow.Nama
		user.Email = userRow.Email
		user.Password = userRow.Password
		user.ProfileFile = file
		user.TanggalLahir = userRow.TanggalLahir
		user.NoHp = userRow.NoHp
		user.Role = userRow.Role
		err = s.repository.Update(&user, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return userNil, err
		}
	}

	return user, nil
}
func (s *service) RegisterUser(ctx context.Context, input RegisterUserInput) (User, error) {

	var user User
	var userEmailAfalabel User
	var userNil User

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return userNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return userNil, ctx.Err()
	default:

		err := s.validate.Struct(input)
		if err != nil {
			return userNil, errors.New("isi form dengan benar!")
		}

		if input.Password != input.PasswordRetype {
			return userNil, errors.New("Password yang anda masukan salah")
		}

		err = s.repository.FindByEmail(&userEmailAfalabel, input.Email, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		if userEmailAfalabel.ID != 0 {
			tx.Rollback()
			return userNil, errors.New("Email sudah pernah digunakan!")
		}

		user.Nama = input.Nama
		user.Email = input.Email
		user.NoHp = input.NoHp
		date, err := helper.StringToDate(input.TanggalLahir)
		if err != nil {
			tx.Rollback()
			return userNil, errors.New("Tangal Lahir tidak sesuai.")
		}
		user.TanggalLahir = date
		user.Password = helper.Sha1ToString(input.Password)
		user.Role = "user"

		err = s.repository.Save(&user, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return userNil, err
		}
	}

	return user, nil
}
func (s *service) Login(ctx context.Context, input LoginInput) (User, error) {
	var user User
	var userNil User

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return userNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return userNil, ctx.Err()
	default:
		err := s.validate.Struct(input)
		if err != nil {
			return userNil, errors.New("isi form dengan benar")
		}

		err = s.repository.FindByEmail(&user, input.Email, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		if user.ID == 0 {
			tx.Rollback()
			return userNil, errors.New("Pengguna dengan email tersebut tidak di temaukan")
		}

		ok, err := helper.VerifySHA1Hash(input.Password, user.Password)
		if err != nil {
			tx.Rollback()
			return userNil, errors.New("Password tidak sesuai")
		}

		if ok == false {
			tx.Rollback()
			return userNil, errors.New("Password tidak sesuai")
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return userNil, err
		}
	}

	return user, nil
}
func (s *service) UpdateUser(ctx context.Context, input UpdateUserInput) (User, error) {
	var user User
	var userRow User
	var userUpdate User
	var userNil User

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return userNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return userNil, ctx.Err()
	default:
		err := s.validate.Struct(input)
		if err != nil {
			return userNil, errors.New("isi form dengan benar")
		}

		err = s.repository.FindByID(&userRow, input.ID, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		user.ID = userRow.ID
		user.Nama = input.Nama
		user.Email = input.Email
		user.Password = userRow.Password
		date, _ := helper.StringToDate(input.TanggalLahir)
		user.TanggalLahir = date
		user.NoHp = strings.Replace(input.NoHp, "_", "", -1)
		user.Role = input.Role
		user.CreatedAt = userRow.CreatedAt
		err = s.repository.Update(&user, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return userNil, err
		}
	}

	return userUpdate, nil
}
func (s *service) UpdatePassword(ctx context.Context, input UpdatePasswordInput) (User, error) {
	var userRow User
	var userUpdate User
	var userNil User

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return userNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return userNil, ctx.Err()
	default:

		err := s.validate.Struct(input)
		if err != nil {
			return userNil, errors.New("isi form dengan benar")
		}

		err = s.repository.FindByID(&userRow, input.ID, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		if input.Password != input.PasswordRetype {
			tx.Rollback()
			return userNil, errors.New("Password salah")
		}

		userUpdate.ID = userRow.ID
		userUpdate.Nama = userRow.Nama
		userUpdate.Email = userRow.Email
		userUpdate.Password = helper.Sha1ToString(input.Password)
		userUpdate.TanggalLahir = userRow.TanggalLahir
		userUpdate.NoHp = userRow.NoHp
		userUpdate.Role = userRow.Role

		err = s.repository.Update(&userUpdate, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return userNil, err
		}
	}

	return userUpdate, nil
}
func (s *service) DeleteUserSoft(ctx context.Context, id uint) (User, error) {
	var user User
	var userRow User

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return userRow, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return userRow, ctx.Err()
	default:
		err := s.repository.DeleteSoft(&user, id, tx)
		if err != nil {
			tx.Rollback()
			return userRow, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return userRow, err
		}
	}

	return user, nil
}
func (s *service) DeleteUser(ctx context.Context, id uint) (User, error) {
	var user User
	var userNil User

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return userNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return userNil, ctx.Err()
	default:

		err := s.repository.Delete(&user, id, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return userNil, err
		}
	}
	return user, nil
}
func (s *service) RestoreUser(ctx context.Context, id uint) (User, error) {
	var user User
	var userRow User
	var userNil User

	tx := s.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return userNil, err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return userNil, ctx.Err()
	default:

		err := s.repository.FindByIDDeletedAt(&user, id, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		err = s.repository.UpdateDeletedAt(&userRow, user.ID, tx)
		if err != nil {
			tx.Rollback()
			return userNil, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return userNil, err
		}
	}

	return userRow, nil
}
