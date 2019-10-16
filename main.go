package main

import (
	"./ws"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	hub := ws.NewHub()
	go hub.Run()

	//hub.WriteBroadcastMsg("sdfsdf")

	// Subscribing to new data from clients
	/*cde := ws.ClientDataEvent{
		Callback: func(v []byte) {
			log.Println("Got data from client in callback", v)
		},
	}
	hub.SubscribeClientEvent(&cde)*/

	//http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
