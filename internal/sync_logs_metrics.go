package internal

import (
	"context"
	"sync"
)

type Callback func(ctx context.Context)

var (
	registryOnce      sync.Once
	registrySingleton *CallbackRegistry
)

type CallbackRegistry struct {
	mu       sync.RWMutex
	registry map[string][]Callback
}

func GetCallbackRegistry() *CallbackRegistry {
	registryOnce.Do(func() {
		registrySingleton = &CallbackRegistry{
			mu:       sync.RWMutex{},
			registry: make(map[string][]Callback),
		}
	})
	return registrySingleton
}

func (c *CallbackRegistry) Register(name string, callback Callback) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.registry[name] = append(c.registry[name], callback)
}

func (c *CallbackRegistry) Purge(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.registry, name)
}

func (c *CallbackRegistry) Event(ctx context.Context, name string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	wg := sync.WaitGroup{}
	for _, callback := range c.registry[name] {
		wg.Add(1)
		go callback(ctx)
	}
	wg.Wait()
}
