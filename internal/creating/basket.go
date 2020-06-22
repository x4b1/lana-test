package creating

import (
	"context"

	checkout "github.com/xabi93/lana-test/internal"
)

//baskets provides access to baskets repository
type baskets interface {
	Add(context.Context, checkout.Basket) error
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
func (bc BasketCreator) Create(ctx context.Context, c checkout.BasketCurrency) (checkout.Basket, error) {
	b := checkout.NewBasket(c)

	if err := bc.baskets.Add(ctx, b); err != nil {
		return checkout.Basket{}, err
	}

	return b, nil
}
