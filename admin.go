package main

//-------------------------------------------------------------//
// ## JRouter Admin ##
//-------------------------------------------------------------//
//
// Let's test the websocket system for the adminHandler
//
//
//
// ""
//

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// adminHandler
func (router *Router) adminHandler(w http.ResponseWriter, r *http.Request) {

	// NOTE: this is just testing the websocket connections

	upgrader.CheckOrigin = func(r *http.Request) bool {
		//TODO: check if origin is allowed
		// Bypass for now
		origin := r.Header.Get("Origin")
		fmt.Println("Websocket Origin: ", origin)
		return true
	}

	connection, err := upgrader.Upgrade(w, r, w.Header())

	if err != nil {
		router.Logger("Error upgrading the admin handler to websocket")
		fmt.Println(err)
		goto NOWEBSOCKET
	}

	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v, user-agent: %v", err, r.Header.Get("User-Agent"))
			}
			return
		}
		fmt.Printf("%s sent: %s\n", connection.RemoteAddr(), string(message))

		// Write message back to browser
		if err = connection.WriteMessage(messageType, message); err != nil {
			fmt.Println(err)
		}

	}

NOWEBSOCKET:
}
