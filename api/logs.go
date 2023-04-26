package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khalidzahra/robot-logging-service/controller"
	"github.com/khalidzahra/robot-logging-service/repository"
	"github.com/khalidzahra/robot-logging-service/service"
)

type LogRoute struct {
}

func (logRoute LogRoute) RegisterRoutes(router *gin.RouterGroup) {
	var (
		logRepository repository.LogRepository  = repository.NewLogRepository()
		logService    service.ILogService       = service.NewLogService(logRepository)
		logController controller.ILogController = controller.NewLogController(logService)
	)

	router.GET("/logs", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, logController.FindAll())
	})

	router.POST("/logs", func(ctx *gin.Context) {
		err := logController.Save(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Log received successfully",
			})
		}
	})
}
