package auth

import (
	"time"

	"github.com/google/uuid"
)

type RobotToken struct {
	Token   string
	Created time.Time
}

type TokenRegistry map[string]RobotToken

// expirationPeriod is the time after which a token will expire.
// Specified in minutes
const expirationPeriod = 30

func NewTokenRegistry() TokenRegistry {
	registry := make(TokenRegistry)
	go registry.expiredTokenPurgeService() // Start purge service with new registry
	return registry
}

func (registry TokenRegistry) GenerateToken() RobotToken {
	token := RobotToken{
		Token:   uuid.NewString(),
		Created: time.Now(),
	}
	registry[token.Token] = token
	return token
}

func (registry TokenRegistry) expiredTokenPurgeService() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		for _, token := range registry {
			if token.Created.Add(expirationPeriod * time.Minute).Before(time.Now()) { // Token has expired
				delete(registry, token.Token)
			}
		}
	}
}

func (registry TokenRegistry) ValidateToken(token string) bool {
	if robotToken, exists := registry[token]; exists {
		delete(registry, robotToken.Token)
		return true
	}
	return false
}
