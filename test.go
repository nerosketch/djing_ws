package main

import (
	"./glob_types"
	"./usock"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	/*ev := glob_types.DeviceNotifyEvent{
		Device: 1,
		Group: 2,
		NotifyType: glob_types.DeviceNotifyEvent_DeviceDown,
	}

	dump, err := proto.Marshal(&ev)
	if err != nil {
		log.Fatal("Error marshalling:", err)
		return
	}

	ioutil.WriteFile("./test_dn.bin", dump, 0644)*/

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, os.Kill, syscall.SIGTERM)

	sc := usock.NewSocket()
	cde := glob_types.DataEvent{
		Callback: func(v []byte) {
			log.Println("Got data from socket in callback", v)
		},
	}
	sc.SubscribeEvent(&cde)

	go func() {
		signalType := <-signalCh
		signal.Stop(signalCh)
		log.Println("Exit command, received", signalType, ". Exiting...")
		sc.Stop()
		os.Exit(0)
	}()

	sc.Listen()
}
