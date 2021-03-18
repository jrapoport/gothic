package events

import (
	"sync"
	"sync/atomic"

	"github.com/gookit/event"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/sirupsen/logrus"
)

// Dispatch for events
type Dispatch struct {
	mgr       *event.Manager
	log       logrus.FieldLogger
	listeners map[chan types.Map]Event
	closed    uint32
	quit      chan struct{}
	mu        sync.RWMutex
}

// NewDispatch creates a new event dispatch
func NewDispatch(name string, l logrus.FieldLogger) *Dispatch {
	d := &Dispatch{
		mgr:       event.NewManager(name),
		listeners: map[chan types.Map]Event{},
		quit:      make(chan struct{}),
		log:       l.WithField("dispatch", name),
	}
	return d
}

// Callback event listener callback function
type Callback func(evt Event, msg types.Map)

// Listen listens for an event and signals the channel
func (d *Dispatch) Listen(evt Event) chan types.Map {
	if atomic.LoadUint32(&d.closed) != 0 {
		d.log.Warn("dispatch closed")
		return nil
	}
	if evt == Unknown {
		d.log.Error("unknown event")
		return nil
	}
	d.mu.Lock()
	d.log.Debugf("listen %s", evt)
	ch := make(chan types.Map)
	d.listeners[ch] = evt
	d.mu.Unlock()
	d.AddListener(evt, func(evt Event, msg types.Map) {
		if atomic.LoadUint32(&d.closed) != 0 {
			return
		}
		select {
		case ch <- msg:
			break
		case <-d.quit:
			break
		}

	})
	return ch
}

// AddListener adds a listener callback for the event
func (d *Dispatch) AddListener(evt Event, cb Callback) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if evt == Unknown || cb == nil {
		return
	}
	if evt == All {
		evt = "*"
	}
	var fnHandler = func(e event.Event) error {
		d.log.Debugf("received event %s", e.Name())
		cb(Event(e.Name()), e.Data())
		return nil
	}
	d.log.Debugf("add listener for: %s", evt)
	d.mgr.AddListener(string(evt), event.ListenerFunc(fnHandler), event.Normal)
}

// Dispatch an event to the listeners
func (d *Dispatch) Dispatch(evt Event, data types.Map) {
	go func(e Event, m types.Map) {
		d.dispatch(e, m)
	}(evt, data)
}

func (d *Dispatch) dispatch(evt Event, data types.Map) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if evt == Unknown || evt == All {
		d.log.Warnf("invalid dispatch event: %s", evt)
		return
	}
	if data == nil {
		data = types.Map{}
	}
	data[key.Event] = evt
	e := event.NewBasic(string(evt), event.M(data))
	d.log.Debugf("send %s ch: %v", evt, data)
	err := d.mgr.FireEvent(e)
	if err != nil {
		d.log.Error(err)
	}
}

// Close the event dispatch
func (d *Dispatch) Close() {
	if atomic.LoadUint32(&d.closed) != 0 {
		d.log.Warn("dispatch closed")
		return
	}
	atomic.StoreUint32(&d.closed, 1)
	d.mu.Lock()
	defer d.mu.Unlock()
	close(d.quit)
	d.mgr.Clear()
	d.log.Debug("close dispatch")
	for ch, evt := range d.listeners {
		d.log.Debugf("close %s ch: %v", evt, ch)
		close(ch)
		delete(d.listeners, ch)
	}
}
