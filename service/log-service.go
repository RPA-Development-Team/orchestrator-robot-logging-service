package service

import (
	"github.com/khalidzahra/robot-logging-service/entity"
	"github.com/khalidzahra/robot-logging-service/repository"
)

type ILogService interface {
	Save(*entity.Log) entity.Log
	Update(log entity.Log)
	Delete(log entity.Log)
	FindAll() []entity.Log
}

type LogService struct {
	logRepository repository.LogRepository
}

func NewLogService() ILogService {
	return &LogService{
		logRepository: repository.NewLogRepository(),
	}
}

func (service *LogService) Save(log *entity.Log) entity.Log {
	service.logRepository.Save(log)
	return *log
}

func (service *LogService) FindAll() []entity.Log {
	return service.logRepository.FindAll()
}

func (service *LogService) Update(log entity.Log) {
	service.logRepository.Update(log)
}

func (service *LogService) Delete(log entity.Log) {
	service.logRepository.Delete(log)
}
