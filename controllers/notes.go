package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dstryaaa/kode-test/models"
	"github.com/dstryaaa/kode-test/utils"
	"github.com/sirupsen/logrus"
)

func PostNote(w http.ResponseWriter, r *http.Request) {
	logger := utils.CreateNewLogger()
	// Заводим новый экземпляр модели Note
	var user models.Note
	// Декодируем в него отправленную записку
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
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
		}).Error("Opening database error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	// Проверяем доступ к пользователя к странице
	tokenCheck, err := utils.ValidateAccessToken(r.Header.Get("Authorization")[7:], []byte("secret"))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Token error")
		return
	}
	//  Из access токена достаем ID пользователя
	claims := tokenCheck.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	// С помощью ID пользователя так же достаем username
	var foundUser models.User
	err = db.QueryRow("SELECT username FROM users WHERE userID = $1", userID).Scan(&foundUser.Username)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Database scan error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	foundUser.ID = userID
	// Проводим проверку и исправление заметки с помощью Яндекс.Спеллер
	checkedNote := utils.SpellCheck(user.Note)
	// Записываем все полученные данные в таблицу notes
	insertQuery := "INSERT INTO notes (userID, username, note) VALUES ($1, $2, $3)"
	_, err = db.Exec(insertQuery, foundUser.ID, foundUser.Username, checkedNote)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Database insertion error")
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	logger.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info("Note created")

	json.NewEncoder(w).Encode(true)
}

func ViewNotes(w http.ResponseWriter, r *http.Request) {
	logger := utils.CreateNewLogger()
	// подключаемся к БД
	db, err := sql.Open("postgres", "postgresql://postgres:qwerty@postgres:5432/postgres?sslmode=disable")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Opening database error")
		return
	}
	defer db.Close()
	// Проверяем доступ к пользователя к странице
	tokenCheck, err := utils.ValidateAccessToken(r.Header.Get("Authorization")[7:], []byte("secret"))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Token error")
		return
	}

	//  Из access токена достаем ID пользователя
	claims := tokenCheck.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	// Получаем все заметки написанные пользователем
	notesStorage, err := GetNotesByUserId(userID, db)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Getting notes error error")
		return
	}
	// Переводим записки в формат JSON
	response, err := json.Marshal(notesStorage)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Error("Marshaling response error")
		return
	}

	logger.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info("Notes shown")

	// Выводим заметки в ответе
	w.Write(response)
}

func GetNotesByUserId(userID string, db *sql.DB) ([]string, error) {
	logger := utils.CreateNewLogger()
	//  Из таблицы notes достаем note сделаный пользователем с userID
	rows, err := db.Query("SELECT note FROM notes WHERE userID = $1", userID)
	if err != nil {
		logger.Error("Database selection error")
		return nil, err
	}
	defer rows.Close()
	// Создаем хранилище для заметок
	var notesStorage []string
	// Записываем все заметки по одной в созданное для них хранилище
	for rows.Next() {
		var singleNote string
		err := rows.Scan(&singleNote)
		if err != nil {
			logger.Error("Scan database rows error")
			return nil, err
		}
		notesStorage = append(notesStorage, singleNote)
	}

	err = rows.Err()
	if err != nil {
		logger.Error("Database rows error")
		return nil, err
	}

	return notesStorage, nil
}
