package privacypool

import (
	"bytes"
	"reflect"

	watcher "github.com/0xBow-io/asp-go-buildkit/core/watcher"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
)

var _ watcher.State = (*State)(nil)

type State struct {
	Sc []byte `cbor:"scope"`
	H  []byte `cbor:"hash"`
	E  []byte `cbor:"e"`
	S  []byte `cbor:"s"`
}

// DeriveState returns the new state derived from the record event
func (s *State) DeriveFrom(scope []byte, r *PrivacyPoolRecord) *State {
	if chainEvent, trans :=
		new(StateTransitionEvent).FromRecord(r); chainEvent != nil && trans != nil {
		s = &State{
			Sc: make([]byte, 32),
			H:  make([]byte, 32),
			E:  chainEvent.Serialize(),
			S:  trans.Serialize(),
		}

		copy(s.Sc, scope)
		copy(s.H, trans.NewRoot)
		return s
	}
	return nil
}

func (s *State) Serialize() []byte {
	if out, err := cbor.Marshal(s); err == nil {
		return out
	}
	return nil
}

func (*State) Deserialize(b []byte) *State {
	s := &State{}
	if err := cbor.Unmarshal(b, s); err == nil {
		return s
	}
	return nil
}

// Scope returns the scope of the state
func (s *State) Scope() []byte { return s.Sc }

// Event returns the deserialized on-chain event associated
// with the state transition
func (s *State) Event() *watcher.Event {
	return new(watcher.Event).Deserialize(s.E)
}

// Cmp compares the state with another state
func (s *State) Cmp(x watcher.State) int { return StateComparatorFunc(s, x) }

// Hash returns the hash of the state
func (s *State) Hash() []byte { return s.H }

// Inner returns the serialized state transition details
func (s *State) Inner() []byte { return s.S }

func (dst *State) copy(src *State) *State {
	dst = &State{
		Sc: make([]byte, 32),
		H:  make([]byte, 32),
		E:  make([]byte, len(src.E)),
		S:  make([]byte, len(src.S)),
	}

	copy(dst.Sc, src.Sc)
	copy(dst.H, src.H)
	copy(dst.S, src.S)
	copy(dst.E, src.E)
	return dst
}

func (s *State) Clone() watcher.State { return new(State).copy(s) }

type TransitionInput struct {
	Src          []byte `cbor:"src"`
	Sink         []byte `cbor:"sink"`
	FeeCollector []byte `cbor:"feeCollector"`
	Fee          []byte `cbor:"fee"`
}

type StateTransitionEvent struct {
	NewRoot []byte `cbor:"newRoot"`
	NewSize []byte `cbor:"newSize"`
	TransitionInput
}

func (trans *StateTransitionEvent) Serialize() []byte {
	if out, err := cbor.Marshal(trans); err == nil {
		return out
	}
	return nil
}

func (*StateTransitionEvent) Deserialize(data []byte) *StateTransitionEvent {
	trans := &StateTransitionEvent{
		TransitionInput: TransitionInput{},
	}
	if err := cbor.Unmarshal(data, trans); err == nil {
		return trans
	}
	return nil
}

func (trans *StateTransitionEvent) InputMatch(e *StateTransitionEvent) bool {
	return reflect.DeepEqual(trans.TransitionInput, e.TransitionInput)
}

func (trans *StateTransitionEvent) RootMatch(e *StateTransitionEvent) bool {
	if bytes.Compare(trans.NewRoot, e.NewRoot) != 0 {
		return false
	}
	if bytes.Compare(trans.NewSize, e.NewSize) != 0 {
		return false
	}
	return true
}

func (trans *StateTransitionEvent) FromRecord(r *PrivacyPoolRecord) (*watcher.Event, *StateTransitionEvent) {
	trans = &StateTransitionEvent{
		TransitionInput: TransitionInput{
			Src:          make([]byte, 20),
			Sink:         make([]byte, 20),
			FeeCollector: make([]byte, 20),
			Fee:          make([]byte, 32),
		},
		NewRoot: make([]byte, 32),
		NewSize: make([]byte, 32),
	}

	// Copy over TransitionInput
	copy(trans.Src, r.R.Src.Bytes())
	copy(trans.Sink, r.R.Sink.Bytes())
	copy(trans.FeeCollector, r.R.FeeCollector.Bytes())
	copy(trans.Fee, common.BigToHash(r.R.Fee).Bytes())

	// parse the record and return the stat
	copy(trans.NewRoot, common.BigToHash(r.StateRoot).Bytes())
	copy(trans.NewSize, common.BigToHash(r.StateSize).Bytes())

	return new(watcher.Event).FromLog(&r.Raw), trans
}

// StateComparatorFunc is a function that compares two states
// returns 1 if the states are inequal, 0 if they are equal
// and -1 if they are not comparable
func StateComparatorFunc(x watcher.State, y watcher.State) int {
	// Scope of States must be equal
	if bytes.Equal(x.Scope(), y.Scope()) {
		// inequal hash = inequal state
		if !bytes.Equal(x.Hash(), y.Hash()) {
			return 1
		}
		return 0
	}
	return -1
}
