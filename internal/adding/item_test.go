package adding_test

import (
	"context"
	"fmt"
	"testing"

	checkout "github.com/xabi93/lana-test/internal"
	"github.com/xabi93/lana-test/internal/adding"

	"github.com/xabi93/lana-test/pkg/errors"
	"github.com/xabi93/lana-test/pkg/money"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type basketsMock struct{ mock.Mock }

func (m *basketsMock) Get(ctx context.Context, id checkout.BasketID) (*checkout.Basket, error) {
	args := m.Called(ctx, id)
	var b *checkout.Basket

	if args.Get(0) != nil {
		b = args.Get(0).(*checkout.Basket)
	}

	return b, args.Error(1)
}

func (m *basketsMock) Save(ctx context.Context, b checkout.Basket) error {
	args := m.Called(ctx, b)

	return args.Error(0)
}

type productsMock struct{ mock.Mock }

func (m *productsMock) Get(ctx context.Context, c checkout.ProductCode) (*checkout.Product, error) {
	args := m.Called(ctx, c)
	var b *checkout.Product

	if args.Get(0) != nil {
		b = args.Get(0).(*checkout.Product)
	}

	return b, args.Error(1)
}

type AddingItemsSuite struct {
	suite.Suite

	ctx context.Context

	service  *adding.BasketItemsAdder
	baskets  *basketsMock
	products *productsMock

	basket  checkout.Basket
	product checkout.Product
}

func (s *AddingItemsSuite) SetupSuite() {
	s.basket = checkout.NewBasket(checkout.BasketCurrency(money.EUR))
	s.product = checkout.Product{
		Code:  checkout.ProductCode("PEN"),
		Name:  checkout.ProductName("Lana Pen"),
		Price: checkout.ProductPrice{Money: money.Eur(99)},
	}
}

func (s *AddingItemsSuite) SetupTest() {
	s.ctx = context.Background()

	s.baskets = new(basketsMock)
	s.products = new(productsMock)

	s.service = adding.NewBasketItemsAdder(s.baskets, s.products)
}

func (s AddingItemsSuite) TestFailsCannotGetBasket() {
	expectedError := fmt.Errorf("unexpected error")

	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(nil, expectedError)

	err := s.service.Add(s.ctx, s.basket.ID, s.product.Code)

	s.Equal(expectedError, err)
}

func (s AddingItemsSuite) TestFailsBasketDoesNotExists() {
	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(nil, nil)

	err := s.service.Add(s.ctx, s.basket.ID, s.product.Code)

	s.True(errors.IsNotFound(err))
	s.Contains(err.Error(), checkout.ErrBasketNotExists.Error())
}

func (s AddingItemsSuite) TestFailsWhenCannotGetProduct() {
	expectedError := fmt.Errorf("unexpected error")

	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(&s.basket, nil)

	s.products.
		On("Get", s.ctx, s.product.Code).
		Return(nil, expectedError)

	err := s.service.Add(s.ctx, s.basket.ID, s.product.Code)

	s.Equal(expectedError, err)
}

func (s AddingItemsSuite) TestFailsWhenProductNotExists() {
	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(&s.basket, nil)

	s.products.
		On("Get", s.ctx, s.product.Code).
		Return(nil, nil)

	err := s.service.Add(s.ctx, s.basket.ID, s.product.Code)

	s.True(errors.IsNotFound(err))
	s.Contains(err.Error(), checkout.ErrProductNotExists.Error())

}

func (s AddingItemsSuite) TestFailsWhenCannotSaveBasket() {
	expectedError := fmt.Errorf("unexpected error")

	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(&s.basket, nil).
		On("Save", s.ctx, mock.AnythingOfType("checkout.Basket")).
		Return(expectedError)

	s.products.
		On("Get", s.ctx, s.product.Code).
		Return(&s.product, nil)

	err := s.service.Add(s.ctx, s.basket.ID, s.product.Code)

	s.Equal(expectedError, err)
}

func (s AddingItemsSuite) TestSuccessWhenAddsProductToBasketAndSaves() {
	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(&s.basket, nil).
		On("Save", s.ctx, mock.AnythingOfType("checkout.Basket")).
		Return(nil)

	s.products.
		On("Get", s.ctx, s.product.Code).
		Return(&s.product, nil)

	err := s.service.Add(s.ctx, s.basket.ID, s.product.Code)

	s.Nil(err)

	s.baskets.AssertCalled(s.T(), "Save", s.ctx, mock.MatchedBy(func(b checkout.Basket) bool {
		if b.ID != s.basket.ID {
			return false
		}
		if len(s.basket.Items) != 1 {
			return false
		}

		return s.basket.Items[s.product.Code].Product == s.product.Code
	}))
}

func TestAddingItemsSuite(t *testing.T) {
	suite.Run(t, new(AddingItemsSuite))
}
