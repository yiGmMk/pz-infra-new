package config

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Configer defines how to get and set value from configuration.
// Each get operation will return a zero value if it’s not found.
// To check if a given key exists, try `IsSet()` .
type Configer interface {
	Get(key string) interface{}
	String(key string) string
	Bool(key string) bool
	Int(key string) int
	Float64(key string) float64
	Time(key string) time.Time
	Map(key string) map[string]interface{}
	Strings(key string) []string
	StringMap(key string) map[string]string

	Set(key string, value interface{})

	IsSet(key string) bool
}

// Provider is the interface that must be implemented by a configuration provider.
type Provider interface {
	// New returns a new configer to the configuration.
	// The option is a provider-specific value used to set option(s) for the provider.
	New(option interface{}) (Configer, error)
}

var (
	providersMu sync.Mutex
	providers   = make(map[string]Provider)
)

// Register makes a configuration provider available by the provided name.
// If Register is called twice with the same name or if provider is nil,
// it panics.
func Register(name string, provider Provider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if provider == nil {
		panic("config: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("config: Register called twice for provider " + name)
	}
	providers[name] = provider
}

// Providers returns a sorted list of the names of the registered providers.
func Providers() []string {
	providersMu.Lock()
	defer providersMu.Unlock()
	var names []string
	for name := range providers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetConfiger returns a provided configer.
func GetConfiger(name string, option interface{}) (Configer, error) {
	provider, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("config: unknown provider[%s] (forgot to Register?)", name)
	}
	return provider.New(option)
}
