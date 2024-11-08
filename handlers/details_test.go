package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	// During tests, the templates dir is in the parent directory
	templatesDir = filepath.Join("..", "templates")
	os.Exit(m.Run())
}

func TestDetailsHandler(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		id            string
		expectedCode  int
		expectedError bool
	}{
		{
			name:          "Valid GET request with valid ID",
			method:        "GET",
			id:            "1",
			expectedCode:  http.StatusOK,
			expectedError: false,
		},
		{
			name:          "Invalid method POST",
			method:        "POST",
			id:            "1",
			expectedCode:  http.StatusMethodNotAllowed,
			expectedError: true,
		},
		{
			name:          "Invalid ID",
			method:        "GET",
			id:            "invalid",
			expectedCode:  http.StatusNotFound,
			expectedError: true,
		},
		{
			name:          "Empty ID",
			method:        "GET",
			id:            "",
			expectedCode:  http.StatusNotFound,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Create a request with the test case parameters
				req := httptest.NewRequest(tt.method, "/details?id="+tt.id, nil)
				w := httptest.NewRecorder()

				// Call the handler
				DetailsHandler(w, req)

				// Check status code
				if w.Code != tt.expectedCode {
					t.Errorf(
						"DetailsHandler returned wrong status code: got %v want %v",
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

// TestDetailsHandlerIntegration performs an integration test
// with the actual API endpoint
func TestDetailsHandlerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	req := httptest.NewRequest("GET", "/details?id=1", nil)
	w := httptest.NewRecorder()

	DetailsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", w.Code)
	}

	// Check if response contains expected HTML elements
	if w.Body.Len() == 0 {
		t.Error("Expected non-empty response body")
	}
}
