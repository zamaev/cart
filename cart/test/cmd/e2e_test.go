package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"route256/cart/internal/app/server"
	"route256/cart/internal/pkg/config"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

const contentType = "application/json"

func TestServer(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx := context.Background()

	app := server.NewApp(ctx, config.NewConfig())
	serverApp := httptest.NewServer(app.Handler)
	defer serverApp.Close()

	t.Run("AddProduct", func(t *testing.T) {
		AddProduct(serverApp, t)
	})
	t.Run("Checkout", func(t *testing.T) {
		Checkout(serverApp, t)
	})
	t.Run("AddProduct", func(t *testing.T) {
		AddProduct(serverApp, t)
	})
	t.Run("GetCart", func(t *testing.T) {
		GetCart(serverApp, t)
	})
	t.Run("RemoveProduct", func(t *testing.T) {
		RemoveProduct(serverApp, t)
	})

	app.Shutdown(ctx)
}

func AddProduct(serverApp *httptest.Server, t *testing.T) {
	addProductRequest := server.AddProductRequest{Count: 1}
	data, err := json.Marshal(addProductRequest)
	require.NoError(t, err)

	resp, err := http.Post(serverApp.URL+"/user/31337/cart/1076963", contentType, bytes.NewBuffer(data))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Checkout(serverApp *httptest.Server, t *testing.T) {
	checkoutRequest := server.CheckoutRequest{UserId: 31337}
	data, err := json.Marshal(checkoutRequest)
	require.NoError(t, err)

	resp, err := http.Post(serverApp.URL+"/cart/checkout", contentType, bytes.NewBuffer(data))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func GetCart(serverApp *httptest.Server, t *testing.T) {
	resp, err := http.Get(serverApp.URL + "/user/31337/cart/list")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var getCartResponse server.GetCartResponse
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(data, &getCartResponse)
	require.NoError(t, err)

	require.Equal(t, uint32(3379), getCartResponse.TotalPrice)
}

func RemoveProduct(serverApp *httptest.Server, t *testing.T) {
	req, err := http.NewRequest("DELETE", serverApp.URL+"/user/31337/cart/1076963", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}
