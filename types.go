package main

type Job struct {
	URL     string
	Title   string
	Company string
}

type JobSet map[string]Job

type Scraper interface {
	Jobs() ([]Job, error)
}

type Store interface {
	Load() ([]string, error)
	Save(jobs []Job) error
}

type Differ interface {
	FindNew(jobs []Job, known []Job) []Job
}

type Notifier interface {
	Send(jobs []Job) error
}
