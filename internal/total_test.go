package checkout_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	checkout "github.com/xabi93/lana-test/internal"
	"github.com/xabi93/lana-test/pkg/money"
)

func TestCalculateTotalReturnsTotal(t *testing.T) {
	product := checkout.Product{
		Code:  checkout.ProductCode("PEN"),
		Name:  checkout.ProductName("Lana Pen"),
		Price: checkout.ProductPrice{money.Eur(500)},
	}

	b := checkout.Basket{
		Currency: checkout.BasketCurrency(money.EUR),
		Items: map[checkout.ProductCode]checkout.Item{
			product.Code: checkout.Item{
				Product:  product.Code,
				Quantity: 2,
			},
		},
	}

	type testCase struct {
		calculator checkout.TotalCalculator
		Expected   money.Money
	}
	testCases := map[string]testCase{
		"with no discounts": testCase{
			calculator: checkout.TotalCalculator{
				Basket: b,
				Products: map[checkout.ProductCode]checkout.Product{
					product.Code: product,
				},
			},
			Expected: money.Eur(1000),
		},
		"with discounts": testCase{
			calculator: checkout.TotalCalculator{
				Basket: b,
				Products: map[checkout.ProductCode]checkout.Product{
					product.Code: product,
				},
				Discounts: []checkout.Discount{new(halfAmountDiscount)},
			},
			Expected: money.Eur(500),
		},
	}
	for name, c := range testCases {
		t.Run(name, func(t *testing.T) {
			amount, err := c.calculator.Calculate()
			require.Nil(t, err)
			require.Equal(t, c.Expected, amount)
		})
	}
}

func TestCalculateNotMatchItemWithProductList(t *testing.T) {
	product := checkout.Product{
		Code:  checkout.ProductCode("PEN"),
		Name:  checkout.ProductName("Lana Pen"),
		Price: checkout.ProductPrice{money.Eur(500)},
	}
	b := checkout.Basket{
		Items: map[checkout.ProductCode]checkout.Item{
			product.Code: checkout.Item{
				Product:  product.Code,
				Quantity: 2,
			},
		},
	}
	calculator := checkout.TotalCalculator{
		Basket: b,
	}

	_, err := calculator.Calculate()

	require.Equal(t, checkout.ErrNotMatchingItemProduct, err)
}

type halfAmountDiscount struct{}

func (d halfAmountDiscount) Apply(itemList map[checkout.ProductCode]checkout.Item, total money.Money) (money.Money, error) {
	return total.Divide(2), nil
}
