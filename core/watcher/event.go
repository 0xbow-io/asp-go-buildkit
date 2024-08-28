package watcher

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/fxamacker/cbor/v2"
)

type Event struct {
	BlockNumber uint64 `cbor:"BlockNumber"`
	BlockHash   []byte `cbor:"BlockHash"`
	TxHash      []byte `cbor:"TxHash"`
	TxIndex     uint   `cbor:"TxIndex"`
	CalLData    []byte `cbor:"CalLData"`
	LogIndex    uint   `cbor:"LogIndex"`
	LogTopics   []byte `cbor:"LogTopics"`
	LogData     []byte `cbor:"LogData"`
	LogAddress  []byte `cbor:"LogAddress"`
}

func (e *Event) Format() map[string]string {
	return map[string]string{
		"BlockNumber": fmt.Sprintf("%d", e.BlockNumber),
		"BlockHash":   common.BytesToHash(e.BlockHash).Hex(),
		"TxHash":      common.BytesToHash(e.TxHash).Hex(),
		"TxIndex":     fmt.Sprintf("%d", e.TxIndex),
		"CalLData":    hex.EncodeToString(e.CalLData),
		"LogIndex":    fmt.Sprintf("%d", e.LogIndex),
		"LogTopics":   hex.EncodeToString(e.LogTopics),
		"LogData":     hex.EncodeToString(e.LogData),
		"LogAddress":  common.BytesToAddress(e.LogAddress).Hex(),
	}
}

func (e *Event) Equal(x *Event) bool {
	return e.BlockNumber == x.BlockNumber &&
		common.BytesToHash(e.BlockHash).Cmp(common.BytesToHash(x.BlockHash)) == 0 &&
		common.BytesToHash(e.TxHash).Cmp(common.BytesToHash(x.TxHash)) == 0 &&
		e.TxIndex == x.TxIndex &&
		bytes.Equal(e.CalLData, x.CalLData) &&
		e.LogIndex == x.LogIndex &&
		bytes.Equal(e.LogTopics, x.LogTopics) &&
		bytes.Equal(e.LogData, x.LogData) &&
		bytes.Equal(e.LogAddress, x.LogAddress)
}

func (e *Event) FromLog(log *types.Log) *Event {
	e = &Event{
		BlockNumber: log.BlockNumber,
		TxIndex:     log.TxIndex,
		LogIndex:    log.Index,

		BlockHash:  make([]byte, 32),
		TxHash:     make([]byte, 32),
		LogData:    make([]byte, len(log.Data)),
		LogAddress: make([]byte, 20),
		LogTopics:  make([]byte, len(log.Topics)*32),
	}

	copy(e.BlockHash, log.BlockHash.Bytes())
	copy(e.TxHash, log.TxHash.Bytes())
	copy(e.LogData, log.Data)
	copy(e.LogAddress, log.Address.Bytes())
	for i, topic := range log.Topics {
		for j, b := range topic.Bytes() {
			e.LogTopics[i*32+j] = b
		}
	}

	return e
}

func (e *Event) Serialize() []byte {
	if out, err := cbor.Marshal(e); err == nil {
		return out
	}
	return nil
}

func (*Event) Deserialize(data []byte) *Event {
	e := &Event{}
	if err := cbor.Unmarshal(data, e); err == nil {
		return e
	}
	return nil
}

// Cmp compares the event with another event.
// returns 1 if x is a later event
// returns 0 if x is the same event
// returns -1 if x is an earlier event
func (e *Event) Cmp(x *Event) int {
	if common.BytesToHash(e.BlockHash).Cmp(common.BytesToHash(x.BlockHash)) == 0 {
		if common.BytesToHash(e.TxHash).Cmp(common.BytesToHash(x.TxHash)) == 0 {
			logIndexDiff := e.LogIndex - x.LogIndex
			if logIndexDiff > 0 {
				return 1
			} else if logIndexDiff < 0 {
				return -1
			}
			return 0
		}
		txIndexDiff := e.TxIndex - x.TxIndex
		if txIndexDiff > 0 {
			return 1
		}
		return -1
	}
	blockNumberDiff := e.BlockNumber - x.BlockNumber
	if blockNumberDiff > 0 {
		return 1
	}
	return -1
}
