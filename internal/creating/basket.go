package creating

import (
	"context"

	basket "github.com/xabi93/lana-test/internal"
)

//baskets provides access to baskets repository
type baskets interface {
	Add(context.Context, basket.Basket) error
}

//NewBasketCreator creates a new BasketCreator instance
func NewBasketCreator(b baskets) *BasketCreator {
	return &BasketCreator{b}
}

//BasketCreator defines the service to create baskets
type BasketCreator struct {
	baskets baskets
}

//Create creates a new basket in system an returns it
func (bc BasketCreator) Create(ctx context.Context) (basket.Basket, error) {
	b := basket.NewBasket()

	if err := bc.baskets.Add(ctx, b); err != nil {
		return basket.Basket{}, err
	}

	return b, nil
}
