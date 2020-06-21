package creating_test

import (
	"context"
	"errors"
	"testing"

	checkout "github.com/xabi93/lana-test/internal"
	"github.com/xabi93/lana-test/internal/creating"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type basketsMock struct{ mock.Mock }

func (m *basketsMock) Add(ctx context.Context, b checkout.Basket) error {
	args := m.Called(ctx, b)

	return args.Error(0)
}

type CreatingBasketSuite struct {
	suite.Suite

	ctx context.Context

	service *creating.BasketCreator
	baskets *basketsMock

	currency checkout.BasketCurrency
}

func (s *CreatingBasketSuite) SetupTest() {
	s.ctx = context.Background()

	s.baskets = new(basketsMock)

	s.service = creating.NewBasketCreator(s.baskets)
}

func (s CreatingBasketSuite) TestFailsWhenCannotSaveNewBasket() {
	expectedError := errors.New("unexpected error")

	s.baskets.
		On("Add", s.ctx, mock.AnythingOfType("checkout.Basket")).
		Return(expectedError)

	_, err := s.service.Create(s.ctx, s.currency)

	s.Equal(expectedError, err)
}

func (s CreatingBasketSuite) TestSuccessWhenSaveAndReturnsNewBasket() {
	s.baskets.
		On("Add", s.ctx, mock.AnythingOfType("checkout.Basket")).
		Return(nil)

	b, err := s.service.Create(s.ctx, s.currency)

	s.Nil(err)
	s.NotNil(b)

	s.baskets.AssertCalled(s.T(), "Add", s.ctx, b)
}

func TestCreatingBasketSuite(t *testing.T) {
	suite.Run(t, new(CreatingBasketSuite))
}
