package modules

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type SubdomainResult struct {
	Subdomain string
	IP        string
	Resolved  bool
}

func LoadSubdomainWordlist(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var subs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sub := strings.TrimSpace(scanner.Text())
		if sub != "" {
			subs = append(subs, sub)
		}
	}
	return subs, scanner.Err()
}

func ResolveSubdomains(target string, wordlist []string, threads int) []SubdomainResult {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	results := []SubdomainResult{}

	sem := make(chan struct{}, threads)

	for _, word := range wordlist {
		wg.Add(1)
		sem <- struct{}{}

		go func(sub string) {
			defer wg.Done()
			defer func() { <-sem }()
			
			fqdn := fmt.Sprintf("%s.%s", sub, target)
			ips, err := net.LookupHost(fqdn)
			if err == nil && len(ips) > 0 {
				mutex.Lock()
				results = append(results, SubdomainResult{
					Subdomain: fqdn,
					IP:        ips[0],
					Resolved:  true,
				})
				mutex.Unlock()
			}
		}(word)
	}

	wg.Wait()
	return results
}
