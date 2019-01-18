package main

import (
	"log"
	"net/http"
)

var router *Router

// ConfigRouting sets up all the server routes
func ConfigRouting() {

	router = &Router{DebugLog: true}
	router.POST("/post", router.LogWrapper(router.testPost))
	router.GET("/get", router.LogWrapper(router.testGet))
	router.GET("/remove", router.LogWrapper(router.testRemove))
	router.GET("/list", router.LogWrapper(router.listHandler))

	router.NotFoundHandler = http.FileServer(http.Dir("public"))
	router.AdminHandler = router.adminHandler

}

func main() {

	ConfigRouting()

	log.Fatal(http.ListenAndServe(":8081", router))

	// sigs := make(chan os.Signal, 1)
	// done := make(chan bool, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// go func() {

	// 	sig := <-sigs
	// 	fmt.Println()
	// 	fmt.Println(sig)

	// 	done <- true
	// }()
	// go func() {
	// 	log.Fatal(http.ListenAndServe(":8081", router))
	// }()

	// fmt.Println("")
	// router.Logger("JRouter running on port :8081")
	// fmt.Println("")

	// fmt.Println("CTRL + C to shutdown")
	// <-done
	// fmt.Println("Shutdown procedure ...")
	// router.WriteToAdminConsole("JRouter Shutdown procedure ...")

}
