package privacypool

import (
	"math/big"
	"reflect"
	"testing"

	watcher "github.com/0xBow-io/asp-go-buildkit/core/watcher"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func Test_StateTransitionEvent(t *testing.T) {

	var (
		req = IPrivacyPoolRequest{
			Src:          common.HexToAddress("0x01"),
			Sink:         common.HexToAddress("0x02"),
			FeeCollector: common.HexToAddress("0x03"),
			Fee:          big.NewInt(1000),
		}
		ste = &StateTransitionEvent{
			TransitionInput: TransitionInput{
				Src:          common.HexToAddress("0x01").Bytes(),
				Sink:         common.HexToAddress("0x02").Bytes(),
				FeeCollector: common.HexToAddress("0x03").Bytes(),
				Fee:          common.BigToHash(big.NewInt(1000)).Bytes(),
			},
			NewRoot: common.BigToHash(big.NewInt(1001)).Bytes(),
			NewSize: common.BigToHash(big.NewInt(1002)).Bytes(),
		}
		ce = watcher.Event{
			BlockNumber: 7,
			BlockHash:   common.HexToHash("0x08").Bytes(),
			TxHash:      common.HexToHash("0x09").Bytes(),
			TxIndex:     10,
			LogIndex:    11,
			LogTopics:   common.HexToHash("0x12").Bytes(),
			LogData:     []byte("0x13"),
			LogAddress:  common.HexToAddress("0x14").Bytes(),
		}
	)

	event, trans := new(StateTransitionEvent).FromRecord(&PrivacyPoolRecord{
		R:         req,
		StateRoot: common.BytesToHash(ste.NewRoot).Big(),
		StateSize: common.BytesToHash(ste.NewSize).Big(),
		Raw: types.Log{
			BlockNumber: ce.BlockNumber,
			BlockHash:   common.BytesToHash(ce.BlockHash),
			TxHash:      common.BytesToHash(ce.TxHash),
			TxIndex:     ce.TxIndex,
			Index:       ce.LogIndex,
			Topics:      []common.Hash{common.BytesToHash(ce.LogTopics)},
			Data:        ce.LogData,
			Address:     common.BytesToAddress(ce.LogAddress),
		},
	})

	require.NotNil(t, event)
	require.NotNil(t, trans)
	require.Equal(t, true, ste.InputMatch(trans))
	require.Equal(t, true, ste.RootMatch(trans))
	require.Equal(t, true, ce.Equal(event))

	serialized := trans.Serialize()
	require.NotNil(t, serialized)

	deserialized := new(StateTransitionEvent).Deserialize(serialized)
	require.NotNil(t, deserialized)

	require.Equal(t, true, reflect.DeepEqual(trans, deserialized))
}

func Test_State(t *testing.T) {

	var (
		req = IPrivacyPoolRequest{
			Src:          common.HexToAddress("0x01"),
			Sink:         common.HexToAddress("0x02"),
			FeeCollector: common.HexToAddress("0x03"),
			Fee:          big.NewInt(1000),
		}
		ste = &StateTransitionEvent{
			TransitionInput: TransitionInput{
				Src:          common.HexToAddress("0x01").Bytes(),
				Sink:         common.HexToAddress("0x02").Bytes(),
				FeeCollector: common.HexToAddress("0x03").Bytes(),
				Fee:          common.BigToHash(big.NewInt(1000)).Bytes(),
			},
			NewRoot: common.BigToHash(big.NewInt(1001)).Bytes(),
			NewSize: common.BigToHash(big.NewInt(1002)).Bytes(),
		}
		ce = watcher.Event{
			BlockNumber: 7,
			BlockHash:   common.HexToHash("0x08").Bytes(),
			TxHash:      common.HexToHash("0x09").Bytes(),
			TxIndex:     10,
			LogIndex:    11,
			LogTopics:   common.HexToHash("0x12").Bytes(),
			LogData:     []byte("0x13"),
			LogAddress:  common.HexToAddress("0x14").Bytes(),
		}
		comparable = &State{
			Sc: common.BigToHash(big.NewInt(1000)).Bytes(),
			H:  ste.NewRoot,
			E:  ce.Serialize(),
		}
		nonComparable = &State{
			// different scope
			Sc: common.BigToHash(big.NewInt(1001)).Bytes(),
			H:  ste.NewRoot,
			E:  ce.Serialize(),
		}
	)
	// test for equality
	equalState := new(State).DeriveFrom(comparable.Sc, &PrivacyPoolRecord{
		R:         req,
		StateRoot: common.BytesToHash(ste.NewRoot).Big(),
		StateSize: common.BytesToHash(ste.NewSize).Big(),
		Raw: types.Log{
			BlockNumber: ce.BlockNumber,
			BlockHash:   common.BytesToHash(ce.BlockHash),
			TxHash:      common.BytesToHash(ce.TxHash),
			TxIndex:     ce.TxIndex,
			Index:       ce.LogIndex,
			Topics:      []common.Hash{common.BytesToHash(ce.LogTopics)},
			Data:        ce.LogData,
			Address:     common.BytesToAddress(ce.LogAddress),
		},
	})
	require.NotNil(t, equalState)
	require.Equal(t, 0, comparable.Cmp(equalState))
	require.Equal(t, true, comparable.Event().Equal(equalState.Event()))

	// Test for uncomparable
	require.Equal(t, -1, nonComparable.Cmp(equalState))

	// test for inequality
	notEqualState := new(State).DeriveFrom(comparable.Sc, &PrivacyPoolRecord{
		R: req,
		// root is different
		StateRoot: big.NewInt(1),
		StateSize: common.BytesToHash(ste.NewSize).Big(),
		Raw: types.Log{
			// block number is different
			BlockNumber: ce.BlockNumber + 1,
			BlockHash:   common.BytesToHash(ce.BlockHash),
			TxHash:      common.BytesToHash(ce.TxHash),
			TxIndex:     ce.TxIndex,
			Index:       ce.LogIndex,
			Topics:      []common.Hash{common.BytesToHash(ce.LogTopics)},
			Data:        ce.LogData,
			Address:     common.BytesToAddress(ce.LogAddress),
		},
	})
	require.NotNil(t, notEqualState)
	require.Equal(t, 1, comparable.Cmp(notEqualState))
	require.Equal(t, false, comparable.Event().Equal(notEqualState.Event()))

}
