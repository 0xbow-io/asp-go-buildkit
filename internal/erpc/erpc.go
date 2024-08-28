package erpc

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ERPC is a wrapper around the geth rpc clientlient struct instead of interface...
// see https://github.com/ethereum/go-ethereum/issues/28267
type ERPC struct {
	connType connType
	client   *ethclient.Client
}

var _ Backend = (*ERPC)(nil)

func getConnType(conn string) connType {
	if strings.Contains(conn, "wss://") {
		return WSS
	} else if strings.Contains(conn, "ws://") {
		return WS
	} else if strings.Contains(conn, "http://") {
		return HTTP
	} else if strings.Contains(conn, "https://") {
		return HTTPS
	}
	return UNSUPPORTED
}

func NewERPC(conn string) (*ERPC, error) {
	connType := getConnType(conn)
	if connType == UNSUPPORTED {
		return nil, errors.New(ErrUnknownConnType)
	}
	client, err := ethclient.Dial(conn)
	if err != nil {
		return nil, err
	}

	return NewERPCFromClient(client), nil
}

func NewERPCFromClient(
	client *ethclient.Client,
) *ERPC {
	return &ERPC{
		client: client,
	}
}

// Generic function used to handle all client calls
func callFunc[T any](
	rpcCall func() (T, error),
	rpcMethodName string,
	erpc *ERPC,
) (value T, err error) {
	result, err := rpcCall()
	return result, nil
}

func (erpc *ERPC) ConnType() connType {
	return erpc.connType
}

// gethClient interface methods
func (erpc *ERPC) ChainID(ctx context.Context) (*big.Int, error) {
	chainID := func() (*big.Int, error) { return erpc.client.ChainID(ctx) }
	id, err := callFunc[*big.Int](chainID, "eth_chainId", erpc)
	return id, err
}

func (erpc *ERPC) BalanceAt(
	ctx context.Context,
	account common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	balanceAt := func() (*big.Int, error) { return erpc.client.BalanceAt(ctx, account, blockNumber) }
	balance, err := callFunc[*big.Int](balanceAt, "eth_getBalance", erpc)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (erpc *ERPC) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	blockByHash := func() (*types.Block, error) { return erpc.client.BlockByHash(ctx, hash) }
	block, err := callFunc[*types.Block](blockByHash, "eth_getBlockByHash", erpc)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (erpc *ERPC) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	blockByNumber := func() (*types.Block, error) { return erpc.client.BlockByNumber(ctx, number) }
	block, err := callFunc[*types.Block](
		blockByNumber,
		"eth_getBlockByNumber",
		erpc,
	)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (erpc *ERPC) BlockNumber(ctx context.Context) (uint64, error) {
	blockNumber := func() (uint64, error) { return erpc.client.BlockNumber(ctx) }
	number, err := callFunc[uint64](blockNumber, "eth_blockNumber", erpc)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (erpc *ERPC) CallContract(
	ctx context.Context,
	call ethereum.CallMsg,
	blockNumber *big.Int,
) ([]byte, error) {
	callContract := func() ([]byte, error) { return erpc.client.CallContract(ctx, call, blockNumber) }
	bytes, err := callFunc[[]byte](callContract, "eth_call", erpc)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (erpc *ERPC) CodeAt(
	ctx context.Context,
	contract common.Address,
	blockNumber *big.Int,
) ([]byte, error) {
	call := func() ([]byte, error) { return erpc.client.CodeAt(ctx, contract, blockNumber) }
	bytes, err := callFunc[[]byte](call, "eth_getCode", erpc)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (erpc *ERPC) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	estimateGas := func() (uint64, error) { return erpc.client.EstimateGas(ctx, call) }
	gas, err := callFunc[uint64](estimateGas, "eth_estimateGas", erpc)
	if err != nil {
		return 0, err
	}
	return gas, nil
}

func (erpc *ERPC) FeeHistory(
	ctx context.Context,
	blockCount uint64,
	lastBlock *big.Int,
	rewardPercentiles []float64,
) (*ethereum.FeeHistory, error) {
	feeHistory := func() (*ethereum.FeeHistory, error) {
		return erpc.client.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
	}
	history, err := callFunc[*ethereum.FeeHistory](
		feeHistory,
		"eth_feeHistory",
		erpc,
	)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (erpc *ERPC) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	filterLogs := func() ([]types.Log, error) { return erpc.client.FilterLogs(ctx, query) }
	logs, err := callFunc[[]types.Log](filterLogs, "eth_getLogs", erpc)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (erpc *ERPC) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	headerByHash := func() (*types.Header, error) { return erpc.client.HeaderByHash(ctx, hash) }
	header, err := callFunc[*types.Header](
		headerByHash,
		"eth_getBlockByHash",
		erpc,
	)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (erpc *ERPC) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	headerByNumber := func() (*types.Header, error) { return erpc.client.HeaderByNumber(ctx, number) }
	header, err := callFunc[*types.Header](
		headerByNumber,
		"eth_getBlockByNumber",
		erpc,
	)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (erpc *ERPC) NonceAt(
	ctx context.Context,
	account common.Address,
	blockNumber *big.Int,
) (uint64, error) {
	nonceAt := func() (uint64, error) { return erpc.client.NonceAt(ctx, account, blockNumber) }
	nonce, err := callFunc[uint64](nonceAt, "eth_getTransactionCount", erpc)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (erpc *ERPC) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	pendingBalanceAt := func() (*big.Int, error) { return erpc.client.PendingBalanceAt(ctx, account) }
	balance, err := callFunc[*big.Int](pendingBalanceAt, "eth_getBalance", erpc)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (erpc *ERPC) PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error) {
	pendingCallContract := func() ([]byte, error) { return erpc.client.PendingCallContract(ctx, call) }
	bytes, err := callFunc[[]byte](pendingCallContract, "eth_call", erpc)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (erpc *ERPC) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	pendingCodeAt := func() ([]byte, error) { return erpc.client.PendingCodeAt(ctx, account) }
	bytes, err := callFunc[[]byte](pendingCodeAt, "eth_getCode", erpc)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (erpc *ERPC) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	pendingNonceAt := func() (uint64, error) { return erpc.client.PendingNonceAt(ctx, account) }
	nonce, err := callFunc[uint64](
		pendingNonceAt,
		"eth_getTransactionCount",
		erpc,
	)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (erpc *ERPC) PendingStorageAt(
	ctx context.Context,
	account common.Address,
	key common.Hash,
) ([]byte, error) {
	pendingStorageAt := func() ([]byte, error) { return erpc.client.PendingStorageAt(ctx, account, key) }
	bytes, err := callFunc[[]byte](pendingStorageAt, "eth_getStorageAt", erpc)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (erpc *ERPC) PendingTransactionCount(ctx context.Context) (uint, error) {
	pendingTransactionCount := func() (uint, error) { return erpc.client.PendingTransactionCount(ctx) }
	count, err := callFunc[uint](
		pendingTransactionCount,
		"eth_getBlockTransactionCountByNumber",
		erpc,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (erpc *ERPC) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	// callFunc takes a function that returns a value and an error
	// so we just wrap the SendTransaction method in a function that returns 0 as its value,
	// which we throw out below
	sendTransaction := func() (int, error) { return 0, erpc.client.SendTransaction(ctx, tx) }
	_, err := callFunc[int](sendTransaction, "eth_sendRawTransaction", erpc)
	return err
}

func (erpc *ERPC) StorageAt(
	ctx context.Context,
	account common.Address,
	key common.Hash,
	blockNumber *big.Int,
) ([]byte, error) {
	storageAt := func() ([]byte, error) { return erpc.client.StorageAt(ctx, account, key, blockNumber) }
	bytes, err := callFunc[[]byte](storageAt, "eth_getStorageAt", erpc)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (erpc *ERPC) SubscribeFilterLogs(
	ctx context.Context,
	query ethereum.FilterQuery,
	ch chan<- types.Log,
) (ethereum.Subscription, error) {
	subscribeFilterLogs := func() (ethereum.Subscription, error) { return erpc.client.SubscribeFilterLogs(ctx, query, ch) }
	subscription, err := callFunc[ethereum.Subscription](
		subscribeFilterLogs,
		"eth_subscribe",
		erpc,
	)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (erpc *ERPC) SubscribeNewHead(
	ctx context.Context,
	ch chan<- *types.Header,
) (ethereum.Subscription, error) {
	subscribeNewHead := func() (ethereum.Subscription, error) { return erpc.client.SubscribeNewHead(ctx, ch) }
	subscription, err := callFunc[ethereum.Subscription](
		subscribeNewHead,
		"eth_subscribe",
		erpc,
	)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (erpc *ERPC) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	suggestGasPrice := func() (*big.Int, error) { return erpc.client.SuggestGasPrice(ctx) }
	gasPrice, err := callFunc[*big.Int](suggestGasPrice, "eth_gasPrice", erpc)
	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

func (erpc *ERPC) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	suggestGasTipCap := func() (*big.Int, error) { return erpc.client.SuggestGasTipCap(ctx) }
	gasTipCap, err := callFunc[*big.Int](
		suggestGasTipCap,
		"eth_maxPriorityFeePerGas",
		erpc,
	)
	if err != nil {
		return nil, err
	}
	return gasTipCap, nil
}

func (erpc *ERPC) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	syncProgress := func() (*ethereum.SyncProgress, error) { return erpc.client.SyncProgress(ctx) }
	progress, err := callFunc[*ethereum.SyncProgress](
		syncProgress,
		"eth_syncing",
		erpc,
	)
	if err != nil {
		return nil, err
	}
	return progress, nil
}

// We write the instrumentation of this function directly because callFunc[] generic fct only takes a single
// return value
func (erpc *ERPC) TransactionByHash(
	ctx context.Context,
	hash common.Hash,
) (tx *types.Transaction, isPending bool, err error) {
	tx, isPending, err = erpc.client.TransactionByHash(ctx, hash)
	return tx, isPending, nil
}

func (erpc *ERPC) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	transactionCount := func() (uint, error) { return erpc.client.TransactionCount(ctx, blockHash) }
	count, err := callFunc[uint](
		transactionCount,
		"eth_getBlockTransactionCountByHash",
		erpc,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (erpc *ERPC) TransactionInBlock(
	ctx context.Context,
	blockHash common.Hash,
	index uint,
) (*types.Transaction, error) {
	transactionInBlock := func() (*types.Transaction, error) { return erpc.client.TransactionInBlock(ctx, blockHash, index) }
	tx, err := callFunc[*types.Transaction](
		transactionInBlock,
		"eth_getTransactionByBlockHashAndIndex",
		erpc,
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (erpc *ERPC) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	transactionReceipt := func() (*types.Receipt, error) { return erpc.client.TransactionReceipt(ctx, txHash) }
	receipt, err := callFunc[*types.Receipt](
		transactionReceipt,
		"eth_getTransactionReceipt",
		erpc,
	)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
