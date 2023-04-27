package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/khalidzahra/robot-logging-service/api"
	"github.com/khalidzahra/robot-logging-service/ws"
)

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Printf("ran")
}

func setupLogFile() {
	f, _ := os.Create("server.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func setupAPIRoutes(server *gin.Engine) {
	apiRouter := server.Group("/api")
	routes := []api.APIRoute{api.LogRoute{}}
	for _, route := range routes {
		route.RegisterRoutes(apiRouter)
	}
}

func setupWebsocketRoute(server *gin.Engine) {
	server.GET("/rtlogs", ws.WebsocketHandler)
}

func main() {
	loadEnvVars()
	setupLogFile()

	server := gin.Default()

	setupAPIRoutes(server)
	setupWebsocketRoute(server)

	server.Run(":8000")
}
