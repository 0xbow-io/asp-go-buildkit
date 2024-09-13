package erpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type connType uint

const (
	UNSUPPORTED connType = iota
	HTTPS
	WSS
	IPC
	HTTP
	WS
)

type Backend interface {
	bind.ContractBackend
	ConnType() connType
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	TransactionByHash(
		ctx context.Context,
		hash common.Hash,
	) (tx *types.Transaction, isPending bool, err error)
}
