package cart

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"route256/cart/internal/pkg/customerror"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/utils"
	"time"
)

const rps = 10

type cartRepository interface {
	AddProduct(ctx context.Context, userId model.UserId, ProductSku model.ProductSku, count uint16) error
	RemoveProduct(ctx context.Context, userId model.UserId, ProductSku model.ProductSku) error
	ClearCart(ctx context.Context, userId model.UserId) error
	GetCart(ctx context.Context, userId model.UserId) (model.Cart, error)
	GetProductCount(ctx context.Context, userId model.UserId, ProductSku model.ProductSku) (uint16, error)
}

type productService interface {
	GetProduct(ctx context.Context, ProductSku model.ProductSku) (*model.Product, error)
}

type lomsService interface {
	OrderCreate(ctx context.Context, user model.UserId, cart model.Cart) (model.OrderId, error)
	StocksInfo(ctx context.Context, sku model.ProductSku) (uint64, error)
}

type CartService struct {
	cartRepository cartRepository
	productService productService
	lomsService    lomsService
}

func NewCartService(cartRepository cartRepository, productService productService, lomsService lomsService) *CartService {
	return &CartService{
		cartRepository: cartRepository,
		productService: productService,
		lomsService:    lomsService,
	}
}

func (r *CartService) AddProduct(ctx context.Context, userId model.UserId, ProductSku model.ProductSku, count uint16) error {
	if userId < 1 || ProductSku < 1 || count < 1 {
		return errors.New("invalid userId or ProductSku or count")
	}
	if _, err := r.productService.GetProduct(ctx, ProductSku); err != nil {
		return fmt.Errorf("r.productService.GetProduct: %w", err)
	}
	cartProductCount, err := r.cartRepository.GetProductCount(ctx, userId, ProductSku)
	if err != nil {
		return fmt.Errorf("r.cartRepository.GetCart: %w", err)
	}
	if stockCount, err := r.lomsService.StocksInfo(ctx, ProductSku); err != nil {
		return fmt.Errorf("r.lomsService.StocksInfo: %w", err)
	} else if stockCount < uint64(cartProductCount+count) {
		return customerror.NewErrStatusCode("not enough products in stock", http.StatusPreconditionFailed)
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

	eg, ctx := utils.NewErrGroup(ctx)
	cartFullMx := model.NewCartFullMx(len(cart))
	
	period := time.NewTicker(time.Second / rps)
	defer period.Stop()

	for productSku, count := range cart {
		eg.Go(func() error {
			<-period.C
			product, err := r.productService.GetProduct(ctx, productSku)
			if err != nil {
				return fmt.Errorf("r.productService.GetProduct: %w", err)
			}
			cartFullMx.Add(*product, count)
			return nil
		})
	}
	
	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("eg.Wait: %w", err)
	}

	return cartFullMx.GetCartFull(), nil
}

func (r *CartService) Checkout(ctx context.Context, userId model.UserId) (model.OrderId, error) {
	if userId < 1 {
		return 0, errors.New("invalid userId")
	}
	cart, err := r.cartRepository.GetCart(ctx, userId)
	if err != nil {
		return 0, fmt.Errorf("r.GetCart: %w", err)
	}
	orderId, err := r.lomsService.OrderCreate(ctx, userId, cart)
	if err != nil {
		return 0, fmt.Errorf("r.lomsService.OrderCreate: %w", err)
	}
	if err := r.cartRepository.ClearCart(ctx, userId); err != nil {
		return 0, fmt.Errorf("r.cartRepository.ClearCart: %w", err)
	}
	return orderId, nil
}
