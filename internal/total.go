package checkout

import (
	"errors"

	"github.com/xabi93/lana-test/pkg/money"
)

var ErrNotMatchingItemProduct = errors.New("Not matching item with product")

//TotalCalculator defines the required data that needs system to calculate total amount of a basket applying discounts
type TotalCalculator struct {
	Basket    Basket
	Discounts []Discount
	Products  map[ProductCode]Product
}

//Calculate calculates the total amount of a basket and then applies the defines discounts
func (tc TotalCalculator) Calculate() (money.Money, error) {
	total, err := tc.total()
	if err != nil {
		return total, err
	}

	return tc.applyDiscounts(total)
}

func (tc TotalCalculator) total() (money.Money, error) {
	var err error
	total := money.Eur(0)

	for _, i := range tc.Basket.Items {
		p, ok := tc.Products[i.Product]
		if !ok {
			return money.Money{}, ErrNotMatchingItemProduct
		}

		total, err = total.Add(p.Price.Multiply(i.Quantity))
		if err != nil {
			return money.Money{}, err
		}
	}

	return total, nil
}

func (tc TotalCalculator) applyDiscounts(total money.Money) (money.Money, error) {
	var err error

	for _, d := range tc.Discounts {
		total, err = d.Apply(tc.Basket.Items, total)
		if err != nil {
			return money.Money{}, err
		}
	}

	return total, nil
}
