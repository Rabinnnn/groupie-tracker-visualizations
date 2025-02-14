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

func Contains(a, b location) bool {
	re := regexp.MustCompile(`(\s|-|_|,|;|:|.)+`)
	locationsInB := re.Split(b, -1)

	for _, locInB := range locationsInB {
		var larger, smaller location
		if len(a) > len(locInB) {
			larger = a
			smaller = locInB
		} else {
			larger = locInB
			smaller = a
		}

		if smaller != "" && (strings.Contains(strings.ToLower(larger), strings.ToLower(smaller))) {
			return true
		}
	}

	return false
}
