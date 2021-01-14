package session

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gyf841010/pz-infra-new/log"
	"github.com/gyf841010/pz-infra-new/redisUtil"
)

// redis session store
type redisSessionStore struct {
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime int //expire time in seconds
}

// set value in redis session
func (rs *redisSessionStore) Set(key, value interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values[key] = value
	return nil
}

// get value in redis session
func (rs *redisSessionStore) Get(key interface{}) interface{} {
	rs.lock.RLock()
	defer rs.lock.RUnlock()
	if v, ok := rs.values[key]; ok {
		return v
	} else {
		return nil
	}
}

// delete value in redis session
func (rs *redisSessionStore) Delete(key interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.values, key)
	return nil
}

// clear all values in redis session
func (rs *redisSessionStore) Flush() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values = make(map[interface{}]interface{})
	return nil
}

// get redis session id
func (rs *redisSessionStore) SessionID() string {
	return rs.sid
}

// save session values to redis
func (rs *redisSessionStore) SessionRelease(httpResponseWriter http.ResponseWriter, sessionKey string) {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	b, err := EncodeGob(rs.values)
	if err != nil {
		return
	}
	values := string(b)

	if err := redisUtil.SetStringWithExpire(rs.sid, values, rs.maxlifetime); err != nil {
		log.Error("Set session values to redis failed", err.Error())
	}
	if err := redisUtil.SetStringWithExpire(getSesstionLifeTimeKey(rs.sid), strconv.Itoa(rs.maxlifetime), rs.maxlifetime); err != nil {
		log.Error("Set session values to redis failed", err.Error())
	}

	if httpResponseWriter != nil {
		httpResponseWriter.Header().Set(sessionKey, rs.sid)
	}
}

// redis session provider
type redisSessionProvider struct {
}

// init redis session
func (rp *redisSessionProvider) SessionInit(savePath string) error {
	return nil
}

// read redis session by sid
func (rp *redisSessionProvider) SessionRead(sid string) (SessionStore, error) {
	values, err := redisUtil.GetString(sid)
	if err != nil {
		return nil, err
	}
	lifeTimeStr, err := redisUtil.GetString(getSesstionLifeTimeKey(sid))
	if err != nil {
		log.Warnf("Get lifetime of session %s failed, use 10 minutes by default", sid)
		lifeTimeStr = strconv.Itoa(60 * 10)
	}
	if lifeTimeStr == "" {
		log.Warnf("Get lifetime of session %s is empty, use 10 minutes by default", sid)
		lifeTimeStr = strconv.Itoa(60 * 10)
	}
	lifeTime, err := strconv.Atoi(lifeTimeStr)
	if err != nil {
		log.Warnf("Get lifetime of session %s failed: %s", sid, err.Error())
		return nil, err
	}

	var kv map[interface{}]interface{}
	if len(values) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = DecodeGob([]byte(values))
		if err != nil {
			return nil, err
		}
	}

	rs := &redisSessionStore{sid: sid, values: kv, maxlifetime: lifeTime}
	return rs, nil
}

func (rp *redisSessionProvider) SessionGenerate(lifeTime int, sid string) (SessionStore, error) {
	kv := make(map[interface{}]interface{})
	rs := &redisSessionStore{sid: sid, values: kv, maxlifetime: lifeTime}
	return rs, nil
}

// check redis session exist by sid
func (rp *redisSessionProvider) SessionExist(sid string) bool {
	return redisUtil.Exists(sid)
}

// delete redis session by id
func (rp *redisSessionProvider) SessionDestroy(sid string) error {
	return redisUtil.Delete(sid)
}

// Impelment method, no used.
func (rp *redisSessionProvider) SessionGC() {
	return
}

// @todo
func (rp *redisSessionProvider) SessionAll() int {
	return 0
}

func getSesstionLifeTimeKey(sid string) string {
	return fmt.Sprintf("%s-lifetime", sid)
}
