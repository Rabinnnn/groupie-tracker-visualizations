package location

import (
	"regexp"
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
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	locationsInA := removeBlank(re.Split(a, -1))
	locationsInB := removeBlank(re.Split(b, -1))

	for _, locA := range locationsInA {
		for _, locB := range locationsInB {
			locA = strings.ToLower(locA)
			locB = strings.ToLower(locB)
			if strings.Contains(locA, locB) || strings.Contains(locB, locA) {
				return true
			}
		}
	}

	return false
}

func removeBlank(a []string) []string {
	for i, b := range a {
		if strings.TrimSpace(b) == "" {
			a = append(a[:i], a[i+1:]...)
		}
	}

	return a
}
