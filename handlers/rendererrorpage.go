package handlers

import (
	"groupie-tracker/api"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

// RenderErrorPage renders an HTML error page with the provided error message and HTTP status code.
//
// Parameters:
//   - w http.ResponseWriter: The response writer to send the rendered HTML page
//   - errorText string: The error message to display on the page
//   - statusCode int: The HTTP status code to set in the response
//
// The function sets the HTTP status code, then renders the error page template (errorPage.html)
// with the provided error message and status code. If there are any errors during template
// parsing or execution, it returns a generic 500 Internal Server Error response.
//
// Example usage:
//
//	RenderErrorPage(w, "Resource not found", http.StatusNotFound)
func RenderErrorPage(w http.ResponseWriter, errorText string, statusCode int) {
	w.WriteHeader(statusCode)

	content := api.ErrorContent{
		Message: errorText,
		Code:    strconv.Itoa(statusCode),
	}

	temp, err := template.ParseFiles(filepath.Join(templatesDir, "errorPage.html"))
	if err != nil {
		http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		log.Printf("Error parsing templates: %v", err)
		return
	}

	err = temp.Execute(w, content)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}
