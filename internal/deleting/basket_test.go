package deleting_test

import (
	"context"
	"fmt"
	"testing"

	checkout "github.com/xabi93/lana-test/internal"
	"github.com/xabi93/lana-test/internal/deleting"
	"github.com/xabi93/lana-test/pkg/errors"

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

func (m *basketsMock) Delete(ctx context.Context, id checkout.BasketID) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

type DeleteBasketSuite struct {
	suite.Suite

	ctx context.Context

	service *deleting.BasketDeleter
	baskets *basketsMock

	basket checkout.Basket
}

func (s *DeleteBasketSuite) SetupSuite() {
	s.basket = checkout.NewBasket()
}

func (s *DeleteBasketSuite) SetupTest() {
	s.ctx = context.Background()

	s.baskets = new(basketsMock)

	s.service = deleting.NewBasketDeleter(s.baskets)
}

func (s DeleteBasketSuite) TestFailsCannotGetBasket() {
	expectedError := fmt.Errorf("unexpected error")

	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(nil, expectedError)

	err := s.service.Delete(s.ctx, s.basket.ID)

	s.Equal(expectedError, err)
}

func (s DeleteBasketSuite) TestFailsBasketDoesNotExists() {
	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(nil, nil)

	err := s.service.Delete(s.ctx, s.basket.ID)

	s.True(errors.IsNotFound(err))
	s.Contains(err.Error(), checkout.ErrBasketNotExists.Error())
}

func (s DeleteBasketSuite) TestFailsWhenCannotDeleteBasket() {
	expectedError := fmt.Errorf("unexpected error")

	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(&s.basket, nil).
		On("Delete", s.ctx, s.basket.ID).
		Return(expectedError)

	err := s.service.Delete(s.ctx, s.basket.ID)

	s.Equal(expectedError, err)
}

func (s DeleteBasketSuite) TestSuccessWhenDeleteBasket() {
	s.baskets.
		On("Get", s.ctx, s.basket.ID).
		Return(&s.basket, nil).
		On("Delete", s.ctx, s.basket.ID).
		Return(nil)

	err := s.service.Delete(s.ctx, s.basket.ID)

	s.Nil(err)
}

func TestDeleteBasketSuite(t *testing.T) {
	suite.Run(t, new(DeleteBasketSuite))
}
