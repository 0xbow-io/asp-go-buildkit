package srv

import (
	"context"
	"fmt"

	r "github.com/0xBow-io/asp-go-buildkit/core/recorder"
	chainalaysis "github.com/0xBow-io/asp-go-buildkit/internal/category/feature/extractors/plugins/chainalysis"
	erpc "github.com/0xBow-io/asp-go-buildkit/internal/erpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	caAPI = chainalaysis.NewAPIV2("")
)

func Categorize(rec []byte, adapter erpc.Backend) {
	var (
		ctx    = context.Background()
		record = r.DeserializeRecord(rec)
	)
	if record == nil {
		fmt.Println("Failed to deserialize record")
		return
	}

	resp := HandleRegistration(ctx, record, adapter)
	externalid := resp.ExternalID

	if externalid == "" {
		fmt.Println("Failed to get external ID from registration")
		return
	}

	result, err := caAPI.GetTransferSummary(ctx, externalid)
	if err != nil {
		fmt.Printf("Error getting transfer summary: %v\n", err)
	} else {
		fmt.Printf("Transfer Summary: %+v\n", result)
	}

	exposure, err := caAPI.GetDirectExposure(ctx, externalid)
	if err != nil {
		fmt.Printf("Error getting direct exposure: %v\n", err)
	} else {
		fmt.Printf("Direct Exposure: %+v\n", exposure)
	}

	alerts, err := caAPI.GetAlerts(ctx, externalid)
	if err != nil {
		fmt.Printf("Error getting alerts: %v\n", err)
	} else {
		for _, alert := range alerts {
			fmt.Printf("Alert: %+v\n", alert)
		}
	}
}

func HandleRegistration(ctx context.Context, record r.Record, adapter erpc.Backend) *chainalaysis.TransferRegisterResp {
	txHash := common.BytesToHash(record.Event().TxHash)

	tx, isPending, err := adapter.TransactionByHash(ctx, txHash)
	if err != nil || isPending {
		fmt.Printf("error %+v", err)
		return nil
	}

	transferVal := tx.Value()
	if transferVal == nil {
		return nil
	}

	toAddress := tx.To().String()
	fromAddress, _ := GetFrom(tx)

	registerResp, err := caAPI.RegisterTransfer(ctx, toAddress, chainalaysis.TransferRegisterReq{
		Network:           "gnosis",
		Asset:             "xDAI",
		TransferReference: fmt.Sprintf("%s:%s", txHash.Hex(), toAddress),
		Direction:         chainalaysis.DirectionSent,
		InputAddresses:    []string{fromAddress},
		OutputAddress:     toAddress,
	})

	return &registerResp

}

func GetFrom(tx *types.Transaction) (string, error) {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	return from.String(), err
}
