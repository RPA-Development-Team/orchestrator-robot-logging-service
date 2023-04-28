package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/khalidzahra/robot-logging-service/auth"
)

type ClientRegistry map[*Client]bool

type Manager struct {
	registry      ClientRegistry
	sync.RWMutex             // To provide thread saftey for the registry
	egress        chan Event // Unbuffered channel for writing to connection to prevent writes being hogged up
	eventHandlers map[string]EventHandler
	TokenRegistry auth.TokenRegistry
}

func NewManager() *Manager {
	m := &Manager{
		registry:      make(ClientRegistry),
		egress:        make(chan Event),
		eventHandlers: make(map[string]EventHandler),
		TokenRegistry: auth.NewTokenRegistry(),
	}
	m.registerHandlers()
	return m
}

func (m *Manager) registerHandlers() {
	m.eventHandlers[EventLogEmit] = LogEmitEventHandler
}

func (m *Manager) HandleSocketConn(ctx *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	token := ctx.Query("token")
	if !m.TokenRegistry.ValidateToken(token) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid token",
		})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		log.Print("Request upgrade error:", err)
		return
	}

	client := NewClient(conn, m)
	m.RegisterClient(client)

	go client.InitReadListener()
	go client.InitWriteListener()
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

func (m *Manager) HandleEvent(e Event, c *Client) error {
	if handler, exists := m.eventHandlers[e.Type]; exists {
		if err := handler(e, c); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid event type: %s", e.Type)
	}
	return nil
}
