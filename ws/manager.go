package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ClientRegistry map[*Client]bool

type Manager struct {
	registry ClientRegistry
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		registry: make(ClientRegistry),
	}
}

func (m *Manager) HandleSocketConn(ctx *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		log.Print("Request upgrade error:", err)
		return
	}

	client := NewClient(conn, m)
	m.RegisterClient(client)

	go client.InitListener()
}

func (m *Manager) RegisterClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.registry[client] = true
}

func (m *Manager) RemoveClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, exists := m.registry[client]; exists {
		client.conn.Close()
		delete(m.registry, client)
	}
}
