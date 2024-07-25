package model

import (
	"time"
)

type Event struct {
	OrderID OrderID
	Status  OrderStatus
	Time    time.Time
}

type Headers struct {
	TraceID string
}

type OutboxItem struct {
	Id          int64
	Topic       string
	Event       []byte
	Headers     []byte
	CreatedAt   time.Time
	CompletedAt *time.Time
}
