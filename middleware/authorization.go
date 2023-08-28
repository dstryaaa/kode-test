package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			token, err := jwt.Parse(authHeader[7:], func(token *jwt.Token) (interface{}, error) {
				// Проверяем, что подпись JWT токена не была изменена
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				// Возвращаем секретный ключ для проверки подписи
				return []byte("secret"), nil
			})
			if err == nil && token.Valid {
				// Пользователь авторизован, передаем управление следующему обработчику
				next.ServeHTTP(w, r)
				return
			}
		}
		// Пользователь не авторизован, возвращаем ошибку 401
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
