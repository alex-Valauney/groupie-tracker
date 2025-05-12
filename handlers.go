package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func CleanRelation(unclean Relation) Relation {
	cleaned := Relation{
		Id:             unclean.Id,
		DatesLocations: make(map[string][]string),
	}

	for location, dates := range unclean.DatesLocations {
		cleanedLocation := strings.ReplaceAll(location, "-", " - ")
		cleanedLocation = strings.ReplaceAll(cleanedLocation, "_", " ")
		cleaned.DatesLocations[cleanedLocation] = dates
	}
	return cleaned
}

func Index(w http.ResponseWriter, r *http.Request, artists []Artists) {
	//Find the template file
	t, err := template.ParseFiles("templates/Index.html")
	//Error if the template is not found
	if err != nil {
		Error(w, r, http.StatusNotFound, "Could not find template.")
		return
	}
	//Error 404 if the user is not on the index page which is the only page accessible
	if r.URL.Path != "/" {
		Error(w, r, http.StatusNotFound, "Page not found")

		return
	}

	err = t.Execute(w, artists)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, "Server internal error. Please try again later")
	}
}

func Error(a http.ResponseWriter, b *http.Request, c int, d string) {
	//Structure to display errors
	type StrucErreur struct {
		CodeErr    int
		MessageErr string
	}
	Erreur := StrucErreur{c, d}
	//Loads the error page if an error happens
	t, err := template.ParseFiles("templates/Error.html")
	if err != nil {
		http.Error(a, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(a, Erreur)
}

func Details(w http.ResponseWriter, r *http.Request, artists []Artists, relations []Relation) {
	//Find the template file
	idStr := r.URL.Path[len("/id/"):] // Get the part after "/id/"
	if idStr[len(idStr)-1:] == "/" {
		idStr = idStr[:len(idStr)-1]
	}

	artistId, err := strconv.Atoi(idStr) // Convert it to an integer
	if err != nil || artistId < 1 || artistId > len(artists) {
		Error(w, r, http.StatusNotFound, "Invalid artist ID")
		return
	}

	// Find the artist by ID (artistId is 1-based, slice index is 0-based)
	artist := artists[artistId-1]
	relation := relations[artistId-1]
	cleanedRelation := CleanRelation(relation)

	// Parse the template for the details page
	t, err := template.ParseFiles("templates/Details.html")
	if err != nil {
		log.Fatal(err)
		Error(w, r, http.StatusNotFound, "Could not find template.")
		return
	}

	data := struct {
		Artist    Artists
		Relations Relation
	}{
		Artist:    artist,
		Relations: cleanedRelation,
	}

	// Render the details page with the artist data
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
		Error(w, r, http.StatusInternalServerError, "Failed to render template")
	}
}
