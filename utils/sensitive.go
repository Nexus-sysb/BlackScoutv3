package utils

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func FindSensitiveData(url string) []string {
	var results []string
	client := http.Client{Timeout: 7 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return results
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return results
	}

	html := string(body)

	patterns := map[string]string{
		"Email":     `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`,
		"IP":        `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`,
		"CPF":       `\d{3}\.\d{3}\.\d{3}-\d{2}`,
		"CNPJ":      `\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}`,
		"Token":     `(?i)token[\s"']*[:=][\s"']*[A-Za-z0-9_\-\.]{10,}`,
		"API Key":   `(?i)(api_key|apikey|key)[\s"']*[:=][\s"']*[A-Za-z0-9_\-\.]{10,}`,
		"Senha":     `(?i)(senha|password)[\s"']*[:=][\s"']*[^"'&<>\s]+`,
	}

	for name, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(html, -1)
		if len(matches) > 0 {
			results = append(results, "["+name+"]")
			for _, m := range matches {
				results = append(results, strings.TrimSpace(m))
			}
		}
	}

	return results
}
