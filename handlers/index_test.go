package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		path          string
		expectedCode  int
		expectedError bool
	}{
		{
			name:          "Valid GET request to root",
			method:        "GET",
			path:          "/",
			expectedCode:  http.StatusOK,
			expectedError: false,
		},
		{
			name:          "Invalid method POST",
			method:        "POST",
			path:          "/",
			expectedCode:  http.StatusMethodNotAllowed,
			expectedError: true,
		},
		{
			name:          "Invalid path",
			method:        "GET",
			path:          "/invalid",
			expectedCode:  http.StatusNotFound,
			expectedError: true,
		},
		{
			name:          "Root path with trailing slash",
			method:        "GET",
			path:          "//",
			expectedCode:  http.StatusNotFound,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Create a request with the test case parameters
				req := httptest.NewRequest(tt.method, tt.path, nil)
				w := httptest.NewRecorder()

				// Call the handler
				IndexHandler(w, req)

				// Check status code
				if w.Code != tt.expectedCode {
					t.Errorf(
						"IndexHandler returned wrong status code: got %v want %v",
						w.Code, tt.expectedCode,
					)
				}

				// Check if response contains error page when expected
				if tt.expectedError {
					if w.Body.Len() == 0 {
						t.Error("Expected error page content, got empty response")
					}
				}
			},
		)
	}
}

// TestIndexHandlerIntegration performs an integration test
// with the actual API endpoint
func TestIndexHandlerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	IndexHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", w.Code)
	}

	// Check if response contains expected HTML elements
	response := w.Body.String()
	expectedElements := []string{
		"<html",
		"<body",
		"</html>",
	}

	for _, element := range expectedElements {
		if !strings.Contains(response, element) {
			t.Errorf("Response missing expected element: %s", element)
		}
	}
}

func TestIndexHandlerNoTemplates(t *testing.T) {
	originalTemplateDir := templatesDir
	templatesDir = ""
	defer func() {
		templatesDir = originalTemplateDir
	}()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	IndexHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %v; got %v", http.StatusInternalServerError, w.Code)
	}
}
