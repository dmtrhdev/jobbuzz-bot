package main

import (
	"log"
	"os"
	"sync"
)

func main() {
	token, chatID := os.Getenv("TELEGRAM_TOKEN"),
		os.Getenv("CHAT_ID")
	if token == "" || chatID == "" {
		log.Fatal("need TELEGRAM_TOKEN and CHAT_ID")
	}

	store := NewFileStore("known_jobs.json")
	differ := NewJobDiffer()
	bot := NewTelegramBot(token, chatID)

	known, err := store.Load()
	if err != nil {
		log.Fatal(err)
	}

	scrappers := []Scraper{
		NewDOUScraper(),
		NewDjinniScraper(),
	}

	jobs := scrape(scrappers)
	newJobs := differ.FindNew(jobs, known)

	if len(newJobs) > 0 {
		if err := bot.Send(newJobs); err != nil {
			log.Println(err)
		}
	}

	for url, job := range jobs {
		known[url] = job
	}

	if err := store.Save(known); err != nil {
		log.Fatal(err)
	}
}

func scrape(scrappers []Scraper) JobSet {
	var wg sync.WaitGroup
	var mu sync.Mutex
	result := make(JobSet)

	for _, s := range scrappers {
		wg.Add(1)
		go func(s Scraper) {
			defer wg.Done()
			jobs, err := s.Jobs()
			if err != nil {
				log.Println(err)
				return
			}
			mu.Lock()
			for _, j := range jobs {
				result[j.URL] = j
			}
			mu.Unlock()
		}(s)
	}

	wg.Wait()
	return result
}
