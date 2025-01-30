package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchHandler(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		query         string
		expectedCode  int
		expectedError bool
	}{
		{
			name:          "Valid GET request",
			method:        "GET",
			query:         "queen",
			expectedCode:  http.StatusOK,
			expectedError: false,
		},
		{
			name:          "Invalid method POST",
			method:        "POST",
			query:         "queen",
			expectedCode:  http.StatusMethodNotAllowed,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Create a request with the test case parameters
				req := httptest.NewRequest(tt.method, "/?q="+tt.query, nil)
				w := httptest.NewRecorder()

				// Call the handler
				SearchHandler(w, req)

				// Check status code
				if w.Code != tt.expectedCode {
					t.Errorf(
						"SearchHandler returned wrong status code: got %v want %v",
						w.Code, tt.expectedCode,
					)
				}
			},
		)
	}
}
