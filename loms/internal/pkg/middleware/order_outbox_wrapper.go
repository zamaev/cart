package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"route256/loms/internal/pkg/inrfa/kafka"
	"route256/loms/internal/pkg/model"
	"route256/loms/pkg/tracing"
	"time"
)

type orderRepository interface {
	Create(context.Context, model.Order, ...func(context.Context, model.OrderID, model.OrderStatus) error) (model.OrderID, error)
	GetById(context.Context, model.OrderID) (model.Order, error)
	SetStatus(context.Context, model.OrderID, model.OrderStatus, ...func(context.Context) error) error
	GetAll(context.Context) ([]model.Order, error)
}

type outboxRepository interface {
	Create(ctx context.Context, topic string, event []byte, headers []byte) (int64, error)
	GetWaitList(ctx context.Context) ([]model.OutboxItem, error)
	SetComplete(ctx context.Context, id int64) (err error)
}

type OrderOutboxWrapper struct {
	orderRepository
	outboxRepository outboxRepository
	config           kafka.Config
}

func NewOrderOutboxWrapper(orderRepository orderRepository, outboxRepository outboxRepository, config kafka.Config) *OrderOutboxWrapper {
	return &OrderOutboxWrapper{
		orderRepository:  orderRepository,
		outboxRepository: outboxRepository,
		config:           config,
	}
}

func (r *OrderOutboxWrapper) Create(ctx context.Context, order model.Order) (_ model.OrderID, err error) {
	ctx, span := tracing.Start(ctx, "OrderOutboxWrapper.Create")
	defer tracing.EndWithCheckError(span, &err)
	traceID := ""
	if span != nil {
		traceID = span.SpanContext().TraceID().String()
	}

	inTx := func(ctx context.Context, orderID model.OrderID, orderStatus model.OrderStatus) error {
		eventData, err := json.Marshal(model.Event{OrderID: orderID, Status: orderStatus, Time: time.Now()})
		if err != nil {
			return fmt.Errorf("json.Marshal: %w", err)
		}
		headersData, err := json.Marshal(model.Headers{TraceID: traceID})
		if err != nil {
			return fmt.Errorf("json.Marshal headers: %w", err)
		}

		if _, err := r.outboxRepository.Create(ctx, r.config.OrderEventsTopic, eventData, headersData); err != nil {
			return fmt.Errorf("r.outboxRepository.Create: %w", err)
		}
		return nil
	}
	return r.orderRepository.Create(ctx, order, inTx)
}

func (r *OrderOutboxWrapper) SetStatus(ctx context.Context, orderID model.OrderID, status model.OrderStatus) (err error) {
	ctx, span := tracing.Start(ctx, "OrderOutboxWrapper.SetStatus")
	defer tracing.EndWithCheckError(span, &err)
	traceID := ""
	if span != nil {
		traceID = span.SpanContext().TraceID().String()
	}

	inTx := func(ctx context.Context) error {
		eventData, err := json.Marshal(model.Event{OrderID: orderID, Status: status, Time: time.Now()})
		if err != nil {
			return fmt.Errorf("json.Marshal event: %w", err)
		}
		headersData, err := json.Marshal(model.Headers{TraceID: traceID})
		if err != nil {
			return fmt.Errorf("json.Marshal headers: %w", err)
		}

		if _, err := r.outboxRepository.Create(ctx, r.config.OrderEventsTopic, eventData, headersData); err != nil {
			return fmt.Errorf("outboxRepository.Create: %w", err)
		}
		return nil
	}
	return r.orderRepository.SetStatus(ctx, orderID, status, inTx)
}
