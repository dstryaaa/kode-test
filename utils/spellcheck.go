package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func SpellCheck(text string) string {
	logger := CreateNewLogger()
	// Создаем запрос на Яндекс.Спеллер
	checkURL := fmt.Sprintf("https://speller.yandex.net/services/spellservice.json/checkText?text=%s", url.QueryEscape(text))
	checkedText := text
	resp, err := http.Get(checkURL)
	if err != nil {
		logger.Error("Response from Yandex.Speller error")
		return ""
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Bode reading error")
		return ""
	}

	var result []struct {
		Word string   `json:"word"`
		S    []string `json:"s"`
	}
	// Записываем полученный ответ в структуру result
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Error("Body unmarshal error")
		return ""
	}
	// Проходим по структуре result заменяя слово с ошибкой
	for _, res := range result {
		for _, r := range res.S {
			checkedText = strings.Replace(checkedText, res.Word, r, 1)
		}
	}

	return checkedText
}
