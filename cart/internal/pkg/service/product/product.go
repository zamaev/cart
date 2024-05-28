package product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"route256/cart/internal/pkg/customerror"
	"route256/cart/internal/pkg/model"
)

const URL = "http://route256.pavl.uk:8080"
const DEV_TOKEN = "testtoken"

type ProductService struct {
	token string
}

func NewProductService(token string) *ProductService {
	if token == "" {
		token = DEV_TOKEN
	}
	return &ProductService{
		token: token,
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

func (ps *ProductService) GetProduct(ProductSku model.ProductSku) (*model.Product, error) {
	url := URL + "/get_product"
	getProductRequest := GetProductRequest{
		Token: ps.token,
		Sku:   uint32(ProductSku),
	}
	body, err := json.Marshal(getProductRequest)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	res, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("http.Post: %w", err)
	}
	defer res.Body.Close()

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
