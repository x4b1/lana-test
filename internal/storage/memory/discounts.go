package memory

import (
	"context"
	"sync"

	checkout "github.com/xabi93/lana-test/internal"
)

func NewDiscountRepository(db *DB) *DiscountRepository {
	return &DiscountRepository{
		db:    db,
		mutex: sync.RWMutex{},
	}
}

type DiscountRepository struct {
	db    *DB
	mutex sync.RWMutex
}

func (pr *DiscountRepository) All(ctx context.Context) ([]checkout.Discount, error) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	return pr.db.discounts, nil
}
