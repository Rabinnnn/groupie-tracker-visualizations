package handlers

import (
	"encoding/json"
	"groupie-tracker/cache"
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

	// Template data
	data := struct{ ArtistsJson string }{}

	artists, _, _, _, err := cache.GetCachedData()
	if err != nil {
		log.Println(err)
		RenderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jsonBytes, _ := json.Marshal(artists)
	data.ArtistsJson = string(jsonBytes)

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
