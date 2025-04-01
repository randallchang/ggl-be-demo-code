package service

import (
	"context"
	"errors"
	"sync"
)

//go:generate mockgen -source=service.go -destination=mock/mock_service.go -package=mock

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Task struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type Service interface {
	ListTasks(ctx context.Context) ([]Task, error)
	CreateTask(ctx context.Context, name string) (*Task, error)
	UpdateTask(ctx context.Context, id int, name string, status int) (*Task, error)
	DeleteTask(ctx context.Context, id int) error
}

type TaskService struct {
	mu     sync.RWMutex
	tasks  map[int]*Task
	nextID int
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks:  make(map[int]*Task),
		nextID: 1,
	}
}

func (s *TaskService) ListTasks(ctx context.Context) ([]Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func (s *TaskService) CreateTask(ctx context.Context, name string) (*Task, error) {
	if name == "" {
		return nil, errors.New("task name cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task := &Task{
		ID:     s.nextID,
		Name:   name,
		Status: 0,
	}
	s.tasks[task.ID] = task
	s.nextID++

	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int, name string, status int) (*Task, error) {
	if name == "" {
		return nil, errors.New("task name cannot be empty")
	}
	if status != 0 && status != 1 {
		return nil, errors.New("invalid status value")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}

	task.Name = name
	task.Status = status

	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}
