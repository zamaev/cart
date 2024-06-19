package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"route256/cart/internal/app/server"
	"testing"

	"github.com/stretchr/testify/require"
)

const contentType = "application/json"

func TestServer(t *testing.T) {
	serverApp := httptest.NewServer(server.NewApp().Handler)
	defer serverApp.Close()

	AddProduct(serverApp, t)
	Checkout(serverApp, t)
	AddProduct(serverApp, t)
	GetCart(serverApp, t)
	RemoveProduct(serverApp, t)
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
