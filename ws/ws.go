package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebsocketHandler(ctx *gin.Context) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		log.Print("Request upgrade error:", err)
		return
	}

	defer ws.Close()

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Incoming: %s", msg)

		err = ws.WriteMessage(msgType, msg)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
