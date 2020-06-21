package calculating

import (
	"context"

	checkout "github.com/xabi93/lana-test/internal"
	"github.com/xabi93/lana-test/pkg/money"
)

//baskets provides access to baskets repository
type baskets interface {
	Get(context.Context, checkout.BasketID) (*checkout.Basket, error)
}

//products provides access to products repository
type products interface {
	All(context.Context) (map[checkout.ProductCode]checkout.Product, error)
}

//discounts provides access to discounts repository
type discounts interface {
	All(context.Context) ([]checkout.Discount, error)
}

//NewTotalBasket creates a new TotalBasket instance
func NewTotalBasket(baskets baskets, products products, discounts discounts) *TotalBasket {
	return &TotalBasket{baskets, products, discounts}
}

//TotalBasket defines the service to calculate total amount of baskets
type TotalBasket struct {
	baskets   baskets
	products  products
	discounts discounts
}

//Total given a basket returns his total amount if it exists
func (t TotalBasket) Total(ctx context.Context, id checkout.BasketID) (money.Money, error) {
	b, err := t.getBasket(ctx, id)
	if err != nil {
		return money.Money{}, err
	}

	p, err := t.products.All(ctx)
	if err != nil {
		return money.Money{}, err
	}

	d, err := t.discounts.All(ctx)
	if err != nil {
		return money.Money{}, err
	}

	calculator := checkout.TotalCalculator{
		Discounts: d,
		Basket:    *b,
		Products:  p,
	}

	return calculator.Calculate()
}

func (t TotalBasket) getBasket(ctx context.Context, id checkout.BasketID) (*checkout.Basket, error) {
	p, err := t.baskets.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, checkout.ErrBasketNotExists
	}

	return p, nil
}