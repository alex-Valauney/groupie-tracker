package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var Magenta = "\033[35m"
var Red = "\033[31m"
var Yellow = "\033[33m"

type Artists struct {
	Id            int      `json:"id"`
	Name          string   `json:"name"`
	Image         string   `json:"image"`
	CreationDates int      `json:"creationDate"`
	FirstAlbum    string   `json:"firstAlbum"`
	Members       []string `json:"members"`
}

type Locations struct {
	Id           int      `json:"id"`
	Localisation []string `json:"dates"`
}

type Dates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type RelationResp struct {
	Index []Relation `json:"index"`
}

type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func FetchJson(url string, i interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(i); err != nil {
		log.Fatal(err)
		fmt.Println(Red+"Failed to decode JSON", err)
		return err
	}

	return nil
}

const port = ":8080"

func main() {

	fmt.Println(Magenta + "Server successfully started")
	fmt.Println(Yellow + "To access the server --> (http://localhost:8080)")
	fmt.Println(Red + "To stop the server, simply press Ctrl+C")

	artist := []Artists{}
	err := FetchJson("https://groupietrackers.herokuapp.com/api/artists", &artist)
	if err != nil {
		fmt.Println(Red+"Failed to fetch artists:", err)
		return
	}

	var relationsResp RelationResp
	err = FetchJson("https://groupietrackers.herokuapp.com/api/relation", &relationsResp)
	if err != nil {
		fmt.Println(Red+"Failed to fetch concerts dates and locations:", err)
		return
	}

	relations := relationsResp.Index

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Index(w, r, artist)
	})
	mux.HandleFunc("/id/", func(w http.ResponseWriter, r *http.Request) {
		Details(w, r, artist, relations)
	})

	server := &http.Server{
		Addr:              port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println(Red + "Error starting server")
		log.Fatal(err)
	}
}
