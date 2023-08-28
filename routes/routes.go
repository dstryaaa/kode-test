package routes

import (
	"github.com/dstryaaa/kode-test/controllers"
	"github.com/dstryaaa/kode-test/middleware"
	"github.com/gorilla/mux"
)

var UserRoutes = func(r *mux.Router) {
	r.HandleFunc("/user/signup", controllers.UserSignUp).Methods("POST") // Обработчик запроса на регистрацию нового пользователя
	r.HandleFunc("/user/login", controllers.UserLogIn).Methods("POST")   // Обработчик запроса на вход пользователя в систему

	// Подмассив для защищенных роутов
	protectedNote := r.PathPrefix("/notes").Subrouter()
	protectedNote.Use(middleware.AuthenticateMiddleware)                   // Middleware для защиты защищенных роутов
	protectedNote.HandleFunc("/new", controllers.PostNote).Methods("POST") // Обработчик запроса на создание новой заметки
	protectedNote.HandleFunc("", controllers.ViewNotes).Methods("GET")     // Обработчик запроса на получение списка заметок
}
