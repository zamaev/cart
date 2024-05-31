package cart

import (
	"context"
	"errors"
	"fmt"
	"route256/cart/internal/pkg/model"
)

type cartRepository interface {
	AddProduct(ctx context.Context, userId model.UserId, ProductSku model.ProductSku, count uint16) error
	RemoveProduct(ctx context.Context, userId model.UserId, ProductSku model.ProductSku) error
	ClearCart(ctx context.Context, userId model.UserId) error
	GetCart(ctx context.Context, userId model.UserId) (model.Cart, error)
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

func (r *CartService) AddProduct(ctx context.Context, userId model.UserId, ProductSku model.ProductSku, count uint16) error {
	if userId < 1 || ProductSku < 1 || count < 1 {
		return errors.New("invalid userId or ProductSku or count")
	}
	if _, err := r.productService.GetProduct(ProductSku); err != nil {
		return fmt.Errorf("r.productService.GetProduct: %w", err)
	}
	if err := r.cartRepository.AddProduct(ctx, userId, ProductSku, count); err != nil {
		return fmt.Errorf("r.cartRepository.AddProduct: %w", err)
	}
	return nil
}

func (r *CartService) RemoveProduct(ctx context.Context, userId model.UserId, ProductSku model.ProductSku) error {
	if userId < 1 || ProductSku < 1 {
		return errors.New("invalid userId or ProductSku")
	}
	if err := r.cartRepository.RemoveProduct(ctx, userId, ProductSku); err != nil {
		return fmt.Errorf("r.cartRepository.RemoveProduct: %w", err)
	}
	return nil
}

func (r *CartService) ClearCart(ctx context.Context, userId model.UserId) error {
	if userId < 1 {
		return errors.New("invalid userId")
	}
	if err := r.cartRepository.ClearCart(ctx, userId); err != nil {
		return fmt.Errorf("r.cartRepository.ClearCart: %w", err)
	}
	return nil
}

func (r *CartService) GetCart(ctx context.Context, userId model.UserId) (model.CartFull, error) {
	if userId < 1 {
		return nil, errors.New("invalid userId")
	}

	cart, err := r.cartRepository.GetCart(ctx, userId)
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
