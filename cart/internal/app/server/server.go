package server

import "route256/cart/internal/pkg/model"

type cartService interface {
	AddProduct(userId model.UserId, ProductSku model.ProductSku, count uint16) error
	RemoveProduct(userId model.UserId, ProductSku model.ProductSku) error
	ClearCart(userId model.UserId) error
	GetCart(userId model.UserId) (model.CartFull, error)
}

type Server struct {
	cartService cartService
}

func NewServer(cartService cartService) *Server {
	return &Server{cartService: cartService}
}
