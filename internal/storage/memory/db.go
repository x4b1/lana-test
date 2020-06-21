package memory

import (
	checkout "github.com/xabi93/lana-test/internal"
	"github.com/xabi93/lana-test/pkg/money"
)

func NewDB() *DB {
	pen := checkout.Product{
		Code:  checkout.ProductCode("PEN"),
		Name:  checkout.ProductName("Lana Pen"),
		Price: checkout.ProductPrice{Money: money.Eur(500)},
	}
	tshirt := checkout.Product{
		Code:  checkout.ProductCode("TSHIRT"),
		Name:  checkout.ProductName("Lana T-Shirt"),
		Price: checkout.ProductPrice{Money: money.Eur(2000)},
	}
	mug := checkout.Product{
		Code:  checkout.ProductCode("MUG"),
		Name:  checkout.ProductName("Lana Coffee Mug"),
		Price: checkout.ProductPrice{Money: money.Eur(750)},
	}

	return &DB{
		baskets: make(map[checkout.BasketID]checkout.Basket),
		products: map[checkout.ProductCode]checkout.Product{
			pen.Code:    pen,
			tshirt.Code: tshirt,
			mug.Code:    mug,
		},
		discounts: []checkout.Discount{
			checkout.BuyXGetXDiscount{
				Product: pen,
				Factor:  2,
			},
			checkout.BulkPurchaseDiscount{
				Product:         tshirt,
				MinQuantity:     3,
				DiscountPercent: 25,
			},
		},
	}
}

type DB struct {
	baskets   map[checkout.BasketID]checkout.Basket
	products  map[checkout.ProductCode]checkout.Product
	discounts []checkout.Discount
}
