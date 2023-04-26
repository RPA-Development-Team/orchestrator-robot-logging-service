package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/khalidzahra/robot-logging-service/entity"
	"github.com/khalidzahra/robot-logging-service/service"
)

type ILogController interface {
	FindAll() []entity.Log
	Save(ctx *gin.Context) error
}

type LogController struct {
	service service.ILogService
}

func NewLogController(service service.ILogService) LogController {
	return LogController{service: service}
}

func (controller LogController) FindAll() []entity.Log {
	return controller.service.FindAll()
}

func (controller LogController) Save(ctx *gin.Context) error {
	var log entity.Log

	err := ctx.ShouldBindJSON(&log)
	if err != nil {
		return err
	}

	controller.service.Save(log)
	return nil
}
