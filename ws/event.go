package ws

import (
	"encoding/json"
	"fmt"

	"github.com/khalidzahra/robot-logging-service/entity"
	"github.com/khalidzahra/robot-logging-service/service"
)

const (
	EventLogEmit      = "logEmitEvent"
	EventErrorMessage = "errorMessageEvent"
	EventLogReceive   = "logReceiveEvent"
)

var LogService service.ILogService

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type LogEmitEvent struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	RobotID   uint64 `json:"robotId"`
}

type ErrorMessageEvent struct {
	Error string `json:"error"`
}

type LogReceiveEvent struct {
	Message string `json:"message"`
}

type EventHandler func(e Event, c *Client) error

func LogEmitEventHandler(e Event, c *Client) error {
	var logEvent LogEmitEvent

	if err := json.Unmarshal(e.Payload, &logEvent); err != nil {
		return fmt.Errorf("invalid payload error:\n %v", err)
	}

	logEntry := entity.Log{
		Timestamp: logEvent.Timestamp,
		Message:   logEvent.Message,
		RobotID:   logEvent.RobotID,
	}

	go LogService.Save(logEntry)

	logReceiveEvent, _ := json.Marshal(LogReceiveEvent{
		Message: "Log entry receieved.",
	})

	c.manager.egress <- Event{
		Type:    EventLogReceive,
		Payload: logReceiveEvent,
	}

	return nil
}
