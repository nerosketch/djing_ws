package usock

import (
	"log"
	"net"
	"os"
	"time"
)

type Socket struct {
	sockFname string
}

func NewSocket() *Socket {
	return &Socket{
		sockFname: "/run/djing_ws.sock",
	}
}

const (
	// Time allowed to write a message to the peer.
	readWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	//pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	//pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	//maxMessageSize = 512
)


func (s *Socket) Listen() {
	// s.addr + ":" + strconv.FormatInt(int64(s.port), 10)
	ln, err := net.Listen("unix", s.sockFname)
	if err != nil {
		log.Fatal("Failed to open socket")
		return
	}
	defer func() {
		log.Println("Free sock")
		if er := ln.Close(); er != nil{
			log.Println("Failed to close socket:", er.Error())
		}
		if er := os.Remove(s.sockFname); er != nil {
			log.Println("Failed to remove socket file:", er.Error())
		}
	}()

	log.Println("Listen...")

	buf := make([]byte, 0xfff)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Failed to listen socket")
			return
		}
		if rderr := conn.SetReadDeadline(time.Now().Add(readWait)); rderr != nil {
			log.Println("Failed to set read deadline")
			continue
		}
		recLen, err := conn.Read(buf)
		if err != nil {
			log.Fatal("Error reading:", err.Error())
		}
		log.Println("Received", recLen, "bytes:", buf[:recLen])
	}
}
