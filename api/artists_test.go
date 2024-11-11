package api

import "testing"

// TestGetArtists_Success tests that GetArtists returns a valid list of artists
// and does not return an error when the API is available and correctly formatted.
func TestGetArtists_Success(t *testing.T) {
	artists, err := GetArtists()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if len(artists) == 0 {
		t.Fatal("Expected non-zero number of artists, but got 0")
	}
}
