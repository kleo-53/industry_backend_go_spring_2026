package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

var TaskNotFound = errors.New("no task with such id in storage")

type Clock interface {
	Now() time.Time
}

type Task struct {
	ID        string
	Title     string
	Done      bool
	UpdatedAt time.Time
}

type TaskRepo interface {
	Create(title string) (Task, error)
	Get(id string) (Task, bool)
	List() []Task
	SetDone(id string, done bool) (Task, error)
}

type Storage struct {
	mu      sync.RWMutex
	tasks   map[string]Task
	clock   Clock
	counter uint64
}

func NewInMemoryTaskRepo(clock Clock) *Storage {
	return &Storage{
		clock: clock,
		tasks: map[string]Task{},
	}
}

func (s *Storage) Create(title string) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprint(s.counter)
	s.counter++
	newTask := Task{
		ID:        id,
		Title:     title,
		Done:      false,
		UpdatedAt: s.clock.Now(),
	}
	s.tasks[id] = newTask
	return newTask, nil
}

func (s *Storage) Get(id string) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if el, ok := s.tasks[id]; ok {
		return el, ok
	} else {
		return Task{}, ok
	}
}

func (s *Storage) List() (res []Task) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, task := range s.tasks {
		res = append(res, task)
	}
	return res
}

func (s *Storage) SetDone(id string, done bool) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if el, ok := s.tasks[id]; ok {
		el.Done = done
		el.UpdatedAt = s.clock.Now()
		return el, nil
	}

	return Task{}, TaskNotFound
}

func NewHTTPHandler(s *Storage) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", s.handleCreateTask)
	mux.HandleFunc("GET /tasks/{id}", s.handleGetTask)
	mux.HandleFunc("GET /tasks", s.handleListTasks)
	mux.HandleFunc("PATCH /tasks/{id}", s.handleUpdateTask)

	return mux
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}

func (s *Storage) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}

	task, err := s.Create(req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (s *Storage) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Field ID cannot be emtpy", http.StatusBadRequest)
		return
	}

	task, ok := s.Get(id)
	if !ok {
		http.Error(w, TaskNotFound.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(task)
}

func (s *Storage) handleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks := s.List()
	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].UpdatedAt.Equal(tasks[j].UpdatedAt) {
			return tasks[i].ID < tasks[j].ID
		}
		return tasks[i].UpdatedAt.After(tasks[j].UpdatedAt)
	})

	json.NewEncoder(w).Encode(tasks)
}

type UpdateTaskRequest struct {
	Done *bool `json:"done"`
}

func (s *Storage) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Field ID cannot be emtpy", http.StatusBadRequest)
		return
	}

	var req UpdateTaskRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Done == nil {
		http.Error(w, "Field 'done' is required", http.StatusBadRequest)
		return
	}

	task, err := s.SetDone(id, *req.Done)
	if err != nil {
		if errors.Is(err, TaskNotFound) {
			http.Error(w, TaskNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(task)
}
