package adding

import (
	"context"

	checkout "github.com/xabi93/lana-test/internal"
	"github.com/xabi93/lana-test/pkg/errors"
)

//baskets provides access to baskets repository
type baskets interface {
	Get(context.Context, checkout.BasketID) (*checkout.Basket, error)
	Save(context.Context, checkout.Basket) error
}

//products provides access to products repository
type products interface {
	Get(context.Context, checkout.ProductCode) (*checkout.Product, error)
}

//NewBasketItemsAdded creates a new BasketItemsAdder instance
func NewBasketItemsAdder(b baskets, p products) *BasketItemsAdder {
	return &BasketItemsAdder{b, p}
}

//BasketItemsAdder defines the service to add products to a basket
type BasketItemsAdder struct {
	baskets  baskets
	products products
}

//Add adds a new item to the basket if the given basket and product exists in the system
func (a BasketItemsAdder) Add(ctx context.Context, id checkout.BasketID, c checkout.ProductCode) error {
	b, err := a.getBasket(ctx, id)
	if err != nil {
		return err
	}

	p, err := a.getProduct(ctx, c)
	if err != nil {
		return err
	}

	b.AddItem(p.Code)

	return a.baskets.Save(ctx, *b)
}

func (s BasketItemsAdder) getBasket(ctx context.Context, id checkout.BasketID) (*checkout.Basket, error) {
	p, err := s.baskets.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.WrapNotFound(checkout.ErrBasketNotExists, "Adding item to basket")
	}

	return p, nil
}

func (s BasketItemsAdder) getProduct(ctx context.Context, c checkout.ProductCode) (*checkout.Product, error) {
	p, err := s.products.Get(ctx, c)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.WrapNotFound(checkout.ErrProductNotExists, "Adding item to basket")
	}

	return p, nil
}
