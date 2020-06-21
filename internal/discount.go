package checkout

import (
	"github.com/xabi93/lana-test/pkg/money"
)

//Discount defines a general discount
type Discount interface {
	Apply(itemList map[ProductCode]Item, total money.Money) (money.Money, error)
}

// BuyXGetXDiscount defines the discount of buying x amount of one product getting x free
type BuyXGetXDiscount struct {
	Product Product
	Factor  int
}

//Apply given a list of items and amount returns the total amount with discount applied
func (d BuyXGetXDiscount) Apply(items map[ProductCode]Item, total money.Money) (money.Money, error) {
	i, ok := items[d.Product.Code]
	if !ok {
		return total, nil
	}

	return total.Substract(d.Product.Price.Multiply(i.Quantity / d.Factor))
}

// BulkPurchaseDiscount defines the discount of getting a minimun quantity of products applies x percent discount to that product
type BulkPurchaseDiscount struct {
	Product         Product
	MinQuantity     int
	DiscountPercent int
}

//Apply given a list of items and amount returns the total amount with discount applied
// if the minimun queantity of the discount product has been reached
func (d BulkPurchaseDiscount) Apply(items map[ProductCode]Item, total money.Money) (money.Money, error) {
	i, ok := items[d.Product.Code]
	if !ok {
		return total, nil
	}

	if i.Quantity < d.MinQuantity {
		return total, nil
	}

	return total.Substract(d.Product.Price.Multiply(i.Quantity).Discount(d.DiscountPercent))
}
