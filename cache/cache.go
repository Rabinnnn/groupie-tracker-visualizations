// Package cache contains utilities to ensure we have a local copy of the Groupie Trackers API data,
// updating the data from the API when necessary
package cache

import (
	"fmt"
	"groupie-tracker/api"
	"sync"
	"sync/atomic"
	"time"
)

var (
	artistCache      []api.Artist
	locationCache    []api.Location
	locationMapCache map[int][]string
	dateCache        []api.Date
	relationCache    []api.Relations
	// cacheTime keeps track of when last the offline cache was updated with online content
	cacheTime          time.Time
	cacheMutex         sync.RWMutex
	isCacheInitialized bool
)

// cacheDuration how long the application will work with offline
// data before getting new data from the external API
const cacheDuration = 2 * time.Hour

// GetCachedData fetches artists data. If available locally, and the cache is still valid,
// the data is returned immediately, else, network request is made to the Groupie Trackers API to get the latest data
func GetCachedData() ([]api.Artist, []api.Location, []api.Date, []api.Relations, error) {
	err := updateCache()
	return artistCache, locationCache, dateCache, relationCache, err
}

// GetCachedLocationsMap returns the map cached locations data
func GetCachedLocationsMap() map[int][]string {
	return locationMapCache
}

func updateCache() error {
	if isCacheInitialized && time.Since(cacheTime) < cacheDuration {
		return nil
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// keep track of how many errors have
	//been encountered by the go routines
	var errCount atomic.Int32
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		artists, err := api.GetArtists()
		if err != nil {
			errCount.Add(1)
			return
		}
		artistCache = artists
	}()

	go func() {
		defer wg.Done()
		locations, err := api.GetAllLocations()
		if err != nil {
			errCount.Add(1)
			return
		}
		locationCache = locations
	}()

	go func() {
		defer wg.Done()
		dates, err := api.GetAllDates()
		if err != nil {
			errCount.Add(1)
			return
		}
		dateCache = dates
	}()

	go func() {
		defer wg.Done()
		relations, err := api.GetAllRelations()
		if err != nil {
			errCount.Add(1)
			return
		}
		relationCache = relations
	}()

	wg.Wait()
	if errCount.Load() != 0 {
		return fmt.Errorf("failed to fetch data from the Groupie Trackers API")
	}

	cacheTime = time.Now()
	isCacheInitialized = true

	// map the locations so that the keys are the id's of the artists
	locationMapCache = make(map[int][]string)
	for _, loc := range locationCache {
		locationMapCache[loc.Id] = loc.Locations
	}

	return nil
}
