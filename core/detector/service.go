package detector

import (
	watcher "github.com/0xBow-io/asp-go-buildkit/core/watcher"
	logging "github.com/ipfs/go-log/v2"
	"github.com/pkg/errors"
)

var (
	log = logging.Logger("detector")
)

type StateBuffer interface {
	// Stash the state into the buffer.
	// This does not commit the state to the buffer
	// And thus can be discarded at any time.
	Stash(watcher.State) error
	PeekLast() (watcher.State, error)
}

type Service struct {
	StateBuffer
}

func NewService(sb StateBuffer) *Service {
	return &Service{
		StateBuffer: sb,
	}
}

func (s *Service) Absorb(in []watcher.State) error {
	var (
		lastKnown, err = s.PeekLast()
	)
	if err != nil {
		return errors.Wrap(err, "failed to peek the last known state")
	}
	for i := 0; i < len(in); i++ {
		if in[i].Cmp(lastKnown) != 1 {
			return errors.New("invalid state transition detected")
		}

		lastEvent := lastKnown.Event()
		if lastEvent == nil {
			return errors.New("last known state has no event")
		}

		currEvent := in[i].Event()
		if currEvent == nil {
			return errors.New("current state has no event")
		}

		if lastEvent.Cmp(currEvent) != 1 {
			return errors.New("state event is not in order")
		}

		if err := s.Stash(in[i]); err != nil {
			return errors.Wrap(err, "failed to push the state into the buffer")
		}
		lastKnown = in[i]
	}
	return nil
}
