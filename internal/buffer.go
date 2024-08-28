package internal

import (
	watcher "github.com/0xBow-io/asp-go-buildkit/core/watcher"
)

type StateBuffer struct {
}

func (s *StateBuffer) Stash(watcher.State) error {
	return nil
}

func (s *StateBuffer) Discard() bool {
	return false
}

func (s *StateBuffer) PeekLast() (watcher.State, error) {
	return nil, nil
}

func (s *StateBuffer) Size() int {
	return 0
}

func (s *StateBuffer) Pop() (watcher.State, error) {
	return nil, nil
}

func (s *StateBuffer) Commit() error {
	return nil
}

func (s *StateBuffer) Snapshot() ([]byte, error) {
	return nil, nil
}

func (s *StateBuffer) Purge() error {
	return nil
}
