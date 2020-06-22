package checkout

import (
	"errors"

	"github.com/xabi93/lana-test/pkg/money"
	"github.com/xabi93/lana-test/pkg/uuid"
)

var (
	ErrBasketNotExists = errors.New("Basket does not exists")
)

//NewBasket returns a new basket initialized with default data
func NewBasket(curr BasketCurrency) Basket {
	return Basket{
		ID:       BasketID(uuid.New()),
		Currency: curr,
		Items:    make(map[ProductCode]Item),
	}
}

//BasketID defines the unique id for a basket
type BasketID string

//BasketCurrency defines the currency of a basket
type BasketCurrency money.Currency

//Basket defines a basket in the system
type Basket struct {
	ID       BasketID
	Currency BasketCurrency
	Items    map[ProductCode]Item
}

//AddItem given a product adds a new item to the basket
func (b *Basket) AddItem(p ProductCode) {
	if i, ok := b.Items[p]; ok {
		i.Add()
		b.Items[p] = i
		return
	}

	b.Items[p] = NewItem(p, b.ID)
}

var ErrNotMatchingItemProduct = errors.New("not matching item with product")

//BasketTotalWithDiscounts calculates the total amount of a basket and then applies the defined discounts
func BasketTotalWithDiscounts(b Basket, products map[ProductCode]Product, discounts []Discount) (money.Money, error) {
	total, err := basketTotal(b, products)
	if err != nil {
		return total, err
	}

	discount, err := calculateTotalDiscount(b, discounts)
	if err != nil {
		return money.Money{}, err
	}

	if discount.IsZero() {
		return total, nil
	}

	return total.Substract(discount)
}

//basketTotal calculates the total amount of a basket
func basketTotal(b Basket, products map[ProductCode]Product) (money.Money, error) {
	var err error
	total := money.FromCurrency(money.Currency(b.Currency))

	for _, i := range b.Items {
		p, ok := products[i.Product]
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

//basketAmountwithDiscounts given the total amount
func calculateTotalDiscount(b Basket, discounts []Discount) (money.Money, error) {
	total := money.FromCurrency(money.Currency(b.Currency))
	var err error

	for _, d := range discounts {
		discount := d.Calculate(b.Items)
		if !discount.IsZero() {
			total, err = total.Add(discount)
			if err != nil {
				return money.Money{}, err
			}
		}
	}

	return total, nil
}
