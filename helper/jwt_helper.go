package helper

import (
	"errors"
	"go-toko-kovan-al/config"

	"github.com/golang-jwt/jwt/v5"
)

func ClaimTokenHelper(token string, config config.Config) (uint, error) {

	// if !strings.Contains(token, "Bearer") {
	// 	return 0, errors.New("token failed")
	// }

	// tokenString := ""
	// arrayToken := strings.Split(token, " ")
	// if len(arrayToken) == 2 {
	// 	tokenString = arrayToken[1]
	// }

	tokenVerif, err := TokenVerifHelper(token, config)
	if err != nil {
		return 0, err
	}

	claim, ok := tokenVerif.Claims.(jwt.MapClaims)
	if !ok || !tokenVerif.Valid {
		return 0, errors.New("token failed")
	}

	idUser := uint(claim["user_id"].(float64))

	return idUser, nil
}

func TokenVerifHelper(encodeToken string, config config.Config) (*jwt.Token, error) {
	token, err := jwt.Parse(encodeToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Invalid token")
		}
		return []byte(config.Enk.Key), nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenGenerateHelper(value jwt.MapClaims, config config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, value)

	signedToken, err := token.SignedString([]byte(config.Enk.Key))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
