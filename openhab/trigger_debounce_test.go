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

	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			trigger.callback(nil)
		}

		time.Sleep(200 * time.Millisecond)
	}

	assert.Equal(t, uint64(3), atomic.LoadUint64(&counter))
}

func TestDebounceConcurrentRun(t *testing.T) {
	t.Parallel()
	var (
		count = 10
		wg    sync.WaitGroup
		flag  uint64
	)

	trigger := &mockTrigger{}
	debounced := Debounce(100*time.Millisecond, trigger)
	err := debounced.activate(nil, func(event.Event) {
		atomic.CompareAndSwapUint64(&flag, 0, 1)
	}, RuleData{})
	require.NoError(t, err)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			trigger.callback(nil)
		}()
	}
	wg.Wait()

	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, uint64(1), atomic.LoadUint64(&flag), "Flag not set")
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

	trigger.callback(nil)

	time.Sleep(300 * time.Millisecond)

	assert.Equal(t, uint64(1), atomic.LoadUint64(&counter))
}

func TestDebounceCancelled(t *testing.T) {
	t.Parallel()
	var counter uint64

	trigger := &mockTrigger{}
	debounced := Debounce(100*time.Millisecond, trigger)
	err := debounced.activate(nil, func(event.Event) {
		atomic.AddUint64(&counter, 1)
	}, RuleData{})
	require.NoError(t, err)

	trigger.callback(nil)

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

	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			trigger1.callback(nil)
			trigger2.callback(nil)
		}

		time.Sleep(200 * time.Millisecond)
	}

	assert.Equal(t, uint64(3), atomic.LoadUint64(&counter))
}

func TestDebounceConcurrentRunOfThreeTriggers(t *testing.T) {
	t.Parallel()
	var (
		count = 10
		wg    sync.WaitGroup
		flag  uint64
	)
	trigger1 := &mockTrigger{}
	trigger2 := &mockTrigger{}
	trigger3 := &mockTrigger{}
	debounced := Debounce(100*time.Millisecond, trigger1, trigger2, trigger3)
	err := debounced.activate(nil, func(event.Event) {
		atomic.CompareAndSwapUint64(&flag, 0, 1)
	}, RuleData{})
	require.NoError(t, err)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			trigger1.callback(nil)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			trigger2.callback(nil)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			trigger3.callback(nil)
		}()
	}
	wg.Wait()

	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, uint64(1), atomic.LoadUint64(&flag), "Flag not set")
}
