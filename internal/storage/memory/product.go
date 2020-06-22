package memory

import (
	"context"
	"sync"

	checkout "github.com/xabi93/lana-test/internal"
)

func NewProductRepository(db *DB) *ProductRepository {
	return &ProductRepository{
		db:    db,
		mutex: sync.RWMutex{},
	}
}

type ProductRepository struct {
	db    *DB
	mutex sync.RWMutex
}

// Get gets the product by code and returns if it exists
func (pr *ProductRepository) Get(ctx context.Context, c checkout.ProductCode) (*checkout.Product, error) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()
	p, exists := pr.db.products[c]
	if !exists {
		return nil, nil
	}

	return &p, nil
}

//All returns all products in the storage
func (pr *ProductRepository) All(ctx context.Context) (map[checkout.ProductCode]checkout.Product, error) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	return pr.db.products, nil
}
