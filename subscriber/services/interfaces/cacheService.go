package interfaces

import (
	"time"
)

type Cache interface {
	// Set sets a key-value pair with duration expiration
	Set(key string, value interface{}, duration time.Duration)

	// Get returns a value by a key
	Get(key string) (interface{}, bool)

	// Delete deletes a key-value pair by a key
	Delete(key string) error

	// Restore tries to restore key-value pairs.
	// We don't know what pairs were stored before.
	// We set {count} first entries from the db to the cache.
	Restore(orderService OrderService, count int)

	// GetAll returns all values from the cache
	GetAll() []interface{}

	// GetDefaultExpiration returns default expiration
	GetDefaultExpiration() time.Duration
}
