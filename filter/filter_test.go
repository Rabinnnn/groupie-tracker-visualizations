package filter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI(t *testing.T) {
	tests := []struct {
		name       string
		request    APIRequestData
		expected   []string
		statusCode int
		method     int
	}{
		{
			name: "Filter artists between creation dates 1995 and 2000",
			request: APIRequestData{
				CreationDateFilterQuery: CreationDateFilterQuery{
					From: 1995,
					To:   2000,
					Type: "range",
				},
				Combinator: "",
			},
			expected:   []string{"SOJA", "Mamonas Assassinas", "Thirty Seconds to Mars", "Nickelback", "NWA", "Gorillaz", "Linkin Park", "Eminem", "Coldplay"},
			statusCode: http.StatusOK,
		},
		{
			name: "Filter artists by first album between 1990 and 1992",
			request: APIRequestData{
				FirstAlbumDateFilterQuery: FirstAlbumDateFilterQuery{
					From: "01-01-1990",
					To:   "31-12-1992",
					Type: "range",
				},
				Combinator: "",
			},
			expected:   []string{"Pearl Jam", "Red Hot Chili Peppers"},
			statusCode: http.StatusOK,
		},
		{
			name: "Filter artists with exactly 6 members",
			request: APIRequestData{
				NumberOfMembersFilterQuery: NumberOfMembersFilterQuery{
					In:   []int{6},
					Type: "in",
				},
				Combinator: "",
			},
			expected:   []string{"Pink Floyd", "Arctic Monkeys", "Linkin Park", "Foo Fighters"},
			statusCode: http.StatusOK,
		},
		{
			name: "Filter artists with concerts in Texas, USA",
			request: APIRequestData{
				LocationsOfConcertsFilterQuery: LocationsOfConcertsFilterQuery{
					In: []string{"Texas, USA"},
				},
				Combinator: "",
			},
			expected:   []string{"R3HAB", "Logic", "Joyner Lucas", "Twenty One Pilots"},
			statusCode: http.StatusOK,
		},
		{
			name: "Filter solo artists between creation dates 1970 and 2000",
			request: APIRequestData{
				CreationDateFilterQuery: NumberOfMembersFilterQuery{
					From: 1970,
					To:   2000,
					Type: "range",
				},
				NumberOfMembersFilterQuery: NumberOfMembersFilterQuery{
					In:   []int{1},
					Type: "in",
				},
				Combinator: "and",
			},
			expected:   []string{"Bobby McFerrins", "Eminem"},
			statusCode: http.StatusOK,
		},
		{
			name: "Filter artists created after 2010 and first album after 2010",
			request: APIRequestData{
				CreationDateFilterQuery: NumberOfMembersFilterQuery{
					From: 2010,
					To:   9999,
					Type: "range",
				},
				FirstAlbumDateFilterQuery: FirstAlbumDateFilterQuery{
					From: "01-01-2010",
					To:   "01-01-9999",
					Type: "range",
				},
				Combinator: "and",
			},
			expected:   []string{"XXXTentacion", "Juice Wrld", "Alec Benjamin", "Post Malone"},
			statusCode: http.StatusOK,
		},
		{
			name: "Filter artists with concerts in Washington, USA and more than 3 members",
			request: APIRequestData{
				LocationsOfConcertsFilterQuery: LocationsOfConcertsFilterQuery{
					In: []string{"Washington, USA"},
				},
				NumberOfMembersFilterQuery: NumberOfMembersFilterQuery{
					From: 4,
					To:   9999,
					Type: "range",
				},
				Combinator: "and",
			},
			expected:   []string{"The Rolling Stones"},
			statusCode: http.StatusOK,
		},
		{
			name: "Filter artists by first albums between 1980 and 1990 with max 4 members",
			request: APIRequestData{
				FirstAlbumDateFilterQuery: FirstAlbumDateFilterQuery{
					From: "01-01-1980",
					To:   "31-12-1990",
					Type: "range",
				},
				NumberOfMembersFilterQuery: NumberOfMembersFilterQuery{
					From: 0,
					To:   4,
					Type: "range",
				},
				Combinator: "and",
			},
			expected:   []string{"Phil Collins", "Bobby McFerrins", "Red Hot Chili Peppers", "Metallica"},
			statusCode: http.StatusOK,
		},
		{
			name: "Method Not Allowed: Filter artists by first albums between 1980 and 1990 with max 4 members",
			request: APIRequestData{
				FirstAlbumDateFilterQuery: FirstAlbumDateFilterQuery{
					From: "01-01-1980",
					To:   "31-12-1990",
					Type: "range",
				},
				NumberOfMembersFilterQuery: NumberOfMembersFilterQuery{
					From: 0,
					To:   4,
					Type: "range",
				},
				Combinator: "and",
			},
			expected:   []string{},
			statusCode: http.StatusMethodNotAllowed,
			method:     1,
		},
	}

	handler := http.HandlerFunc(API)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal the request payload to JSON
			payload, err := json.Marshal(tc.request)
			if err != nil {
				t.Fatalf("Failed to marshal request payload: %v", err)
			}

			Method := http.MethodPost
			if tc.method == 1 {
				Method = http.MethodGet
			}

			// Create a new HTTP POST request with the JSON payload
			req, err := http.NewRequest(Method, "/api/filter", bytes.NewBuffer(payload))
			if err != nil {
				t.Fatalf("Failed to create HTTP request: %v", err)
			}

			// Set the content type header to JSON
			req.Header.Set("Content-Type", "application/json")

			// Record the response using httptest
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			// Assert the status code
			if rec.Code != tc.statusCode {
				t.Errorf("Expected status code %d, but got %d", tc.statusCode, rec.Code)
			}

			// Decode the JSON response
			var resp struct {
				Status  int `json:"status"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
			}

			err = json.Unmarshal(rec.Body.Bytes(), &resp)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			// Extract artist names from the response
			var artistNames []string
			for _, artist := range resp.Artists {
				artistNames = append(artistNames, artist.Name)
			}

			if !slicesEqual(tc.expected, artistNames) {
				t.Errorf("Expected artists %v, but got %v", tc.expected, artistNames)
			}
		})
	}
}

// slicesEqual checks if two slices of strings contain exactly the same elements,
// not necessarily in the same order
func slicesEqual(s1, s2 []string) bool {
	// If slices are of different lengths, they cannot be equal
	if len(s1) != len(s2) {
		return false
	}

	// Create a map to count occurrences of each string in the first slice
	counts := make(map[string]int)
	for _, str := range s1 {
		counts[str]++
	}

	// Decrement counts based on the second slice
	for _, str := range s2 {
		if _, exists := counts[str]; !exists {
			return false
		}
		counts[str]--
		if counts[str] < 0 {
			return false
		}
	}

	// Ensure all counts are zero
	for _, count := range counts {
		if count != 0 {
			return false
		}
	}

	return true
}
