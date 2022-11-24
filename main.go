package main

import (
	// import our encoding/json package

	"fmt"
	"net/http"
)

var (
	groupies  Groupies
	bands     []Artist
	locations Locations
	dates     Dates
	relations Relations
)

func main() {
	/* Now, we can parse the response body into the "Groupies" struct.
	First, we need to read the response body before we can do a mapping to the struct.
	The result is a byte slice if there's no error.
	After that, we just need to unmarshal from the slice of bytes to an object with the type "Groupies".
	*/

	// Routing

	http.HandleFunc("/", artistHandler)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/search", searchHandler)
	// Web server that listens for requests. (Without it, the main func never executes)
	fmt.Println("Starting the server on " + portNumber)
	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		fmt.Println("\nCannot start server")
	}
}
