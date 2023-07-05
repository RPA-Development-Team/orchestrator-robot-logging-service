package api

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khalidzahra/robot-logging-service/entity"
	"github.com/khalidzahra/robot-logging-service/ws"
)

var (
	authURL              string
	keycloakClientId     string
	keycloakClientSecret string
	Manager              *ws.Manager
)

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
	keycloakClientId = os.Getenv("KC_CLIENT_ID")
	keycloakClientSecret = os.Getenv("KC_CLIENT_SECRET")
}

func handleAuthRequest(user entity.User) bool {
	payload := url.Values{}
	payload.Set("grant_type", "password")
	payload.Set("client_id", keycloakClientId)
	payload.Set("client_secret", keycloakClientSecret)
	payload.Set("username", user.Username)
	payload.Set("password", user.Password)

	r, err := http.NewRequest(http.MethodPost, authURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return false
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return false
	}

	defer res.Body.Close()

	return res.StatusCode == http.StatusOK
}
