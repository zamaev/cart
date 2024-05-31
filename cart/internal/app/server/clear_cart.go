package server

import (
	"fmt"
	"net/http"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/utils"
)

func (s *Server) ClearCart(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "application/json")

	userId, err := utils.GetIntPahtValue(r, "user_id")
	if err != nil {
		return fmt.Errorf("utils.GetIntPahtValue: %w", err)
	}

	if err = s.cartService.ClearCart(r.Context(), model.UserId(userId)); err != nil {
		return fmt.Errorf("s.cartService.ClearCart: %w", err)
	}

	w.WriteHeader(http.StatusNoContent)
	utils.SuccessReponse(w)
	return nil
}
