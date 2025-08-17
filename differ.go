package main

type JobDiffer struct{}

func NewJobDiffer() *JobDiffer {
	return &JobDiffer{}
}

func (d *JobDiffer) FindNew(jobs, known JobSet) []Job {
	var newJobs []Job
	for _, j := range jobs {
		if _, exists := known[j.URL]; !exists {
			newJobs = append(newJobs, j)
		}
	}
	return newJobs
}
