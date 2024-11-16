package handlers

import (
	"errors"
	"groupie-tracker/api"
	"groupie-tracker/xerrors"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

// DetailsHandler handles HTTP GET requests for artist details.
//
// It processes requests with an "id" query parameter, validates the ID,
// fetches all details for the specified artist, and renders them using
// the detailsPage.html template.
//
// The handler performs the following steps:
// 1. Validates that the request method is GET
// 2. Extracts and validates the "id" query parameter
// 3. Fetches all artist details using api.GetAllDetails
// 4. Renders the details using the detailsPage.html template
//
// If any error occurs during these steps, it renders an appropriate error page
// with the corresponding HTTP status code.
//
// Parameters:
//   - w http.ResponseWriter: The response writer to send the HTTP response
//   - r *http.Request: The HTTP request containing the artist ID in query parameters
//
// The handler returns appropriate HTTP status codes:
//   - 200 OK: Successfully rendered artist details
//   - 405 Method Not Allowed: Request method is not GET
//   - 404 Not Found: Invalid or non-existent artist ID
//   - 500 Internal Server Error: Server-side processing errors
func DetailsHandler(w http.ResponseWriter, r *http.Request) {
	handlerTemplate := "detailsPage.html"
	if r.Method != "GET" {
		RenderErrorPage(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	data, err := api.GetAllDetails(id)
	log.Printf("Found err: %v\n", err)
	if errors.Is(err, xerrors.ErrNotFound) {
		RenderErrorPage(w, "Not Found!", http.StatusNotFound)
		log.Printf("Error is NOT Found: %v\n", err)
		return
	} else if err != nil {
		RenderErrorPage(w, "Internal Server Error!", http.StatusInternalServerError)
		log.Printf("Bad Result: Error is noooot NOT Found: %v\n", err)
		return
	}

	// Define an add function for the Go templates
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	temp, err := template.New(handlerTemplate).Funcs(funcMap).ParseFiles(filepath.Join(templatesDir, handlerTemplate))
	if err != nil {
		RenderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v\n", err)
		return
	}

	err = temp.Execute(w, data)
	if err != nil {
		RenderErrorPage(w, "Internal Server error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v\n", err)
		return
	}
}
