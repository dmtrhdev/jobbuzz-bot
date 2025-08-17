package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type DOUScraper struct {
	client *http.Client
}

func NewDOUScraper() *DOUScraper {
	return &DOUScraper{
		client: &http.Client{},
	}
}

func (s *DOUScraper) Jobs() ([]Job, error) {
	resp, err := s.client.Get("https://jobs.dou.ua/vacancies/?category=Front%20End")
	if err != nil {
		return nil, fmt.Errorf("dou: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dou returned %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("dou: %w", err)
	}

	var jobs []Job
	walk(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" && hasClass(n, "l-vacancy") {
			var title, url, company string

			walk(n, func(a *html.Node) {
				if a.Type == html.ElementNode && a.Data == "a" {
					if hasClass(a, "vt") {
						title = text(a)
						url = attr(a, "href")
					}
					if hasClass(a, "company") {
						company = text(a)
					}
				}
			})

			if title != "" && url != "" && company != "" {
				jobs = append(jobs, Job{URL: url, Title: title, Company: company})
			}
		}
	})

	return jobs, nil
}

type DjinniScraper struct {
	client *http.Client
}

func NewDjinniScraper() *DjinniScraper {
	return &DjinniScraper{
		client: &http.Client{},
	}
}

func (s *DjinniScraper) Jobs() ([]Job, error) {
	req, err := http.NewRequest("GET", "https://djinni.co/jobs/?primary_keyword=JavaScript&primary_keyword=Angular&primary_keyword=React.js&primary_keyword=Svelte&primary_keyword=Vue.js&primary_keyword=Markup", nil)
	if err != nil {
		return nil, fmt.Errorf("djinni: %w", err)
	}

	// Seems like these headers are not set, Djinni might return a different page or block the request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Referer", "https://djinni.co/")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("djinni: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("djinni returned %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("djinni: %w", err)
	}

	var jobs []Job
	walk(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" && strings.HasPrefix(attr(n, "id"), "job-item") {
			var title, company, url string

			walk(n, func(a *html.Node) {
				if a.Type == html.ElementNode && a.Data == "a" {
					if hasClass(a, "job-item__title-link") {
						title = text(a)
						url = fmt.Sprintf("https://djinni.co%s", attr(a, "href"))
					}
					if strings.Contains(attr(a, "href"), "company-") {
						company = text(a)
					}
				}
			})

			if title != "" && url != "" && company != "" {
				jobs = append(jobs, Job{URL: url, Title: title, Company: company})
			}
		}
	})

	return jobs, nil
}

func hasClass(n *html.Node, class string) bool {
	for _, a := range n.Attr {
		if a.Key == "class" && strings.Contains(a.Val, class) {
			return true
		}
	}
	return false
}

func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	var b strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b.WriteString(text(c))
	}
	return strings.TrimSpace(b.String())
}

func walk(n *html.Node, fn func(*html.Node)) {
	fn(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walk(c, fn)
	}
}
