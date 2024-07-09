package loms

import (
	"context"
	"fmt"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/utils/metrics"
	"route256/cart/pkg/api/loms/v1"
	"route256/cart/pkg/tracing"
	"time"
)

type LomsService struct {
	client loms.LomsClient
}

func NewLomsService(client loms.LomsClient) *LomsService {
	return &LomsService{
		client: client,
	}
}

func (s *LomsService) OrderCreate(ctx context.Context, user model.UserId, cart model.Cart) (_ model.OrderId, err error) {
	ctx, span := tracing.Start(ctx, "LomsService.OrderCreate")
	defer tracing.EndWithCheckError(span, &err)

	metrics.ExternalRequestCounter("loms.OrderCreate")
	defer func(start time.Time) {
		metrics.ExternalRequestDurationWithError("loms.OrderCreate", err, time.Since(start).Seconds())
	}(time.Now())

	req := loms.OrderCreateRequest{
		User:  int64(user),
		Items: make([]*loms.OrderItem, 0, len(cart)),
	}
	for productSku, count := range cart {
		req.Items = append(req.Items, &loms.OrderItem{
			Sku:   uint32(productSku),
			Count: uint32(count),
		})
	}
	res, err := s.client.OrderCreate(ctx, &req)
	if err != nil {
		return model.OrderId(0), fmt.Errorf("lomsClient.OrderCreate: %w", err)
	}
	return model.OrderId(res.OrderId), nil
}

func (s *LomsService) StocksInfo(ctx context.Context, sku model.ProductSku) (_ uint64, err error) {
	ctx, span := tracing.Start(ctx, "LomsService.StocksInfo")
	defer tracing.EndWithCheckError(span, &err)

	metrics.ExternalRequestCounter("loms.StocksInfo")
	defer func(start time.Time) {
		metrics.ExternalRequestDurationWithError("loms.StocksInfo", err, time.Since(start).Seconds())
	}(time.Now())

	res, err := s.client.StocksInfo(ctx, &loms.StocksInfoRequest{
		Sku: uint32(sku),
	})
	if err != nil {
		return 0, fmt.Errorf("lomsClient.StocksInfo: %w", err)
	}
	return res.Count, nil
}
