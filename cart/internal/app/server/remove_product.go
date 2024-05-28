package server

import (
	"fmt"
	"net/http"
	"route256/cart/internal/pkg/model"
	"strconv"
)

func (s *Server) RemoveProduct(w http.ResponseWriter, r *http.Request) error {
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

	err = s.cartService.RemoveProduct(
		model.UserId(userId),
		model.ProductSku(skuId),
	)
	if err != nil {
		return fmt.Errorf("s.cartService.RemoveProduct: %w", err)
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("{}"))
	return nil
}
