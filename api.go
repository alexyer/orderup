package main

import "net/http"

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
