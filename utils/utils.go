package utils

import (
	"log"
	"net/url"
	"regexp"
	"strings"

	"ssl-checker/cache"
)

type CertStats struct {
	TotalCertificates        int
	ValidCertificates        int
	ExpiringSoonCertificates int
	CriticalCertificates     int
	ErrorCertificates        int
}

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
	browserRe := regexp.MustCompile(`(Chrome|Firefox|Safari|Edge|Opera)/[\d.]+`)
	browserMatches := browserRe.FindString(userAgent)

	osRe := regexp.MustCompile(`\(([^)]+)\)`)
	osMatches := osRe.FindStringSubmatch(userAgent)

	var osInfo string
	if len(osMatches) > 1 {
		osInfo = osMatches[1]
	} else {
		osInfo = "Unknown OS"
	}

	if browserMatches == "" {
		browserMatches = "Unknown Browser"
	}

	return browserMatches, osInfo
}

func CalculateCertificateStats() (CertStats, error) {
	var stats CertStats

	domains, err := cache.LoadCache()
	if err != nil {
		log.Printf("Error loading cache: %v", err)
		return stats, err
	}

	stats.TotalCertificates = len(domains)

	for _, domain := range domains {
		if domain.Message != "" {
			stats.ErrorCertificates++
		} else if domain.DaysLeft > 7 {
			stats.ValidCertificates++
		} else if domain.DaysLeft <= 7 && domain.DaysLeft > 3 {
			stats.ExpiringSoonCertificates++
		} else if domain.DaysLeft <= 3 {
			stats.CriticalCertificates++
		}
	}

	return stats, nil
}
