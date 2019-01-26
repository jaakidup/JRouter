package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func (router *Router) listHandler(w http.ResponseWriter, r *http.Request, params Params) {
	w.Header().Set("Content-Type", "application/json")

	var results []interface{}

	results = append(results, router.Routes["GET"].ListKeys("GET"))
	results = append(results, router.Routes["POST"].ListKeys("POST"))

	json.NewEncoder(w).Encode(results)
}

func (router *Router) getPerson(w http.ResponseWriter, r *http.Request, params Params) {
	// w.WriteHeader(200)

	// send it through to the logic circuits

	var results []interface{}
	results = append(results, "Cool Beans!")
	results = append(results, params)

	json.NewEncoder(w).Encode(results)

	// io.WriteString(w, output)
}

func (router *Router) testGet(w http.ResponseWriter, r *http.Request, params Params) {

	// output := "TestGet function"
	// io.WriteString(w, output)

	reply := struct {
		Request string
		Params  Params
	}{
		Request: r.RequestURI,
		Params:  params,
	}

	json.NewEncoder(w).Encode(reply)

}

func (router *Router) testRemove(w http.ResponseWriter, r *http.Request, params Params) {

	router.Unregister("GET", "/get")

	io.WriteString(w, "Removed a live handler")
}

func defaultHeaders(h Handle) Handle {
	return func(w http.ResponseWriter, r *http.Request, p Params) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")
		h(w, r, p)
	}
}
