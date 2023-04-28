package api

import "github.com/gin-gonic/gin"

type APIRoute interface {
	RegisterRoutes(router *gin.RouterGroup)
	LoadEnvVariables()
}
