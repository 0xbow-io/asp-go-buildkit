package core

import (
	. "github.com/0xBow-io/asp-go-buildkit/core/watcher"
	logging "github.com/ipfs/go-log/v2"
)

var (
	log = logging.Logger("observer")
)

type StateBuffer interface {
	Stash(State) error
	Discard() bool
	PeekLast() (State, error)
	Size() int
	Pop() (State, error)
	Commit() error
	Snapshot() ([]byte, error)
	Purge() error
}

type Watcher interface {
	Watch(_ Observable) ([]State, error)
}

type Detector interface {
	Absorb(_ []State) error
}

type Recorder interface {
	Record(_ chan<- []byte) error
}

type Observer struct {
	StateBuffer
	Watcher
	Detector
	Recorder
}

func (o *Observer) Observe(obs Observable) (chan<- []byte, error) {
	stream := make(chan<- []byte)
	go func() {
		if err := o.Record(stream); err != nil {
			log.Error(err)
			return
		}
	}()
	go func() {
		for {
			states, err := o.Watch(obs)
			if err != nil {
				log.Error(err)
				continue
			}

			if err := o.Absorb(states); err != nil {
				log.Error(err)
				continue
			}
		}
	}()
	return stream, nil
}
