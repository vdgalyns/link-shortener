package entities

import "strings"

func ValidateURL(originalURL string) bool {
	_, domain, _ := strings.Cut(originalURL, "//")
	if len(domain) == 0 {
		return false
	}
	name, zone, _ := strings.Cut(domain, ".")
	if len(name) == 0 || len(zone) < 2 {
		return false
	}
	return true
}
