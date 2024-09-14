package utils

import (
	"net/url"
	"strings"
)

func ExtractDomain(domainURL string) string {
	// Удаляем пробелы
	domainURL = strings.TrimSpace(domainURL)

	// Проверяем, если URL начинается с http:// или https://
	if strings.HasPrefix(domainURL, "http://") || strings.HasPrefix(domainURL, "https://") {
		parsedURL, err := url.Parse(domainURL)
		if err == nil {
			host := parsedURL.Hostname()
			// Удаляем порт, если он есть
			return strings.Split(host, ":")[0]
		}
	}

	// Если это просто домен, возвращаем его
	// Также обрабатываем случаи, когда домен может содержать пробелы
	return strings.TrimSpace(domainURL)
}
