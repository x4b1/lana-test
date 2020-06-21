package memory

import (
	"context"
	"fmt"
	"sync"

	checkout "github.com/xabi93/lana-test/internal"
)

func NewBasketRepository(db *DB) *BasketRepository {
	return &BasketRepository{
		db:    db,
		mutex: sync.RWMutex{},
	}
}

type BasketRepository struct {
	db    *DB
	mutex sync.RWMutex
}

func (br *BasketRepository) Add(_ context.Context, b checkout.Basket) error {
	br.mutex.Lock()
	defer br.mutex.Unlock()
	if _, ok := br.db.baskets[b.ID]; ok {
		return fmt.Errorf("duplicated basket %s", b.ID)
	}

	br.db.baskets[b.ID] = b

	return nil
}

func (br *BasketRepository) Get(_ context.Context, id checkout.BasketID) (*checkout.Basket, error) {
	br.mutex.Lock()
	defer br.mutex.Unlock()
	b, exists := br.db.baskets[id]
	if !exists {
		return nil, nil
	}

	return &b, nil
}

func (br *BasketRepository) Delete(_ context.Context, id checkout.BasketID) error {
	br.mutex.Lock()
	defer br.mutex.Unlock()

	delete(br.db.baskets, id)

	return nil
}

func (br *BasketRepository) Save(_ context.Context, b checkout.Basket) error {
	br.mutex.Lock()
	defer br.mutex.Unlock()
	_, exists := br.db.baskets[b.ID]
	if !exists {
		return fmt.Errorf("basket not exists")
	}

	br.db.baskets[b.ID] = b

	return nil
}
