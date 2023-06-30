package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"fmt"

	"github.com/khalidzahra/robot-logging-service/entity"
)

type ILogService interface {
	Save(entity.Log) bool
}

type LogService struct {
	elasticEndpoint string
}

func NewLogService() ILogService {
	return &LogService{
		elasticEndpoint: os.Getenv("ELASTIC_ENDPOINT"),
	}
}

func (service *LogService) Save(log entity.Log) bool {
	payload, err := json.Marshal(log)
	if err != nil {
		return false
	}

	r, err := http.NewRequest("POST", service.elasticEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return false
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return false
	}

	defer res.Body.Close()

	return res.StatusCode == http.StatusOK
}