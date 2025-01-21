package api

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/fileio"
	"groupie-tracker/xerrors"
	"net/http"
	"sync"
	"io"
)

// GetLocation retrieves location data for a specific artist from the groupietrackers API.
//
// The function makes an HTTP GET request to the locations endpoint of the groupietrackers API
// using the provided artist ID. It then decodes the JSON response into a Location struct.
//
// Parameters:
//   - id string: The unique identifier of the artist whose locations are being requested
//
// Returns:
//   - Location: A struct containing an array of location strings
//   - error: An error if the request fails, the response is invalid, or the JSON decoding fails
//
// The function handles the following error cases:
//   - Network connectivity issues
//   - Invalid HTTP responses
//   - JSON decoding errors
//
// Example usage:
//
//	locations, err := GetLocation("1")
//	if err != nil {
//	    log.Printf("Failed to get locations: %v", err)
//	    return
//	}
func GetLocation(id string) (Location, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations/" + id)
	if err != nil {
		return Location{}, err
	} else if resp.StatusCode == 404 {
		return Location{}, xerrors.ErrNotFound
	} else if resp.StatusCode != 200 {
		return Location{}, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	defer fileio.Close(resp.Body)

	var data Location
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Location{}, err
	}
	return data, nil
}

// GetDates retrieves concert dates data for a specific artist from the groupietrackers API.
//
// The function makes an HTTP GET request to the dates endpoint of the groupietrackers API
// using the provided artist ID. It then decodes the JSON response into a Date struct.
//
// Parameters:
//   - id string: The unique identifier of the artist whose concert dates are being requested
//
// Returns:
//   - Date: A struct containing an array of concert date strings
//   - error: An error if the request fails, the response is invalid, or the JSON decoding fails
//
// The function handles the following error cases:
//   - Network connectivity issues
//   - Invalid HTTP responses
//   - JSON decoding errors
//
// Example usage:
//
//	dates, err := GetDates("1")
//	if err != nil {
//	    log.Printf("Failed to get concert dates: %v", err)
//	    return
//	}
func GetDates(id string) (Date, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/dates/" + id)
	if err != nil {
		return Date{}, err
	} else if resp.StatusCode == 404 {
		return Date{}, xerrors.ErrNotFound
	} else if resp.StatusCode != 200 {
		return Date{}, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	defer fileio.Close(resp.Body)

	var data Date
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Date{}, err
	}
	return data, nil
}

// GetRelations retrieves the relationship data between concert dates and locations
// for a specific artist from the groupietrackers API.
//
// The function makes an HTTP GET request to the relation endpoint of the groupietrackers API
// using the provided artist ID. It then decodes the JSON response into a Relations struct,
// which contains a map of dates to their corresponding locations.
//
// Parameters:
//   - id string: The unique identifier of the artist whose relations are being requested
//
// Returns:
//   - Relations: A struct containing a map of dates to location arrays
//   - error: An error if the request fails, the response is invalid, or the JSON decoding fails
//
// The function handles the following error cases:
//   - Network connectivity issues
//   - Invalid HTTP responses
//   - JSON decoding errors
//
// Example usage:
//
//	relations, err := GetRelations("1")
//	if err != nil {
//	    log.Printf("Failed to get relations: %v", err)
//	    return
//	}
func GetRelations(id string) (Relations, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + id)
	if err != nil {
		return Relations{}, err
	} else if resp.StatusCode == 404 {
		return Relations{}, xerrors.ErrNotFound
	} else if resp.StatusCode != 200 {
		return Relations{}, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	defer fileio.Close(resp.Body)

	var data Relations
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Relations{}, err
	}
	return data, nil
}

// GetDetails retrieves detailed information about a specific artist from the groupietrackers API.
//
// The function makes an HTTP GET request to the artists endpoint of the groupietrackers API
// using the provided artist ID. It then decodes the JSON response into a Details struct containing
// basic information about the artist such as name, image, members, creation date, and first album.
//
// Parameters:
//   - id string: The unique identifier of the artist whose details are being requested
//
// Returns:
//   - Details: A struct containing the artist's basic information
//   - error: An error if the request fails, the response is invalid, or the JSON decoding fails
//
// The function handles the following error cases:
//   - Network connectivity issues
//   - Invalid HTTP responses
//   - JSON decoding errors
//
// Example usage:
//
//	details, err := GetDetails("1")
//	if err != nil {
//	    log.Printf("Failed to get artist details: %v", err)
//	    return
//	}
func GetDetails(id string) (Details, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + id)
	if err != nil {
		return Details{}, err
	} else if resp.StatusCode == 404 {
		return Details{}, xerrors.ErrNotFound
	} else if resp.StatusCode != 200 {
		return Details{}, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	defer fileio.Close(resp.Body)

	var data Details
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Details{}, xerrors.ErrNotFound
	} else if data.ID == 0 {
		return Details{}, xerrors.ErrNotFound
	}
	return data, nil
}

// GetAllDetails retrieves comprehensive information about a specific artist by aggregating
// data from multiple endpoints of the groupietrackers API.
//
// The function makes multiple API calls to collect:
// - Basic artist details (name, image, members, etc.)
// - Concert dates
// - Location information
// - Relations between dates and locations
//
// It combines all this information into a single AllDetails struct for convenient access.
//
// Parameters:
//   - id string: The unique identifier of the artist whose complete information is being requested
//
// Returns:
//   - AllDetails: A struct containing all available information about the artist, including:
//   - Details: Basic artist information (name, image, members, etc.)
//   - Dates: Concert dates
//   - Location: Performance locations
//   - Relations: Mapping between dates and locations
//   - error: An error if any of the API calls fail or if the data cannot be processed
//
// The function will return an error if any of the following occurs:
//   - Network connectivity issues
//   - Invalid responses from any endpoint
//   - JSON decoding errors
//   - Invalid or non-existent artist ID
//
// Example usage:
//
//	allDetails, err := GetAllDetails("1")
//	if err != nil {
//	    log.Printf("Failed to get complete artist information: %v", err)
//	    return
//	}
//	// Access combined information
//	fmt.Printf("Artist: %s\n", allDetails.Details.Name)
//	fmt.Printf("Upcoming concerts: %v\n", allDetails.Dates.Dates)
//
// Note: This function makes multiple API calls, so it may take longer to complete
// than individual endpoint calls.
func GetAllDetails(id string) (AllDetails, error) {
	details, err := GetDetails(id)
	if err != nil {
		return AllDetails{}, err
	}

	data := AllDetails{
		Details: details,
	}

	// Speed up the other fetch with goroutines
	wg := sync.WaitGroup{}
	var errs = [3]error{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		data.Dates, errs[0] = GetDates(id)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data.Relations, errs[1] = GetRelations(id)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data.Location, errs[2] = GetLocation(id)
	}()

	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return AllDetails{}, err
		}
	}

	return data, nil
}




var (
	LocationsURL = "https://groupietrackers.herokuapp.com/api/locations"
	DatesURL     = "https://groupietrackers.herokuapp.com/api/dates"
	RelationURL  = "https://groupietrackers.herokuapp.com/api/relation"
)

// FetchData makes an HTTP GET request to the given URL and returns the response body
func FetchData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}


// GetLocations fetches the location data from the API and returns a slice of Location structs
func GetAllLocations() ([]Location, error) {
	body, err := FetchData(LocationsURL)
	if err != nil {
		return nil, err
	}

	var locations struct {
		Index []Location `json:"index"`
	}

	if err := json.Unmarshal(body, &locations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal locations: %v", err)
	}
	return locations.Index, nil
}

// GetDates fetches the date data from the API and returns a slice of Date structs
func GetAllDates() ([]Date, error) {
	body, err := FetchData(DatesURL)
	if err != nil {
		return nil, err
	}

	var dates struct {
		Index []Date `json:"index"`
	}

	if err := json.Unmarshal(body, &dates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dates: %v", err)
	}

	return dates.Index, nil
}

// GetRelations fetches the relation data from the API and returns a slice of Relation structs
func GetAllRelations() ([]Relations, error) {
	body, err := FetchData(RelationURL)
	if err != nil {
		return nil, err
	}

	var relations struct {
		Index []Relations `json:"index"`
	}
	if err := json.Unmarshal(body, &relations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal relations: %v", err)
	}

	return relations.Index, nil
}
