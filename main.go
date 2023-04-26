package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/khalidzahra/robot-logging-service/api"
)

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

func main() {
	setupLogFile()

	server := gin.Default()

	setupAPIRoutes(server)

	server.Run(":8000")
}
