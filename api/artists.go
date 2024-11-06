package api

import (
	"encoding/json"
	"groupie-tracker/fileio"
	"net/http"
)

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

func GetLength() (int, error) {
	artists, err := GetArtists()
	if err != nil {
		return 0, err
	}

	return len(artists), nil
}
