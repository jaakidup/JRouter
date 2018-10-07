package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
)

// Handle ...
type Handle func(w http.ResponseWriter, r *http.Request)

// Router ...
type Router struct {
	AutoOptions bool
	Routes      *DigitalTree
	DebugLog    bool
}

// New ...
func New() *Router {
	return &Router{
		AutoOptions: true,
		Routes:      NewDigitalTree(),
		DebugLog:    false,
	}
}

func (router *Router) Error(w http.ResponseWriter, r *http.Request) {
	if router.DebugLog {
		io.WriteString(w, "Whoopsie!!!!")
		log.Println("Oooops.")
	}
}

// ServerHTTP ...
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if router.DebugLog {
		io.WriteString(w, html.EscapeString(r.Method))
		io.WriteString(w, html.EscapeString(r.URL.Path))
	}

	method := html.EscapeString(r.Method)
	path := html.EscapeString(r.URL.Path)

	if method == "GET" {
		if router.DebugLog {
			log.Println("Received GET request on path, ", path)
		}
		found, jfunc := router.Routes.Find(path)
		if found {
			jfunc(w, r)
		} else {
			router.Error(w, r)
		}

	} else if method == "POST" {
		if router.DebugLog {
			log.Println("Received Post request on path, ", path)
		}
		found, jfunc := router.Routes.Find(path)
		if found {
			jfunc(w, r)
		} else {
			router.Error(w, r)
		}

	} else if method == "OPTIONS" {
		if router.DebugLog {
			log.Println("Received OPTIONS request on path, ", path)
		}
		// TODO: handle OPTIONS automatically
	}

}

// Register ...
func (router *Router) Register(method string, path string, handle Handle) {
	if router.DebugLog {
		fmt.Println("Register: ", method, path)
		fmt.Printf("function to register is type: %T\n", handle)
	}
	if router.Routes == nil {
		if router.DebugLog {
			fmt.Println("Router tree is nil!!")
		}
		router.Routes = NewDigitalTree()
	}

	router.Routes.Add("/tester", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	})

	router.Routes.Add(path, handle)

}

// POST ...
func (router *Router) POST(path string, handle Handle) {
	router.Register("POST", path, handle)
}

// GET ...
func (router *Router) GET(path string, handle Handle) {
	router.Register("GET", path, handle)
}