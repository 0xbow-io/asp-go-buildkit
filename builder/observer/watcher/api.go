/*
Watches the network for signals (e.g. event emissions) by protocols that indicate a state-change has occured.
*/

package watcher

import (
	core "github.com/0xBow-io/asp-go-buildkit/core/watcher"
)

var _ Module = (*API)(nil)

type Module interface {
	// Watch returns a set of state
	// of an observable instance.
	Watch(_ core.Observable) ([]core.State, error)
}

type API struct {
	Internal struct {
		Watch func(
			core.Observable,
		) ([]core.State, error)
	}
}

func (a *API) Watch(obs core.Observable) ([]core.State, error) {
	return a.Internal.Watch(obs)
}
