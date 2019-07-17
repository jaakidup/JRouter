
# jaakitrouter

  

Rest api router with GET, POST, DELETE, OPTIONS, PUT and PATCH wrappers

Currently supports parameters in the form of http://localhost:8081/person/@id/@other_param/@etc/

  
  

  

    import (
	    "encoding/json"
	    "log"
	    "net/http"
    )	
    
    var router *Router

    func main() {
	    router = &Router{DebugLog: true}

	    router.GET("/person/@something/@spectacular", router.LogWrapper(router.getUsers))	    
	    router.GET("/list", router.listHandler)
	    router.GET("/", router.index)
	    
	    router.Logger("Starting Server on :8081")
	    log.Fatal(http.ListenAndServe(":8081", router))
    }
    
      
    
    func (router *Router) index(w http.ResponseWriter, r *http.Request, _ NamedParams) {
	    w.Write([]byte("<div align='center'>Index page, use GET /users</div>"))    
    }
    
      
    
    func (router *Router) getUsers(w http.ResponseWriter, r *http.Request, params NamedParams) {
	    w.Header().Set("content-type", "application/json")
	    json.NewEncoder(w).Encode(params)
    }
    
      
    
    func (router *Router) listHandler(w http.ResponseWriter, r *http.Request, params NamedParams) {
	    w.Header().Set("Content-Type", "application/json")
	    var results []interface{}
	    results = append(results, router.Routes["GET"].ListKeys("GET"))
	    results = append(results, router.Routes["POST"].ListKeys("POST"))
	    json.NewEncoder(w).Encode(results)
    }





