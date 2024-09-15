package utils

import (
	"net/url"
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
