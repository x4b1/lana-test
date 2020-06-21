package main

import (
	"os"

	"github.com/xabi93/lana-test/internal/adding"
	"github.com/xabi93/lana-test/internal/calculating"
	"github.com/xabi93/lana-test/internal/creating"
	"github.com/xabi93/lana-test/internal/deleting"
	"github.com/xabi93/lana-test/internal/server"
	"github.com/xabi93/lana-test/internal/storage/memory"
	"github.com/xabi93/lana-test/pkg/log/logrus"
)

func main() {
	db := memory.NewDB()

	baskets := memory.NewBasketRepository(db)
	products := memory.NewProductRepository(db)
	discounts := memory.NewDiscountRepository(db)

	s := server.New(
		logrus.New(),
		creating.NewBasketCreator(baskets),
		adding.NewBasketItemsAdder(baskets, products),
		deleting.NewBasketDeleter(baskets),
		calculating.NewTotalBasket(baskets, products, discounts),
	)

	s.Run(os.Getenv("PORT"))
}
