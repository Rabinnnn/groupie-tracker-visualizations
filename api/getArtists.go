package api

import (
	"net/http"
	"encoding/json"
)


func GetArtists()([]Artist,error){
	results, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil{
		return nil, err
	}
	defer results.Body.Close()

	var artists []Artist
	if err := json.NewDecoder(results.Body).Decode(&artists); err != nil {
		return nil, err
	}
	return artists, nil
}