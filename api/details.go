package api

import (
	"encoding/json"
	"groupie-tracker/fileio"
	"groupie-tracker/xerrors"
	"net/http"
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
	}
	defer fileio.Close(resp.Body)

	var data Details
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Details{}, err
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

	dates, err := GetDates(id)
	if err != nil {
		return AllDetails{}, err
	}

	relations, err := GetRelations(id)
	if err != nil {
		return AllDetails{}, err
	}

	location, err := GetLocation(id)
	if err != nil {
		return AllDetails{}, err
	}

	data := AllDetails{
		Details:   details,
		Dates:     dates,
		Location:  location,
		Relations: relations,
	}
	return data, nil
}
