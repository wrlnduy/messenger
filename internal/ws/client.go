package ws

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte

	UserID uuid.UUID
}

func NewClient(hub *Hub, conn *websocket.Conn, userID uuid.UUID) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 100),
		UserID: userID,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("User:%v failed on reading message: %v\n", c.UserID, err)
			break
		}
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("User:%v failed on writing message: %v\n", c.UserID, err)
			break
		}
	}
}

func (c *Client) PostMessage(msg []byte) {
	select {
	case c.send <- msg:
	default:
		c.hub.unregister <- c
	}
}
