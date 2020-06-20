package checkout

import (
	"errors"

	"github.com/xabi93/lana-test/pkg/uuid"
)

var (
	ErrBasketNotExists = errors.New("Basket does not exists")
)

//NewBasket returns a new basket initialized with default data
func NewBasket() Basket {
	return Basket{
		ID:    BasketID(uuid.New()),
		Items: make(map[ProductCode]Item),
	}
}

//BasketID defines the unique id for a basket
type BasketID string

//Basket defines a basket in the system
type Basket struct {
	ID    BasketID
	Items map[ProductCode]Item
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
