package checkout_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	checkout "github.com/xabi93/lana-test/internal"
	"github.com/xabi93/lana-test/pkg/money"
)

func TestNewBasketReturnsBasketInitialized(t *testing.T) {
	curr := checkout.BasketCurrency(money.EUR)
	b := checkout.NewBasket(curr)

	require.NotEmpty(t, b.ID)
	require.Empty(t, b.Items)
	require.Equal(t, curr, b.Currency)
}

func TestAddMultipleItemShouldAddCorrectly(t *testing.T) {
	b := checkout.NewBasket(checkout.BasketCurrency(money.EUR))

	pen := checkout.ProductCode("PEN")
	mug := checkout.ProductCode("MUG")

	b.AddItem(pen)
	b.AddItem(pen)
	b.AddItem(mug)

	require.Len(t, b.Items, 2)
	require.Equal(t, 2, b.Items[pen].Quantity)
	require.Equal(t, 1, b.Items[mug].Quantity)
}

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
		basket    checkout.Basket
		products  map[checkout.ProductCode]checkout.Product
		discounts []checkout.Discount
		Expected  money.Money
	}
	testCases := map[string]testCase{
		"with no discounts": testCase{
			basket: b,
			products: map[checkout.ProductCode]checkout.Product{
				product.Code: product,
			},
			Expected: money.Eur(1000),
		},
		"with discounts": testCase{
			basket: b,
			products: map[checkout.ProductCode]checkout.Product{
				product.Code: product,
			},
			discounts: []checkout.Discount{new(oneEurDiscount)},
			Expected:  money.Eur(900),
		},
	}
	for name, c := range testCases {
		t.Run(name, func(t *testing.T) {
			amount, err := checkout.BasketTotalWithDiscounts(c.basket, c.products, c.discounts)
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

	_, err := checkout.BasketTotalWithDiscounts(b, nil, nil)

	require.Equal(t, checkout.ErrNotMatchingItemProduct, err)
}

type oneEurDiscount struct{}

func (d oneEurDiscount) Calculate(itemList map[checkout.ProductCode]checkout.Item) money.Money {
	return money.Eur(100)
}
