package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetSubdomains(domain string) []string {
	var subdomains []string

	apiURL := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return subdomains
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return subdomains
	}

	var certs []map[string]interface{}
	err = json.Unmarshal(body, &certs)
	if err != nil {
		return subdomains
	}

	seen := make(map[string]bool)
	for _, cert := range certs {
		if name, ok := cert["name_value"].(string); ok {
			for _, sub := range strings.Split(name, "\n") {
				clean := strings.TrimSpace(sub)
				if !seen[clean] && strings.Contains(clean, domain) {
					seen[clean] = true
					subdomains = append(subdomains, clean)
				}
			}
		}
	}

	return subdomains
}
