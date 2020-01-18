package conn

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kirileec/ot/ot"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var defaultSession = NewSession(`package main

import "fmt"

func main() {
	fmt.Println("Hello, playground")
}`)

func InitOtServer() {
	ot.TextEncoding = ot.TextEncodingTypeUTF8

	r := mux.NewRouter()
	r.HandleFunc("/ws", serveWs)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	go defaultSession.HandleEvents()

	fmt.Println("Listening on port :2020")

	err := http.ListenAndServe(":2020", r)
	if err != nil {
		fmt.Println("Error: ", err)
	}

}

func serveWs(w http.ResponseWriter, r *http.Request) {
	var err error
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			fmt.Println(err)
		}
		return
	}

	NewConnection(defaultSession, conn).Handle()
}
