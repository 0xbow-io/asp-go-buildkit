package detector

import (
	"fmt"
	"math/big"

	watcher "github.com/0xBow-io/asp-go-buildkit/core/watcher"
	logging "github.com/ipfs/go-log/v2"
	"github.com/cockroachdb/errors"
)

var (
	log = logging.Logger("detector")
)

type StateBuffer interface {
	Stash(_ []byte) (*big.Int, error)
}

type Service struct {
	stashed   <-chan watcher.State
	lastKnown watcher.State
	StateBuffer
}

func NewService(sb StateBuffer) *Service {
	return &Service{
		lastKnown:   nil,
		StateBuffer: sb,
	}
}

func (s *Service) Absorb(in []watcher.State) (*big.Int, error) {
	var (
		root  = big.NewInt(0)
		err   error
		index int = 0
	)
	if s.lastKnown == nil {
		s.lastKnown = in[0]
		index++
	}
	for index < len(in) {
		if cmp := in[index].Cmp(s.lastKnown); cmp != 1 {
			return nil, errors.New(fmt.Sprintf("invalid state transition detected %d", cmp))
		}

		// Get the last cached state transition event
		lastEvent := s.lastKnown.Event()
		if lastEvent == nil {
			return nil, errors.New("last known state has no event")
		}

		// Get the current state transition event
		currEvent := in[index].Event()
		if currEvent == nil {
			return nil, errors.New("current state has no event")
		}

		// Check if the current state transition event is greater
		//  than the last known state transition event
		if lastEvent.Cmp(currEvent) != 1 {
			return nil, errors.New("state events are out of order")
		}

		if root, err = s.StateBuffer.Stash(in[index].Serialize()); err != nil {
			return nil, errors.Wrap(err, "failed to push the state into the buffer")
		}
		s.lastKnown = in[index].Clone()
		index++
	}
	return root, nil
}
