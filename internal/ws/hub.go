package ws

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c.UserID] = c
		case c := <-h.unregister:
			if _, ok := h.clients[c.UserID]; ok {
				delete(h.clients, c.UserID)
				close(c.send)
			}
		}
	}
}

func (h *Hub) Broadcast(msg []byte) {
	for _, c := range h.clients {
		c.PostMessage(msg)
	}
}
