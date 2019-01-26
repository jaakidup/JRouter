package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// // Handle ...
// type Handle func(w http.ResponseWriter, r *http.Request)

// Handle ...
type Handle func(w http.ResponseWriter, r *http.Request, params Params)

// Params ...
type Params map[int]string

// Router ...
type Router struct {
	AutoOptions     bool
	Routes          map[string]*DigitalTree
	DebugLog        bool
	NotFoundHandler http.Handler
	AdminHandler    http.HandlerFunc
	AdminConnected  bool
	AdminConnection *websocket.Conn
}

// New ...
func New() *Router {
	return &Router{
		AutoOptions: true,
		Routes:      make(map[string]*DigitalTree),
		DebugLog:    false,
	}
}

// ErrorHandler is the default error handler
func (router *Router) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(420)
	router.Logger("Enhance your calm")
}

// OptionsHandler will handle OPTIONS request if no other OPTIONS handler is declared
func (router *Router) OptionsHandler(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// ServerHTTP ...
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	method := html.EscapeString(r.Method)
	path := html.EscapeString(r.URL.Path)

	if method == "OPTIONS" {
		router.OptionsHandler(w, r)
		goto SKIPFOROPTIONS
	}

	if path == "/" {
		router.NotFoundHandler.ServeHTTP(w, r)
	} else if path == "/admin" {
		router.AdminHandler(w, r)
	} else {

		pathobject := router.disectPath(path)

		found, jfunc := router.Routes[method].Find(pathobject[0])
		if found {
			jfunc(w, r, pathobject)
		} else {
			fmt.Println("Didn't find route")
			w.WriteHeader(404)
		}
	}

SKIPFOROPTIONS:
}

// ==================================================== //
func (*Router) splitPath(path string) []string {
	fmt.Println("Splitter", path, strings.Split(path, "/"))
	return strings.Split(path, "/")
}

func (router *Router) disectPath(path string) Params {
	paramkey := "@"

	sections := router.splitPath(path)
	params := make(map[int]string)

	// if strings.ContainsAny(path, paramkey) {
	// fmt.Println("Path contains parameter, strip them out and store them in the DigitalTree")
	for i, entry := range sections {
		index := i - 1
		if index == -1 {
			continue
		}
		if strings.Contains(entry, paramkey) {
			params[index] = strings.Split(entry, paramkey)[1]
		} else {
			params[index] = entry
		}
	}
	return params
}

// ==================================================== //

// Register ...
func (router *Router) Register(method string, path string, handle Handle) {

	pathobject := router.disectPath(path)
	fmt.Println(len(pathobject), pathobject)

	// TODO: register the

	if router.Routes == nil {
		router.Routes = make(map[string]*DigitalTree)
	}

	if router.Routes[method] == nil {
		router.Routes[method] = NewDigitalTree()
	}

	router.Logger("Registering: " + method + " " + pathobject[0])
	router.Routes[method].Add(pathobject[0], handle)
}

// Unregister ...
func (router *Router) Unregister(method string, path string) {
	router.Routes[method].Delete(path)
}

// Logger logs the message[s] on a single line if debug:true
func (router *Router) Logger(message interface{}) {
	if router.DebugLog {
		log.Println(message)
		router.WriteToAdminConsole(message)
	}
}

// LogWrapper is wrapped around the handler to send logs to admin console
func (router *Router) LogWrapper(h Handle) Handle {
	return func(w http.ResponseWriter, r *http.Request, p Params) {
		router.WriteToAdminConsole("Received request[" + r.Method + "][" + r.RemoteAddr + "][" + r.RequestURI + "]")

		var report interface{}

		json.NewDecoder(r.Body).Decode(&report)

		router.WriteToAdminConsole(report)

		h(w, r, p)
	}
}

// POST ...
func (router *Router) POST(path string, handle Handle) {
	router.Register("POST", path, handle)
}

// DELETE ...
func (router *Router) DELETE(path string, handle Handle) {
	router.Register("DELETE", path, handle)
}

// PUT ...
func (router *Router) PUT(path string, handle Handle) {
	router.Register("PUT", path, handle)
}

// PATCH ...
func (router *Router) PATCH(path string, handle Handle) {
	router.Register("PATCH", path, handle)
}

// HEAD ...
func (router *Router) HEAD(path string, handle Handle) {
	router.Register("HEAD", path, handle)
}

// GET ...
func (router *Router) GET(path string, handle Handle) {
	router.Register("GET", path, handle)
}

// OPTIONS ...
func (router *Router) OPTIONS(path string, handle Handle) {
	router.Register("OPTIONS", path, handle)
}
