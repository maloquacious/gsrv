package gsrv

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthHandler(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	server.started = time.Now().Add(-5 * time.Minute) // simulate 5 minutes uptime

	handler := server.HealthHandler()

	t.Run("GET request returns health status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
		}

		var response map[string]string
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if _, ok := response["uptime"]; !ok {
			t.Error("expected uptime field in response")
		}
	})

	t.Run("non-GET request returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/health", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestShutdownHandler(t *testing.T) {
	server, err := New(WithShutdownKey("test-key"))
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	handler := server.ShutdownHandler()

	t.Run("POST with correct key returns 202", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/shutdown/test-key", nil)
		req.SetPathValue("key", "test-key")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusAccepted {
			t.Errorf("expected status %d, got %d", http.StatusAccepted, w.Code)
		}

		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
		}

		var response map[string]string
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response["status"] != "server shutting down" {
			t.Errorf("expected status 'server shutting down', got %s", response["status"])
		}
	})

	t.Run("POST with wrong key returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/shutdown/wrong-key", nil)
		req.SetPathValue("key", "wrong-key")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("non-POST request returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shutdown/test-key", nil)
		req.SetPathValue("key", "test-key")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestShutdownHandlerWithoutKey(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	server.admin.keys.shutdown = "" // simulate no shutdown key

	handler := server.ShutdownHandler()

	req := httptest.NewRequest(http.MethodPost, "/shutdown/any-key", nil)
	req.SetPathValue("key", "any-key")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
