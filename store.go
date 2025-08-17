package main

import (
	"encoding/json"
	"os"
)

type FileStore struct {
	filename string
}

func NewFileStore(filename string) *FileStore {
	return &FileStore{filename: filename}
}

func (s *FileStore) Load() (JobSet, error) {
	data, err := os.ReadFile(s.filename)
	if os.IsNotExist(err) {
		return make(JobSet), nil
	}
	if err != nil {
		return nil, err
	}
	var jobs JobSet
	err = json.Unmarshal(data, &jobs)
	return jobs, err
}

func (s *FileStore) Save(jobs JobSet) error {
	data, err := json.Marshal(jobs)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, data, 0644)
}
