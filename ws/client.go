package ws

import (
	"encoding/json"
	"log"

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

// InitReadListener starts listening for incoming messages from the client to the server.
func (c *Client) InitReadListener() {
	defer c.manager.RemoveClient(c)
	for {
		_, payload, err := c.conn.ReadMessage()

		// Connection was closed somehow
		if err != nil {
			break
		}

		var event Event
		if err := json.Unmarshal(payload, &event); err != nil {
			log.Println(err)
			return
		}

		if err := c.manager.HandleEvent(event, c); err != nil {
			// Sending error to client in an ErrorMessageEvent
			errMsg, _ := json.Marshal(ErrorMessageEvent{
				Error: err.Error(),
			})
			c.manager.egress <- Event{
				Type:    EventErrorMessage,
				Payload: errMsg,
			}
		}
	}
}

// InitWriteListener listens for new messages in the channel associated with the client's manager
func (c *Client) InitWriteListener() {
	defer c.manager.RemoveClient(c)

	for {
		select {
		case event, exists := <-c.manager.egress:
			if !exists { // Channel was closed somehow so exit routine
				return
			}

			payload, err := json.Marshal(event)

			if err == nil {
				if err := c.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
					log.Println(err)
				}
			}
		}
	}

}
