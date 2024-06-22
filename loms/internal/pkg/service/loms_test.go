package service

import (
	"context"
	"errors"
	"route256/loms/internal/pkg/model"
	"route256/loms/internal/pkg/service/mock"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestOrderCreate(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	stockRepositoryMock := mock.NewStockRepositoryMock(ctrl)
	orderRepositoryMock := mock.NewOrderRepositoryMock(ctrl)
	service := NewLomsService(stockRepositoryMock, orderRepositoryMock)

	type mocks struct {
		stockRepositoryMock *mock.StockRepositoryMock
		orderRepositoryMock *mock.OrderRepositoryMock
	}
	testMocks := mocks{
		stockRepositoryMock,
		orderRepositoryMock,
	}

	testData := []struct {
		name    string
		order   model.Order
		prepare func(mocks *mocks)
		test    func(orderID model.OrderID, err error)
	}{
		{
			name: "valid params",
			order: model.Order{
				User: 1,
				Items: []model.OrderItem{
					{
						Sku:   1,
						Count: 1,
					},
				},
			},
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.CreateMock.Return(1, nil)
				mocks.stockRepositoryMock.ReserveMock.Expect(ctx, 1, 1).Return(nil)
				mocks.orderRepositoryMock.SetStatusMock.Expect(ctx, 1, model.OrderStatusAwaitingPayment).Return(nil)
			},
			test: func(orderID model.OrderID, err error) {
				assert.Equal(t, model.OrderID(1), orderID)
				assert.NoError(t, err)
			},
		},
		{
			name: "stock reserve error",
			order: model.Order{
				User: 1,
				Items: []model.OrderItem{
					{
						Sku:   1,
						Count: 1,
					},
				},
			},
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.CreateMock.Return(1, nil)
				mocks.stockRepositoryMock.ReserveMock.Expect(ctx, 1, 1).Return(errors.New("reserve error"))
				mocks.orderRepositoryMock.SetStatusMock.Expect(ctx, 1, model.OrderStatusFailed).Return(nil)
			},
			test: func(orderID model.OrderID, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "order set status error",
			order: model.Order{
				User: 1,
				Items: []model.OrderItem{
					{
						Sku:   1,
						Count: 1,
					},
				},
			},
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.CreateMock.Return(1, nil)
				mocks.stockRepositoryMock.ReserveMock.Expect(ctx, 1, 1).Return(nil)
				mocks.orderRepositoryMock.SetStatusMock.Expect(ctx, 1, model.OrderStatusAwaitingPayment).Return(errors.New("set status error"))
			},
			test: func(orderID model.OrderID, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "reserver cancel error",
			order: model.Order{
				User: 1,
				Items: []model.OrderItem{
					{
						Sku:   1,
						Count: 1,
					},
					{
						Sku:   2,
						Count: 1,
					},
				},
			},
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.CreateMock.Return(1, nil)
				mocks.stockRepositoryMock.ReserveMock.When(ctx, 1, 1).Then(nil)
				mocks.stockRepositoryMock.ReserveMock.When(ctx, 2, 1).Then(errors.New("reserve error"))
				mocks.orderRepositoryMock.SetStatusMock.Expect(ctx, 1, model.OrderStatusFailed).Return(nil)
				mocks.stockRepositoryMock.ReserveCancelMock.Expect(ctx, 1, 1).Return(errors.New("reserve cancel error"))
			},
			test: func(orderID model.OrderID, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(&testMocks)
			orderID, err := service.OrderCreate(ctx, tt.order)
			tt.test(orderID, err)
		})
	}
}

func TestOrderInfo(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	stockRepositoryMock := mock.NewStockRepositoryMock(ctrl)
	orderRepositoryMock := mock.NewOrderRepositoryMock(ctrl)
	service := NewLomsService(stockRepositoryMock, orderRepositoryMock)

	type mocks struct {
		stockRepositoryMock *mock.StockRepositoryMock
		orderRepositoryMock *mock.OrderRepositoryMock
	}
	testMocks := mocks{
		stockRepositoryMock,
		orderRepositoryMock,
	}

	testData := []struct {
		name    string
		orderID model.OrderID
		prepare func(mocks *mocks)
		test    func(orderID model.Order, err error)
	}{
		{
			name:    "valid params",
			orderID: 1,
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.GetByIdMock.Expect(ctx, 1).Return(model.Order{
					User: 1,
					Items: []model.OrderItem{
						{
							Sku:   1,
							Count: 1,
						},
					},
				}, nil)
			},
			test: func(order model.Order, err error) {
				assert.Equal(t, model.Order{
					User: 1,
					Items: []model.OrderItem{
						{
							Sku:   1,
							Count: 1,
						},
					},
				}, order)
				assert.NoError(t, err)
			},
		},
		{
			name:    "invalid order id",
			orderID: 1,
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.GetByIdMock.Expect(ctx, 1).Return(model.Order{}, errors.New("invalid order id: 1"))
			},
			test: func(order model.Order, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(&testMocks)
			order, err := service.OrderInfo(ctx, tt.orderID)
			tt.test(order, err)
		})
	}
}

func TestOrderPay(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	stockRepositoryMock := mock.NewStockRepositoryMock(ctrl)
	orderRepositoryMock := mock.NewOrderRepositoryMock(ctrl)
	service := NewLomsService(stockRepositoryMock, orderRepositoryMock)

	type mocks struct {
		stockRepositoryMock *mock.StockRepositoryMock
		orderRepositoryMock *mock.OrderRepositoryMock
	}
	testMocks := mocks{
		stockRepositoryMock,
		orderRepositoryMock,
	}

	testData := []struct {
		name    string
		orderID model.OrderID
		prepare func(mocks *mocks)
		test    func(err error)
	}{
		{
			name:    "valid params",
			orderID: 1,
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.GetByIdMock.Expect(ctx, 1).Return(model.Order{
					User:   1,
					Status: model.OrderStatusAwaitingPayment,
					Items: []model.OrderItem{
						{
							Sku:   1,
							Count: 1,
						},
					},
				}, nil)
				mocks.stockRepositoryMock.ReserveRemoveMock.Expect(ctx, 1, 1).Return(nil)
				mocks.orderRepositoryMock.SetStatusMock.Expect(ctx, 1, model.OrderStatusPaid).Return(nil)
			},
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:    "invalid order id",
			orderID: 1,
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.GetByIdMock.Expect(ctx, 1).Return(model.Order{}, errors.New("invalid order id: 1"))
			},
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "invalid status",
			orderID: 1,
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.GetByIdMock.Expect(ctx, 1).Return(model.Order{
					User:   1,
					Status: model.OrderStatusFailed,
					Items: []model.OrderItem{
						{
							Sku:   1,
							Count: 1,
						},
					},
				}, nil)
			},
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "reserve remove error",
			orderID: 1,
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.GetByIdMock.Expect(ctx, 1).Return(model.Order{
					User:   1,
					Status: model.OrderStatusAwaitingPayment,
					Items: []model.OrderItem{
						{
							Sku:   1,
							Count: 1,
						},
					},
				}, nil)
				mocks.stockRepositoryMock.ReserveRemoveMock.Expect(ctx, 1, 1).Return(errors.New("reserve remove error"))
			},
			test: func(err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(&testMocks)
			err := service.OrderPay(ctx, tt.orderID)
			tt.test(err)
		})
	}
}

func TestOrderCancel(t *testing.T) {
	ctx := context.Background()
	ctrl := minimock.NewController(t)
	stockRepositoryMock := mock.NewStockRepositoryMock(ctrl)
	orderRepositoryMock := mock.NewOrderRepositoryMock(ctrl)
	service := NewLomsService(stockRepositoryMock, orderRepositoryMock)

	type mocks struct {
		stockRepositoryMock *mock.StockRepositoryMock
		orderRepositoryMock *mock.OrderRepositoryMock
	}
	testMocks := mocks{
		stockRepositoryMock,
		orderRepositoryMock,
	}

	testData := []struct {
		name    string
		orderID model.OrderID
		prepare func(mocks *mocks)
		test    func(err error)
	}{
		{
			name:    "valid params",
			orderID: 1,
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.GetByIdMock.Expect(ctx, 1).Return(model.Order{
					User:   1,
					Status: model.OrderStatusAwaitingPayment,
					Items: []model.OrderItem{
						{
							Sku:   1,
							Count: 1,
						},
					},
				}, nil)
				mocks.stockRepositoryMock.ReserveCancelMock.Expect(ctx, 1, 1).Return(nil)
				mocks.orderRepositoryMock.SetStatusMock.Expect(ctx, 1, model.OrderStatusCancelled).Return(nil)
			},
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:    "invalid order id",
			orderID: 1,
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.GetByIdMock.Expect(ctx, 1).Return(model.Order{}, errors.New("invalid order id: 1"))
			},
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "invalid status",
			orderID: 1,
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.GetByIdMock.Expect(ctx, 1).Return(model.Order{
					User:   1,
					Status: model.OrderStatusFailed,
					Items: []model.OrderItem{
						{
							Sku:   1,
							Count: 1,
						},
					},
				}, nil)
			},
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "reserve cancel error",
			orderID: 1,
			prepare: func(mocks *mocks) {
				mocks.orderRepositoryMock.GetByIdMock.Expect(ctx, 1).Return(model.Order{
					User:   1,
					Status: model.OrderStatusAwaitingPayment,
					Items: []model.OrderItem{
						{
							Sku:   1,
							Count: 1,
						},
					},
				}, nil)
				mocks.stockRepositoryMock.ReserveCancelMock.Expect(ctx, 1, 1).Return(errors.New("reserve cancel error"))
			},
			test: func(err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(&testMocks)
			err := service.OrderCancel(ctx, tt.orderID)
			tt.test(err)
		})
	}
}
