package ws

import (
	"../glob_types"
	"container/list"
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

func (h *Hub) WriteBroadcastMsg(msg []byte) {
	// Converts from binary to json
	hdr := glob_types.MessageHeader{}
	//TODO: сделать перевод из бинаря в JSON
	// получать тип сообщения уже получилось, в hdr
	if err := proto.Unmarshal(msg, &hdr); err != nil {
		log.Println("Error unmarshalling broadcast message:", err)
	}
	h.broadcast <- msg
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
