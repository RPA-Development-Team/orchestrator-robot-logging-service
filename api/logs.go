package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogRoute struct {
}

func (logRoute LogRoute) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/logs", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Success",
		})
	})

	router.POST("/logs", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Success",
		})
	})
}
