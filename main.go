package main

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/khalidzahra/robot-logging-service/api"
)

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
	routes := []api.APIRoute{api.LogRoute{}}
	for _, route := range routes {
		route.RegisterRoutes(apiRouter)
	}
}

func main() {
	loadEnvVars()
	setupLogFile()

	server := gin.Default()

	setupAPIRoutes(server)

	server.Run(":8000")
}
