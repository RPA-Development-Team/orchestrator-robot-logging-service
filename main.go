package main

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/khalidzahra/robot-logging-service/api"
	"github.com/khalidzahra/robot-logging-service/service"
	"github.com/khalidzahra/robot-logging-service/ws"
)

var logService service.ILogService

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func setupLogFile() {
	f, _ := os.Create("server.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func setupAPIRoutes(server *gin.Engine) {
	apiRouter := server.Group("/api")
	routes := []api.APIRoute{api.RobotRoute{}}
	for _, route := range routes {
		route.LoadEnvVariables()
		route.RegisterRoutes(apiRouter)
	}
}

func setupWebsocketRoute(server *gin.Engine) {
	manager := ws.NewManager()
	api.Manager = manager
	server.GET("/rtlogs", manager.HandleSocketConn)
}

func initServices() {
	logService = service.NewLogService()
	ws.LogService = logService
}

func main() {
	loadEnvVars()
	setupLogFile()
	initServices()

	server := gin.Default()

	setupWebsocketRoute(server)
	setupAPIRoutes(server)

	server.Run(":8000")
}
