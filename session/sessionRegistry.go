package session

import (
	"sync"

	. "github.com/yiGmMk/pz-infra-new/logging"
	"github.com/yiGmMk/pz-infra-new/redisUtil"
)

const (
	SESSION_REGISTRY_KEY_PREFIX = "SESSION_REG_PREFIX_"
)

// Provider contains global session Registry Method,
// Reversing with user id and user client id
type SessionRegistry interface {
	// Registry User Session When Create Session
	SessionRegistry(userId, clientId, sessionId string, lifeTime int) error
	// Find Registry User Session
	GetUserSession(userId string) (sessionId, clientId string, err error)
}

var registry SessionRegistry
var registryOnce = &sync.Once{}

func Registry() SessionRegistry {
	registryOnce.Do(func() {
		if registry == nil {
			registry = new(redisSessionRegistry)
		}
	})
	return registry
}

type registryObject struct {
	UserId    string `json:"user_id"`
	ClientId  string `json:"client_id"`
	SessionId string `json:"session_id"`
}

type redisSessionRegistry struct {
}

// Registry User Session When Create Session
func (r *redisSessionRegistry) SessionRegistry(userId, clientId, sessionId string, lifeTime int) error {
	registry := registryObject{
		UserId:    userId,
		ClientId:  clientId,
		SessionId: sessionId,
	}
	return redisUtil.SetObjectWithExpire(r.getRegistryKey(userId), &registry, lifeTime)
}

// Find Registry User Session
func (r *redisSessionRegistry) GetUserSession(userId string) (string, string, error) {
	redisKey := r.getRegistryKey(userId)
	var sessionRegistry registryObject
	if err := redisUtil.GetObject(redisKey, &sessionRegistry); err != nil {
		Log.Error("Failed to Get User Session Registry", With("userId", userId), WithError(err))
		return "", "", err
	}
	if sessionRegistry.SessionId == "" {
		Log.Debug("No User Session Registry Found", With("userId", userId))
		return "", "", nil
	}
	return sessionRegistry.SessionId, sessionRegistry.ClientId, nil
}

func (r *redisSessionRegistry) getRegistryKey(userId string) string {
	return SESSION_REGISTRY_KEY_PREFIX + userId
}
