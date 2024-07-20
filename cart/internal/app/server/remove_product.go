package server

import (
	"fmt"
	"net/http"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/utils"
	"route256/cart/pkg/tracing"
)

func (s *Server) RemoveProduct(w http.ResponseWriter, r *http.Request) (err error) {
	ctx, span := tracing.Start(r.Context(), "server.RemoveProduct")
	defer tracing.EndWithCheckError(span, &err)

	w.Header().Add("Content-Type", "application/json")

	userId, err := utils.GetIntPahtValue(r, "user_id")
	if err != nil {
		return fmt.Errorf("utils.GetIntPahtValue: %w", err)
	}

	skuId, err := utils.GetIntPahtValue(r, "sku_id")
	if err != nil {
		return fmt.Errorf("utils.GetIntPahtValue: %w", err)
	}

	err = s.cartService.RemoveProduct(
		ctx,
		model.UserId(userId),
		model.ProductSku(skuId),
	)
	if err != nil {
		return fmt.Errorf("s.cartService.RemoveProduct: %w", err)
	}

	w.WriteHeader(http.StatusNoContent)
	utils.SuccessReponse(w)
	return nil
}
