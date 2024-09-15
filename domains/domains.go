package domains_utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Domain struct {
	Name       string
	ValidUntil time.Time
	DaysLeft   int
	Message    string
}

var Domains = []string{}

func LoadDomains() ([]string, error) {
	data, err := os.ReadFile("data/domains.list")
	if err != nil {
		log.Println("Error reading domains file:", err)
		return nil, err
	}

	domains := FilterEmptyLines(strings.Split(strings.TrimSpace(string(data)), "\n"))
	return domains, nil
}

func SaveDomains(domains []string) error {
	data := strings.Join(domains, "\n")
	err := os.WriteFile("data/domains.list", []byte(data), 0644)
	if err != nil {
		log.Println("Error writing domains file:", err)
		return err
	}
	return nil
}

func FilterEmptyLines(lines []string) []string {
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return result
}

func CheckCertificate(domain string) (time.Time, string) {
	conn, err := tls.Dial("tcp", domain+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return time.Time{}, fmt.Sprintf("Error: %v", err)
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	validUntil := cert.NotAfter
	return validUntil, ""
}

func UpdateDomainTable(domains []string) ([]Domain, error) {
	var result []Domain

	for _, domainName := range domains {
		validUntil, message := CheckCertificate(domainName)
		daysLeft := int(validUntil.Sub(time.Now()).Hours() / 24)

		domainMessage := ""
		if message != "" {
			domainMessage = message
		}

		domain := Domain{
			Name:       domainName,
			ValidUntil: validUntil,
			DaysLeft:   daysLeft,
			Message:    domainMessage,
		}

		result = append(result, domain)
	}

	return result, nil
}
