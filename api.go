package main

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// General API definitions.
// Independent from the API version.

const API_PREFIX = "/api"

type Route struct {
	Path        string           // URI endpoint
	HandlerFunc http.HandlerFunc // Handler function
	Methods     []string         // Allowed HTTP methods list
}

type API struct {
	Routes []Route
}

type APIErrorResponse struct {
	Response string   `json:"response"`
	Errors   []string `json:"errors"`
}

// Basic Auth decorate.
// Decorate API route handler to protect endpoint with token.
func (o *Orderup) BasicAuth(pass http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If Authorization and password aren't set - allow everything.
		if r.Header["Authorization"] == nil {
			if o.validate("", "") {
				pass(w, r)
				return
			} else {
				http.Error(w, "Authorization failed", http.StatusUnauthorized)
				return
			}
		}

		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "bad syntax", http.StatusBadRequest)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !o.validate(pair[0], pair[1]) {
			http.Error(w, "Authorization failed", http.StatusUnauthorized)
			return
		}

		pass(w, r)
	}
}

func (o *Orderup) validate(username, password string) bool {
	// If password is not set allow everything.
	if o.password == "" {
		return true
	}

	if password == o.password {
		return true
	}

	return false
}
