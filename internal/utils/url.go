package utils

import "strings"

func ValidateURL(url string) bool {
	_, domain, _ := strings.Cut(url, "//")
	if len(domain) == 0 {
		return false
	}
	name, zone, _ := strings.Cut(domain, ".")
	if len(name) == 0 || len(zone) < 2 {
		return false
	}
	return true
}
