package ws

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	manager *Manager
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		conn:    conn,
		manager: manager,
	}
}

// InitListener starts listening for incoming messages from the client to the server.
func (c *Client) InitListener() {
	defer c.manager.RemoveClient(c)
	for {
		msgType, payload, err := c.conn.ReadMessage()

		// Connection was closed somehow
		if err != nil {
			break
		}

		fmt.Println(string(payload))
		c.conn.WriteMessage(msgType, payload)
	}
}
