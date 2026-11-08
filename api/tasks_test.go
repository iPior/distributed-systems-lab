package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestConcurrentCreateTasks(t *testing.T) {
	const count = 10
	store := newTaskStore()
	created := make([]Task, count)
	start := make(chan struct{})
	var wg sync.WaitGroup

	// Register all workers before starting them so Wait cannot finish early.
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(index int) {
			defer wg.Done()
			<-start
			// Each worker writes to its own slice element.
			created[index] = store.create("Task " + strconv.Itoa(index))
		}(i)
	}

	// Release the workers together, then wait before reading their results.
	close(start)
	wg.Wait()

	seen := make(map[int]bool, count)
	for i, task := range created {
		if task.ID < 1 || task.ID > count {
			t.Errorf("worker %d: ID %d outside range 1..%d", i, task.ID, count)
		}
		if seen[task.ID] {
			t.Errorf("duplicate task ID: %d", task.ID)
		}
		seen[task.ID] = true

		wantTitle := "Task " + strconv.Itoa(i)
		if task.Title != wantTitle {
			t.Errorf("worker %d: title got %q, want %q", i, task.Title, wantTitle)
		}
		got, found := store.get(task.ID)
		if !found || got != task {
			t.Errorf("lookup ID %d: found=%t, got %+v, want %+v", task.ID, found, got, task)
		}
	}
	if got := len(store.tasks); got != count {
		t.Errorf("stored task count: got %d, want %d", got, count)
	}
}

func TestCreateAndGetTask(t *testing.T) {
	store := newTaskStore()
	created := store.create("Learn Go")
	got, found := store.get(created.ID)
	if !found {
		t.Fatal("task not found")
	}
	if got != created {
		t.Fatalf("got %+v, want %+v", got, created)
	}
}

func TestGetMissingTask(t *testing.T) {
	store := newTaskStore()
	_, found := store.get(999)
	if found {
		t.Fatal("expected not found")
	}
}

func TestCreateTaskHTTP(t *testing.T) {
	store := newTaskStore()
	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(`{"title":"Learn Go"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	store.createTask(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusCreated)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type: got %q, want application/json", got)
	}

	var got Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	want := Task{ID: 1, Title: "Learn Go"}
	if got != want {
		t.Fatalf("body: got %+v, want %+v", got, want)
	}
	stored, found := store.get(got.ID)
	if !found || stored != want {
		t.Fatalf("stored: found=%t, got %+v, want %+v", found, stored, want)
	}
}

func TestCreateTaskInvalidInputHTTP(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "whitespace title", body: `{"title":"   "}`},
		{name: "empty title", body: `{"title":""}`},
		{name: "missing title", body: `{}`},
		{name: "empty body", body: ""},
		{name: "malformed JSON", body: `{"title":`},
		{name: "wrong title type", body: `{"title":42}`},
		{name: "extra JSON value", body: `{"title":"Learn Go"} {"title":"Extra"}`},
		{name: "trailing garbage", body: `{"title":"Learn Go"} garbage`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newTaskStore()
			req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			store.createTask(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status: got %d, want %d; body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
			}
			if len(store.tasks) != 0 {
				t.Fatalf("rejected input stored %d tasks, want 0", len(store.tasks))
			}
			if store.nextID != 1 {
				t.Fatalf("next ID: got %d, want 1 after rejected input", store.nextID)
			}
		})
	}
}

func TestGetTaskHTTP(t *testing.T) {
	store := newTaskStore()
	want := store.create("Learn Go")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", store.getTask)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+strconv.Itoa(want.ID), nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type: got %q, want application/json", got)
	}
	var got Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got != want {
		t.Fatalf("body: got %+v, want %+v", got, want)
	}
}

func TestGetTaskErrorsHTTP(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		wantStatus int
		wantBody   string
	}{
		{name: "missing task", id: "999", wantStatus: http.StatusNotFound, wantBody: "task not found\n"},
		{name: "non-numeric ID", id: "hello", wantStatus: http.StatusBadRequest, wantBody: "invalid id\n"},
		{name: "zero ID", id: "0", wantStatus: http.StatusBadRequest, wantBody: "invalid id\n"},
		{name: "negative ID", id: "-1", wantStatus: http.StatusBadRequest, wantBody: "invalid id\n"},
		{name: "overflowing ID", id: "999999999999999999999999999999", wantStatus: http.StatusBadRequest, wantBody: "invalid id\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newTaskStore()
			mux := http.NewServeMux()
			mux.HandleFunc("GET /tasks/{id}", store.getTask)
			req := httptest.NewRequest(http.MethodGet, "/tasks/"+tt.id, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d, want %d; body: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if got := rec.Body.String(); got != tt.wantBody {
				t.Fatalf("body: got %q, want %q", got, tt.wantBody)
			}
		})
	}
}
