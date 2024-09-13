package srv

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	. "github.com/0xBow-io/asp-go-buildkit/internal"
	erpc "github.com/0xBow-io/asp-go-buildkit/internal/erpc"

	"encoding/hex"

	"github.com/0xBow-io/asp-go-buildkit/core/detector"
	r "github.com/0xBow-io/asp-go-buildkit/core/recorder"
	watcher "github.com/0xBow-io/asp-go-buildkit/core/watcher"
)

type Watcher interface {
	Watch(_ watcher.Observable) ([]watcher.State, error)
}

type Detector interface {
	Absorb(in []watcher.State) (*big.Int, error)
}

func InitBuff(sink chan []byte, initroot *big.Int) Buffer {
	buff := NewBuffer(initroot)
	// kick off the sink go-routine
	go func() {
		buff.Sink(sink)
	}()

	return buff
}

func Observe(
	ctx context.Context,
	obs watcher.Observable,
	adapter erpc.Backend,
	startBlock uint64,
	maxWindowSize uint64,
	waitTimeMs time.Duration,
	des watcher.StateDeserializer,
) (<-chan r.Record, error) {
	var (
		stream     = make(chan []byte)
		wg         = new(sync.WaitGroup)
		buff       = InitBuff(stream, big.NewInt(0))
		detector   = detector.NewService(buff)
		recorder   = r.NewService()
		watcher    = watcher.NewService(adapter)
		window     = [2]uint64{startBlock, startBlock + maxWindowSize}
		recordChan = make(chan r.Record)
	)

	go func() {
		for {

			// get the latest block number
			// from the adapter
			latestBlock, err := adapter.BlockNumber(ctx)
			if err != nil {
				fmt.Printf("adapter failure: %s \n", err.Error())
				continue
			}
			if latestBlock == window[0] {
				time.Sleep(waitTimeMs)
				continue
			}
			if latestBlock-window[0] < maxWindowSize {
				window[1] = latestBlock
			} else {
				window[1] = window[0] + maxWindowSize
			}

			// watch the osbservable states
			// and return the observations
			observations, err := watcher.Watch(obs, window)
			if err != nil {
				fmt.Printf("watcher failure: %s \n", err.Error())
				os.Exit(1)
			}

			fmt.Printf("Watched Window [%d %d] .. Received %d observations\n",
				window[0], window[1], len(observations))

			if len(observations) == 0 {
				window[0] = window[1]
				continue
			}

			// absorb the observations into the buffer
			// detector will verify that the observations are of valid state transitions
			if root, err := detector.Absorb(observations); err != nil {
				fmt.Printf("detector failure: %s \n", err.Error())
				os.Exit(1)
			} else {
				// gurantee that all the
				// observed states has been stashed into the buffer
				// the root calculated by the buffer should
				// be the same as the root returned by the detector
				if root.Cmp(buff.Root()) != 0 {
					fmt.Printf("root mismatch, got: %d expected: %d \n", root, buff.Root())
					os.Exit(1)
				}
				fmt.Printf("Buffer Root: %+v Cnt: %d\n", root, buff.Cnt())
			}
			window[0] = window[1]
		}
	}()

	for ss := range stream {
		// deserialize the state
		if state := des(ss); state != nil {
			event := state.Event()
			if event == nil {
				fmt.Println("failed to extract event !")
				os.Exit(1)
			}
			fmt.Printf("New State --> hash: %+v, event: %+v \n",
				hex.EncodeToString(state.Hash()),
				event.Format())

			rec, err := recorder.Record(state)
			if err != nil {
				fmt.Printf("recorder failure: %s \n", err.Error())
				os.Exit(1)
			}
			if rec != nil {
				Categorize(rec.Serialize(), adapter)
			}
		} else {
			fmt.Println("failed to deserialize state !")
			os.Exit(1)
		}
	}

	wg.Wait()
	return recordChan, nil
}
