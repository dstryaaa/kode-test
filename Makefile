.PHONY: build clean

build:
	go build -o kode-auth .\cmd\main\main.go

migrate:
	migrate -path database/migration/ -database "postgresql://postgres:qwerty@localhost:5432/postgres?sslmode=disable" -verbose up

clean:
	rm -f kode-auth