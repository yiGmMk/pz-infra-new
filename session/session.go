package session

import (
	"net/http"
	"sync"

	. "github.com/yiGmMk/pz-infra-new/logging"

	"github.com/astaxie/beego"
)

// SessionStore contains all data for one session process with specific id.
type SessionStore interface {
	Set(key, value interface{}) error                 //set session value
	Get(key interface{}) interface{}                  //get session value
	Delete(key interface{}) error                     //delete session value
	SessionID() string                                //back current sessionID
	SessionRelease(w http.ResponseWriter, key string) // release the resource & save data to provider & return the data
	Flush() error                                     //delete all data
}

// Provider contains global session methods and saved SessionStores.
// it can operate a SessionStore by its id.
type SessionProvider interface {
	//gclifetime -1: session persistence
	SessionInit(config string) error
	SessionRead(sid string) (SessionStore, error)
	SessionExist(sid string) bool
	SessionGenerate(lifeTime int, sid string) (SessionStore, error)
	SessionDestroy(sid string) error
	SessionAll() int //get all active session
	SessionGC()
}

type session struct {
	sid    string
	values map[interface{}]interface{}
}

var provider SessionProvider
var initializer = &sync.Once{}

func initialize() {
	providerName := beego.AppConfig.String("sessionProvider")
	switch providerName {
	case "redis":
		provider = new(redisSessionProvider)
		provider.SessionInit("")
	default:
		Log.Info("Missing session provider configuration. Use redis instead")
		provider = new(redisSessionProvider)
		provider.SessionInit("")
	}
}

func Provider() SessionProvider {
	initializer.Do(initialize)
	return provider
}
