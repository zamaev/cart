package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"route256/cart/internal/pkg/model"
	"route256/cart/pkg/tracing"
)

type CheckoutRequest struct {
	UserId int64 `json:"user" validate:"required"`
}

type CheckoutResponse struct {
	OrderId int64 `json:"order_id"`
}

func (s *Server) Checkout(w http.ResponseWriter, r *http.Request) (err error) {
	ctx, span := tracing.Start(r.Context(), "server.Checkout")
	defer tracing.EndWithCheckError(span, &err)

	w.Header().Add("Content-Type", "application/json")

	reqData, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	var checkoutRequest CheckoutRequest
	err = json.Unmarshal(reqData, &checkoutRequest)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	orderId, err := s.cartService.Checkout(ctx, model.UserId(checkoutRequest.UserId))
	if err != nil {
		return fmt.Errorf("s.cartService.Checkout: %w", err)
	}

	checkoutRespone := CheckoutResponse{
		OrderId: int64(orderId),
	}

	resData, err := json.Marshal(checkoutRespone)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	w.Write(resData)
	return nil
}
