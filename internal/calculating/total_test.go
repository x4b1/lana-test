package calculating_test

import (
	"context"
	"fmt"
	"testing"

	checkout "github.com/xabi93/lana-test/internal"
	"github.com/xabi93/lana-test/internal/calculating"

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

type productsMock struct{ mock.Mock }

func (m *productsMock) All(ctx context.Context) (map[checkout.ProductCode]checkout.Product, error) {
	args := m.Called(ctx)
	var p map[checkout.ProductCode]checkout.Product

	if args.Get(0) != nil {
		p = args.Get(0).(map[checkout.ProductCode]checkout.Product)
	}

	return p, args.Error(1)
}

type discountsMock struct{ mock.Mock }

func (m *discountsMock) All(ctx context.Context) ([]checkout.Discount, error) {
	args := m.Called(ctx)
	var d []checkout.Discount

	if args.Get(0) != nil {
		d = args.Get(0).([]checkout.Discount)
	}

	return d, args.Error(1)
}

type CalculatingBasketTotalSuite struct {
	suite.Suite

	ctx context.Context

	service   *calculating.TotalBasket
	baskets   *basketsMock
	products  *productsMock
	discounts *discountsMock

	basket       checkout.Basket
	productList  map[checkout.ProductCode]checkout.Product
	discountList []checkout.Discount
}

func (s *CalculatingBasketTotalSuite) SetupSuite() {
	product := checkout.Product{
		Code:  checkout.ProductCode("PEN"),
		Name:  checkout.ProductName("Lana Pen"),
		Price: checkout.ProductPrice{Money: money.Eur(1500)},
	}

	s.productList = map[checkout.ProductCode]checkout.Product{product.Code: product}

	s.basket = checkout.NewBasket(checkout.BasketCurrency(money.EUR))
	s.basket.AddItem(product.Code)
	s.basket.AddItem(product.Code)
	s.basket.AddItem(product.Code)

	s.discountList = []checkout.Discount{new(oneEurDiscount)}
}

func (s *CalculatingBasketTotalSuite) SetupTest() {
	s.ctx = context.Background()

	s.baskets = new(basketsMock)
	s.products = new(productsMock)
	s.discounts = new(discountsMock)

	s.service = calculating.NewTotalBasket(s.baskets, s.products, s.discounts)
}

func (s CalculatingBasketTotalSuite) TestFailsCannotGetBasket() {
	expectedError := fmt.Errorf("unexpected error")

	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(nil, expectedError)

	_, err := s.service.Total(s.ctx, s.basket.ID)

	s.Equal(expectedError, err)
}

func (s CalculatingBasketTotalSuite) TestFailsBasketDoesNotExists() {
	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(nil, nil)

	_, err := s.service.Total(s.ctx, s.basket.ID)

	s.True(errors.IsNotFound(err))
	s.Contains(err.Error(), checkout.ErrBasketNotExists.Error())
}

func (s CalculatingBasketTotalSuite) TestFailsWhenCannotGetProducts() {
	expectedError := fmt.Errorf("unexpected error")

	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(&s.basket, nil)

	s.products.
		On("All", s.ctx).
		Return(nil, expectedError)

	_, err := s.service.Total(s.ctx, s.basket.ID)

	s.Equal(expectedError, err)
}

func (s CalculatingBasketTotalSuite) TestFailsWhenCannotGetDiscounts() {
	expectedError := fmt.Errorf("unexpected error")

	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(&s.basket, nil)

	s.products.
		On("All", s.ctx).
		Return(s.productList, nil)

	s.discounts.
		On("All", s.ctx).
		Return(nil, expectedError)

	_, err := s.service.Total(s.ctx, s.basket.ID)

	s.Equal(expectedError, err)
}

func (s CalculatingBasketTotalSuite) TestReturnsTotalAmountCalculated() {
	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(&s.basket, nil)

	s.products.
		On("All", s.ctx).
		Return(s.productList, nil)

	s.discounts.
		On("All", s.ctx).
		Return(s.discountList, nil)

	total, err := s.service.Total(s.ctx, s.basket.ID)

	s.NoError(err)
	s.Equal(money.Eur(4400), total)
}

func TestCalculatingBasketTotalSuite(t *testing.T) {
	suite.Run(t, new(CalculatingBasketTotalSuite))
}

type oneEurDiscount struct{}

func (d oneEurDiscount) Calculate(itemList map[checkout.ProductCode]checkout.Item) money.Money {
	return money.Eur(100)
}
