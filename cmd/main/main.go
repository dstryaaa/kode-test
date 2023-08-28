package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql"

	"github.com/dstryaaa/kode-test/routes"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:qwerty@postgres:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		fmt.Println("ggs")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///go/src/kode/database/migration",
		"postgres", instance)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	routes.UserRoutes(r)
	http.Handle("/", r)
	fmt.Println("Server started at port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
