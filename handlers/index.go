package handlers

import (
	"groupie-tracker/api"
	"html/template"
	"net/http"
)

// IndexHandler handles HTTP GET requests for the main index page.
//
// It serves as the main entry point of the application, displaying a list of all artists.
// The handler performs the following steps:
// 1. Validates that the request method is GET
// 2. Ensures the request path is exactly "/"
// 3. Fetches the list of all artists using api.GetArtists
// 4. Renders the artists data using the index.html template
//
// If any error occurs during these steps, it renders an appropriate error page
// with the corresponding HTTP status code.
//
// Parameters:
//   - w http.ResponseWriter: The response writer to send the HTTP response
//   - r *http.Request: The HTTP request
//
// The handler returns appropriate HTTP status codes:
//   - 200 OK: Successfully rendered the index page
//   - 405 Method Not Allowed: Request method is not GET
//   - 404 Not Found: URL path is not "/"
//   - 500 Internal Server Error: Server-side processing errors
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		renderErrorPage(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		renderErrorPage(w, "Page Not Found!", http.StatusNotFound)
		return
	}

	artists, err := api.GetArtists()
	if err != nil {
		renderErrorPage(w, "Internal Server Error!", http.StatusInternalServerError)
		return
	}

	temp, err := template.ParseFiles("static/templates/index.html")
	if err != nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = temp.Execute(w, artists)
	if err != nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
