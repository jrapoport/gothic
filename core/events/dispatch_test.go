package events

import (
	"sync"
	"testing"
	"time"

	"github.com/jrapoport/gothic/log"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
)

func allEvents() []Event {
	return []Event{
		Login,
		Signup,
		Confirmed,
	}
}

func testDispatch(t *testing.T) *Dispatch {
	n := utils.RandomUsername()
	d := NewDispatch(n, log.Log)
	t.Cleanup(func() {
		d.Close()
	})
	return d
}

func TestDispatch_AddListener(t *testing.T) {
	t.Parallel()
	testListen(t, func(d *Dispatch, evt Event, cb Callback) {
		d.AddListener(evt, cb)
	})
}

func TestDispatch_Listen(t *testing.T) {
	t.Parallel()
	testListen(t, func(d *Dispatch, evt Event, cb Callback) {
		ch := d.Listen(evt)
		go func() {
			for {
				select {
				case msg, ok := <-ch:
					if !ok {
						return
					}
					e, ok := msg[key.Event].(Event)
					assert.True(t, ok)
					assert.NotEqual(t, e, All)
					cb(e, msg)
				}
			}
		}()
	})
	d := testDispatch(t)
	bad := d.Listen(Unknown)
	assert.Nil(t, bad)
	d.closed = 1
	bad = d.Listen(Logout)
	assert.Nil(t, bad)
}

func testListen(t *testing.T, lis func(d *Dispatch, evt Event, cb Callback)) {
	type WaitAssertionFunc = func(t assert.TestingT,
		condition func() bool,
		waitFor time.Duration,
		tick time.Duration,
		msgAndArgs ...interface{}) bool
	tests := []struct {
		e    Event
		fire bool
		Err  WaitAssertionFunc
	}{
		//	{Signup, true, assert.Eventually},
		//	{Login, false, assert.Never},
		{All, true, assert.Eventually},
	}
	d := testDispatch(t)
	for _, test := range tests {
		var wait sync.WaitGroup
		t.Run(string(test.e), func(t *testing.T) {
			var m types.Map
			var mu sync.RWMutex
			var recv int
			var expect = 1
			em := tutils.RandomEmail()
			lis(d, test.e, func(evt Event, msg types.Map) {
				mu.Lock()
				defer mu.Unlock()
				m = msg
				recv++
				wait.Done()
			})
			m = nil
			if test.fire && test.e == All {
				expect = len(allEvents())
				wait.Add(expect)
				for _, e := range allEvents() {
					d.Dispatch(e, types.Map{
						key.Email: em,
					})
				}
			} else if test.fire {
				wait.Add(expect)
				d.Dispatch(test.e, types.Map{
					key.Email: em,
				})
			}
			test.Err(t, func() bool {
				mu.RLock()
				defer mu.RUnlock()
				return recv == expect
			}, 1*time.Second, 10*time.Millisecond)
			if test.fire {
				assert.NotNil(t, m)
				assert.Equal(t, em, m[key.Email])
			}
		})
		wait.Wait()
		d.mu.Lock()
		d.mgr.Clear()
		d.mu.Unlock()
	}
}

func TestDispatch_Dispatch(t *testing.T) {
	t.Parallel()
	d := testDispatch(t)
	ch := make(chan types.Map)
	d.listeners[ch] = Unknown
	d.Dispatch(Unknown, nil)
	assert.Never(t, func() bool {
		var ok bool
		go func() {
			_, ok = <-ch
		}()
		return ok
	}, 1*time.Second, 10*time.Millisecond)
	var m types.Map
	var mu sync.RWMutex
	d.AddListener(Login, func(evt Event, msg types.Map) {
		mu.Lock()
		m = msg
		mu.Unlock()
	})
	d.Dispatch(Login, nil)
	assert.Eventually(t, func() bool {
		mu.RLock()
		defer mu.RUnlock()
		return m != nil
	}, 1*time.Second, 10*time.Millisecond)
	assert.Equal(t, Login, m[key.Event].(Event))
}

func TestDispatch_Close(t *testing.T) {
	t.Parallel()
	em := tutils.RandomEmail()
	d := testDispatch(t)
	ch := d.Listen(Signup)
	go func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				e, ok := msg[key.Event].(Event)
				assert.True(t, ok)
				assert.NotEqual(t, e, All)
			}
		}
	}()
	d.Dispatch(Signup, types.Map{
		key.Email: em,
	})
	d.Close()
	// already closed
	assert.NotPanics(t, func() {
		d.Dispatch(Signup, types.Map{
			key.Email: em,
		})
	})
}
