package main

import (
	"log"
	"net/http"
)



func main() {
	hub := NewHub()
	go hub.Run()

	//handling login request
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		serveLogin(hub, w, r)
	})
	//handing index request
	http.HandleFunc("/", serveIndex)
	//handing websocket request
	//serveWs defined in client.go
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
