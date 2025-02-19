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
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	locationsInA := removeBlank(re.Split(a, -1))
	locationsInB := removeBlank(re.Split(b, -1))

	slices.Reverse(locationsInA)
	slices.Reverse(locationsInB)

	la := len(locationsInA)
	lb := len(locationsInB)

	if la > 0 && lb > 0 {
		if strings.Contains(locationsInA[0], locationsInB[0]) || strings.Contains(locationsInB[0], locationsInA[0]) {
			return true
		}
	}

	if la > 1 && lb > 1 {
		if strings.Contains(locationsInA[1], locationsInB[1]) || strings.Contains(locationsInB[1], locationsInA[1]) {
			return true
		}

		//if strings.Contains(locationsInA[0], locationsInB[1]) || strings.Contains(locationsInB[0], locationsInA[1]) {
		//	return true
		//}
	}

	//for _, locA := range locationsInA {
	//	for _, locB := range locationsInB {
	//		locA = strings.ToLower(locA)
	//		locB = strings.ToLower(locB)
	//		if strings.Contains(locA, locB) || strings.Contains(locB, locA) {
	//			return true
	//		}
	//	}
	//}

	return false
}

// Contains2 checks if a given location is contained in the other location.
// For example Seattle, Washington, USA is part of Washington, USA
func Contains2(a, b location) bool {
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	return strings.Contains(a, b) || strings.Contains(b, a)
}

func removeBlank(a []string) []string {
	for i, b := range a {
		if strings.TrimSpace(b) == "" {
			a = append(a[:i], a[i+1:]...)
		}
	}

	return a
}

func getCityCountry(s string) (city, country string) {
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	locations := removeBlank(re.Split(s, -1))

	slices.Reverse(locations)

	if len(locations) > 0 {
		country = locations[0]
	}

	if len(locations) > 1 {
		city = locations[1]
	}

	return
}
