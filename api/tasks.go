package main

import (
	"sync"
	"net/http"
	"encoding/json"
	"strings"
	"log"
	"io"
)

type Task struct {
	ID int `json:"id"`
	Title string `json:"title"`
}

type taskStore struct {
	mu sync.Mutex
	tasks map[int]Task
	nextID int
}

func newTaskStore() *taskStore {
    return &taskStore{
        tasks:  make(map[int]Task),
        nextID: 1,
    }
}

func (s *taskStore) create(title string) Task {
    s.mu.Lock()
	defer s.mu.Unlock()
	task := Task{
		ID:    s.nextID,
		Title: title,
	}
	s.tasks[s.nextID] = task
	s.nextID++
	return task
}

func (s *taskStore) createTask(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Title string `json:"title"`
    }

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err = decoder.Decode(new(any)); err != io.EOF {
		http.Error(w, "body must contain exactly one JSON value", http.StatusBadRequest)
		return
	}

	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		http.Error(w, "title cannot be empty", http.StatusBadRequest)
		return
	}

	task := s.create(input.Title)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Printf("write task response: %v", err)
		return
	}
}