package main

// JROUTER usage

// package main

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// )

// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	output := "Index Route     / "
// 	io.WriteString(w, output)
// }

// func testPost(w http.ResponseWriter, r *http.Request) {
// 	output := "Test function"
// 	io.WriteString(w, output)
// }

// func test(w http.ResponseWriter, r *http.Request) {
// 	output := "Test function"
// 	io.WriteString(w, output)
// }

// func main() {
// 	sigs := make(chan os.Signal, 1)
// 	done := make(chan bool, 1)
// 	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

// 	router := &Router{DebugLog: true}

// 	router.POST("/post", testPost)
// 	router.GET("/get", test)
// 	router.GET("/", indexHandler)

// 	go func() {

// 		sig := <-sigs
// 		fmt.Println()
// 		fmt.Println(sig)

// 		done <- true
// 	}()
// 	go func() {
// 		log.Fatal(http.ListenAndServe(":8080", router))
// 	}()

// 	fmt.Println("Awaiting Signal")
// 	<-done
// 	fmt.Println("Shutdown procedure")

// }

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
	// Routes      *DigitalTree
	Routes   map[string]*DigitalTree
	DebugLog bool
	NotFound http.Handler
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
	if router.DebugLog {
		io.WriteString(w, "Whoopsie!!!!")
	}
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
	router.Logger("Serving", r.RemoteAddr, " on ", method, path)

	if path == "/" {
		router.NotFound.ServeHTTP(w, r)
	} else {
		found, jfunc := router.Routes[method].Find(path)
		if found {
			jfunc(w, r)
		} else {
			if router.NotFound == nil {
				router.ErrorHandler(w, r)
			} else {
				router.NotFound.ServeHTTP(w, r)
			}
		}
	}

}

// Register ...
func (router *Router) Register(method string, path string, handle Handle) {
	// if router.Routes == nil {
	// 	router.Logger("Router tree is nil, let's create")
	// 	router.Routes = NewDigitalTree()
	// }

	router.Logger("Let's check router.Routes")

	if router.Routes == nil {
		router.Routes = make(map[string]*DigitalTree)
	}

	if router.Routes[method] == nil {
		router.Logger("Routes for Method", method, "is nil, let's create")
		router.Routes[method] = NewDigitalTree()
	}

	fmt.Println(router.Routes)

	// router.Routes[method] = getRoutes
	// router.Routes["GET"] = NewDigitalTree()
	// fmt.Println(router.Routes)
	// router.Logger("Create router.Routes[POST]")
	// router.Routes["POST"] = NewDigitalTree()

	router.Logger("Registering: ", method, path)
	router.Routes[method].Add(path, handle)
}

// Logger logs the message[s] on a single line if debug:true
func (router *Router) Logger(message ...interface{}) {
	if router.DebugLog {
		log.Println(message...)
	}
}

// POST ...
func (router *Router) POST(path string, handle Handle) {
	router.Register("POST", path, handle)
}

// GET ...
func (router *Router) GET(path string, handle Handle) {
	router.Register("GET", path, handle)
}

// OPTIONS ...
func (router *Router) OPTIONS(path string, handle Handle) {
	router.Register("OPTIONS", path, handle)
}
