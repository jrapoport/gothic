package core

import (
	"github.com/jrapoport/gothic/core/events"
	"github.com/jrapoport/gothic/store/types"
)

// Listen listens to an event.
func (a *API) Listen(e events.Event) <-chan types.Map {
	return a.evt.Listen(e)
}

// AddListener adds a listener for an event.
func (a *API) AddListener(e events.Event, listener events.Callback) {
	a.evt.AddListener(e, listener)
}

func (a *API) dispatchEvent(e events.Event, data types.Map) {
	a.evt.Dispatch(e, data)
}

func (a *API) closeDispatch() {
	a.evt.Close()
}
