package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var router *Router

// ConfigRouting sets up all the server routes
func ConfigRouting() {

	router = &Router{DebugLog: true}
	router.POST("/post", router.testPost)
	router.GET("/get", router.testGet)
	router.GET("/remove", router.testRemove)
	router.GET("/list", router.listHandler)

	router.NotFoundHandler = http.FileServer(http.Dir("public"))
	router.AdminHandler = router.adminHandler

}

func main() {

	ConfigRouting()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)

		done <- true
	}()
	go func() {
		log.Fatal(http.ListenAndServe(":8081", router))
	}()

	fmt.Println("")
	router.Logger("JRouter running on port :8081")
	fmt.Println("")

	fmt.Println("CTRL + C to shutdown")
	<-done
	fmt.Println("Shutdown procedure")

}
