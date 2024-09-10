package recorder

import (
	"errors"
	"sync"

	"github.com/0xBow-io/asp-go-buildkit/core/watcher"
	logging "github.com/ipfs/go-log/v2"
)

var (
	log = logging.Logger("recorder")
)

type Service struct {
	wg       *sync.WaitGroup
	preState watcher.State
}

func NewService() *Service {
	return &Service{
		preState: nil,
		wg:       new(sync.WaitGroup),
	}
}

func (s *Service) Record(postState watcher.State) (Record, error) {
	var rec Record = nil
	if s.preState != nil {
		if rec := new(_Record).Build(postState, s.preState); rec == nil {
			return nil, errors.New("failed to build a new record")
		} else {
			s.preState = postState.Clone()
			return rec, nil
		}
	}

	s.preState = postState.Clone()
	return rec, nil
}
