package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthHandler tests the /health endpoint
func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	// ステータスコードの確認
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Content-Type の確認
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	// JSON レスポンスのパース
	var resp HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// レスポンス内容の確認
	if resp.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", resp.Status)
	}

	if resp.Message != "Server is running" {
		t.Errorf("Expected message 'Server is running', got '%s'", resp.Message)
	}
}

// TestRootHandler tests the / endpoint
func TestRootHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	rootHandler(w, req)

	// ステータスコードの確認
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// JSON レスポンスのパース
	var resp HelloResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// メッセージの確認
	if resp.Message == "" {
		t.Error("Expected non-empty message")
	}
}

// TestHelloHandler tests the /hello endpoint with various inputs
func TestHelloHandler(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		expectedName string
		expectedMsg  string
	}{
		{
			name:         "with name parameter",
			query:        "?name=Alice",
			expectedName: "Alice",
			expectedMsg:  "Hello, Alice!",
		},
		{
			name:         "without name parameter",
			query:        "",
			expectedName: "Guest",
			expectedMsg:  "Hello, Guest!",
		},
		{
			name:         "with CodeDeploy name",
			query:        "?name=CodeDeploy",
			expectedName: "CodeDeploy",
			expectedMsg:  "Hello, CodeDeploy!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/hello"+tt.query, nil)
			w := httptest.NewRecorder()

			helloHandler(w, req)

			// ステータスコードの確認
			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			// JSON レスポンスのパース
			var resp HelloResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			// Name フィールドの確認
			if resp.Name != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, resp.Name)
			}

			// Message フィールドの確認
			if resp.Message != tt.expectedMsg {
				t.Errorf("Expected message '%s', got '%s'", tt.expectedMsg, resp.Message)
			}
		})
	}
}

// TestLoggingMiddleware tests the logging middleware
func TestLoggingMiddleware(t *testing.T) {
	// ダミーハンドラー
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// ミドルウェアを適用
	wrappedHandler := loggingMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	// ステータスコードの確認
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
