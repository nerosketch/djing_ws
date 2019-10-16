package main

import "./usock"

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

	sc := usock.NewSocket()
	sc.Listen()
}
