package services

import (
	"context"
	"errors"
	"log"
	"subscriber/services/interfaces"
	"sync"
	"time"
)

// src
// https://habr.com/ru/articles/359078/

type RWMCache struct {
	sync.RWMutex
	items             map[string]Item
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

// Item struct cache item
type Item struct {
	Value      interface{}
	Expiration int64
	Created    time.Time
}

// New. Initializing a new memory cache
func NewRWMCache(defaultExpiration, cleanupInterval time.Duration) interfaces.Cache {
	items := make(map[string]Item)
	// cache item
	cache := RWMCache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}
	if cleanupInterval > 0 {
		cache.StartGC()
	}
	return &cache
}

// Set setting a cache by key
func (c *RWMCache) Set(key string, value interface{}, duration time.Duration) {
	var expiration int64
	if duration == 0 {
		duration = c.defaultExpiration
	}
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	c.Lock()
	defer c.Unlock()
	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

// Get getting a cache by key
func (c *RWMCache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	item, found := c.items[key]
	// cache not found
	if !found {
		return nil, false
	}
	if item.Expiration > 0 {
		// cache expired
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	return item.Value, true
}

// Delete cache by key
// Return false if key not found
func (c *RWMCache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()
	if _, found := c.items[key]; !found {
		return errors.New("key not found")
	}
	delete(c.items, key)
	return nil
}

// StartGC start Garbage Collection
func (c *RWMCache) StartGC() {
	go c.GC()
}

// GC Garbage Collection
func (c *RWMCache) GC() {
	for {
		<-time.After(c.cleanupInterval)
		if c.items == nil {
			return
		}
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}

}

// expiredKeys returns key list which are expired.
func (c *RWMCache) expiredKeys() (keys []string) {
	c.RLock()
	defer c.RUnlock()
	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}
	return
}

// clearItems removes all the items which key in keys.
func (c *RWMCache) clearItems(keys []string) {
	c.Lock()
	defer c.Unlock()
	for _, k := range keys {
		delete(c.items, k)
	}
}

func (c *RWMCache) Restore(orderService interfaces.OrderService, count int) {
	c.Lock()
	defer c.Unlock()
	orders, err := orderService.GetMany(count, context.Background())
	if err != nil {
		log.Fatalln("cannot get all orders")
		return
	}
	for _, value := range orders {
		expiration := time.Now().Add(c.defaultExpiration).UnixNano()
		c.items[value.OrderUID] = Item{
			Value:      value,
			Expiration: expiration,
			Created:    time.Now(),
		}
	}
}

func (c *RWMCache) GetAll() []interface{} {
	c.Lock()
	defer c.Unlock()
	resultSlice := make([]any, 0)
	for _, item := range c.items {
		if item.Expiration > 0 && time.Now().UnixNano() < item.Expiration {
			resultSlice = append(resultSlice, item.Value)
		}
	}
	return resultSlice
}

func (c *RWMCache) GetDefaultExpiration() time.Duration {
	return c.defaultExpiration
}
