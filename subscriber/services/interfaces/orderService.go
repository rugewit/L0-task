package interfaces

import (
	"context"
	"subscriber/model"
)

type OrderService interface {
	// Insert inserts one order into a storage
	Insert(order *model.Order, ctx context.Context) error

	// InsertMany inserts many orders into a storage
	InsertMany(order *model.Order, ctx context.Context) error

	// Get returns an order from a storage
	Get(id string, ctx context.Context, cache Cache) (model.Order, error)

	// GetMany returns {count} orders from a storage
	GetMany(count int, ctx context.Context) ([]model.Order, error)

	// GetAll returns all orders from a storage
	GetAll(ctx context.Context) ([]model.Order, error)
}
