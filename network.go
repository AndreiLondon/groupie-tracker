package main

import (
	"io"
	"net/http"
)

var api = "https://groupietrackers.herokuapp.com/api"

func getData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	//Converting an HTTP response body to a string
	// We can convert []byte to a string here
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
