package watcher

import (
	"context"

	"github.com/0xBow-io/asp-go-buildkit/internal/erpc"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	logging "github.com/ipfs/go-log/v2"
	"github.com/pkg/errors"
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
	getRange func(context.Context) (uint64, uint64)
	adapter  erpc.Backend
}

func NewService(
	getRange func(context.Context) (uint64, uint64),
	adapter erpc.Backend,
) *Service {
	return &Service{
		getRange: getRange,
		adapter:  adapter,
	}
}

func (s *Service) Watch(obs Observable) ([]State, error) {
	var (
		err        error
		states     []State
		start, end = s.getRange(context.Background())
	)

	stream, err := obs.Play(s.adapter, &bind.FilterOpts{
		Context: context.Background(),
		Start:   start,
		End:     &end,
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
