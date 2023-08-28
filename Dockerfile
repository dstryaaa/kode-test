# Определение базового образа
FROM golang:1.21.0-alpine

# Рабочая директория внутри контейнера
WORKDIR /go/src/kode

# Копирование файлов зависимостей Go в контейнер
COPY go.mod go.sum ./

# Установка зависимостей Go из файлов go.mod и go.sum
RUN go mod download

# Копирование исходных файлов приложения в контейнер
COPY . .

# RUN migrate -path database/migration -database "postgresql://postgres:qwerty@localhost:5432/postgres?sslmode=disable"

# Сборка приложения внутри контейнера
RUN go build -o main ./cmd/main/main.go

# Определение команды запуска приложения
CMD ["./main"]