package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/utils"
	"route256/cart/pkg/tracing"
	"sort"
)

type GetCartResponseProduct struct {
	SkuId int64  `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

type GetCartResponse struct {
	Items      []GetCartResponseProduct `json:"items"`
	TotalPrice uint32                   `json:"total_price"`
}

func (s *Server) GetCart(w http.ResponseWriter, r *http.Request) (err error) {
	ctx, span := tracing.Start(r.Context(), "server.GetCart")
	defer tracing.EndWithCheckError(span, &err)

	w.Header().Add("Content-Type", "application/json")

	userId, err := utils.GetIntPahtValue(r, "user_id")
	if err != nil {
		return fmt.Errorf("utils.GetIntPahtValue: %w", err)
	}

	cartFull, err := s.cartService.GetCart(ctx, model.UserId(userId))
	if err != nil {
		return fmt.Errorf("s.cartService.ClearCart: %w", err)
	}

	if len(cartFull) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
		return nil
	}

	items := make([]GetCartResponseProduct, 0, len(cartFull))
	totalPrice := uint32(0)
	for product, count := range cartFull {
		items = append(items, GetCartResponseProduct{
			SkuId: int64(product.Sku),
			Name:  product.Name,
			Count: count,
			Price: product.Price,
		})
		totalPrice += product.Price * uint32(count)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].SkuId < items[j].SkuId
	})

	getCartResponse := GetCartResponse{
		Items:      items,
		TotalPrice: uint32(totalPrice),
	}

	data, err := json.Marshal(getCartResponse)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	w.Write(data)
	return nil
}
