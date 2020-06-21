package money_test

import (
	"testing"

	"github.com/xabi93/lana-test/pkg/money"

	"github.com/stretchr/testify/require"
)

func TestAddMoneySameCurrencySucess(t *testing.T) {
	m, err := money.Eur(1000).Add(money.Eur(50))

	require.Nil(t, err)
	require.Equal(t, 1050, m.Amount)
	require.Equal(t, money.EUR, m.Currency)
}

func TestAddMoneyDifferentCurrencyFails(t *testing.T) {
	_, err := money.Eur(1000).Add(money.Usd(50))

	require.Error(t, err)
}

func TestSubstractMoneySameCurrencySucess(t *testing.T) {
	m, err := money.Eur(1000).Substract(money.Eur(50))

	require.Nil(t, err)
	require.Equal(t, 950, m.Amount)
	require.Equal(t, money.EUR, m.Currency)
}

func TestSubstractMoneySameCurrencyFails(t *testing.T) {
	_, err := money.Eur(1000).Substract(money.Usd(50))

	require.Error(t, err)
}

func TestMultiplyMoneyReturnsCorrectValue(t *testing.T) {
	m := money.Eur(1000).Multiply(2)

	require.Equal(t, 2000, m.Amount)
	require.Equal(t, money.EUR, m.Currency)
}

func TestDivideMoneyReturnsCorrectValue(t *testing.T) {
	m := money.Eur(1050).Divide(2)

	require.Equal(t, 525, m.Amount)
	require.Equal(t, money.EUR, m.Currency)
}

func TestAppliesCorrecltyDiscountToAmount(t *testing.T) {
	m := money.Eur(1000).Discount(25)

	require.Equal(t, 250, m.Amount)
	require.Equal(t, money.EUR, m.Currency)
}

func TestPrintsAmountCorrecltyFormated(t *testing.T) {
	require.Equal(t, "10.50\u20ac", money.Eur(1050).String())
}
