package cart

import (
	"errors"
	"fmt"
	"route256/cart/internal/pkg/model"
)

type cartRepository interface {
	AddProduct(userId model.UserId, ProductSku model.ProductSku, count uint16) error
	RemoveProduct(userId model.UserId, ProductSku model.ProductSku) error
	ClearCart(userId model.UserId) error
	GetCart(userId model.UserId) (model.Cart, error)
}

type productService interface {
	GetProduct(ProductSku model.ProductSku) (*model.Product, error)
}

type CartService struct {
	cartRepository cartRepository
	productService productService
}

func NewCartService(cartRepository cartRepository, productService productService) *CartService {
	return &CartService{
		cartRepository: cartRepository,
		productService: productService,
	}
}

func (r *CartService) AddProduct(userId model.UserId, ProductSku model.ProductSku, count uint16) error {
	if userId < 1 || ProductSku < 1 || count < 1 {
		return errors.New("invalid userId or ProductSku or count")
	}
	if _, err := r.productService.GetProduct(ProductSku); err != nil {
		return fmt.Errorf("r.productService.GetProduct: %w", err)
	}
	if err := r.cartRepository.AddProduct(userId, ProductSku, count); err != nil {
		return fmt.Errorf("r.cartRepository.AddProduct: %w", err)
	}
	return nil
}

func (r *CartService) RemoveProduct(userId model.UserId, ProductSku model.ProductSku) error {
	if userId < 1 || ProductSku < 1 {
		return errors.New("invalid userId or ProductSku")
	}
	if err := r.cartRepository.RemoveProduct(userId, ProductSku); err != nil {
		return fmt.Errorf("r.cartRepository.RemoveProduct: %w", err)
	}
	return nil
}

func (r *CartService) ClearCart(userId model.UserId) error {
	if userId < 1 {
		return errors.New("invalid userId")
	}
	if err := r.cartRepository.ClearCart(userId); err != nil {
		return fmt.Errorf("r.cartRepository.ClearCart: %w", err)
	}
	return nil
}

func (r *CartService) GetCart(userId model.UserId) (model.CartFull, error) {
	if userId < 1 {
		return nil, errors.New("invalid userId")
	}

	cart, err := r.cartRepository.GetCart(userId)
	if err != nil {
		return nil, fmt.Errorf("r.cartRepository.ClearCart: %w", err)
	}

	cartFull := make(model.CartFull, len(cart))
	for productSku, count := range cart {
		product, err := r.productService.GetProduct(productSku)
		if err != nil {
			return nil, fmt.Errorf("r.productService.GetProduct: %w", err)
		}
		cartFull[*product] = count
	}

	return cartFull, nil
}
