package ws

import (
	"container/list"
)

type Hub struct {
	// Registered clients.
	clients                  map[*Client]bool
	eventFromClientListeners list.List

	broadcast chan []byte

	register chan *Client

	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) WriteBroadcastMsg(msg []byte) {
	h.broadcast <- msg
}

func (h *Hub) SubscribeClientEvent(fn *ClientDataEvent) {
	h.eventFromClientListeners.PushBack(fn)
}
func (h *Hub) UnsubscribeClientEvent(fn *ClientDataEvent) {
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
				ev := e.Value.(*ClientDataEvent)
				go ev.Callback(message)
			}

			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
