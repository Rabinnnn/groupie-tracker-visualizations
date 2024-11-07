package api

import (
	"encoding/json"
	"groupie-tracker/fileio"
	"net/http"
)

// GetArtists fetches the list of artists from the external API.
// It sends an HTTP GET request to the API endpoint
// and decodes the response body into a slice of Artist structs.
//
// Returns:
// - A slice of Artist structs representing the list of artists fetched from the API.
// - An error if the network request fails or if the data cannot be decoded into the Artist struct.
func GetArtists() ([]Artist, error) {
	results, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer fileio.Close(results.Body)

	var artists []Artist
	if err := json.NewDecoder(results.Body).Decode(&artists); err != nil {
		return nil, err
	}
	return artists, nil
}
