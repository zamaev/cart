package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/utils"

	"github.com/go-playground/validator/v10"
)

type AddProductRequest struct {
	Count uint16 `json:"count" validate:"gt=0"`
}

func (s *Server) AddProduct(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "application/json")

	userId, err := utils.GetIntPahtValue(r, "user_id")
	if err != nil {
		return fmt.Errorf("utils.GetIntPahtValue: %w", err)
	}

	skuId, err := utils.GetIntPahtValue(r, "sku_id")
	if err != nil {
		return fmt.Errorf("utils.GetIntPahtValue: %w", err)
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

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(addProductRequest)
	if err != nil {
		return fmt.Errorf("validation addProductRequest: %w", err)
	}

	err = s.cartService.AddProduct(
		r.Context(),
		model.UserId(userId),
		model.ProductSku(skuId),
		addProductRequest.Count,
	)
	if err != nil {
		return fmt.Errorf("s.cartService.AddProduct: %w", err)
	}

	utils.SuccessReponse(w)
	return nil
}
