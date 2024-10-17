package internal

import (
	"math/big"
	"sync"

	posiedon "github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/cockroachdb/errors"
)

/*
	StateBuffer is a simple channel based ring buffer
	which is used to quickly
	cache the serialized states of the observable.
*/

type Buffer interface {
	Stash(_ []byte) (*big.Int, error)
	Sink(sink chan<- []byte)
	Size() int
	Root() *big.Int
	Cnt() uint64
	Purge() bool
}

type _Buffer struct {
	stash   chan []byte
	root    *big.Int
	counter uint64
	lock    *sync.RWMutex
}

func NewBuffer(initRoot *big.Int) Buffer {
	return &_Buffer{
		stash:   make(chan []byte),
		counter: 0,
		root:    big.NewInt(0).Set(initRoot),
		lock:    &sync.RWMutex{},
	}
}

func (s *_Buffer) Root() *big.Int {
	return s.root
}

func (s *_Buffer) Stash(ss []byte) (*big.Int, error) {

	var (
		hash    *big.Int
		newRoot *big.Int
		err     error
	)
	if hash, err = posiedon.HashBytes(ss); err == nil {
		if newRoot, err = posiedon.HashWithState([]*big.Int{hash}, s.root); err == nil {
			s.stash <- ss
			s.counter++
			return s.root.Set(newRoot), nil
		}
	}
	return nil, errors.Wrap(err, "failed to stash incoming")
}

func (s *_Buffer) Sink(sink chan<- []byte) {
	for ss := range s.stash {
		sink <- ss
	}
}

func (s *_Buffer) Size() int {
	return len(s.stash)
}

func (s *_Buffer) Cnt() uint64 {
	return s.counter
}

// Purge removes all the stashed & to-be sinked
// states from the buffer
// By simply dumping the items into the ether
func (s *_Buffer) Purge() bool {
	defer s.lock.Unlock()
	s.lock.Lock()
	if s.counter == 0 {
		return false
	}
	for {
		select {
		case <-s.stash:
			continue
		default:
			s.counter = 0
			s.root = nil
			return true
		}
	}
}
