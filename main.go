package main

import (
	"./glob_types"
	"./usock"
	"./ws"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()


	// Web socket
	hub := ws.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})


	go func() {
		if err := http.ListenAndServe(*addr, nil); err != nil {
			log.Fatal("Error listen and serve http: ", err.Error())
		}
	}()





	// Local socket
	localSocket := usock.NewSocket()
	cde := glob_types.DataEvent{
		Callback: func(v []byte) {
			//log.Println("Got data from socket in callback", v)
			hub.WriteBroadcastMsg(v)
		},
	}
	localSocket.SubscribeEvent(&cde)

	// Catch system signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		signalType := <-signalCh
		signal.Stop(signalCh)
		log.Println("Exit command, received", signalType, "Exiting...")
		localSocket.Stop()
		os.Exit(0)
	}()


	localSocket.Listen()
}
