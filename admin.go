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

	var err error
	router.AdminConnection, err = upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		router.Logger("Error upgrading the admin handler to websocket")
		fmt.Println(err)
		goto NOWEBSOCKET
	} else {
		router.AdminConnected = true
	}

	// for {
	// 	// messageType, message, err := router.AdminConnection.ReadJSON()
	// 	messageType, message, err := router.AdminConnection.ReadMessage() //text
	// 	fmt.Println("Message type", messageType)
	// 	if err != nil {
	// 		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
	// 			log.Printf("error: %v, user-agent: %v", err, r.Header.Get("User-Agent"))
	// 			router.AdminConnected = false
	// 		}
	// 		return
	// 	}
	// 	fmt.Printf("%s sent: %s\n", router.AdminConnection.RemoteAddr(), string(message))

	// 	// Write message back to browser
	// 	if err = router.AdminConnection.WriteMessage(messageType, message); err != nil {
	// 		fmt.Println(err)
	// 	}

	// }

NOWEBSOCKET:
}

// WriteToAdminConsole ...
func (router *Router) WriteToAdminConsole(message interface{}) {
	// func (router *Router) WriteToAdminConsole(message ...string) {
	// var messages []string
	// messages = append(messages, message...)

	// output := struct {
	// 	Type    string   `json:"type,omitempty"`
	// 	Message []string `json:"message,omitempty"`
	// }{
	// 	Type:    "test Json message",
	// 	Message: messages,
	// }

	if router.AdminConnected {
		router.AdminConnection.WriteJSON(message)
	}
}
