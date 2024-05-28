package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"route256/cart/internal/pkg/model"
	"strconv"
)

type AddProductRequest struct {
	Count uint16 `json:"count"`
}

func (s *Server) AddProduct(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "application/json")

	userIdRaw := r.PathValue("user_id")
	userId, err := strconv.ParseInt(userIdRaw, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt userIdRaw: %w", err)
	}

	skuIdRaw := r.PathValue("sku_id")
	skuId, err := strconv.ParseInt(skuIdRaw, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt skuIdRaw: %w", err)
	}

	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	var addProductRequest AddProductRequest
	err = json.Unmarshal(data, &addProductRequest)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	err = s.cartService.AddProduct(
		model.UserId(userId),
		model.ProductSku(skuId),
		addProductRequest.Count,
	)
	if err != nil {
		return fmt.Errorf("s.cartService.AddProduct: %w", err)
	}

	w.Write([]byte("{}"))
	return nil
}
