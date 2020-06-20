package money_test

import (
	"testing"

	"github.com/xabi93/lana-test/pkg/money"

	"github.com/stretchr/testify/require"
)

func TestAddMoneySameCurrencySucess(t *testing.T) {
	m, err := money.Eur(1000).Add(money.Eur(50))

	require.Nil(t, err)
	require.Equal(t, "10.50 ?", m.String())
}

func TestAddMoneyDifferentCurrencyFails(t *testing.T) {
	_, err := money.Eur(1000).Add(money.Usd(50))

	require.Error(t, err)
}

func TestSubstractMoneySameCurrencySucess(t *testing.T) {
	m, err := money.Eur(1000).Substract(money.Eur(50))

	require.Nil(t, err)
	require.Equal(t, "9.50 ?", m.String())
}

func TestSubstractMoneySameCurrencyFails(t *testing.T) {
	_, err := money.Eur(1000).Substract(money.Usd(50))

	require.Error(t, err)
}

func TestMultiplyMoneyReturnsCorrectValue(t *testing.T) {
	m := money.Eur(1000).Multiply(2)

	require.Equal(t, "20.00 ?", m.String())
}

func TestDivideMoneyReturnsCorrectValue(t *testing.T) {
	m := money.Eur(1050).Divide(2)

	require.Equal(t, "5.25 ?", m.String())
}
