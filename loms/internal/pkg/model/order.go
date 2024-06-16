package model

type OrderID int64

type OrderStatus string

const (
	OrderStatusNone            OrderStatus = ""
	OrderStatusNew             OrderStatus = "new"
	OrderStatusAwaitingPayment OrderStatus = "awaiting payment"
	OrderStatusFailed          OrderStatus = "failed"
	OrderStatusPaid            OrderStatus = "payed"
	OrderStatusCancelled       OrderStatus = "cancelled"
)

type OrderItem struct {
	Sku   ProductSku
	Count uint16
}

type Order struct {
	Status OrderStatus
	User   UserID
	Items  []OrderItem
}
