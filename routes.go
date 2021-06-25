package main

import (
	"io"
	"log"
	"net/http"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	log.Print("Start handler")
	io.WriteString(w, generateHTML())
}
