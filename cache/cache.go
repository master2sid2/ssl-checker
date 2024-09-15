package cache

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"ssl-checker/domains"
)

var CacheFile = "data/cache.json"
var CacheMutex sync.Mutex

func LoadCache() ([]domains_utils.Domain, error) {
	CacheMutex.Lock()
	defer CacheMutex.Unlock()

	file, err := os.Open(CacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []domains_utils.Domain{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var domains []domains_utils.Domain
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&domains)
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func SaveCache(domains []domains_utils.Domain) error {
	CacheMutex.Lock()
	defer CacheMutex.Unlock()

	file, err := os.Create(CacheFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(domains)
}

func UpdateCache() {
	for {
		domains, err := domains_utils.LoadDomains()
		if err != nil {
			log.Println("Error loading domains:", err)
			continue
		}

		domainData, err := domains_utils.UpdateDomainTable(domains)
		if err != nil {
			log.Println("Error updating domain table:", err)
			continue
		}

		err = SaveCache(domainData)
		if err != nil {
			log.Println("Error saving cache:", err)
		}

		time.Sleep(1 * time.Hour)
	}
}
