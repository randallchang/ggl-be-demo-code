package service

import (
	"context"
	"sync"
	"testing"
)

func TestListTasks(t *testing.T) {
	t.Run("empty task list", func(t *testing.T) {
		// Given
		svc := NewTaskService()

		// When
		tasks, err := svc.ListTasks(context.Background())

		// Then
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(tasks) != 0 {
			t.Errorf("expected empty task list, got %d tasks", len(tasks))
		}
	})

	t.Run("multiple tasks", func(t *testing.T) {
		// Given
		svc := NewTaskService()
		task1, _ := svc.CreateTask(context.Background(), "Task 1")
		task2, _ := svc.CreateTask(context.Background(), "Task 2")

		// When
		tasks, err := svc.ListTasks(context.Background())

		// Then
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(tasks) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(tasks))
		}
		if tasks[0].ID != task1.ID || tasks[1].ID != task2.ID {
			t.Error("tasks not in expected order")
		}
	})
}

func TestCreateTask(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		// Given
		svc := NewTaskService()

		// When
		task, err := svc.CreateTask(context.Background(), "New Task")

		// Then
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if task.Name != "New Task" {
			t.Errorf("expected task name 'New Task', got '%s'", task.Name)
		}
		if task.Status != 0 {
			t.Errorf("expected status 0, got %d", task.Status)
		}
		if task.ID != 1 {
			t.Errorf("expected ID 1, got %d", task.ID)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		// Given
		svc := NewTaskService()

		// When
		task, err := svc.CreateTask(context.Background(), "")

		// Then
		if err == nil {
			t.Error("expected error for empty name, got nil")
		}
		if task != nil {
			t.Error("expected nil task, got non-nil")
		}
	})
}

func TestUpdateTask(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		// Given
		svc := NewTaskService()
		task, _ := svc.CreateTask(context.Background(), "Original")

		// When
		updated, err := svc.UpdateTask(context.Background(), task.ID, "Updated", 1)

		// Then
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if updated.Name != "Updated" {
			t.Errorf("expected name 'Updated', got '%s'", updated.Name)
		}
		if updated.Status != 1 {
			t.Errorf("expected status 1, got %d", updated.Status)
		}
	})

	t.Run("task not found", func(t *testing.T) {
		// Given
		svc := NewTaskService()

		// When
		updated, err := svc.UpdateTask(context.Background(), 999, "Updated", 1)

		// Then
		if err == nil {
			t.Error("expected error for non-existent task, got nil")
		}
		if updated != nil {
			t.Error("expected nil task, got non-nil")
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		// Given
		svc := NewTaskService()
		task, _ := svc.CreateTask(context.Background(), "Original")

		// When
		updated, err := svc.UpdateTask(context.Background(), task.ID, "Updated", 2)

		// Then
		if err == nil {
			t.Error("expected error for invalid status, got nil")
		}
		if updated != nil {
			t.Error("expected nil task, got non-nil")
		}
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		// Given
		svc := NewTaskService()
		task, _ := svc.CreateTask(context.Background(), "To Delete")

		// When
		err := svc.DeleteTask(context.Background(), task.ID)

		// Then
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		tasks, _ := svc.ListTasks(context.Background())
		if len(tasks) != 0 {
			t.Error("task was not deleted")
		}
	})

	t.Run("task not found", func(t *testing.T) {
		// Given
		svc := NewTaskService()

		// When
		err := svc.DeleteTask(context.Background(), 999)

		// Then
		if err == nil {
			t.Error("expected error for non-existent task, got nil")
		}
	})
}

func TestConcurrency(t *testing.T) {
	t.Run("concurrent reads", func(t *testing.T) {
		// Given
		svc := NewTaskService()
		svc.CreateTask(context.Background(), "Task")

		// When
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				tasks, err := svc.ListTasks(context.Background())
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(tasks) != 1 {
					t.Errorf("expected 1 task, got %d", len(tasks))
				}
			}()
		}
		wg.Wait()
	})

	t.Run("concurrent writes", func(t *testing.T) {
		// Given
		svc := NewTaskService()

		// When
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				_, err := svc.CreateTask(context.Background(), "Task")
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}(i)
		}
		wg.Wait()

		// Then
		tasks, _ := svc.ListTasks(context.Background())
		if len(tasks) != 10 {
			t.Errorf("expected 10 tasks, got %d", len(tasks))
		}
	})

	t.Run("concurrent read/write", func(t *testing.T) {
		// Given
		svc := NewTaskService()

		// When
		var wg sync.WaitGroup
		for i := 0; i < 5; i++ {
			// Add reader
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _ = svc.ListTasks(context.Background())
			}()

			// Add writer
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				_, _ = svc.CreateTask(context.Background(), "Task")
			}(i)
		}
		wg.Wait()

		// Then
		tasks, _ := svc.ListTasks(context.Background())
		if len(tasks) != 5 {
			t.Errorf("expected 5 tasks, got %d", len(tasks))
		}
	})
}
