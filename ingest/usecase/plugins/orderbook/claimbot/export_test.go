package claimbot

import (
	"context"

	"github.com/osmosis-labs/sqs/domain"
	sqstx "github.com/osmosis-labs/sqs/domain/cosmos/tx"
	"github.com/osmosis-labs/sqs/domain/keyring"
	"github.com/osmosis-labs/sqs/domain/mvc"
	orderbookdomain "github.com/osmosis-labs/sqs/domain/orderbook"

	txfeestypes "github.com/osmosis-labs/osmosis/v27/x/txfees/types"

	"github.com/osmosis-labs/osmosis/osmomath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// ProcessedOrderbook is order alias data structure for testing purposes.
type ProcessedOrderbook = processedOrderbook

// ProcessOrderbooksAndGetClaimableOrders is test wrapper for processOrderbooksAndGetClaimableOrders.
// This function is exported for testing purposes.
func ProcessOrderbooksAndGetClaimableOrders(
	ctx context.Context,
	orderbookusecase mvc.OrderBookUsecase,
	fillThreshold osmomath.Dec,
	orderbooks []domain.CanonicalOrderBooksResult,
) ([]ProcessedOrderbook, error) {
	return processOrderbooksAndGetClaimableOrders(ctx, orderbookusecase, fillThreshold, orderbooks)
}

// SendBatchClaimTx a test wrapper for sendBatchClaimTx.
// This function is used only for testing purposes.
func SendBatchClaimTx(
	ctx context.Context,
	keyring keyring.Keyring,
	txfeesClient txfeestypes.QueryClient,
	gasCalculator sqstx.GasCalculator,
	txServiceClient txtypes.ServiceClient,
	chainID string,
	account *authtypes.BaseAccount,
	contractAddress string,
	claims orderbookdomain.Orders,
) (*sdk.TxResponse, error) {
	return sendBatchClaimTx(ctx, keyring, txfeesClient, gasCalculator, txServiceClient, chainID, account, contractAddress, claims)
}

// PrepareBatchClaimMsg is a test wrapper for prepareBatchClaimMsg.
// This function is exported for testing purposes.
func PrepareBatchClaimMsg(claims orderbookdomain.Orders) ([]byte, error) {
	return prepareBatchClaimMsg(claims)
}

// GetOrderbooks is a test wrapper for getOrderbooks.
// This function is exported for testing purposes.
func GetOrderbooks(poolsUsecase mvc.PoolsUsecase, metadata domain.BlockPoolMetadata) ([]domain.CanonicalOrderBooksResult, error) {
	return getOrderbooks(poolsUsecase, metadata)
}
