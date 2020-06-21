package deleting

import (
	"context"

	checkout "github.com/xabi93/lana-test/internal"
)

type baskets interface {
	Get(context.Context, checkout.BasketID) (*checkout.Basket, error)
	Delete(context.Context, checkout.BasketID) error
}

//NewBasketDeleter creates a new BasketDeleter instance
func NewBasketDeleter(b baskets) *BasketDeleter {
	return &BasketDeleter{b}
}

//BasketDeleter defines the service to delete baskets
type BasketDeleter struct {
	baskets baskets
}

//Delete checks if the given basket exists and deletes it from system
func (s BasketDeleter) Delete(ctx context.Context, id checkout.BasketID) error {
	p, err := s.baskets.Get(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return checkout.ErrBasketNotExists
	}

	return s.baskets.Delete(ctx, id)
}
