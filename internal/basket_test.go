package checkout_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	checkout "github.com/xabi93/lana-test/internal"
)

func TestNewBasketReturnsBasketInitialized(t *testing.T) {
	b := checkout.NewBasket()

	require.NotEmpty(t, b.ID)
	require.Empty(t, b.Items)
}

func TestAddMultipleItemShouldAddCorrectly(t *testing.T) {
	b := checkout.NewBasket()

	pen := checkout.ProductCode("PEN")
	mug := checkout.ProductCode("MUG")

	b.AddItem(pen)
	b.AddItem(pen)
	b.AddItem(mug)

	require.Len(t, b.Items, 2)
	require.Equal(t, 2, b.Items[pen].Quantity)
	require.Equal(t, 1, b.Items[mug].Quantity)
}
