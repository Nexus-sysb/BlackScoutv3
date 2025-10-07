package utils

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

var TotalRequests int64

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
	"Mozilla/5.0 (X11; Linux x86_64)",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X)",
	"Mozilla/5.0 (iPad; CPU OS 13_6 like Mac OS X)",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B)",
	"curl/7.68.0",
	"Wget/1.21.1 (linux-gnu)",
	"Googlebot/2.1 (+http://www.google.com/bot.html)",
	"Bingbot/2.0 (+http://www.bing.com/bingbot.htm)",
	"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (Windows NT 10.0; rv:102.0) Gecko/20100101 Firefox/102.0",
	"Mozilla/5.0 (Linux; Android 12; SM-G998B)",
	"Mozilla/5.0 (compatible; PetalBot/1.0; +https://petalsearch.com/bot)",
}

func RandomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	return userAgents[rand.Intn(len(userAgents))]
}

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Banner() {
	fmt.Println(`
 ▄▄▄▄    ██▓    ▄▄▄       ▄████▄   ██ ▄█▀  ██████  ▄████▄   ▒█████   █    ██ ▄▄▄█████▓
▓█████▄ ▓██▒   ▒████▄    ▒██▀ ▀█   ██▄█▒ ▒██    ▒ ▒██▀ ▀█  ▒██▒  ██▒ ██  ▓██▒▓  ██▒ ▓▒
▒██▒ ▄██▒██░   ▒██  ▀█▄  ▒▓█    ▄ ▓███▄░ ░ ▓██▄   ▒▓█    ▄ ▒██░  ██▒▓██  ▒██░▒ ▓██░ ▒░
▒██░█▀  ▒██░   ░██▄▄▄▄██ ▒▓▓▄ ▄██▒▓██ █▄   ▒   ██▒▒▓▓▄ ▄██▒▒██   ██░▓▓█  ░██░░ ▓██▓ ░
░▓█  ▀█▓░██████▒▓█   ▓██▒▒ ▓███▀ ░▒██▒ █▄▒██████▒▒▒ ▓███▀ ░░ ████▓▒░▒▒█████▓   ▒██▒ ░
░▒▓███▀▒░ ▒░▓  ░▒▒   ▓▒█░░ ░▒ ▒  ░▒ ▒▒ ▓▒▒ ▒▓▒ ▒ ░░ ░▒ ▒  ░░ ▒░▒░▒░ ░▒▓▒ ▒ ▒   ▒ ░░
▒░▒   ░ ░ ░ ▒  ░ ▒   ▒▒ ░  ░  ▒   ░ ░▒ ▒░░ ░▒  ░ ░  ░  ▒     ░ ▒ ▒░ ░░▒░ ░ ░     ░
 ░    ░   ░ ░    ░   ▒   ░        ░ ░░ ░ ░  ░  ░  ░        ░ ░ ░ ▒   ░░░ ░ ░   ░
 ░          ░  ░     ░  ░░ ░      ░  ░         ░  ░ ░          ░ ░     ░
      ░                  ░                        ░
       BlackScoutV3 - Coded By Nexus | Team The Project Nexus
`)
}

func ReadInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func RandomHeaders() map[string]string {
	return map[string]string{
		"Accept":          "*/*",
		"Accept-Language": "en-US,en;q=0.9",
		"Connection":      "keep-alive",
	}
}

func Delay(ms int) {
	if ms <= 0 {
		return
	}
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func IncrementRequests() {
	atomic.AddInt64(&TotalRequests, 1)
}

func ShowLiveProgress(ptr *int64, start time.Time) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		reqs := atomic.LoadInt64(ptr)
		elapsed := time.Since(start).Seconds()
		if elapsed <= 0 {
			elapsed = 1
		}
		rps := float64(reqs) / elapsed
		fmt.Printf("\r[+] Requisições: %d | Tempo: %.0fs | RPS: %.2f", reqs, elapsed, rps)
	}
}

var linkRegex = regexp.MustCompile(`(?i)(?:href|src)=["']([^"']+)["']|https?://[^\s"'<>]+|/[-A-Za-z0-9_./?=&%]+`)

func ExtractLinks(content string, base string) []string {
	matches := linkRegex.FindAllStringSubmatch(content, -1)
	var out []string
	for _, m := range matches {
		for i := 1; i < len(m); i++ {
			if m[i] == "" {
				continue
			}
			u := strings.TrimSpace(m[i])
			parsed, err := url.Parse(u)
			if err != nil {
				continue
			}
			if !parsed.IsAbs() {
				baseURL, err := url.Parse(base)
				if err != nil {
					continue
				}
				u = baseURL.ResolveReference(parsed).String()
			}
			out = append(out, u)
		}
	}
	return uniqueStrings(out)
}

var jsURLRegex = regexp.MustCompile(`https?://[^\s"'<>]+|/(?:[^\s"'<>]+\.js)`)

func ExtractJSLinks(jsContent string, origin string) []string {
	matches := jsURLRegex.FindAllString(jsContent, -1)
	var out []string
	for _, m := range matches {
		u := strings.TrimSpace(m)
		if strings.HasPrefix(u, "/") {
			baseURL, err := url.Parse(origin)
			if err == nil {
				u = baseURL.ResolveReference(&url.URL{Path: u}).String()
			}
		}
		out = append(out, u)
	}
	return uniqueStrings(out)
}

func uniqueStrings(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	var out []string
	for _, s := range in {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

func IsLikelyAPI(u string) bool {
	l := strings.ToLower(u)
	if strings.Contains(l, "/api/") || strings.HasSuffix(l, ".json") || strings.Contains(l, "graphql") {
		return true
	}
	return false
}
