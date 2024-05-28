package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"route256/cart/internal/pkg/model"
	"sort"
	"strconv"
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

func (s *Server) GetCart(w http.ResponseWriter, r *http.Request) error {
	userIdRaw := r.PathValue("user_id")
	userId, err := strconv.ParseInt(userIdRaw, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt: %w", err)
	}

	cartFull, err := s.cartService.GetCart(model.UserId(userId))
	if err != nil {
		return fmt.Errorf("s.cartService.ClearCart: %w", err)
	}

	if len(cartFull) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
		return nil
	}

	var getCartResponse GetCartResponse
	for product, count := range cartFull {
		getCartResponse.Items = append(getCartResponse.Items, GetCartResponseProduct{
			SkuId: int64(product.Sku),
			Name:  product.Name,
			Count: count,
			Price: product.Price,
		})
		getCartResponse.TotalPrice += product.Price * uint32(count)
	}

	sort.Slice(getCartResponse.Items, func(i, j int) bool {
		return getCartResponse.Items[i].SkuId < getCartResponse.Items[j].SkuId
	})

	data, err := json.Marshal(getCartResponse)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	w.Write(data)
	return nil
}
