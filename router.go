package main

import (
	"encoding/json"
	"html"
	"net/http"

	"github.com/gorilla/websocket"
)

// Handle ...
type Handle func(w http.ResponseWriter, r *http.Request)

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
	router.Logger("Some error, still needs to be decided on.")
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
		found, jfunc := router.Routes[method].Find(path)
		if found {
			jfunc(w, r)
		} else {
			if router.NotFoundHandler == nil {
				router.ErrorHandler(w, r)
			} else {
				router.NotFoundHandler.ServeHTTP(w, r)
			}
		}
	}
SKIPFOROPTIONS:
}

// Register ...
func (router *Router) Register(method string, path string, handle Handle) {

	if router.Routes == nil {
		router.Routes = make(map[string]*DigitalTree)
	}

	if router.Routes[method] == nil {
		router.Logger("Routes for Method: " + method + "is nil, let's create")
		router.Routes[method] = NewDigitalTree()
	}

	router.Logger("Registering: " + method + path)
	router.Routes[method].Add(path, handle)
}

// Unregister ...
func (router *Router) Unregister(method string, path string) {
	router.Routes[method].Delete(path)
}

// Logger logs the message[s] on a single line if debug:true
// TODO: send logs over channel to admin section
func (router *Router) Logger(message interface{}) {
	if router.DebugLog {
		// log.Println(message...)
		router.WriteToAdminConsole(message)
	}
}

// LogWrapper is wrapped around the handler to send logs to admin console
func (router *Router) LogWrapper(h Handle) Handle {
	return func(w http.ResponseWriter, r *http.Request) {
		router.WriteToAdminConsole("Received request[" + r.Method + "][" + r.RemoteAddr + "][" + r.RequestURI + "]")

		var report interface{}

		json.NewDecoder(r.Body).Decode(&report)

		router.WriteToAdminConsole(report)

		h(w, r)
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

// UPDATE ...
func (router *Router) UPDATE(path string, handle Handle) {
	router.Register("UPDATE", path, handle)
}

// GET ...
func (router *Router) GET(path string, handle Handle) {
	router.Register("GET", path, handle)
}

// OPTIONS ...
func (router *Router) OPTIONS(path string, handle Handle) {
	router.Register("OPTIONS", path, handle)
}
