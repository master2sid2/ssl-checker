package utils

import (
	"net/url"
	"regexp"
	"strings"
)

func ExtractDomain(domainURL string) string {
	domainURL = strings.TrimSpace(domainURL)

	if strings.HasPrefix(domainURL, "http://") || strings.HasPrefix(domainURL, "https://") {
		parsedURL, err := url.Parse(domainURL)
		if err == nil {
			host := parsedURL.Hostname()
			return strings.Split(host, ":")[0]
		}
	}
	return strings.TrimSpace(domainURL)
}

func ParseUserAgent(userAgent string) (string, string) {
	// Регулярное выражение для извлечения информации о браузере
	browserRe := regexp.MustCompile(`(Chrome|Firefox|Safari|Edge|Opera)/[\d.]+`)
	browserMatches := browserRe.FindString(userAgent)

	// Регулярное выражение для извлечения информации об операционной системе
	osRe := regexp.MustCompile(`\(([^)]+)\)`)
	osMatches := osRe.FindStringSubmatch(userAgent)

	var osInfo string
	if len(osMatches) > 1 {
		osInfo = osMatches[1] // Возьмем строку, которая внутри скобок
	} else {
		osInfo = "Unknown OS"
	}

	if browserMatches == "" {
		browserMatches = "Unknown Browser"
	}

	return browserMatches, osInfo
}
