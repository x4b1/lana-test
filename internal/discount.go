package checkout

import (
	"github.com/xabi93/lana-test/pkg/money"
)

//Discount defines a general discount
type Discount interface {
	Calculate(itemList map[ProductCode]Item) money.Money
}

// BuyXGetXDiscount defines the discount of buying x amount of one product getting x free
type BuyXGetXDiscount struct {
	Product Product
	Factor  int
}

//Calculate given a list of items and amount returns the total amount with discount applied
func (d BuyXGetXDiscount) Calculate(items map[ProductCode]Item) money.Money {
	i, ok := items[d.Product.Code]
	if !ok {
		return money.Money{}
	}

	return d.Product.Price.Multiply(i.Quantity / d.Factor)
}

// BulkPurchaseDiscount defines the discount of getting a minimun quantity of products applies x percent discount to that product
type BulkPurchaseDiscount struct {
	Product         Product
	MinQuantity     int
	DiscountPercent int
}

//Calculate given a list of items and amount returns the total amount with discount applied
// if the minimun quantity of the discount product has been reached
func (d BulkPurchaseDiscount) Calculate(items map[ProductCode]Item) money.Money {
	i, ok := items[d.Product.Code]
	if !ok {
		return money.Money{}
	}

	if i.Quantity < d.MinQuantity {
		return money.Money{}
	}

	return d.Product.Price.Multiply(i.Quantity).Discount(d.DiscountPercent)
}
