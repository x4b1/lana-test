package checkout

import (
	"errors"

	"github.com/xabi93/lana-test/pkg/money"
)

var (
	ErrProductNotExists = errors.New("Product does not exists")
)

//ProductCode defines the unique id for a product
type ProductCode string

//ProductName defines the name of a product
type ProductName string

//ProductPrice defines the price of a product
type ProductPrice money.Money

type Product struct {
	Code  ProductCode
	Name  ProductName
	Price ProductPrice
}
