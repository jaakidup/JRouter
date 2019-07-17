package jaakitrouter

// ================================================================== //
//
// 	JRouter
//
//	HTTP Router for REST Apis
//
//
//	@Jaakit @Jaakidup
//
// ================================================================== //

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// )

// var router *Router

// func main() {

// 	router = &Router{DebugLog: true}

// 	router.GET("/person/@something/@spectacular", router.LogWrapper(router.getUsers))
// 	router.GET("/list", router.listHandler)
// 	router.GET("/", router.index)

// 	router.Logger("Starting Server on :8081")
// 	log.Fatal(http.ListenAndServe(":8081", router))

// }

// func (router *Router) index(w http.ResponseWriter, r *http.Request, _ NamedParams) {
// 	w.Write([]byte("<div align='center'>Index page, use GET /users</div>"))
// }

// func (router *Router) getUsers(w http.ResponseWriter, r *http.Request, params NamedParams) {
// 	w.Header().Set("content-type", "application/json")
// 	json.NewEncoder(w).Encode(params)
// }

// func (router *Router) listHandler(w http.ResponseWriter, r *http.Request, params NamedParams) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var results []interface{}

// 	results = append(results, router.Routes["GET"].ListKeys("GET"))
// 	results = append(results, router.Routes["POST"].ListKeys("POST"))

// 	json.NewEncoder(w).Encode(results)
// }

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"time"
)

// Handle ...
// type Handle func(w http.ResponseWriter, r *http.Request, params Params)
type Handle func(w http.ResponseWriter, r *http.Request, params NamedParams)

// Params ...
type Params map[int]string

// NamedParams ...
type NamedParams map[string]string

// Router ...
type Router struct {
	AutoOptions     bool
	Routes          map[string]*DigitalTree
	DebugLog        bool
	NotFoundHandler http.Handler
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

	pathobject := router.disectPath(path)

	// TODO: add check for OPTIONS handlers, otherwise ahndle OPTION if set to true
	// if method == "OPTIONS" {
	// 	router.OptionsHandler(w, r)
	// }

	found, jfunc, params, _ := router.Routes[method].Find(pathobject[1])
	if found {
		namedParams := make(map[string]string)
		for index, key := range params {
			namedParams[key] = pathobject[index+1]
		}
		jfunc(w, r, namedParams)

	} else {
		fmt.Println("Didn't find route")
		w.WriteHeader(404)
	}
}

func (router *Router) disectPath(path string) Params {

	sections := strings.Split(path, "/")
	pathlen := len(sections)
	if pathlen == 1 || pathlen == 0 || len(sections[1]) == 0 {
		return nil
	}

	params := make(map[int]string)

	for i := 0; i < len(sections); i++ {
		params[i] = sections[i]
	}

	return params
}

func (router *Router) pathObjectForRegister(path string) Params {
	paramkey := "@"

	sections := strings.Split(path, "/")

	fmt.Println("register", sections)

	params := make(map[int]string)

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

// Register ...
func (router *Router) Register(method string, path string, handle Handle) {

	pathobject := router.pathObjectForRegister(path)

	if router.Routes == nil {
		router.Routes = make(map[string]*DigitalTree)
	}

	if router.Routes[method] == nil {
		router.Routes[method] = NewDigitalTree()
	}

	router.Routes[method].Add(pathobject[0], handle, pathobject)
}

// Unregister ...
func (router *Router) Unregister(method string, path string) {
	router.Routes[method].Delete(path)
}

// Logger logs the message[s] on a single line if debug:true
func (router *Router) Logger(message interface{}) {
	if router.DebugLog {
		log.Println(message)
	}
}

// LogWrapper is wrapped around the handler to send logs to admin console
func (router *Router) LogWrapper(h Handle) Handle {
	return func(w http.ResponseWriter, r *http.Request, p NamedParams) {
		start := time.Now()
		h(w, r, p)
		duration := time.Since(start)
		router.Logger(r.Method + " " + r.RemoteAddr + " " + duration.String() + " " + r.RequestURI)
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
