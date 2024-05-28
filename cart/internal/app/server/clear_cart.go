package server

import (
	"fmt"
	"net/http"
	"route256/cart/internal/pkg/model"
	"strconv"
)

func (s *Server) ClearCart(w http.ResponseWriter, r *http.Request) error {
	userIdRaw := r.PathValue("user_id")
	userId, err := strconv.ParseInt(userIdRaw, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt: %w", err)
	}

	if err = s.cartService.ClearCart(model.UserId(userId)); err != nil {
		return fmt.Errorf("s.cartService.ClearCart: %w", err)
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("{}"))
	return nil
}
