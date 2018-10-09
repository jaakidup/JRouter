package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func (router *Router) listHandler(w http.ResponseWriter, r *http.Request) {
	jenc := json.NewEncoder(w)

	var results []interface{}

	results = append(results, router.Routes["GET"].ListKeys("GET"))
	results = append(results, router.Routes["POST"].ListKeys("POST"))
	// results = append(results, router.Routes["OPTIONS"].ListKeys("OPTIONS"))

	jenc.Encode(results)
}

func (router *Router) testPost(w http.ResponseWriter, r *http.Request) {
	output := "TestPost function"
	io.WriteString(w, output)
}
func (router *Router) testGet(w http.ResponseWriter, r *http.Request) {
	output := "TestGet function"
	io.WriteString(w, output)
}

func (router *Router) testRemove(w http.ResponseWriter, r *http.Request) {

	router.Unregister("GET", "/get")

	io.WriteString(w, "Removed a live handler")
}
