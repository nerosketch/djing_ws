package usock

import (
	"../glob_types"
	"container/list"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Socket struct {
	sockFname string
	ln net.Listener
	eventSubscribers list.List
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

func (s *Socket) Stop() {
	log.Println("Free sock")
	if er := s.ln.Close(); er != nil{
		log.Println("Failed to close socket:", er.Error())
	}
	if er := os.Remove(s.sockFname); er != nil {
		log.Println("Failed to remove socket file:", er.Error())
	}
}

func (s *Socket) SubscribeEvent(ev *glob_types.DataEvent) {
	s.eventSubscribers.PushBack(ev)
}
func (s *Socket) UnsubscribeEvent(ev *glob_types.DataEvent) {
	for e := s.eventSubscribers.Front(); e != nil; e = e.Next() {
		if ev == e.Value {
			s.eventSubscribers.Remove(e)
			break
		}
	}
}


func (s *Socket) Listen() {
	var err error
	s.ln, err = net.Listen("unix", s.sockFname)
	if err != nil {
		log.Println("Failed to open socket")
		return
	}
	defer s.Stop()

	log.Println("Listen...")

	buf := make([]byte, 0xfff)
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println("Failed to listen socket")
			return
		}
		if rderr := conn.SetReadDeadline(time.Now().Add(readWait)); rderr != nil {
			log.Println("Failed to set read deadline")
			continue
		}
		recLen, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("Received empty data, continue")
			}else{
				log.Println("Error reading:", err.Error())
			}
			continue
		}

		// send data to listener callbacks
		for e := s.eventSubscribers.Front(); e != nil; e = e.Next() {
			ev := e.Value.(*glob_types.DataEvent)
			go ev.Callback(buf[:recLen])
		}
	}
}
