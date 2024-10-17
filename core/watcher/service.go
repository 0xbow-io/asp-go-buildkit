package watcher

import (
	"context"

	"github.com/0xBow-io/asp-go-buildkit/internal/erpc"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	logging "github.com/ipfs/go-log/v2"
	"github.com/cockroachdb/errors"
)

var (
	log = logging.Logger("watcher")
)

type Observable interface {
	ID() string
	Scope() []byte
	ChainID() int
	Address() common.Address
	Deserialize(data []byte) State
	Play(_ erpc.Backend, _ *bind.FilterOpts) (<-chan []byte, error)
}

type Service struct {
	adapter erpc.Backend
}

func NewService(
	adapter erpc.Backend,
) *Service {
	return &Service{
		adapter: adapter,
	}
}

func (s *Service) Watch(obs Observable, blockRange [2]uint64) ([]State, error) {
	if blockRange[0] > blockRange[1] || blockRange[0] == 0 || blockRange[1] == 0 {
		return nil, errors.New("invalid block range")
	}
	var (
		err    error
		states []State
	)

	stream, err := obs.Play(s.adapter, &bind.FilterOpts{
		Context: context.Background(),
		Start:   blockRange[0],
		End:     &blockRange[1],
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to observe the given instance")
	}

	for {
		s, ok := <-stream
		if !ok {
			log.Debugw("watcher/Watch: stream closed", "instance", obs.ID())
			break
		}
		if s == nil {
			return nil, errors.New("received a nil state from the observable instance")
		}
		states = append(states, obs.Deserialize(s))
	}
	return states, nil
}
