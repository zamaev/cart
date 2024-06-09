package cart

import (
	"context"
	"route256/cart/internal/pkg/customerror"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/service/cart/mock"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestAddProduct(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	cartRepositoryMock := mock.NewCartRepositoryMock(ctrl)
	productServiceMock := mock.NewProductServiceMock(ctrl)
	cartService := NewCartService(cartRepositoryMock, productServiceMock)

	type mocks struct {
		cartRepositoryMock *mock.CartRepositoryMock
		productServiceMock *mock.ProductServiceMock
	}
	testMocks := mocks{
		cartRepositoryMock,
		productServiceMock,
	}

	testData := []struct {
		name       string
		userId     model.UserId
		productSku model.ProductSku
		count      uint16
		prepare    func(mocks *mocks)
		test       func(err error)
	}{
		{
			name:       "valid params",
			userId:     1,
			productSku: 1,
			count:      1,
			prepare: func(mocks *mocks) {
				mocks.cartRepositoryMock.AddProductMock.Expect(ctx, 1, 1, 1).Return(nil)
				mocks.productServiceMock.GetProductMock.Expect(1).Return(&model.Product{
					Sku:   1111,
					Name:  "Book",
					Price: 100,
				}, nil)
			},
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:       "product not found",
			userId:     1,
			productSku: 30, // invalid sku
			count:      1,
			prepare: func(mocks *mocks) {
				mocks.productServiceMock.GetProductMock.Expect(30).Return(nil, customerror.ErrStatusCode{})
			},
			test: func(err error) {
				assert.ErrorAs(t, err, &customerror.ErrStatusCode{})
			},
		},
		{
			name:       "invalid userId",
			userId:     0,
			productSku: 1,
			count:      1,
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			name:       "invalid productSku",
			userId:     1,
			productSku: 0,
			count:      1,
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			name:       "invalid count",
			userId:     1,
			productSku: 1,
			count:      0,
			test: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(&testMocks)
			}
			err := cartService.AddProduct(ctx, tt.userId, tt.productSku, tt.count)
			tt.test(err)
		})
	}
}

func TestRemoveProduct(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	cartRepositoryMock := mock.NewCartRepositoryMock(ctrl)
	productServiceMock := mock.NewProductServiceMock(ctrl)
	cartService := NewCartService(cartRepositoryMock, productServiceMock)

	testData := []struct {
		name       string
		userId     model.UserId
		productSku model.ProductSku
		test       func(err error)
	}{
		{
			name:       "valid params",
			userId:     1,
			productSku: 3,
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:       "invalid userId",
			userId:     0,
			productSku: 3,
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			name:       "invalid productSku",
			userId:     1,
			productSku: 0,
			test: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			cartRepositoryMock.RemoveProductMock.Expect(ctx, tt.userId, tt.productSku).Return(nil)
			err := cartService.RemoveProduct(ctx, tt.userId, tt.productSku)
			tt.test(err)
		})
	}
}

func TestClearCart(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	cartRepositoryMock := mock.NewCartRepositoryMock(ctrl)
	productServiceMock := mock.NewProductServiceMock(ctrl)
	cartService := NewCartService(cartRepositoryMock, productServiceMock)

	testData := []struct {
		name   string
		userId model.UserId
		test   func(err error)
	}{
		{
			name:   "valid params",
			userId: 1,
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "invalid userId",
			userId: 0,
			test: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			cartRepositoryMock.ClearCartMock.Expect(ctx, tt.userId).Return(nil)
			err := cartService.ClearCart(ctx, tt.userId)
			tt.test(err)
		})
	}
}

func TestGetCart(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	cartRepositoryMock := mock.NewCartRepositoryMock(ctrl)
	productServiceMock := mock.NewProductServiceMock(ctrl)
	cartService := NewCartService(cartRepositoryMock, productServiceMock)

	type mocks struct {
		cartRepositoryMock *mock.CartRepositoryMock
		productServiceMock *mock.ProductServiceMock
	}
	testMocks := mocks{
		cartRepositoryMock,
		productServiceMock,
	}

	testData := []struct {
		name    string
		userId  model.UserId
		prepare func(mocks *mocks)
		test    func(cart model.CartFull, err error)
	}{
		{
			name:   "valid params",
			userId: 1,
			prepare: func(mocks *mocks) {
				mocks.cartRepositoryMock.GetCartMock.Expect(ctx, 1).Return(model.Cart{
					1: 3,
				}, nil)
				mocks.productServiceMock.GetProductMock.Expect(1).Return(&model.Product{
					Sku:   1111,
					Name:  "Book",
					Price: 100,
				}, nil)
			},
			test: func(cart model.CartFull, err error) {
				assert.Equal(t, cart[model.Product{
					Sku:   1111,
					Name:  "Book",
					Price: 100,
				}], uint16(3))
				assert.NoError(t, err)
			},
		},
		{
			name:   "product not found",
			userId: 1,
			prepare: func(mocks *mocks) {
				mocks.cartRepositoryMock.GetCartMock.Expect(ctx, 1).Return(model.Cart{
					1: 1,
				}, nil)
				mocks.productServiceMock.GetProductMock.Expect(1).Return(nil, customerror.ErrStatusCode{})
			},
			test: func(cart model.CartFull, err error) {
				assert.Nil(t, cart)
				assert.ErrorAs(t, err, &customerror.ErrStatusCode{})
			},
		},
		{
			name:   "invalid userId",
			userId: 0,
			test: func(cart model.CartFull, err error) {
				assert.Nil(t, cart)
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(&testMocks)
			}
			cart, err := cartService.GetCart(ctx, tt.userId)
			tt.test(cart, err)
		})
	}
}
