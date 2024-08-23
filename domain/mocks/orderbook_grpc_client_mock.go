package mocks

import (
	"context"

	orderbookdomain "github.com/osmosis-labs/sqs/domain/orderbook"
	orderbookgrpcclientdomain "github.com/osmosis-labs/sqs/domain/orderbook/grpcclient"
	orderbookplugindomain "github.com/osmosis-labs/sqs/domain/orderbook/plugin"
)

var _ orderbookgrpcclientdomain.OrderBookClient = (*OrderbookGRPCClientMock)(nil)

// OrderbookGRPCClientMock is a mock struct that implements orderbookplugindomain.OrderbookGRPCClient.
type OrderbookGRPCClientMock struct {
	MockGetOrdersByTickCb          func(ctx context.Context, contractAddress string, tick int64) ([]orderbookplugindomain.Order, error)
	MockGetActiveOrdersCb          func(ctx context.Context, contractAddress string, ownerAddress string) (orderbookdomain.Orders, uint64, error)
	MockGetTickUnrealizedCancelsCb func(ctx context.Context, contractAddress string, tickIDs []int64) ([]orderbookgrpcclientdomain.UnrealizedTickCancels, error)
	MockQueryTicksCb               func(ctx context.Context, contractAddress string, ticks []int64) ([]orderbookdomain.Tick, error)
}

func (o *OrderbookGRPCClientMock) GetOrdersByTick(ctx context.Context, contractAddress string, tick int64) ([]orderbookplugindomain.Order, error) {
	if o.MockGetOrdersByTickCb != nil {
		return o.MockGetOrdersByTickCb(ctx, contractAddress, tick)
	}

	return nil, nil
}

func (o *OrderbookGRPCClientMock) GetActiveOrders(ctx context.Context, contractAddress string, ownerAddress string) (orderbookdomain.Orders, uint64, error) {
	if o.MockGetActiveOrdersCb != nil {
		return o.MockGetActiveOrdersCb(ctx, contractAddress, ownerAddress)
	}

	return nil, 0, nil
}

func (o *OrderbookGRPCClientMock) GetTickUnrealizedCancels(ctx context.Context, contractAddress string, tickIDs []int64) ([]orderbookgrpcclientdomain.UnrealizedTickCancels, error) {
	if o.MockGetTickUnrealizedCancelsCb != nil {
		return o.MockGetTickUnrealizedCancelsCb(ctx, contractAddress, tickIDs)
	}

	return nil, nil
}

func (o *OrderbookGRPCClientMock) QueryTicks(ctx context.Context, contractAddress string, ticks []int64) ([]orderbookdomain.Tick, error) {
	if o.MockQueryTicksCb != nil {
		return o.MockQueryTicksCb(ctx, contractAddress, ticks)
	}

	return nil, nil
}