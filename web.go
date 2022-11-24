package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

const portNumber = ":8080"

func artistHandler(w http.ResponseWriter, r *http.Request) {
	c := make(chan Error)
	go fetchData(c)
	e := <-c
	if e.Code != http.StatusOK {
		errorHandler(w, e.Code, e.Message)
		return
	}

	path := r.URL.Path
	if path != "/" {
		idString := path[1:]
		id, err := strconv.Atoi(idString)
		if err == nil {
			t, e := template.ParseFiles("templates/artist.html")
			if e == nil {
				for _, artist := range bands {
					if artist.ID == id {
						summary := createArtistSummary(artist)
						t.Execute(w, summary)
						return
					}
				}
				errorHandler(w, http.StatusNotFound, "404 NOT FOUND")
				return
			}

		}

	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		errorHandler(w, http.StatusNotFound, "404 NOT FOUND")
	}
	err = t.Execute(w, bands)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, "500 SERVER ERROR")
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	artistName := r.FormValue("artist")
	for _, artist := range bands {
		if artist.Name == artistName {
			http.Redirect(w, r, "/"+fmt.Sprint(artist.ID), http.StatusFound)
			return
		}
	}
	errorHandler(w, http.StatusNotFound, "404 NOT FOUND: "+artistName)
}

func createArtistSummary(artist Artist) ArtistSummary {
	summary := ArtistSummary{}
	summary.ID = artist.ID
	summary.Name = artist.Name
	summary.Image = artist.Image
	summary.CreationDate = artist.CreationDate
	summary.FirstAlbum = artist.FirstAlbum
	summary.Members = artist.Members

	if len(relations.Index) > 0 {
		for _, rel := range relations.Index {
			if artist.ID == rel.ID {
				concerts := []Concert{}
				for k, v := range rel.DatesLocations {
					constert := Concert{}
					constert.Location = formatLocations(k)
					constert.Dates = v
					concerts = append(concerts, constert)
				}
				summary.Concerts = concerts
			}
		}
		return summary
	}

	locationsSlice := []string{}
	datesSlice := []string{}

	for _, locat := range locations.Index {
		if locat.ID == artist.ID {
			locationsSlice = locat.Locations
			break
		}
	}
	for _, d := range dates.Index {
		if d.ID == artist.ID {
			datesSlice = d.Dates
			break
		}
	}

	datesString := strings.Join(datesSlice, " ")
	datesByStar := strings.Split(datesString, "*")[1:]

	concerts := []Concert{}

	for i, locat := range locationsSlice {
		concert := Concert{}
		locationFormatted := formatLocations(locat)
		concert.Location = locationFormatted
		if i < len(datesByStar) {
			concert.Dates = strings.Fields(datesByStar[i])
			concerts = append(concerts, concert)
		}
	}
	summary.Concerts = concerts

	return summary
}

func errorHandler(w http.ResponseWriter, code int, message string) {
	go w.WriteHeader(code)
	template, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, message, code)
		return
	}
	err = template.Execute(w, Error{Message: message, Code: code})
	if err != nil {
		http.Error(w, message, code)
	}
}

func fetchData(c chan Error) {
	data, err := getData(api)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	err = json.Unmarshal([]byte(data), &groupies)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}

	// Get bands
	artists, err := getData(groupies.Artists)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	err = json.Unmarshal([]byte(artists), &bands)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}

	// Get Relations
	isDataAvailable := true
	relat, err := getData(groupies.Relation)
	if err != nil {
		isDataAvailable = false
	}
	if isDataAvailable {
		err = json.Unmarshal([]byte(relat), &relations)
		if err != nil {
			isDataAvailable = false
		}
	}
	if len(relations.Index) == 0 {
		isDataAvailable = false
	}
	if isDataAvailable {
		c <- Error{Code: http.StatusOK, Message: ""}
		return
	}

	// Get Locations if not avalable
	locat, err := getData(groupies.Locations)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	err = json.Unmarshal([]byte(locat), &locations)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}

	// Get Dates
	d, err := getData(groupies.Dates)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	err = json.Unmarshal([]byte(d), &dates)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	c <- Error{Code: http.StatusOK, Message: ""}
}

func make500Error(message string) Error {
	return Error{Code: http.StatusInternalServerError, Message: "500 INTERNAL SERVER ERROR: " + message}
}
