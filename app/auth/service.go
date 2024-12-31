package auth

import (
	"go-toko-kovan-al/config"
	"go-toko-kovan-al/helper"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	GenerateToken(userID uint, config config.Config) (string, error)
	ValidateToken(token string, config config.Config) (*jwt.Token, error)
}

type jwtService struct {
}

func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userID uint, config config.Config) (string, error) {
	day := time.Hour * 24
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(day * 1).Unix(),
	}

	token, err := helper.TokenGenerateHelper(claims, config)
	if err != nil {
		return token, err
	}
	return token, nil
}

func (s *jwtService) ValidateToken(encodeToken string, config config.Config) (*jwt.Token, error) {
	token, err := helper.TokenVerifHelper(encodeToken, config)

	if err != nil {
		return token, err
	}

	return token, err
}
