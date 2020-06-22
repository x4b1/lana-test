package checkout_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	checkout "github.com/xabi93/lana-test/internal"
	"github.com/xabi93/lana-test/pkg/money"
)

func TestBuyXGetXDiscountGivenAListOfItemsCalculatesTotalDiscount(t *testing.T) {
	product := checkout.Product{
		Code:  checkout.ProductCode("PEN"),
		Name:  checkout.ProductName("Lana Pen"),
		Price: checkout.ProductPrice{money.Eur(500)},
	}
	type testCase struct {
		items    map[checkout.ProductCode]checkout.Item
		amount   money.Money
		Expected money.Money
	}
	testCases := map[string]testCase{
		"matched criteria once": testCase{
			items: map[checkout.ProductCode]checkout.Item{
				product.Code: checkout.Item{
					Product:  product.Code,
					Quantity: 3,
				},
				"TSHIRT": checkout.Item{
					Product:  "TSHIRT",
					Quantity: 1,
				},
			},
			amount:   money.Eur(3500),
			Expected: money.Eur(500),
		},
		"matched criteria twice": testCase{
			items: map[checkout.ProductCode]checkout.Item{
				product.Code: checkout.Item{
					Product:  product.Code,
					Quantity: 4,
				},
				"TSHIRT": checkout.Item{
					Product:  "TSHIRT",
					Quantity: 1,
				},
			},
			amount:   money.Eur(4000),
			Expected: money.Eur(1000),
		},
		"not match criteria": testCase{
			items: map[checkout.ProductCode]checkout.Item{
				"TSHIRT": checkout.Item{
					Product:  "TSHIRT",
					Quantity: 1,
				},
			},
			amount:   money.Eur(2000),
			Expected: money.Money{},
		},
	}
	d := checkout.BuyXGetXDiscount{
		Product: product,
		Factor:  2,
	}
	for name, c := range testCases {
		t.Run(name, func(t *testing.T) {
			amount := d.Calculate(c.items)

			require.Equal(t, c.Expected, amount)
		})
	}
}

func TestBulkPurchaseDiscountGivenAListOfItemsCalculatesTotalDiscount(t *testing.T) {
	product := checkout.Product{
		Code:  checkout.ProductCode("TSHIRT"),
		Name:  checkout.ProductName("Lana T-Shirt"),
		Price: checkout.ProductPrice{money.Eur(2000)},
	}
	type testCase struct {
		items    map[checkout.ProductCode]checkout.Item
		amount   money.Money
		Expected money.Money
	}
	testCases := map[string]testCase{
		"matched minimun quantity": testCase{
			items: map[checkout.ProductCode]checkout.Item{
				product.Code: checkout.Item{
					Product:  product.Code,
					Quantity: 4,
				},
				"PEN": checkout.Item{
					Product:  "Lana pen",
					Quantity: 1,
				},
			},
			amount:   money.Eur(8500),
			Expected: money.Eur(2000),
		},
		"not match criteria": testCase{
			items: map[checkout.ProductCode]checkout.Item{
				"PEN": checkout.Item{
					Product:  "PEN",
					Quantity: 1,
				},
			},
			amount:   money.Eur(500),
			Expected: money.Money{},
		},
	}
	d := checkout.BulkPurchaseDiscount{
		Product:         product,
		MinQuantity:     3,
		DiscountPercent: 25,
	}
	for name, c := range testCases {
		t.Run(name, func(t *testing.T) {
			amount := d.Calculate(c.items)

			require.Equal(t, c.Expected, amount)
		})
	}
}
