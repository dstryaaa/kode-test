package utils

import (
	// "encoding/base64"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	// "github.com/google/uuid"
)

var secretString []byte = []byte("secret")

// Создаем новый access токен указывая в нем ID пользователя и используя secretString
func GenerateAccessToken(userID string) string {
	logger := CreateNewLogger()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	tokenString, err := token.SignedString(secretString)
	if err != nil {
		logger.Error("Creating JWT error")
		return ""
	}
	return tokenString
}
func ValidateAccessToken(tokenString string, sec []byte) (*jwt.Token, error) {
	// Распарсиваем токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Возвращаем секретную строку для валидации
		return sec, nil
	})

	// Возвращаем токен или ошибку
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
