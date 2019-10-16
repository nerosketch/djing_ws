package main

import (
	"./glob_types"
	"flag"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
)

func main() {
	/*ev := glob_types.MessageNotifyEvent{
		ChatId: 13,
		AuthorId: 2,
		ParticipantsLen: 4,
		Participants: []uint32{1, 2, 3, 4},
		Length: 234,
		Text: "TestText",
	}*/
	ev := glob_types.DeviceNotifyEvent{
		MessageType: glob_types.MessageType_DeviceNotify,
		Device: 9,
		Group: 0,
		NotifyType: glob_types.DeviceNotifyEvent_Unknown,
	}

	dump, err := proto.Marshal(&ev)
	if err != nil {
		log.Fatal("Error marshalling:", err)
		return
	}

	ioutil.WriteFile("./test_dn.bin", dump, 0644)

	oev := glob_types.MessageHeader{}
	proto.Unmarshal(dump, &oev)

	flag.Parse()


	/*// Catch system signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, os.Kill, syscall.SIGTERM)



	// Web socket
	hub := ws.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Listen and serve http: ", err.Error())
	}





	// Local socket
	localSocket := usock.NewSocket()
	cde := glob_types.DataEvent{
		Callback: func(v []byte) {
			//log.Println("Got data from socket in callback", v)
			//FIXME: надо сделать переход от бинарного формата в JSON
			hub.WriteBroadcastMsg(v)
		},
	}
	localSocket.SubscribeEvent(&cde)
	go func() {
		signalType := <-signalCh
		signal.Stop(signalCh)
		log.Println("Exit command, received", signalType, "Exiting...")
		localSocket.Stop()
		os.Exit(0)
	}()
	localSocket.Listen()*/
}
