package cmd

import (
	"context"
	"route256/cart/internal/pkg/config"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/service/cart"
	"route256/cart/internal/pkg/service/cart/mock"
	"route256/cart/internal/pkg/service/product"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestGetCart(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx := context.Background()
	config := config.NewConfig()
	ctrl := minimock.NewController(t)

	lomsServiceMock := mock.NewLomsServiceMock(ctrl)
	cartRepository := mock.NewCartRepositoryMock(ctrl)
	productService := product.NewProductService(config)
	cartService := cart.NewCartService(cartRepository, productService, lomsServiceMock)

	testData := []struct {
		name string
		cart model.Cart
		test func(error)
	}{
		{
			name: "valid params",
			cart: model.Cart{
				1076963: 1,
				1148162: 1,
				1625903: 1,
			},
			test: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "invalid ProductSku in the begining",
			cart: model.Cart{
				1:       1, // wrong
				1076963: 1,
				1148162: 1,
				1625903: 1,
				2618151: 1,
				2956315: 1,
				2958025: 1,
				3596599: 1,
				3618852: 1,
				4288068: 1,
				4465995: 1,
				4487693: 1,
				4669069: 1,
				4678287: 1,
				4678816: 1,
				4679011: 1,
				4687693: 1,
				4996014: 1,
				5097510: 1,
				5415913: 1,
				5647362: 1,
			},
			test: func(err error) {
				require.Error(t, err)
			},
		},
		{
			name: "invalid ProductSku in the middle",
			cart: model.Cart{
				1076963: 1,
				1148162: 1,
				1625903: 1,
				2618151: 1,
				2956315: 1,
				2958025: 1,
				3596599: 1,
				3618852: 1,
				4288068: 1,
				4465995: 1,
				1:       1, // wrong
				4487693: 1,
				4669069: 1,
				4678287: 1,
				4678816: 1,
				4679011: 1,
				4687693: 1,
				4996014: 1,
				5097510: 1,
				5415913: 1,
				5647362: 1,
			},
			test: func(err error) {
				require.Error(t, err)
			},
		},
		{
			name: "invalid ProductSku in the end",
			cart: model.Cart{
				1076963: 1,
				1148162: 1,
				1625903: 1,
				2618151: 1,
				2956315: 1,
				2958025: 1,
				3596599: 1,
				3618852: 1,
				4288068: 1,
				4465995: 1,
				4487693: 1,
				4669069: 1,
				4678287: 1,
				4678816: 1,
				4679011: 1,
				4687693: 1,
				4996014: 1,
				5097510: 1,
				5415913: 1,
				5647362: 1,
				1:       1, // wrong
			},
			test: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			cartRepository.GetCartMock.Return(tt.cart, nil)
			_, err := cartService.GetCart(ctx, 1)
			tt.test(err)
		})
	}
}
