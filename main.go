package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	datasource := datasource{}
	badge, err := datasource.GetBadge("jx")
	if err != nil {
		log.Fatalf("error getting badge %s", err)
	}

	log.Printf("Got Badge %s", badge)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(badge)
}

func main() {
	log.Printf("Starting...")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
