package watcher

type StateDeserializer func(_ []byte) State

type State interface {
	// State Hash function
	// Returns the hash of the state
	Event() *Event
	Scope() []byte
	Inner() []byte
	Hash() []byte
	Cmp(_ State) int
	Clone() State
	Serialize() []byte
	Deserialize(_ []byte) State
}
