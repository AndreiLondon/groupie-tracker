package main

import "strings"

func formatLocations(location string) string {
	location = strings.Replace(location, "-", ", ", -1)
	location = strings.Replace(location, "_", " ", -1)
	location = strings.Title(location)
	location = strings.Replace(location, ", Usa", ", USA", -1)
	location = strings.Replace(location, ", Uk", ", UK", -1)
	return location
}
