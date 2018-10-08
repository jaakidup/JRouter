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
	output := "TestPost function"
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
	router.NotFound = http.FileServer(http.Dir("public"))

	go func() {

		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)

		done <- true
	}()
	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	fmt.Println("")
	router.Logger("JRouter running on port :8080")
	fmt.Println("")

	fmt.Println("CTRL + C to shutdown")
	<-done
	fmt.Println("Shutdown procedure")

}
