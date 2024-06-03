package server

import (
	"context"
	"route256/cart/internal/pkg/model"
)

type cartService interface {
	AddProduct(ctx context.Context, userId model.UserId, ProductSku model.ProductSku, count uint16) error
	RemoveProduct(ctx context.Context, userId model.UserId, ProductSku model.ProductSku) error
	ClearCart(ctx context.Context, userId model.UserId) error
	GetCart(ctx context.Context, userId model.UserId) (model.CartFull, error)
}

type Server struct {
	cartService cartService
}

func NewServer(cartService cartService) *Server {
	return &Server{cartService: cartService}
}
