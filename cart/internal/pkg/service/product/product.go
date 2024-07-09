package product

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"route256/cart/internal/pkg/config"
	"route256/cart/internal/pkg/customerror"
	"route256/cart/internal/pkg/middleware"
	"route256/cart/internal/pkg/model"
	"route256/cart/pkg/tracing"

	"route256/cart/internal/pkg/utils/metrics"
	"strconv"
	"time"
)

type ProductService struct {
	url   string
	token string
}

func NewProductService(config config.Config) *ProductService {
	return &ProductService{
		url:   config.ProductServiceUrl,
		token: config.ProductServiceToken,
	}
}

type GetProductRequest struct {
	Token string `json:"token"`
	Sku   uint32 `json:"sku"`
}

type GetProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func (ps *ProductService) GetProduct(ctx context.Context, ProductSku model.ProductSku) (_ *model.Product, err error) {
	ctx, span := tracing.Start(ctx, "ProductService.GetProduct")
	defer tracing.EndWithCheckError(span, &err)

	url := ps.url + "/get_product"
	getProductRequest := GetProductRequest{
		Token: ps.token,
		Sku:   uint32(ProductSku),
	}
	body, err := json.Marshal(getProductRequest)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	metrics.ExternalRequestCounter(url)
	start := time.Now()

	res, err := middleware.NewRetryClient().Post(ctx, url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("http.Post: %w", err)
	}
	defer res.Body.Close()

	metrics.ExternalRequestDuration(url, strconv.Itoa(res.StatusCode), time.Since(start).Seconds())

	if res.StatusCode != http.StatusOK {
		return nil, customerror.NewErrStatusCode(
			fmt.Sprintf("sku %v not found", ProductSku),
			http.StatusPreconditionFailed,
		)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	var getProductResponse GetProductResponse
	err = json.Unmarshal(data, &getProductResponse)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return &model.Product{
		Sku:   ProductSku,
		Name:  getProductResponse.Name,
		Price: getProductResponse.Price,
	}, nil
}
