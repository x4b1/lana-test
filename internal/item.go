package checkout

//NewItem given a product and basket creates a new item type with one
func NewItem(p ProductCode, basket BasketID) Item {
	return Item{
		Product:  p,
		Basket:   basket,
		Quantity: 1,
	}
}

//Item defines a item product of a basket in the system
type Item struct {
	Product  ProductCode
	Basket   BasketID
	Quantity int
}

//Add adds 1 to the quantity of the item
func (i *Item) Add() {
	i.Quantity++
}
