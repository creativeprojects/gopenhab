package openhab

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTrigger struct {
	callback       func(e event.Event)
	OnActivation   func()
	OnDeactivation func()
}

func (t *mockTrigger) activate(client subscriber, run func(ev event.Event), ruleData RuleData) error {
	t.callback = run
	if t.OnActivation != nil {
		t.OnActivation()
	}
	return nil
}

func (t *mockTrigger) deactivate(client subscriber) {
	t.callback = nil
	if t.OnDeactivation != nil {
		t.OnDeactivation()
	}
}

func (t *mockTrigger) match(e event.Event) bool {
	return true
}

func TestDebounce(t *testing.T) {
	t.Parallel()
	var counter uint64

	trigger := &mockTrigger{}
	debounced := Debounce(50*time.Millisecond, trigger)
	err := debounced.activate(nil, func(event.Event) {
		atomic.AddUint64(&counter, 1)
	}, RuleData{})
	require.NoError(t, err)

	for range 3 {
		// run 10 times but will only be counted once after a 50ms without activity
		for range 10 {
			trigger.callback(event.NewGenericEvent("Test", "Topic", "Payload"))
		}
		time.Sleep(200 * time.Millisecond)
	}

	// managed 3 full run with the help of a 200ms sleep in between
	assert.Equal(t, uint64(3), atomic.LoadUint64(&counter))
}

func TestDebounceConcurrentRun(t *testing.T) {
	t.Parallel()
	var (
		count   = 10
		wg      sync.WaitGroup
		counter uint64
	)

	trigger := &mockTrigger{}
	debounced := Debounce(100*time.Millisecond, trigger)
	err := debounced.activate(nil, func(event.Event) {
		atomic.AddUint64(&counter, 1)
	}, RuleData{})
	require.NoError(t, err)

	for range count {
		wg.Add(1)
		go func() {
			defer wg.Done()
			trigger.callback(event.NewGenericEvent("Test", "Topic", "Payload"))
		}()
	}
	wg.Wait()

	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, uint64(1), atomic.LoadUint64(&counter))
}

func TestDebounceDelayed(t *testing.T) {
	t.Parallel()
	var counter uint64

	trigger := &mockTrigger{}
	debounced := Debounce(100*time.Millisecond, trigger)
	err := debounced.activate(nil, func(event.Event) {
		atomic.AddUint64(&counter, 1)
	}, RuleData{})
	require.NoError(t, err)

	time.Sleep(110 * time.Millisecond)

	trigger.callback(event.NewGenericEvent("Test", "Topic", "Payload"))

	time.Sleep(300 * time.Millisecond)

	assert.Equal(t, uint64(1), atomic.LoadUint64(&counter))
}

func TestDebounceCancelled(t *testing.T) {
	t.Parallel()
	var counter uint64

	trigger := &mockTrigger{}
	debounced := Debounce(50*time.Millisecond, trigger)
	err := debounced.activate(nil, func(event.Event) {
		atomic.AddUint64(&counter, 1)
	}, RuleData{})
	require.NoError(t, err)

	trigger.callback(event.NewGenericEvent("Test", "Topic", "Payload"))

	time.Sleep(10 * time.Millisecond)
	debounced.deactivate(nil)

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, uint64(0), atomic.LoadUint64(&counter))
}

func TestDebounceTwoTriggers(t *testing.T) {
	t.Parallel()
	var counter uint64

	trigger1 := &mockTrigger{}
	trigger2 := &mockTrigger{}
	debounced := Debounce(50*time.Millisecond, trigger1, trigger2)
	err := debounced.activate(nil, func(event.Event) {
		atomic.AddUint64(&counter, 1)
	}, RuleData{})
	require.NoError(t, err)

	for range 3 {
		for range 10 {
			trigger1.callback(event.NewGenericEvent("Test", "Topic", "Payload"))
			trigger2.callback(event.NewGenericEvent("Test", "Topic", "Payload"))
		}

		time.Sleep(200 * time.Millisecond)
	}

	assert.Equal(t, uint64(3), atomic.LoadUint64(&counter))
}

func TestDebounceConcurrentRunOfThreeTriggers(t *testing.T) {
	t.Parallel()
	var (
		count   = 10
		wg      sync.WaitGroup
		counter uint64
	)
	trigger1 := &mockTrigger{}
	trigger2 := &mockTrigger{}
	trigger3 := &mockTrigger{}
	debounced := Debounce(100*time.Millisecond, trigger1, trigger2, trigger3)
	err := debounced.activate(nil, func(event.Event) {
		atomic.AddUint64(&counter, 1)
	}, RuleData{})
	require.NoError(t, err)

	for range count {
		wg.Add(1)
		go func() {
			defer wg.Done()
			trigger1.callback(event.NewGenericEvent("Test", "Topic", "Payload"))
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			trigger2.callback(event.NewGenericEvent("Test", "Topic", "Payload"))
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			trigger3.callback(event.NewGenericEvent("Test", "Topic", "Payload"))
		}()
	}
	wg.Wait()

	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, uint64(1), atomic.LoadUint64(&counter))
}
