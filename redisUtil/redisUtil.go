package redisUtil

import (
	"encoding/base64"
	"sync"
	"time"

	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gyf841010/pz-infra-new/encryptUtil"
	. "github.com/gyf841010/pz-infra-new/errorUtil"
	. "github.com/gyf841010/pz-infra-new/logging"
	"github.com/gyf841010/pz-infra-new/slackUtil"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

var once sync.Once

func initPool() {
	url := beego.AppConfig.String("redisUrl")
	Log.Debug("-------------- redis initPool with URL ", With("url", url))

	ciphertext := beego.AppConfig.String("redisPass")
	if len(ciphertext) > 0 {
		str, err := base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			panic(err)
		}
		password, err := encryptUtil.AesDecrypt(encryptUtil.INTERNAL_KEY, string(str))
		if err != nil {
			panic(err)
		}

		pool = newPool(url, string(password))
	} else {
		Log.Debug("-------------- not configure password ")
		pool = newPool(url, "")
	}
}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				Log.Error("redis dial error:", WithError(err))
				return nil, err
			}
			if len(password) > 0 {
				_, err = c.Do("AUTH", password)
				if err != nil {
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func GetPool() *redis.Pool {
	return getPool()
}

func getPool() *redis.Pool {
	if pool == nil {
		once.Do(initPool)
	}

	if pool == nil {
		Log.Error("redis pool is nil, init fail.")
	}
	return pool
}

var SetObject = func(key string, value interface{}) error {
	conn := getPool().Get()
	defer conn.Close()

	if _, err := do(conn, "HMSET", redis.Args{}.Add(key).AddFlat(value)...); err != nil {
		Log.Error("redisUtil SetObject error:", WithError(err))
		return err
	}
	return nil
}

var GetObject = func(key string, value interface{}) (err error) {
	conn := getPool().Get()
	defer conn.Close()
	v, err := redis.Values(do(conn, "HGETALL", key))
	if err != nil {
		Log.Error("redisUtil GetObject error:", WithError(err))
		return err
	}

	Log.Debug("redisUtil GetObject v:", With("object", v))

	if err := redis.ScanStruct(v, value); err != nil {
		Log.Error("Redist Util Error for getting", With("key", key), WithError(err))
		return err
	}
	Log.Debug("redisUtil GetObject value:", With("value", value))
	return nil
}

//设置key多少秒后超时
func Expire(key string, seconds int) error {
	conn := getPool().Get()
	defer conn.Close()
	_, err := do(conn, "EXPIRE", key, seconds)
	if err != nil {
		Log.Error("redisUtil Expire error: ", WithError(err))
	}
	return nil
}

func Delete(key string) error {
	conn := getPool().Get()
	defer conn.Close()
	if _, err := do(conn, "DEL", key); err != nil {
		Log.Error("redisUtil Delete error:", WithError(err))
		return err
	}
	return nil
}

func SetComplexObject(key string, value interface{}) error {
	bytes, _ := json.Marshal(value)
	if err := SetString(key, string(bytes)); err != nil {
		return err
	}
	return nil
}

func SetComplexObjectExpire(key string, value interface{}, expire int) error {
	bytes, _ := json.Marshal(value)
	if err := SetStringWithExpire(key, string(bytes), expire); err != nil {
		return err
	}
	return nil
}

func GetComplexObject(key string, value interface{}) error {
	str, err := GetString(key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(str), value); err != nil {
		return err
	}
	return nil
}

func SetObjectWithExpire(key string, value interface{}, expire int) error {
	conn := getPool().Get()
	defer conn.Close()
	if _, err := do(conn, "HMSET", redis.Args{}.Add(key).AddFlat(value)...); err != nil {
		Log.Error("redisUtil SetObject error:", WithError(err))
		return err
	}
	Expire(key, expire)
	return nil
}

func SetStringWithExpire(key string, value string, expire int) error {
	conn := getPool().Get()
	defer conn.Close()
	if _, err := do(conn, "SET", key, value, "EX", expire); err != nil {
		Log.Error("SetStringWithExpire error ", WithError(err))
		return err
	}
	return nil
}

func SetString(key string, value string) error {
	conn := getPool().Get()
	defer conn.Close()
	if _, err := do(conn, "SET", key, value); err != nil {
		Log.Error("SetString error ", WithError(err))
		return err
	}
	return nil
}

func GetString(key string) (string, error) {
	conn := getPool().Get()
	defer conn.Close()
	v, err := redis.String(do(conn, "GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return "", nil
		}
		return "", err
	}
	return v, nil
}

func Exists(key string) bool {
	conn := getPool().Get()
	defer conn.Close()
	v, err := redis.Bool(do(conn, "EXISTS", key))
	if err != nil {
		Log.Error("redisUtil Exist error:", WithError(err))
		return false
	}
	return v
}

func AddGeoIndex(indexName string, geoKey string, latitude float32, longitude float32) error {
	conn := getPool().Get()
	defer conn.Close()
	_, err := do(conn, "GEOADD", indexName, latitude, longitude, geoKey)
	if err != nil {
		Log.Error("add geo index error:", WithError(err))
		return err
	}
	return nil
}

func GetKeysByPrefix(prefix string) ([]string, error) {
	conn := getPool().Get()
	defer conn.Close()
	keys, err := redis.Strings(do(conn, "KEYS", prefix+"*"))
	if err != nil {
		return nil, err
	}
	return keys, nil
}

//往redis里面插入键值，如果键已存在，返回False，不执行
//如果键不存在，插入键值,返回成功
func SetStringIfNotExist(key, value string, expire int) (bool, error) {
	conn := getPool().Get()
	defer conn.Close()
	result, err := redis.String(do(conn, "SET", key, value, "EX", expire, "NX"))
	if err != nil {
		if err == redis.ErrNil {
			return false, nil
		} else {
			Log.Error("SetStringIfNotExist error ", WithError(err))
			return false, err
		}
	}
	if result == "OK" {
		return true, nil
	} else {
		return false, nil
	}
}

var (
	errKeyIsBlank        = NewHErrorCustom(ERROR_CODE_REDIS_KEY_NULL)
	errValueIsNotPointer = NewHErrorCustom(ERROR_CODE_REDIS_VALUE_NULL_PTR)
	errValueIsNil        = NewHErrorCustom(ERROR_CODE_REDIS_VALUE_NULL)
	ErrKeyNotFound       = NewHErrorCustom(ERROR_CODE_REDIS_KEY_NOT_EXIST)
)

// GetValue key should not be blank and value must be non-nil pointer to int/bool/string/struct ...
// because a value can set only if it is addressable
// if key not found will return ErrKeyNotFound
func GetValue(key string, value interface{}) (err error) {
	if len(key) == 0 {
		return errKeyIsBlank
	}
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return errValueIsNotPointer
	}
	if v.IsNil() {
		return errValueIsNil
	}

	conn := getPool().Get()
	defer conn.Close()

	reply, err := redis.Values(do(conn, "MGET", key))
	if err != nil {
		Log.Error("redis: get Error", With("key", key), WithError(err))
		return err
	}
	if len(reply) == 0 || reply[0] == nil {
		return ErrKeyNotFound
	}

	if v.Elem().Kind() == reflect.Struct {
		err = json.Unmarshal(reply[0].([]byte), value)
	} else {
		_, err = redis.Scan(reply, value)
	}
	if err != nil {
		Log.Error("redis: scan error", With("key", key), WithError(err))
		return err
	}
	Log.Debug("redis: get success", With("key", key))
	return nil
}

// SetValue key should not be blank and value should not be nil
// Struct or pointer to struct values will encoding as JSON objects
func SetValue(key string, value interface{}, seconds ...int) (err error) {
	if len(key) == 0 {
		return errKeyIsBlank
	}
	v := reflect.ValueOf(value)
	isPtr := (v.Kind() == reflect.Ptr)
	if value == nil || (isPtr && v.IsNil()) {
		return errValueIsNil
	}

	conn := getPool().Get()
	defer conn.Close()

	if v.Kind() == reflect.Struct || (isPtr && v.Elem().Kind() == reflect.Struct) {
		bs, err := json.Marshal(value)
		if err != nil {
			return err
		}
		value = string(bs)
	} else {
		if isPtr { // *int/*bool/*string ...
			value = v.Elem()
		}
	}
	exArgs := []interface{}{}
	if len(seconds) > 0 {
		exArgs = append(exArgs, "EX", seconds[0])
	}
	args := append([]interface{}{key, value}, exArgs...)
	_, err = do(conn, "SET", args...)
	if err != nil {
		Log.Error("redis: set error ", With("key", key), WithError(err))
		return err
	}
	Log.Debug("redis: set success", With("key", key))
	return nil
}

func do(c redis.Conn, commandName string, args ...interface{}) (reply interface{}, err error) {
	reply, err = c.Do(commandName, args...)
	if err != nil {
		handleAlertError(err)
	}
	return reply, err
}

func getRedisUrl() string {
	return beego.AppConfig.String("redisUrl")
}

func handleAlertError(err error) {
	if err == nil {
		return
	}
	Log.Debug("send slack alert")
	message := fmt.Sprintf("redis error with redisUrl : %s , error is %s", getRedisUrl(), err.Error())
	slackUtil.SendMessage(message)
}

func SetStrings(key string, ss []string, seconds ...int) error {
	if len(key) == 0 {
		return errKeyIsBlank
	}

	conn := getPool().Get()
	defer conn.Close()
	var err error
	for _, s := range ss {
		if err = conn.Send("SADD", key, s); err != nil {
			Log.Error("redis: Send Error", WithError(err))
			return err
		}
	}
	if err = conn.Flush(); err != nil {
		Log.Error("redis: Flush Error", WithError(err))
		return err
	}
	_, err = conn.Do("")
	if err != nil {
		Log.Error("redis: SADD %s %v", With("key", key), With("strings", ss), WithError(err))
		return err
	}
	if len(seconds) > 0 {
		_, err = conn.Do("EXPIRE", key, seconds[0])
		if err != nil {
			Log.Error("redis: EXPIRE Key", With("key", key), With("seconds", seconds[0]), WithError(err))
			return err
		}
	}
	Log.Debug("redis: SetStrings success", With("key", key))
	return nil
}

func GetStrings(key string) ([]string, error) {
	if len(key) == 0 {
		return nil, errKeyIsBlank
	}

	conn := getPool().Get()
	defer conn.Close()
	ss, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		Log.Error("redis: SMEMBERS Error ", With("key", key), WithError(err))
		return nil, err
	}
	Log.Debug("redis: GetStrings success", With("key", key))
	return ss, nil
}

func Incr(key string) (*int, error) {
	if len(key) == 0 {
		return nil, errKeyIsBlank
	}

	conn := getPool().Get()
	defer conn.Close()

	id, err := redis.Int(conn.Do("INCR", key))
	if err != nil {
		Log.Error("redisUtil INCR error:", WithError(err))
		return nil, err
	}
	Log.Debug("redis: Incr success", With("key", key), With("id", id))
	return &id, nil
}

//SET if Not exists
func SetValueNX(key string, value interface{}, seconds ...int) (err error) {
	if Exists(key) {
		return nil
	}
	if len(key) == 0 {
		return errKeyIsBlank
	}
	v := reflect.ValueOf(value)
	isPtr := (v.Kind() == reflect.Ptr)
	if value == nil || (isPtr && v.IsNil()) {
		return errValueIsNil
	}

	conn := getPool().Get()
	defer conn.Close()

	if v.Kind() == reflect.Struct || (isPtr && v.Elem().Kind() == reflect.Struct) {
		bs, err := json.Marshal(value)
		if err != nil {
			return err
		}
		value = string(bs)
	} else {
		if isPtr { // *int/*bool/*string ...
			value = v.Elem()
		}
	}
	exArgs := []interface{}{}
	if len(seconds) > 0 {
		exArgs = append(exArgs, "EX", seconds[0])
	}
	args := append([]interface{}{key, value}, exArgs...)
	_, err = do(conn, "SET", args...)
	if err != nil {
		Log.Error("redis: set error ", With("key", key), WithError(err))
		return err
	}
	Log.Debug("redis: set success", With("key", key))
	return nil
}

//SET Hash string
func SetHashStringWithExpire(key, field, value string, seconds ...int) (err error) {
	if len(key) == 0 {
		return errKeyIsBlank
	}

	conn := getPool().Get()
	defer conn.Close()

	if err = conn.Send("HSET", key, field, value); err != nil {
		Log.Error("redis: Send Error", WithError(err))
		return err
	}

	if err = conn.Flush(); err != nil {
		Log.Error("redis: Flush Error", WithError(err))
		return err
	}
	_, err = conn.Do("")
	if err != nil {
		Log.Error("redis: HSET %s %v", With("key", key), With("field", field), With("value", value), WithError(err))
		return err
	}
	if len(seconds) > 0 {
		_, err = conn.Do("EXPIRE", key, seconds[0])
		if err != nil {
			Log.Error("redis: EXPIRE Key", With("key", key), With("seconds", seconds[0]), WithError(err))
			return err
		}
	}
	Log.Debug("redis: HSET success", With("key", key), With("field", field), With("value", value))
	return nil
}

func GetHashStrings(key string) ([]string, error) {
	if len(key) == 0 {
		return nil, errKeyIsBlank
	}

	conn := getPool().Get()
	defer conn.Close()
	ss, err := redis.Strings(conn.Do("HGETALL", key))
	if err != nil {
		Log.Error("redis: HGETALL Error ", With("key", key), WithError(err))
		return nil, err
	}
	Log.Debug("redis: GetHashStrings success", With("key", key), With("ss", ss))
	return ss, nil
}

//GET Hash string
func GetHashString(key, field string) (string, error) {
	conn := getPool().Get()
	defer conn.Close()
	v, err := redis.String(do(conn, "HGET", key, field))
	if err != nil {
		if err == redis.ErrNil {
			return "", nil
		}
		return "", err
	}
	return v, nil
}

func LpushString(key, value string) error {
	conn := getPool().Get()
	defer conn.Close()
	if _, err := do(conn, "LPUSH", key, value); err != nil {
		Log.Error("redisUtil Lpush Object error:", WithError(err))
		return err
	}
	return nil
}
