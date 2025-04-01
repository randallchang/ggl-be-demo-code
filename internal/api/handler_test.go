package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/randallchang/ggl-be-demo-code/internal/service"
	"github.com/randallchang/ggl-be-demo-code/internal/service/mock"
)

func setupTest(t *testing.T) (*mock.MockService, *Handler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockService := mock.NewMockService(ctrl)
	handler := NewHandler(mockService)
	router := gin.New()
	SetupRoutes(router, handler)
	return mockService, handler, router
}

func TestListTasks(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		// Given
		mockService, _, router := setupTest(t)
		expectedTasks := []service.Task{{ID: 1, Name: "Task 1", Status: 0}}
		mockService.EXPECT().ListTasks(gomock.Any()).Return(expectedTasks, nil)

		// When
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		router.ServeHTTP(w, req)

		// Then
		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []service.Task
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if len(response) != 1 || response[0].ID != expectedTasks[0].ID {
			t.Error("unexpected response data")
		}
	})
}

func TestCreateTask(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		// Given
		mockService, _, router := setupTest(t)
		expectedTask := &service.Task{ID: 1, Name: "New Task", Status: 0}
		mockService.EXPECT().CreateTask(gomock.Any(), "New Task").Return(expectedTask, nil)

		// When
		w := httptest.NewRecorder()
		body := bytes.NewBufferString(`{"name":"New Task"}`)
		req, _ := http.NewRequest(http.MethodPost, "/tasks", body)
		router.ServeHTTP(w, req)

		// Then
		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response service.Task
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if response.Name != expectedTask.Name {
			t.Error("unexpected response data")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		// Given
		_, _, router := setupTest(t)

		// When
		w := httptest.NewRecorder()
		body := bytes.NewBufferString(`{"name":""}`)
		req, _ := http.NewRequest(http.MethodPost, "/tasks", body)
		router.ServeHTTP(w, req)

		// Then
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestUpdateTask(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		// Given
		mockService, _, router := setupTest(t)
		expectedTask := &service.Task{ID: 1, Name: "Updated Task", Status: 1}
		mockService.EXPECT().UpdateTask(gomock.Any(), 1, "Updated Task", 1).Return(expectedTask, nil)

		// When
		w := httptest.NewRecorder()
		body := bytes.NewBufferString(`{"name":"Updated Task","status":1}`)
		req, _ := http.NewRequest(http.MethodPut, "/tasks/1", body)
		router.ServeHTTP(w, req)

		// Then
		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response service.Task
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if response.Name != expectedTask.Name || response.Status != expectedTask.Status {
			t.Error("unexpected response data")
		}
	})

	t.Run("task not found", func(t *testing.T) {
		// Given
		mockService, _, router := setupTest(t)
		mockService.EXPECT().UpdateTask(gomock.Any(), 999, "Updated Task", 1).
			Return(nil, service.ErrTaskNotFound)

		// When
		w := httptest.NewRecorder()
		body := bytes.NewBufferString(`{"name":"Updated Task","status":1}`)
		req, _ := http.NewRequest(http.MethodPut, "/tasks/999", body)
		router.ServeHTTP(w, req)

		// Then
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		// Given
		mockService, _, router := setupTest(t)
		mockService.EXPECT().DeleteTask(gomock.Any(), 1).Return(nil)

		// When
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
		router.ServeHTTP(w, req)

		// Then
		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}
	})

	t.Run("task not found", func(t *testing.T) {
		// Given
		mockService, _, router := setupTest(t)
		mockService.EXPECT().DeleteTask(gomock.Any(), 999).Return(service.ErrTaskNotFound)

		// When
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/tasks/999", nil)
		router.ServeHTTP(w, req)

		// Then
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestPing(t *testing.T) {
	t.Run("ping endpoint", func(t *testing.T) {
		// Given
		_, _, router := setupTest(t)

		// When
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		router.ServeHTTP(w, req)

		// Then
		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if response["message"] != "pong" {
			t.Error("unexpected response message")
		}
	})
}
