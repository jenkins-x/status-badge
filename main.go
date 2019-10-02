package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	datasource := datasource{}
	vars := mux.Vars(r)
	repo := vars["repo"]
	org := vars["org"]

	badge, err := datasource.GetBadge(org, repo)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error getting badge %s", err)
		return
	}

	if badge == nil {
		log.Printf("Could not find badge for org '%s' and repo '%s'", org, repo)

		w.WriteHeader(http.StatusNotFound)
	} else {
		log.Printf("Got Badge '%s' for org '%s' and repo '%s'", badge, org, repo)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(badge)
	}
}

func main() {
	log.Printf("Starting...")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{repo}", handler)
	router.HandleFunc("/{org}/{repo}", handler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
