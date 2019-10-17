package ws

import (
	"../glob_types"
	"bytes"
	"container/list"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"log"
)

type Hub struct {
	// Registered clients.
	clients						map[*Client]bool
	eventFromClientListeners	list.List

	broadcast 					chan []byte

	register 					chan *Client

	unregister					chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:				make(chan []byte),
		register:				make(chan *Client),
		unregister:				make(chan *Client),
		clients:				make(map[*Client]bool),
	}
}


func (h *Hub) sendPB2WS_JSON(v []byte, ifs proto.Message) bool {
	if msgErr := proto.Unmarshal(v, ifs); msgErr != nil {
		log.Println("Error unmarshalling message:", msgErr)
		return false
	}
	mrsh := jsonpb.Marshaler{}

	//w := &hubBroadcastWriter{hb: h,}
	var buf bytes.Buffer
	if err := mrsh.Marshal(&buf, ifs); err != nil {
		log.Println("error marshalling message to JSON:", err)
		return false
	}
	h.broadcast <- buf.Bytes()
	return true
}


func (h *Hub) WriteBroadcastMsg(v []byte) {
	// Converts from binary to json
	hdr := glob_types.MessageHeader{}

	if err := proto.Unmarshal(v, &hdr); err != nil {
		log.Println("Error unmarshalling broadcast message header:", err)
		return
	}

	// Converts PB binary message to JSON
	var isOk = false
	switch hdr.MessageType {
	case glob_types.MessageType_NewTask:
		msg := glob_types.NewTaskEvent{}
		isOk = h.sendPB2WS_JSON(v, &msg)
	case glob_types.MessageType_MessageNotify:
		msg := glob_types.MessageNotifyEvent{}
		isOk = h.sendPB2WS_JSON(v, &msg)
	case glob_types.MessageType_DeviceNotify:
		msg := glob_types.DeviceNotifyEvent{}
		isOk = h.sendPB2WS_JSON(v, &msg)
	default:
		log.Println("Unknown hdr message type:", hdr.MessageType)
		return
	}

	if !isOk {
		log.Println("Failed marshaling PB 2 JSON")
	}

}

func (h *Hub) SubscribeClientEvent(fn *glob_types.DataEvent) {
	h.eventFromClientListeners.PushBack(fn)
}
func (h *Hub) UnsubscribeClientEvent(fn *glob_types.DataEvent) {
	for e := h.eventFromClientListeners.Front(); e != nil; e = e.Next() {
		if fn == e.Value {
			h.eventFromClientListeners.Remove(e)
			break
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			// Got messages from client
			// Send data to listener callbacks
			//log.Println("Got message from client in hub", message)
			for e := h.eventFromClientListeners.Front(); e != nil; e = e.Next() {
				ev := e.Value.(*glob_types.DataEvent)
				go ev.Callback(message)
			}

			for client := range h.clients {
				select {
				case client.send <- message:

					//mr := jsonpb.Unmarshal()

				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
