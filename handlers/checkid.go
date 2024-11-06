package handlers

import (
	"groupie-tracker/api"
	"strconv"
)

func CheckId(id string) bool {
	num, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	allArtists, err := api.GetLength()
	if err != nil {
		return false
	}

	if !(num > 0 && num <= allArtists) {
		return false
	}
	return true
}
