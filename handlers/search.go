package handlers

import (
	"encoding/json"
	"groupie-tracker/cache"
	"net/http"
	"strconv"
	"strings"
)

type SearchHandlerResponse struct {
	Suggestion string `json:"suggestion"`
	From       string `json:"from"`
}

// SearchHandler exposes a GET request API that accepts a query for a search for an artist,
// album, or concert location.
//
//				Example usage:
//
//				Request query: `/?q=queen`
//
//				Response:
//
//				```json
//				[
//		 			{
//		 			  "suggestion": "Queen",
//		 			  "from": "artist/band"
//		 			},
//		 			{
//		 			  "suggestion": "queensland-australia",
//		 			  "from": "location"
//		 			}
//				]
//				```
//
//				Making a request with the `init` query, allows for initialization functions to make blank queries, that
//				may be helpful when initializing search suggestions.
//
//			 Below is an example of how to initialize all search suggestions:
//
//				Request query: `/?init=true`
//
//				Response:
//
//				```json
//				[
//	    			{
//	    			  "suggestion": "Queen",
//	    			  "from": "artist/band"
//	    			},
//	    			{
//	    			  "suggestion": "queensland-australia",
//	    			  "from": "location"
//	    			}
//	  		]
//	  		```
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	// whether this api request should return all suggestions if the query, q, is blank
	initSuggestions := false
	{
		initQuery := r.URL.Query().Get("init")
		initSuggestions = initQuery == "true"
	}

	// ignore empty search queries, return an empty suggestion list
	if strings.TrimSpace(query) == "" && !initSuggestions {
		_ = json.NewEncoder(w).Encode([]SearchHandlerResponse{})
		return
	}

	var suggestions []SearchHandlerResponse
	artists, locations, _, _, err := cache.GetCachedData()
	if err != nil {
		_ = json.NewEncoder(w).Encode(suggestions)
		return
	}

	for _, artist := range artists {
		// Artist/band name
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			suggestions = append(
				suggestions, SearchHandlerResponse{
					Suggestion: artist.Name,
					From:       "artist/band",
				},
			)
		}

		// Members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), strings.ToLower(query)) {
				suggestions = append(
					suggestions, SearchHandlerResponse{
						Suggestion: member,
						From:       "member (" + artist.Name + ")",
					},
				)
			}
		}

		// First album date
		if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(query)) {
			suggestions = append(
				suggestions, SearchHandlerResponse{
					Suggestion: artist.FirstAlbum,
					From:       "first album date (" + artist.Name + ")",
				},
			)
		}

		// Creation date
		if strings.Contains(strconv.Itoa(artist.CreationDate), query) {
			suggestions = append(
				suggestions, SearchHandlerResponse{
					Suggestion: strconv.Itoa(artist.CreationDate),
					From:       "creation date (" + artist.Name + ")",
				},
			)
		}

	}

	for _, location := range locations {
		for _, loc := range location.Locations {
			if strings.Contains(strings.ToLower(loc), strings.ToLower(query)) {
				suggestions = append(
					suggestions, SearchHandlerResponse{
						Suggestion: loc,
						From:       "location",
					},
				)
			}
		}
	}

	_ = json.NewEncoder(w).Encode(suggestions)
}
