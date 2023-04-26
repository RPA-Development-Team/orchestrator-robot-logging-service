package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogRoute struct {
}

func (logRoute LogRoute) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/log", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Success",
		})
	})

	router.POST("/log", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Success",
		})
	})
}
