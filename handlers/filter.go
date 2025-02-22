package handlers

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

func Filter(w http.ResponseWriter, r *http.Request) {
	handlerTemplate := "filter.html"
	if r.Method != "GET" {
		RenderErrorPage(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}

	data := struct{}{}

	ID := r.URL.Query().Get("id")
	log.Printf(ID)

	temp, err := template.New(handlerTemplate).ParseFiles(filepath.Join(templatesDir, handlerTemplate))
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
