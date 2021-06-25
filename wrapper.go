package main

import (
	"log"
	"net/http"
)

func GetOnly(h handler) handler {
	log.Print("GET ONLY")
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h(w, r)
			return
		}
		http.Error(w, "get only", http.StatusMethodNotAllowed)
	}
}
