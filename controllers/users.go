package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/dstryaaa/kode-test/models"
	"github.com/dstryaaa/kode-test/utils"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func UserSignUp(w http.ResponseWriter, r *http.Request) {
	logger := utils.CreateNewLogger()
	var newUser models.User

	// Декодируем тело запроса в структуру newUser
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		// Если произошла ошибка при декодировании, логируем ошибку и возвращаем статус BadRequest
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Body decode error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Хешируем пароль нового пользователя
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Password hash generation error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Создаем ID для нового пользователя
	newUser.ID = uuid.New().String()
	newUser.Password = string(hashedPassword)

	// подключаемся к БД
	db, err := sql.Open("postgres", "postgresql://postgres:qwerty@postgres:5432/postgres?sslmode=disable")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Opening database error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Выполняем запрос на вставку нового пользователя в таблицу users
	insertQuery := "INSERT INTO users (userID, username, email, password_hash) VALUES ($1, $2, $3, $4)"
	_, err = db.Exec(insertQuery, newUser.ID, newUser.Username, newUser.Email, newUser.Password)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Database insertion error")
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	// Логируем успешное выполнение запроса на регистрацию нового пользователя
	logger.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info("Login request received")
}

func UserLogIn(w http.ResponseWriter, r *http.Request) {
	logger := utils.CreateNewLogger()

	// Декодируем тело запроса в структуру user
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Если произошла ошибка при декодировании тела запроса, логируем ошибку
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Body decode error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// подключаемся к БД
	db, err := sql.Open("postgres", "postgresql://postgres:qwerty@postgres:5432/postgres?sslmode=disable")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Batabase opening error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var foundUser models.User
	// Ищем пользователя в базе данных по email
	err = db.QueryRow("SELECT userid, username, email, password_hash FROM users WHERE email = $1", user.Email).Scan(&foundUser.ID, &foundUser.Username, &foundUser.Email, &foundUser.Password)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Database scan error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Сравниваем хеш пароля, полученный из запроса, с хешем пароля, хранящимся в базе данных
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Passwords doesn't match")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Генерируем JWT-токен для пользователя
	accessToken := utils.GenerateAccessToken(foundUser.ID)
	if accessToken == "" {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Create access token error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Формируем ответ с JWT-токеном
	response := map[string]string{
		"access_token": accessToken,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {

		return
	}
	// Логируем успешное выполнение запроса на вход пользователя
	logger.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info("Login request received")

	// В ответ возвращаем access токен
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
