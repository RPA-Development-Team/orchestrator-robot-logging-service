package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func setupLogFile() {
	f, _ := os.Create("server.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {
	setupLogFile()

	server := gin.Default()

	server.Run(":8000")
}
