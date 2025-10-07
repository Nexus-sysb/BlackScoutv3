package core

import (
	"blackscout/utils"
	"context"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Crawler struct {
	BaseURL     *url.URL
	Visited     sync.Map
	Results     []string
	ThreadLimit chan struct{}
	DelayMs     int
}

func NewCrawler(target string, threads, delay int) (*Crawler, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	return &Crawler{
		BaseURL:     u,
		ThreadLimit: make(chan struct{}, threads),
		DelayMs:     delay,
	}, nil
}

func (c *Crawler) normalize(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}
	if u.IsAbs() {
		if u.Host == c.BaseURL.Host {
			return u.String()
		}
		return ""
	}
	return c.BaseURL.ResolveReference(u).String()
}

func (c *Crawler) fetch(link string, apiScan, jsScan bool) []string {
	var found []string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := &http.Client{}
	req, _ := http.NewRequestWithContext(ctx, "GET", link, nil)
	req.Header.Set("User-Agent", utils.RandomUserAgent())
	for k, v := range utils.RandomHeaders() {
		req.Header.Set(k, v)
	}

	time.Sleep(time.Duration(rand.Intn(c.DelayMs)) * time.Millisecond)

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		return found
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	content := string(body)

	if jsScan && strings.HasSuffix(link, ".js") {
		utils.ExtractJSLinks(content, link)
	}

	links := utils.ExtractLinks(content, link)
	for _, l := range links {
		norm := c.normalize(l)
		if norm == "" {
			continue
		}
		if apiScan && !utils.IsLikelyAPI(norm) {
			continue
		}
		if _, loaded := c.Visited.LoadOrStore(norm, true); !loaded {
			found = append(found, norm)
		}
	}

	return found
}

func (c *Crawler) Start(apiScan, jsScan bool) []string {
	var wg sync.WaitGroup
	var mux sync.Mutex

	queue := make(chan string, 1000)
	queue <- c.BaseURL.String()
	c.Visited.Store(c.BaseURL.String(), true)

	workers := cap(c.ThreadLimit)
	if workers <= 0 {
		workers = 10
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for link := range queue {
				c.ThreadLimit <- struct{}{}
				res := c.fetch(link, apiScan, jsScan)
				mux.Lock()
				c.Results = append(c.Results, link)
				mux.Unlock()
				for _, r := range res {
					queue <- r
				}
				utils.IncrementRequests()
				<-c.ThreadLimit
			}
		}()
	}

	time.AfterFunc(30*time.Second, func() {
		close(queue)
	})
	wg.Wait()
	return c.Results
}
