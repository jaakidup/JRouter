package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	output := "Index Route     / "
	io.WriteString(w, output)
}

func testPost(w http.ResponseWriter, r *http.Request) {
	output := "Test function"
	io.WriteString(w, output)
}

func test(w http.ResponseWriter, r *http.Request) {
	output := "Test function"
	io.WriteString(w, output)
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	router := &Router{DebugLog: true}

	router.POST("/post", testPost)
	router.GET("/get", test)
	router.GET("/", indexHandler)

	router.Routes.ListKeys()

	go func() {

		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)

		done <- true
	}()
	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	fmt.Println("Awaiting Signal")
	<-done
	fmt.Println("Shutdown procedure")

}
