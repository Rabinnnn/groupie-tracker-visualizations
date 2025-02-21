package filter

import (
	"encoding/json"
	"errors"
	"fmt"
	"groupie-tracker/api"
	"groupie-tracker/cache"
	"groupie-tracker/location"
	"groupie-tracker/xtime"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

type CreationDateFilterQuery = NumberOfMembersFilterQuery

type NumberOfMembersFilterQuery struct {
	From int   `json:"from"`
	To   int   `json:"to"`
	In   []int `json:"in"`
	// Type of the filter query. One of `range`, `in`, or `or`.
	//If the query type is `range`, then, the query result will contain
	//values that are within the inclusive range [ From, To ].
	//If the query type is `in`, then, the query result will contain
	//values that are among the values specified in the array, In.
	//If the query type is `or`, then, the query result will contain
	//values that either satisfy the `range` or `in` query types, described above.
	Type string `json:"type"`
}

type FirstAlbumDateFilterQuery struct {
	From string   `json:"from"`
	To   string   `json:"to"`
	In   []string `json:"in"`
	// Type of the filter query. One of `range`, `in`, or `or`.
	//If the query type is `range`, then, the query result will contain
	//values that are within the inclusive range [ From, To ].
	//If the query type is `in`, then, the query result will contain
	//values that are among the values specified in the array, In.
	//If the query type is `or`, then, the query result will contain
	//values that either satisfy the `range` or `in` query types, described above.
	Type string `json:"type"`
}

type LocationsOfConcertsFilterQuery struct {
	In []string `json:"in"`
}

type APIRequestData struct {
	CreationDateFilterQuery        `json:"creation_date"`
	FirstAlbumDateFilterQuery      `json:"first_album_date"`
	LocationsOfConcertsFilterQuery `json:"locations_of_concerts"`
	NumberOfMembersFilterQuery     `json:"number_of_members"`
	Combinator                     string `json:"combinator"`
}

type APIResponseData struct {
	Status  int          `json:"status"`
	Artists []api.Artist `json:"artists"`
}

type APIErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func makeAPIErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	response := APIErrorResponse{
		Status:  statusCode,
		Message: message,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func API(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		makeAPIErrorResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Set content-type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Read JSON from the request body
	var requestData APIRequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		makeAPIErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	AllArtists, _, _, _, err := cache.GetCachedData()
	if err != nil {
		makeAPIErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	filteredArtistsIds := make(map[int]bool)
	filteredArtists := make([]api.Artist, 0)

	// add the given artists to the list of filtered artists if they haven't been included yet
	addArtists := func(artists []api.Artist) {
		for _, artist := range artists {
			// add this artist if it's ID doesn't yet exist
			_, ok := filteredArtistsIds[artist.ID]
			if !ok {
				filteredArtistsIds[artist.ID] = true
				filteredArtists = append(filteredArtists, artist)
			}
		}
	}

	isAnd := strings.TrimSpace(strings.ToLower(requestData.Combinator)) == "and"

	// Filter by creation date
	if requestData.CreationDateFilterQuery.Type != "" {
		matchedArtists, err := filterByCreationDate(AllArtists, requestData.CreationDateFilterQuery)
		if err != nil {
			makeAPIErrorResponse(w, http.StatusBadRequest, "Invalid JSON for creation_date query")
			return
		}

		if isAnd {
			AllArtists = matchedArtists
		} else {
			addArtists(matchedArtists)
		}
	}

	// Filter by first album date
	if requestData.FirstAlbumDateFilterQuery.Type != "" {
		matchedArtists, err := filterByFirstAlbumDate(AllArtists, requestData.FirstAlbumDateFilterQuery)
		if err != nil {
			makeAPIErrorResponse(w, http.StatusBadRequest, "Invalid JSON for first_album_date query: "+err.Error())
			return
		}

		if isAnd {
			AllArtists = matchedArtists
		} else {
			addArtists(matchedArtists)
		}
	}

	// Filter by number_of_members
	if requestData.NumberOfMembersFilterQuery.Type != "" {
		matchedArtists, err := filterByNumberOfMembers(AllArtists, requestData.NumberOfMembersFilterQuery)
		if err != nil {
			makeAPIErrorResponse(w, http.StatusBadRequest, "Invalid JSON for number_of_members query")
			return
		}

		if isAnd {
			AllArtists = matchedArtists
		} else {
			addArtists(matchedArtists)
		}
	}

	// Filter by locations_of_concerts
	if len(requestData.LocationsOfConcertsFilterQuery.In) > 0 {
		matchedArtists, err := filterByLocationsOfConcerts(AllArtists, requestData.LocationsOfConcertsFilterQuery)
		if err != nil {
			makeAPIErrorResponse(w, http.StatusBadRequest, "InternalServerError: Cache Map error")
			return
		}

		if isAnd {
			AllArtists = matchedArtists
		} else {
			addArtists(matchedArtists)
		}
	}

	if isAnd {
		addArtists(AllArtists)
	}

	// Create a response
	responseData := APIResponseData{
		Status:  200,
		Artists: filteredArtists,
	}

	// Encode the response data as JSON and send it
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func filterByCreationDate(artists []api.Artist, q CreationDateFilterQuery) (result []api.Artist, err error) {
	if IsBlank(q.Type) {
		return
	}

	if !slices.Contains([]string{"range", "in", "or"}, q.Type) {
		return []api.Artist{}, errors.New("invalid query type")
	}

	if IsBlank(q.Type) {
		return
	}

	for _, artist := range artists {
		if q.Type == "range" && (artist.CreationDate >= q.From && artist.CreationDate <= q.To) {
			result = append(result, artist)
		}

		if q.Type == "in" && slices.Contains(q.In, artist.CreationDate) {
			result = append(result, artist)
		}

		if q.Type == "or" && ((artist.CreationDate >= q.From && artist.CreationDate <= q.To) ||
			slices.Contains(q.In, artist.CreationDate)) {
			result = append(result, artist)
		}
	}

	return
}

func filterByLocationsOfConcerts(artists []api.Artist, q LocationsOfConcertsFilterQuery) (
	result []api.Artist, err error,
) {
	if len(q.In) == 0 {
		return
	}

	locationsMap := cache.GetCachedLocationsMap()
	if locationsMap == nil {
		return nil, errors.New("locationsMap is nil")
	}

	for _, artist := range artists {
		locations, ok := locationsMap[artist.ID]
		if !ok {
			log.Printf("artist with ID: %d has no concert location data in cache map", artist.ID)
			continue
		}

		broke := false
		for _, hyphenatedLocation := range locations {
			loc := hyphenatedLocation
			city, country := location.Parse(loc)
			loc = fmt.Sprintf("%s, %s", city, country)
			for _, in := range q.In {
				// The user may have entered the location in the hyphenated location format
				possibleHyphenatedLocation := in
				city, country := location.GetCityCountry(in)
				in = fmt.Sprintf("%s, %s", city, country)

				if location.Contains(loc, in) || location.Contains(hyphenatedLocation, possibleHyphenatedLocation) {
					result = append(result, artist)
					// break from the outer loop as well
					broke = true
					break
				}
			}
			// This artist has already been included in the result, other checks are unnecessary
			if broke {
				break
			}
		}
	}

	return
}

func filterByNumberOfMembers(artists []api.Artist, q NumberOfMembersFilterQuery) (result []api.Artist, err error) {
	if IsBlank(q.Type) {
		return
	}

	if !slices.Contains([]string{"range", "in", "or"}, q.Type) {
		return []api.Artist{}, errors.New("invalid query type")
	}

	if q.Type == "or" {
		for _, artist := range artists {
			numberOfMembers := len(artist.Members)
			if (numberOfMembers >= q.From && numberOfMembers <= q.To) ||
				slices.Contains(q.In, numberOfMembers) {
				result = append(result, artist)
			}
		}
	} else if q.Type == "range" {
		for _, artist := range artists {
			numberOfMembers := len(artist.Members)
			if numberOfMembers >= q.From && numberOfMembers <= q.To {
				result = append(result, artist)
			}
		}
	} else if q.Type == "in" {
		for _, artist := range artists {
			numberOfMembers := len(artist.Members)
			if slices.Contains(q.In, numberOfMembers) {
				result = append(result, artist)
			}
		}
	}

	return
}

func filterByFirstAlbumDate(artists []api.Artist, q FirstAlbumDateFilterQuery) (result []api.Artist, err error) {
	if IsBlank(q.Type) {
		return
	}

	if !slices.Contains([]string{"range", "in", "or"}, q.Type) {
		return []api.Artist{}, errors.New("invalid query type")
	}

	compare := func(a time.Time, comparator string, b time.Time) bool {
		if comparator == ">=" {
			return a.Equal(b) || a.After(b)
		} else if comparator == "<=" {
			return a.Equal(b) || a.Before(b)
		} else {
			return false
		}
	}

	for _, artist := range artists {
		typeRange := func() error {
			qFrom, err := xtime.Parse(q.From)
			if err != nil {
				return fmt.Errorf("invalid query from: %s", q.From)
			}

			qTo, err := xtime.Parse(q.To)
			if err != nil {
				return fmt.Errorf("invalid query to: %s", q.To)
			}

			artistFirstAlbumDate, err := xtime.Parse(artist.FirstAlbum)
			if err != nil {
				return fmt.Errorf("invalid first album date format: %v, for artist: %s", artist.FirstAlbum, artist.Name)
			}

			if compare(artistFirstAlbumDate, ">=", qFrom) &&
				compare(artistFirstAlbumDate, "<=", qTo) {
				result = append(result, artist)
			}

			return nil
		}

		typeIn := func() error {
			artistFirstAlbumDate, err := xtime.Parse(artist.FirstAlbum)
			if err != nil {
				return fmt.Errorf("invalid first album date format: %v, for artist: %s", artist.FirstAlbum, artist.Name)
			}

			var qIn []time.Time
			for _, in := range q.In {
				currentIn, err := xtime.Parse(in)
				if err != nil {
					return err
				}

				qIn = append(qIn, currentIn)
			}

			if slices.Contains(qIn, artistFirstAlbumDate) {
				result = append(result, artist)
			}

			return nil
		}

		if q.Type == "range" {
			err := typeRange()
			if err != nil {
				return result, err
			}
		} else if q.Type == "in" {
			err := typeIn()
			if err != nil {
				return result, err
			}
		} else if q.Type == "or" {
			err := typeRange()
			if err != nil {
				return result, err
			}

			err = typeIn()
			if err != nil {
				return result, err
			}
		}
	}

	return
}

func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
