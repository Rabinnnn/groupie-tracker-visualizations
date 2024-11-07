package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetLocation(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		mockResponse Location
		mockStatus   int
		expectError  bool
	}{
		{
			name: "Location id: 3",
			id:   "3",
			mockResponse: Location{
				Locations: []string{
					"london-uk",
					"lausanne-switzerland",
					"lyon-france",
				},
			},
			mockStatus:  http.StatusOK,
			expectError: false,
		},

		{
			name: "Location id: 23",
			id:   "23",
			mockResponse: Location{
				Locations: []string{
					"riyadh-saudi_arabia",
					"rio_de_janeiro-brazil",
					"canton-usa",
					"quebec-canada",
					"new_york-usa",
					"california-usa",
					"las_vegas-usa",
					"mexico_city-mexico",
					"monterrey-mexico",
					"del_mar-usa",
					"berlin-germany",
					"lisbon-portugal",
				},
			},
			mockStatus:  http.StatusOK,
			expectError: false,
		},

		{
			name: "Nonexistent Location id: 9999999",
			id:   "9999999",
			mockResponse: Location{
				Locations: nil,
			},
			mockStatus:  http.StatusNotFound,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Create a test server
				server := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							// Set response status
							w.WriteHeader(tt.mockStatus)

							// If expecting success, write mock response
							if !tt.expectError {
								err := json.NewEncoder(w).Encode(tt.mockResponse)
								if err != nil {
									t.Errorf("failed to encode response: %s", err)
								}
							}
						},
					),
				)
				defer server.Close()

				location, err := GetLocation(tt.id)

				// Check error expectation
				if tt.expectError && err == nil {
					t.Error("Expected error but got none")
				}
				if !tt.expectError && err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// For successful cases, verify the response data
				if !tt.expectError {
					if len(location.Locations) != len(tt.mockResponse.Locations) {
						t.Errorf(
							"Expected %d locations, got %d",
							len(tt.mockResponse.Locations),
							len(location.Locations),
						)
					}
				}
			},
		)
	}
}

func TestGetDates(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		mockResponse Date
		mockStatus   int
		expectError  bool
	}{
		{
			name: "Valid Date id: 4",
			id:   "4",
			mockResponse: Date{
				Dates: []string{
					"*19-02-2020",
					"*22-02-2020",
					"*24-02-2020",
					"*27-02-2020",
					"*01-03-2019",
					"*20-11-2019",
					"*18-11-2019",
					"*15-11-2019",
				},
			},
			mockStatus:  http.StatusOK,
			expectError: false,
		},
		{
			name:         "Empty ID",
			id:           "",
			mockResponse: Date{},
			mockStatus:   http.StatusNotFound,
			expectError:  true,
		},
		{
			name:         "Invalid ID",
			id:           "invalid",
			mockResponse: Date{},
			mockStatus:   http.StatusNotFound,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Create a test server
				server := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							// Verify request method
							if r.Method != http.MethodGet {
								t.Errorf("Expected GET request, got %s", r.Method)
							}

							// Set response status
							w.WriteHeader(tt.mockStatus)

							// If expecting success, write mock response
							if !tt.expectError {
								if err := json.NewEncoder(w).Encode(tt.mockResponse); err != nil {
									t.Fatalf("Failed to encode mock response: %v", err)
								}
							}
						},
					),
				)
				defer server.Close()

				// Temporarily override the API base URL
				// Note: This would require modifying the original function to accept a base URL
				// or use a configuration pattern

				dates, err := GetDates(tt.id)

				// Check error expectation
				if tt.expectError && err == nil {
					t.Error("Expected error but got none")
				}
				if !tt.expectError && err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// For successful cases, verify the response data
				if !tt.expectError {
					if len(dates.Dates) != len(tt.mockResponse.Dates) {
						t.Errorf(
							"Expected %d dates, got %d",
							len(tt.mockResponse.Dates),
							len(dates.Dates),
						)
					}

					// Verify dates content if present
					if len(tt.mockResponse.Dates) > 0 {
						if dates.Dates[0] != tt.mockResponse.Dates[0] {
							t.Errorf(
								"Expected first date %s, got %s",
								tt.mockResponse.Dates[0],
								dates.Dates[0],
							)
						}
					}
				}
			},
		)
	}
}

// TestGetDatesIntegration performs an integration test with the actual API
func TestGetDatesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test with a known valid ID
	dates, err := GetDates("1")
	if err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}

	// Verify that we got a valid response
	if dates.Dates == nil {
		t.Error("Expected non-nil Dates slice")
	}

	// Verify that dates are in expected format
	for i, date := range dates.Dates {
		if date == "" {
			t.Errorf("Empty date string at index %d", i)
		}
	}
}

// TestGetDatesConcurrent tests the function under concurrent access
func TestGetDatesConcurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	concurrentRequests := 10
	errors := make(chan error, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			_, err := GetDates("1")
			errors <- err
		}()
	}

	// Collect all errors
	for i := 0; i < concurrentRequests; i++ {
		if err := <-errors; err != nil {
			t.Errorf("Concurrent request failed: %v", err)
		}
	}
}

func TestGetRelations(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		mockResponse Relations
		mockStatus   int
		expectError  bool
	}{

		{
			name: "Valid relations id: 1",
			id:   "1",
			mockResponse: Relations{
				DatesLocation: map[string][]string{
					"dunedin-new_zealand": {"10-02-2020"},
					"georgia-usa":         {"22-08-2019"},
					"los_angeles-usa":     {"20-08-2019"},
					"nagoya-japan":        {"30-01-2019"},
					"north_carolina-usa":  {"23-08-2019"},
					"osaka-japan":         {"28-01-2020"},
					"penrose-new_zealand": {"07-02-2020"},
					"saitama-japan":       {"26-01-2020"},
				},
			},
			mockStatus:  http.StatusOK,
			expectError: false,
		},
		{
			name:         "Empty ID",
			id:           "",
			mockResponse: Relations{},
			mockStatus:   http.StatusNotFound,
			expectError:  true,
		},
		{
			name:         "Invalid ID",
			id:           "invalid",
			mockResponse: Relations{},
			mockStatus:   http.StatusNotFound,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Create a test server
				server := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							// Verify request method
							if r.Method != http.MethodGet {
								t.Errorf("Expected GET request, got %s", r.Method)
							}

							// Set response status
							w.WriteHeader(tt.mockStatus)

							// If expecting success, write mock response
							if !tt.expectError {
								if err := json.NewEncoder(w).Encode(tt.mockResponse); err != nil {
									t.Fatalf("Failed to encode mock response: %v", err)
								}
							}
						},
					),
				)
				defer server.Close()

				relations, err := GetRelations(tt.id)

				// Check error expectation
				if tt.expectError && err == nil {
					t.Error("Expected error but got none")
				}
				if !tt.expectError && err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// For successful cases, verify the response data
				if !tt.expectError {
					if !reflect.DeepEqual(relations.DatesLocation, tt.mockResponse.DatesLocation) {
						t.Errorf(
							"Expected relations %v, got %v",
							tt.mockResponse.DatesLocation,
							relations.DatesLocation,
						)
					}

					// Verify map content if present
					for date, locations := range tt.mockResponse.DatesLocation {
						gotLocations, exists := relations.DatesLocation[date]
						if !exists {
							t.Errorf("Expected date %s not found in response", date)
							continue
						}
						if !reflect.DeepEqual(locations, gotLocations) {
							t.Errorf(
								"For date %s, expected locations %v, got %v",
								date, locations, gotLocations,
							)
						}
					}
				}
			},
		)
	}
}

// TestGetRelationsIntegration performs an integration test with the actual API
func TestGetRelationsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test with a known valid ID
	relations, err := GetRelations("1")
	if err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}

	// Verify that we got a valid response
	if relations.DatesLocation == nil {
		t.Error("Expected non-nil DatesLocation map")
	}

	// Verify map structure
	for date, locations := range relations.DatesLocation {
		if date == "" {
			t.Error("Empty date string in relations map")
		}
		if len(locations) == 0 {
			t.Errorf("Empty locations array for date %s", date)
		}
		for _, location := range locations {
			if location == "" {
				t.Errorf("Empty location string for date %s", date)
			}
		}
	}
}

// TestGetRelationsConcurrent tests the function under concurrent access
func TestGetRelationsConcurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	concurrentRequests := 10
	errors := make(chan error, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			_, err := GetRelations("1")
			errors <- err
		}()
	}

	// Collect all errors
	for i := 0; i < concurrentRequests; i++ {
		if err := <-errors; err != nil {
			t.Errorf("Concurrent request failed: %v", err)
		}
	}
}

func TestGetDetails(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		mockResponse Details
		mockStatus   int
		expectError  bool
	}{
		{
			name: "Valid artist id: 10",
			id:   "10",
			mockResponse: Details{
				ID:    10,
				Name:  "Pearl Jam",
				Image: "https://groupietrackers.herokuapp.com/api/images/pearljam.jpeg",
				Members: []string{
					"Eddie Vedder",
					"Mike McCready",
					"Stone Gossard",
					"Jeff Ament",
					"Matt Cameron",
				},
				CreationDate: 1990,
				FirstAlbum:   "03-07-1992",
			},
			mockStatus:  http.StatusOK,
			expectError: false,
		},

		{
			name:         "Empty ID",
			id:           "",
			mockResponse: Details{},
			mockStatus:   http.StatusNotFound,
			expectError:  true,
		},
		{
			name:         "Invalid ID",
			id:           "invalid",
			mockResponse: Details{},
			mockStatus:   http.StatusNotFound,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Create a test server
				server := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							// Verify request method
							if r.Method != http.MethodGet {
								t.Errorf("Expected GET request, got %s", r.Method)
							}

							// Set response status
							w.WriteHeader(tt.mockStatus)

							// If expecting success, write mock response
							if !tt.expectError {
								if err := json.NewEncoder(w).Encode(tt.mockResponse); err != nil {
									t.Fatalf("Failed to encode mock response: %v", err)
								}
							}
						},
					),
				)
				defer server.Close()

				details, err := GetDetails(tt.id)

				// Check error expectation
				if tt.expectError && err == nil {
					t.Error("Expected error but got none")
				}
				if !tt.expectError && err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// For successful cases, verify the response data
				if !tt.expectError {
					// Compare all fields
					if details.ID != tt.mockResponse.ID {
						t.Errorf("Expected ID %d, got %d", tt.mockResponse.ID, details.ID)
					}
					if details.Name != tt.mockResponse.Name {
						t.Errorf("Expected Name %s, got %s", tt.mockResponse.Name, details.Name)
					}
					if details.Image != tt.mockResponse.Image {
						t.Errorf("Expected Image %s, got %s", tt.mockResponse.Image, details.Image)
					}
					if !reflect.DeepEqual(details.Members, tt.mockResponse.Members) {
						t.Errorf("Expected Members %v, got %v", tt.mockResponse.Members, details.Members)
					}
					if details.CreationDate != tt.mockResponse.CreationDate {
						t.Errorf(
							"Expected CreationDate %d, got %d",
							tt.mockResponse.CreationDate, details.CreationDate,
						)
					}
					if details.FirstAlbum != tt.mockResponse.FirstAlbum {
						t.Errorf(
							"Expected FirstAlbum %s, got %s",
							tt.mockResponse.FirstAlbum, details.FirstAlbum,
						)
					}
				}
			},
		)
	}
}

// TestGetDetailsIntegration performs an integration test with the actual API
func TestGetDetailsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test with a known valid ID
	details, err := GetDetails("1")
	if err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}

	// Verify required fields
	if details.ID <= 0 {
		t.Error("Expected positive ID")
	}
	if details.Name == "" {
		t.Error("Expected non-empty Name")
	}
	if details.Image == "" {
		t.Error("Expected non-empty Image URL")
	}
	if len(details.Members) == 0 {
		t.Error("Expected at least one member")
	}
	if details.CreationDate <= 0 {
		t.Error("Expected valid CreationDate")
	}
	if details.FirstAlbum == "" {
		t.Error("Expected non-empty FirstAlbum")
	}
}

// TestGetDetailsConcurrent tests the function under concurrent access
func TestGetDetailsConcurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	concurrentRequests := 10
	errors := make(chan error, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			_, err := GetDetails("1")
			errors <- err
		}()
	}

	// Collect all errors
	for i := 0; i < concurrentRequests; i++ {
		if err := <-errors; err != nil {
			t.Errorf("Concurrent request failed: %v", err)
		}
	}
}

func TestGetAllDetails(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		mockDetails  Details
		mockDates    Date
		mockLocation Location
		mockRelation Relations
		mockStatus   int
		expectError  bool
	}{
		{
			name: "Valid complete artist data",
			id:   "12",
			mockDetails: Details{
				ID:           12,
				Name:         "Rihanna",
				Image:        "https://groupietrackers.herokuapp.com/api/images/rihanna.jpeg",
				Members:      []string{"Robyn Rihanna Fenty"},
				CreationDate: 2003,
				FirstAlbum:   "10-09-2005",
			},
			mockDates: Date{
				Dates: []string{
					"*27-11-2016",
					"*24-09-2016",
					"*03-09-2016",
					"04-09-2016",
				},
			},
			mockLocation: Location{
				Locations: []string{
					"abu_dhabi-united_arab_emirates",
					"new_york-usa",
					"pennsylvania-usa",
				},
			},
			mockRelation: Relations{
				DatesLocation: map[string][]string{
					"abu_dhabi-united_arab_emirates": {"27-11-2016"},
					"new_york-usa":                   {"24-09-2016"},
					"pennsylvania-usa":               {"03-09-2016", "04-09-2016"},
				},
			},
			mockStatus:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "Invalid ID",
			id:          "invalid",
			mockStatus:  http.StatusNotFound,
			expectError: true,
		},
		{
			name:        "Empty ID",
			id:          "",
			mockStatus:  http.StatusNotFound,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				allDetails, err := GetAllDetails(tt.id)

				// Check error expectation
				if tt.expectError && err == nil {
					t.Error("Expected error but got none")
					return
				}
				if !tt.expectError && err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				// For successful cases, verify all components of the response
				if !tt.expectError {
					// Verify Details
					if !reflect.DeepEqual(allDetails.Details, tt.mockDetails) {
						t.Errorf(
							"Details mismatch: expected %+v, got %+v",
							tt.mockDetails, allDetails.Details,
						)
					}

					// Verify Dates
					if !reflect.DeepEqual(allDetails.Dates, tt.mockDates) {
						t.Errorf(
							"Dates mismatch: expected %+v, got %+v",
							tt.mockDates, allDetails.Dates,
						)
					}

					// Verify Location
					if !reflect.DeepEqual(allDetails.Location, tt.mockLocation) {
						t.Errorf(
							"Location mismatch: expected %+v, got %+v",
							tt.mockLocation, allDetails.Location,
						)
					}

					// Verify Relations
					if !reflect.DeepEqual(allDetails.Relations, tt.mockRelation) {
						t.Errorf(
							"Relations mismatch: expected %+v, got %+v",
							tt.mockRelation, allDetails.Relations,
						)
					}
				}
			},
		)
	}
}

// TestGetAllDetailsIntegration performs an integration test with the actual API
func TestGetAllDetailsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test with a known valid ID
	allDetails, err := GetAllDetails("1")
	if err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}

	// Verify Details
	if allDetails.Details.ID <= 0 {
		t.Error("Expected positive ID in Details")
	}
	if allDetails.Details.Name == "" {
		t.Error("Expected non-empty Name in Details")
	}
	if len(allDetails.Details.Members) == 0 {
		t.Error("Expected at least one member in Details")
	}

	// Verify Dates
	if len(allDetails.Dates.Dates) == 0 {
		t.Error("Expected at least one date")
	}

	// Verify Location
	if len(allDetails.Location.Locations) == 0 {
		t.Error("Expected at least one location")
	}

	// Verify Relations
	if len(allDetails.Relations.DatesLocation) == 0 {
		t.Error("Expected at least one relation")
	}
}

// TestGetAllDetailsConcurrent tests the function under concurrent access
func TestGetAllDetailsConcurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	concurrentRequests := 5 // Reduced due to multiple API calls per request
	errors := make(chan error, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			_, err := GetAllDetails("1")
			errors <- err
		}()
	}

	// Collect all errors
	for i := 0; i < concurrentRequests; i++ {
		if err := <-errors; err != nil {
			t.Errorf("Concurrent request failed: %v", err)
		}
	}
}
