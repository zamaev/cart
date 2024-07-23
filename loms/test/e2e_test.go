package test

import (
	"context"
	"net"
	"route256/cart/pkg/logger"
	"route256/loms/internal/app/server"
	"route256/loms/internal/pkg/config"
	"route256/loms/pkg/api/loms/v1"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestServer(t *testing.T) {
	config, err := config.NewConfig()
	require.NoError(t, err)
	app := server.NewApp(context.Background(), config)

	lis, err := net.Listen("tcp", config.GrpcUrl)
	require.NoError(t, err)
	defer lis.Close()

	go func() {
		logger.Infow(context.Background(), "starting server app", "url", config.GrpcUrl)
		require.NoError(t, app.GrpcServer.Serve(lis))
	}()

	grpcClient, err := grpc.NewClient(config.GrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	lomsClient := loms.NewLomsClient(grpcClient)

	OrderCreate(lomsClient, t)
	time.Sleep(time.Millisecond)
	OrderInfo(lomsClient, t)
	time.Sleep(time.Millisecond)
	OrderPay(lomsClient, t)
	time.Sleep(time.Millisecond)
	OrderCreate2(lomsClient, t)
	time.Sleep(time.Millisecond)
	OrderCancel(lomsClient, t)
	time.Sleep(time.Millisecond)
	OrderInfo2(lomsClient, t)
	time.Sleep(time.Millisecond)
	StocksInfo(lomsClient, t)

	app.GrpcServer.Stop()
}

func OrderCreate(lomsClient loms.LomsClient, t *testing.T) {
	res, err := lomsClient.OrderCreate(context.Background(), &loms.OrderCreateRequest{
		User: 1,
		Items: []*loms.OrderItem{
			{
				Sku:   1,
				Count: 1,
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), res.OrderId)
}

func OrderInfo(lomsClient loms.LomsClient, t *testing.T) {
	res, err := lomsClient.OrderInfo(context.Background(), &loms.OrderInfoRequest{
		OrderId: 1,
	})
	require.Equal(t, "awaiting payment", res.Status)
	require.Equal(t, int64(1), res.User)
	require.Len(t, res.Items, 1)
	require.Equal(t, uint32(1), res.Items[0].Sku)
	require.Equal(t, uint32(1), res.Items[0].Count)
	require.NoError(t, err)
}

func OrderPay(lomsClient loms.LomsClient, t *testing.T) {
	_, err := lomsClient.OrderPay(context.Background(), &loms.OrderPayRequest{
		OrderId: 1,
	})
	require.NoError(t, err)

	res, err := lomsClient.OrderInfo(context.Background(), &loms.OrderInfoRequest{
		OrderId: 1,
	})
	require.NoError(t, err)
	require.Equal(t, "payed", res.Status)
}

func OrderCreate2(lomsClient loms.LomsClient, t *testing.T) {
	res, err := lomsClient.OrderCreate(context.Background(), &loms.OrderCreateRequest{
		User: 1,
		Items: []*loms.OrderItem{
			{
				Sku:   1,
				Count: 1,
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, int64(2), res.OrderId)
}

func OrderCancel(lomsClient loms.LomsClient, t *testing.T) {
	_, err := lomsClient.OrderCancel(context.Background(), &loms.OrderCancelRequest{
		OrderId: 2,
	})
	require.NoError(t, err)
}

func OrderInfo2(lomsClient loms.LomsClient, t *testing.T) {
	res, err := lomsClient.OrderInfo(context.Background(), &loms.OrderInfoRequest{
		OrderId: 2,
	})
	require.NoError(t, err)
	require.Equal(t, "cancelled", res.Status)

}

func StocksInfo(lomsClient loms.LomsClient, t *testing.T) {
	_, err := lomsClient.StocksInfo(context.Background(), &loms.StocksInfoRequest{
		Sku: 1,
	})
	require.NoError(t, err)
}
