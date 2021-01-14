// Shamelessly copied from http://hjr265.github.io/redsync.go/
package lockUtil

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/gyf841010/pz-infra-new/redisUtil"

	"github.com/garyburd/redigo/redis"
)

const (
	// DefaultExpiry is used when Mutex Duration is 0
	DefaultExpiry = 8 * time.Second
	// DefaultTries is used when Mutex Duration is 0
	DefaultTries = 16
	// DefaultDelay is used when Mutex Delay is 0
	DefaultDelay = 512 * time.Millisecond
	// DefaultFactor is used when Mutex Factor is 0
	DefaultFactor = 0.01
)

var (
	// ErrFailed is returned when lock cannot be acquired
	ErrFailed = errors.New("failed to acquire lock")
)

type lockConfig struct {
	Name   string        // Resouce name
	Expiry time.Duration // Duration for which the lock is valid, DefaultExpiry if 0

	Tries int           // Number of attempts to acquire lock before admitting failure, DefaultTries if 0
	Delay time.Duration // Delay between two attempts to acquire lock, DefaultDelay if 0

	Factor float64 // Drift factor, DefaultFactor if 0

	Quorum int // Quorum for the lock, set to len(addrs)/2+1 by NewMutex()

	value string // value is used in order to release the lock in a safe way
	until time.Time

	nodes []*redis.Pool
}

func NewLockConfig(name, value string) *lockConfig {
	if name == "" {
		panic("redsync: name is empty")
	}
	if value == "" {
		panic("redsync: value is empty")
	}
	nodes := []*redis.Pool{redisUtil.GetPool()}
	if len(nodes) == 0 {
		panic("redsync: nodes is empty")
	}
	return &lockConfig{
		Name:   name,
		Quorum: len(nodes)/2 + 1,
		value:  value,
		nodes:  nodes,
	}
}

func AcquireLock(c *lockConfig) error {
	var value string
	if c.value == "" {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			return err
		}
		value = base64.StdEncoding.EncodeToString(b)
	} else {
		value = c.value
	}

	expiry := c.Expiry
	if expiry == 0 {
		expiry = DefaultExpiry
	}

	retries := c.Tries
	if retries == 0 {
		retries = DefaultTries
	}

	for i := 0; i < retries; i++ {
		n := 0
		start := time.Now()
		for _, node := range c.nodes {
			if node == nil {
				continue
			}

			conn := node.Get()
			reply, err := redis.String(conn.Do("set", c.Name, value, "nx", "px", int(expiry/time.Millisecond)))
			conn.Close()
			if err != nil {
				continue
			}
			if reply != "OK" {
				continue
			}
			n++
		}

		factor := c.Factor
		if factor == 0 {
			factor = DefaultFactor
		}

		until := time.Now().Add(expiry - time.Now().Sub(start) - time.Duration(int64(float64(expiry)*factor)) + 2*time.Millisecond)
		if n >= c.Quorum && time.Now().Before(until) {
			c.value = value
			c.until = until
			return nil
		}

		// no need to clean up

		delay := c.Delay
		if delay == 0 {
			delay = DefaultDelay
		}
		time.Sleep(delay)
	}

	return ErrFailed
}

func TouchLock(c *lockConfig) bool {
	value := c.value
	if value == "" {
		panic("redsync: touch of unlocked mutex")
	}

	expiry := c.Expiry
	if expiry == 0 {
		expiry = DefaultExpiry
	}
	reset := int(expiry / time.Millisecond)

	n := 0
	for _, node := range c.nodes {
		if node == nil {
			continue
		}

		conn := node.Get()
		reply, err := touchScript.Do(conn, c.Name, value, reset)
		conn.Close()
		if err != nil {
			continue
		}
		if reply != "OK" {
			continue
		}
		n++
	}
	if n >= c.Quorum {
		return true
	}
	return false
}

func ReleaseLock(c *lockConfig) bool {
	value := c.value
	if value == "" {
		panic("redsync: unlock of unlocked mutex")
	}

	c.value = ""
	c.until = time.Unix(0, 0)

	n := 0
	for _, node := range c.nodes {
		if node == nil {
			continue
		}

		conn := node.Get()
		status, err := delScript.Do(conn, c.Name, value)
		conn.Close()
		if err != nil {
			continue
		}
		if status == 0 {
			continue
		}
		n++
	}
	if n >= c.Quorum {
		return true
	}
	return false
}

var delScript = redis.NewScript(1, `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`)

var touchScript = redis.NewScript(1, `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("set", KEYS[1], ARGV[1], "xx", "px", ARGV[2])
else
	return "ERR"
end`)
