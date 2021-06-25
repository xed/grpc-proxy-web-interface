package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

type handler func(w http.ResponseWriter, r *http.Request)

func BasicAuth(pass handler) handler {
	log.Print("ENABLE BASIC AUTH")
	return func(w http.ResponseWriter, r *http.Request) {
		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !validate(pair[0], pair[1]) {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		pass(w, r)
	}
}

func validate(username, password string) bool {
	if username == getLogin() && password == getPass() {
		return true
	}
	return false
}
