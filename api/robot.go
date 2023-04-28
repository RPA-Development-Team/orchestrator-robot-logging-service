package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/khalidzahra/robot-logging-service/entity"
	"github.com/khalidzahra/robot-logging-service/ws"
)

var authURL string

var Manager *ws.Manager

type RobotRoute struct {
}

func (robotRoute RobotRoute) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/robot/login", func(ctx *gin.Context) {
		var user entity.User

		err := ctx.ShouldBindJSON(&user)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if handleAuthRequest(user) {
			token := Manager.TokenRegistry.GenerateToken()
			ctx.JSON(http.StatusOK, gin.H{
				"token": token.Token,
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid credentials",
			})
		}

	})
}

func (robotRoute RobotRoute) LoadEnvVariables() {
	authURL = os.Getenv("AUTH_URL")
}

func handleAuthRequest(user entity.User) bool {
	payload, err := json.Marshal(user)
	if err != nil {
		return false
	}

	r, err := http.NewRequest("POST", authURL, bytes.NewBuffer(payload))
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
