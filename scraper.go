package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

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
	base := "https://djinni.co/jobs/"
	qs := url.Values{}
	qs.Add("primary_keyword", "JavaScript")
	qs.Add("primary_keyword", "Angular")
	qs.Add("primary_keyword", "React.js")
	qs.Add("primary_keyword", "Svelte")
	qs.Add("primary_keyword", "Vue.js")
	qs.Add("primary_keyword", "Markup")
	qs.Add("_", fmt.Sprint(time.Now().Unix()))

	req, err := http.NewRequest("GET", base+"?"+qs.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("djinni: %w", err)
	}
	// Djinni shows old jobs if you don't look like a browser.
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "uk-UA,uk;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://djinni.co/jobs/")

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
