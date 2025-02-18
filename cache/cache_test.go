package cache

import (
	"groupie-tracker/api"
	"reflect"
	"slices"
	"testing"
)

func TestGetCachedDataIntegration(t *testing.T) {
	// Invalidate all cache
	artistCache = nil
	locationCache = nil
	locationMapCache = nil
	dateCache = nil
	relationCache = nil
	isCacheInitialized = false

	artists, locations, _, _, _ := GetCachedData()
	ValidateArtistsData(artists, t)

	// Validate the integrity of the locations data
	{
		if locations == nil {
			t.Fatalf("GetCachedData() locations is nil")
		}

		containsLocation := func(id int, locations []string) func(api.Location) bool {
			return func(loc api.Location) bool {
				return loc.Id == id && reflect.DeepEqual(loc.Locations, locations)
			}
		}

		tests := []struct {
			Name      string
			Id        int
			Locations []string
		}{
			{
				Name: "First",
				Id:   1,
				Locations: []string{
					"north_carolina-usa",
					"georgia-usa",
					"los_angeles-usa",
					"saitama-japan",
					"osaka-japan",
					"nagoya-japan",
					"penrose-new_zealand",
					"dunedin-new_zealand",
				},
			},

			{
				Name: "Random",
				Id:   48,
				Locations: []string{
					"texas-usa",
					"oklahoma-usa",
					"california-usa",
					"illinois-usa",
					"scheessel-germany",
					"st_gallen-switzerland",
					"gdynia-poland",
					"arras-france",
				},
			},

			{
				Name: "Last",
				Id:   52,
				Locations: []string{
					"oregon-usa",
					"vancouver-canada",
					"nevada-usa",
					"colorado-usa",
					"munich-germany",
					"prague-czechia",
					"milan-italy",
				},
			},
		}
		for _, tt := range tests {
			t.Run(
				tt.Name, func(t *testing.T) {
					if !slices.ContainsFunc(locations, containsLocation(tt.Id, tt.Locations)) {
						t.Fatalf("expected to find location with id %d; but none matched", tt.Id)
					}
				},
			)
		}
	}

	// Fetch new data before cache expiration
	artistsNew, _, _, _, _ := GetCachedData()

	// We expect the new data to also pass the artists data validation
	ValidateArtistsData(artistsNew, t)
}

// ValidateArtistsData validates the integrity of the given artists data
func ValidateArtistsData(artists []api.Artist, t *testing.T) {
	if artists == nil {
		t.Fatalf("GetCachedData() artists is nil")
	}

	containsArtist := func(id int, name string) func(artist api.Artist) bool {
		return func(artist api.Artist) bool {
			return artist.ID == id && artist.Name == name
		}
	}

	tests := []struct {
		id   int
		name string
	}{
		{
			id:   1,
			name: "Queen",
		},

		{
			id:   52,
			name: "The Chainsmokers",
		},

		{
			id:   49,
			name: "The Rolling Stones",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if !slices.ContainsFunc(artists, containsArtist(tt.id, tt.name)) {
					t.Fatalf("expected to find artist %s with id %d; but none matched", tt.name, tt.id)
				}
			},
		)
	}

}
