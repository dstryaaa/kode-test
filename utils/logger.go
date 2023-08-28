package utils

import (
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// Функция создает и конфигурирует логгер
func CreateNewLogger() *log.Logger {
	// Создаем новый логгер с помощью библиотеки logrus
	logger := logrus.New()

	// Устанавливаем конфигурацию логгера для использования текстового формата с цветным выводом
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	// Устанавливаем вывод логгера в стандартный вывод
	logger.SetOutput(os.Stdout)
	return logger
}
