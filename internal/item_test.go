package checkout_test

import (
	"testing"

	"github.com/xabi93/lana-test/pkg/uuid"

	"github.com/stretchr/testify/require"

	checkout "github.com/xabi93/lana-test/internal"
)

func TestNewItemReturnsItemInitialized(t *testing.T) {
	product := checkout.ProductCode("PEN")
	basket := checkout.BasketID(uuid.New())

	i := checkout.NewItem(product, basket)

	require.Equal(t, basket, i.Basket)
	require.Equal(t, product, i.Product)
	require.Equal(t, 1, i.Quantity)
}

func TestAddAddsOneToQuantity(t *testing.T) {
	product := checkout.ProductCode("PEN")
	basket := checkout.BasketID(uuid.New())

	i := checkout.NewItem(product, basket)

	i.Add()

	require.Equal(t, 2, i.Quantity)

}
