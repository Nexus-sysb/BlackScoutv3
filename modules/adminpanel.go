package modules

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var adminPaths = []string{
	"/admin", "/administrator", "/admin/login", "/adminpanel", "/cpanel", "/admin_area",
	"/panel", "/controlpanel", "/admin1", "/admin2", "/admin/login.php", "/admin/index.php",
}

type AdminResult struct {
	URL        string
	StatusCode int
	Found      bool
}

func DetectAdminPanels(target string) []AdminResult {
	results := []AdminResult{}

	client := &http.Client{Timeout: 5 * time.Second}

	for _, path := range adminPaths {
		fullURL := fmt.Sprintf("%s%s", target, path)
		resp, err := client.Get(fullURL)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		isAdmin := strings.Contains(resp.Request.URL.Path, "admin") || resp.StatusCode == 200
		results = append(results, AdminResult{
			URL:        fullURL,
			StatusCode: resp.StatusCode,
			Found:      isAdmin,
		})
	}

	return results
}
