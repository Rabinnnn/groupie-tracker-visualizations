package location

import (
	"regexp"
	"slices"
	"strings"
)

type location = string

// Parse the given `hyphenated location` and return the extracted city and country locations.
//
// For example, given the `hyphenated location` (`philadelphia-usa`) returns the city (philadelphia) and country (usa)
func Parse(s string) (city, country string) {
	data := strings.Split(s, "-")

	process := func(s string) string {
		return strings.ReplaceAll(s, "_", " ")
	}

	if len(data) >= 2 {
		country = process(data[1])
	}

	if len(data) >= 1 {
		city = process(data[0])
	}

	return
}

// Contains checks if a given location is contained in the other location.
// For example Seattle, Washington, USA is part of Washington, USA
func Contains(a, b location) bool {
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	return strings.Contains(a, b) || strings.Contains(b, a)
}

// removeBlank returns a new slice from a with all element strings that are blank removed
func removeBlank(a []string) (noBlanks []string) {
	for _, current := range a {
		if strings.TrimSpace(current) != "" {
			noBlanks = append(noBlanks, current)
		}
	}
	return
}

// GetCityCountry parses the given location, and extracts the name of the country and the target city of the location,
// trimming any space at the end and start of the names.
//
// Note: This assumes the name of the city is always the first then the second is the name of the country.
func GetCityCountry(s location) (city, country string) {
	s = rev(s)
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	locations := removeBlank(re.Split(s, 2))

	slices.Reverse(locations)

	if len(locations) > 0 {
		city = locations[0]
	}

	if len(locations) > 1 {
		country = locations[1]
	}

	city = strings.TrimSpace(city)
	country = strings.TrimSpace(country)

	return rev(city), rev(country)
}

// rev returns the reverse of the string s
func rev(s string) string {
	chars := strings.Split(s, "")
	slices.Reverse(chars)
	return strings.Join(chars, "")
}
