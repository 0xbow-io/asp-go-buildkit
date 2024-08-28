package recorder

import (
	"bytes"

	"github.com/0xBow-io/asp-go-buildkit/core/watcher"
	"github.com/fxamacker/cbor/v2"
)

type Record interface {
	Hash() []byte
	Scope() []byte
	// Serialized Event
	Event() *watcher.Event
	// Hashes of the pre and post states
	PreState() []byte
	PostState() []byte
	Serialize() []byte
	Deserialize([]byte) Record
}

type _Record struct {
	Sc    []byte `cbor:"scope"`
	H     []byte `cbor:"hash"`
	E     []byte `cbor:"event"`
	PreS  []byte `cbor:"prestate"`
	PostS []byte `cbor:"poststate"`
}

func (*_Record) Build(post watcher.State, pre watcher.State) Record {
	event := post.Event().Serialize()
	if event != nil {
		return nil
	}
	r := &_Record{
		Sc:    make([]byte, 32),
		H:     make([]byte, 32),
		PreS:  make([]byte, 32),
		PostS: make([]byte, 32),
		E:     make([]byte, len(event)),
	}
	copy(r.Sc, post.Scope())
	copy(r.H, post.Hash())
	copy(r.PreS, pre.Hash())
	copy(r.PostS, post.Hash())
	copy(r.E, event)

	// verify the hashes of the pre and post states
	// are different
	if bytes.Equal(r.PreS, r.PostS) {
		return nil
	}

	return r
}

func (r _Record) Hash() []byte          { return r.H }
func (r _Record) Scope() []byte         { return r.Sc }
func (r _Record) Event() *watcher.Event { return new(watcher.Event).Deserialize(r.E) }
func (r _Record) PreState() []byte      { return r.PreS }
func (r _Record) PostState() []byte     { return r.PostS }
func (r _Record) Serialize() []byte {
	if out, err := cbor.Marshal(r); err == nil {
		return out
	}
	return nil
}

func (r *_Record) Deserialize(data []byte) Record {
	_r := &_Record{}

	if err := cbor.Unmarshal(data, _r); err == nil {
		return _r
	}
	return nil
}

var _ Record = (*_Record)(nil)