package recorder

import (
	watcher "github.com/0xBow-io/asp-go-buildkit/core/watcher"
	logging "github.com/ipfs/go-log/v2"
	"github.com/pkg/errors"
)

var (
	log = logging.Logger("recorder")
)

type StateBuffer interface {
	Size() int
	Pop() (watcher.State, error)
	PeekLast() (watcher.State, error)
}

type Service struct {
	StateBuffer
}

func (s *Service) Record(out chan<- []byte) error {
	for {
		if s.Size() >= 2 {
			post, err := s.Pop()
			if err != nil {
				return errors.Wrap(err, "failed to pop the latest state from the buffer")
			}

			pre, err := s.PeekLast()
			if err != nil {
				return errors.Wrap(err, "failed to peek the last known state")
			}
			out <- new(_Record).Build(post, pre).Serialize()
		}
	}
}
