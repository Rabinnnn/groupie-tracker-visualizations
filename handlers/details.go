package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"groupie-tracker/api"
	"groupie-tracker/xerrors"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
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

	ID := r.URL.Query().Get("id")
	data, err := api.GetAllDetails(ID)
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



func SearchHandler(w http.ResponseWriter, r *http.Request) {
	initCache()
	query := r.URL.Query().Get("q")
	if query == " " {
		json.NewEncoder(w).Encode([]string{})
		return
	}
	suggestions := []string{}
	artists, locations, _, _ := getCachedData()

	for _, artist := range artists {
		// Artist/band name
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			suggestions = append(suggestions, fmt.Sprintf("%s - artist/band", artist.Name))
		}

		// Members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), strings.ToLower(query)) {
				suggestions = append(suggestions, fmt.Sprintf("%s - member", member))
			}
		}

		// First album date
		if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(query)) {
			suggestions = append(suggestions, fmt.Sprintf("%s - first album date", artist.FirstAlbum))
		}

		// Creation date
		if strings.Contains(strconv.Itoa(artist.CreationDate), query) {
			suggestions = append(suggestions, fmt.Sprintf("%d - creation date", artist.CreationDate))
		}

	}

	for _, location := range locations {
		for _, loc := range location.Locations{
			if strings.Contains(strings.ToLower(loc), strings.ToLower(query)) {
				suggestions = append(suggestions, fmt.Sprintf("%s - location", loc))
			}
		}
	}

	json.NewEncoder(w).Encode(suggestions)
}




var (
	artistCache        []api.Artist
	locationCache      []api.Location
	dateCache          []api.Date
	relationCache      []api.Relations
	cacheTime          time.Time
	cacheMutex         sync.RWMutex
	isCacheInitialized bool
)

const cacheDuration = 10 * time.Minute

func initCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if !isCacheInitialized {
		updateCache()
		isCacheInitialized = true
	}
}

func updateCache() {

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		artists, err := api.GetArtists()
		if err == nil {
			artistCache = artists
		}
	}()
	go func() {
		defer wg.Done()
		locations, err := api.GetAllLocations()
		if err == nil {
			locationCache = locations
		}
	}()

	go func() {
		defer wg.Done()
		dates, err := api.GetAllDates()
		if err == nil {
			dateCache = dates
		}
	}()

	go func() {
		defer wg.Done()
		relations, err := api.GetAllRelations()
		if err == nil {
			relationCache = relations
		}
	}()

	wg.Wait()
	cacheTime = time.Now()
}

func getCachedData() ([]api.Artist, []api.Location, []api.Date, []api.Relations) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	if time.Since(cacheTime) > cacheDuration {
		go updateCache()
	}
	
	return artistCache, locationCache, dateCache, relationCache
}