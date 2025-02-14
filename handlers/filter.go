package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"groupie-tracker/api"
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

func Filter(w http.ResponseWriter, r http.Request) {
	//creation_date
	//first_album_date
	//number_of_members
	//locations_of_concerts

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//query := r.URL.Query().Get("q")
	//creation_date := r.URL.Query().Get("creation_date")
	//first_album_date := r.URL.Query().Get("first_album_date")
	//number_of_members := r.URL.Query().Get("number_of_members")
	//locations_of_concerts := r.URL.Query().Get("locations_of_concerts")
	//
	//data, locations, dates, relations := getCachedData()

}

type FilterAPIRequestData struct {
	CreationDateFilterQuery        `json:"creation_date"`
	FirstAlbumDateFilterQuery      `json:"first_album_date"`
	LocationsOfConcertsFilterQuery `json:"locations_of_concerts"`
	NumberOfMembersFilterQuery     `json:"number_of_members"`
}

type FilterAPIResponseData struct {
	Status  int          `json:"status"`
	Artists []api.Artist `json:"artists"`
}

func FilterAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content-type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Read JSON from the request body
	var requestData FilterAPIRequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Print the received data (for demonstration)
	fmt.Printf("Received data: %#v\n", requestData)

	// Create a response
	responseData := FilterAPIResponseData{
		Status:  200,
		Artists: nil,
	}

	// Filter by creation date
	{
		matchedArtists, err := filterByCreationDate(requestData.CreationDateFilterQuery)
		if err != nil {
			http.Error(w, "Invalid JSON for creation_date query", http.StatusBadRequest)
			return
		}
		responseData.Artists = append(responseData.Artists, matchedArtists...)
	}

	// Filter by first album date
	{
		matchedArtists, err := filterByFirstAlbumDate(requestData.FirstAlbumDateFilterQuery)
		if err != nil {
			http.Error(w, "Invalid JSON for first_album_date query: "+err.Error(), http.StatusBadRequest)
			return
		}
		responseData.Artists = append(responseData.Artists, matchedArtists...)
	}

	// Filter by number_of_members
	{
		matchedArtists, err := filterByNumberOfMembers(requestData.NumberOfMembersFilterQuery)
		if err != nil {
			http.Error(w, "Invalid JSON for number_of_members query", http.StatusBadRequest)
			return
		}
		responseData.Artists = append(responseData.Artists, matchedArtists...)
	}

	// Filter by locations_of_concerts
	{
		matchedArtists := filterByLocationsOfConcerts(requestData.LocationsOfConcertsFilterQuery)
		responseData.Artists = append(responseData.Artists, matchedArtists...)
	}

	// Encode the response data as JSON and send it
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func filterByCreationDate(q CreationDateFilterQuery) (result []api.Artist, err error) {
	if IsBlank(q.Type) {
		return
	}

	if !slices.Contains([]string{"range", "in", "or"}, q.Type) {
		return []api.Artist{}, errors.New("invalid query type")
	}

	if IsBlank(q.Type) {
		return
	}

	artists, _, _, _ := getCachedData()
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

func filterByLocationsOfConcerts(q LocationsOfConcertsFilterQuery) (result []api.Artist) {
	if len(q.In) == 0 {
		return
	}

	artists, _, _, _ := getCachedData()
	for _, artist := range artists {
		locations, ok := locationMapCache[artist.ID]
		if !ok {
			log.Printf("artist with ID: %d has no concert location data in cache map", artist.ID)
			continue
		}

		broke := false
		for _, loc := range locations {
			city, country := location.Parse(loc)
			loc = fmt.Sprintf("%s, %s", city, country)
			for _, in := range q.In {
				if location.Contains(loc, in) {
					result = append(result, artist)
					broke = true
					break
				}
			}
			if broke {
				break
			}
		}
	}

	return
}

func filterByNumberOfMembers(q NumberOfMembersFilterQuery) (result []api.Artist, err error) {
	if IsBlank(q.Type) {
		return
	}

	if !slices.Contains([]string{"range", "in", "or"}, q.Type) {
		return []api.Artist{}, errors.New("invalid query type")
	}

	artists, _, _, _ := getCachedData()
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

func filterByFirstAlbumDate(q FirstAlbumDateFilterQuery) (result []api.Artist, err error) {
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

	artists, _, _, _ := getCachedData()
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
