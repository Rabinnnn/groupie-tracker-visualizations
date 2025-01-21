package handlers

import (
	//"bytes"
	//"fmt"
	//"fmt"
	//"fmt"
	"encoding/json"
	"fmt"
	"groupie-tracker/api"
	"html/template"
	//"log"

	//"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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

type TemplateData struct {
	Artists   []api.Artist
	Query     string
	NoResults bool
}
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	initCache()
	if r.Method != "GET" {
		RenderErrorPage(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		RenderErrorPage(w, "Page Not Found!", http.StatusNotFound)
		return
	}

	//artists, err := api.GetArtists()
	query := r.URL.Query().Get("query") // Get the query parameter

	artists, locations, _, _ := getCachedData()

	// if err != nil {
	// 	RenderErrorPage(w, "Internal Server Error!", http.StatusInternalServerError)
	// 	return
	// }


	// Update the Query field for each Artist
	// for i := range artists {
	// 	artists[i].Query = query
	// }	
	for i := range artists {
		// Ensure we don't exceed the length of the locations slice
		if i < len(locations) {
			// Convert locations[i].Locations (a []string) into a JSON string
			locationData, err := json.Marshal(locations[i].Locations)
			if err != nil {
				fmt.Printf("Error marshalling locations for artist %d: %v\n", artists[i].ID, err)
				continue
			}
			// Assign the serialized JSON string to the Locations field of the artist
			artists[i].Locations = string(locationData)
		}
	}
	
	//for i := range artists {
	//	artists[i].Locations = locations
		// Format the URL with the value of i
		//url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%d", i+1)

		// Fetch artist locations using the formatted URL
		//locations, err := api.FetchArtistLocations(url)
		//if err == nil {
			//artists[i].Locations = strings.Join(locations, ", ")
		//} else {
		//	log.Printf("Error fetching location for artist %d: %v", artists[i].ID, err)
		//}
//	}
	// for i := range artists{
	// 	fmt.Println(artists[i].Locations)
	// }
	filteredArtists := filterArtists(artists, query)
	//fmt.Println(filteredArtists)
	//fmt.Printf("leng:%d\n", len(filteredArtists))
	//fmt.Println(query)
	// if len(filteredArtists) == 0 && query != "" {
	// 	RenderErrorPage(w, "No Result Found for this search!", http.StatusNotFound)
	// 	return
	// }

	data := TemplateData{
		Artists:   filteredArtists,
		Query:     query,
		NoResults: len(filteredArtists) == 0 && query != "",
	}

	temp, err := template.ParseFiles(filepath.Join(templatesDir, "index.html"))
	if err != nil {
		RenderErrorPage(w, "Internal Server Error333", http.StatusInternalServerError)
		return
	}

	err = temp.Execute(w, data)
	//var buf bytes.Buffer
	//err = temp.Execute(&buf, data)
	if err != nil {
		RenderErrorPage(w, "Internal Server Error222", http.StatusInternalServerError)
		return
	}
}


// filterArtists filters the list of artists based on the search query
func filterArtists(artists []api.Artist, query string) []api.Artist {
	if query == "" {
		return artists
	}

	query = strings.ToLower(query)
	var result []api.Artist

	for _, a := range artists {
		// Artist/band name matches
		if strings.Contains(strings.ToLower(a.Name), query) {
			result = append(result, a)
			continue
		}

		// Members
		for _, member := range a.Members {
			if strings.Contains(strings.ToLower(member), query) {
				result = append(result, a)
				break
			}
		}

		// First album dates
		if strings.Contains(strings.ToLower(a.FirstAlbum), query) {
			result = append(result, a)
			continue
		}

		// creation dates
		if strings.Contains(strconv.Itoa(a.CreationDate), query) {
			result = append(result, a)
			continue
		}

		// locations
		if strings.Contains(strings.ToLower(a.Locations), strings.ToLower(query)) {
			result = append(result, a)
			continue
		}
	}
	return result
}
