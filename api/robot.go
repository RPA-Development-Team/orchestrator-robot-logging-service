package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/khalidzahra/robot-logging-service/entity"
	"github.com/khalidzahra/robot-logging-service/ws"
)

var (
	authURL              string
	keycloakClientId     string
	keycloakClientSecret string
	Manager              *ws.Manager
)

type KeycloakResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

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

		if ok, userId := handleAuthRequest(user); ok {
			token := Manager.TokenRegistry.GenerateToken()
			ctx.JSON(http.StatusOK, gin.H{
				"token":  token.Token,
				"userId": userId,
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

func handleAuthRequest(user entity.User) (bool, string) {
	payload := url.Values{}
	payload.Set("grant_type", "password")
	payload.Set("client_id", keycloakClientId)
	payload.Set("client_secret", keycloakClientSecret)
	payload.Set("username", user.Username)
	payload.Set("password", user.Password)

	r, err := http.NewRequest(http.MethodPost, authURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return false, ""
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return false, ""
	}

	defer res.Body.Close()

	kcRes := &KeycloakResponse{}

	err = json.NewDecoder(res.Body).Decode(kcRes)

	if err != nil {
		return false, ""
	}

	decodedToken, _, err := new(jwt.Parser).ParseUnverified(kcRes.AccessToken, jwt.MapClaims{})

	if err != nil {
		fmt.Print(err)
		return false, ""
	}

	if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok {
		return res.StatusCode == http.StatusOK, claims["sub"].(string)
	}

	return false, ""
}
