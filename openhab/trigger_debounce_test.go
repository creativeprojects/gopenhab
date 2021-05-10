package openhab

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/stretchr/testify/assert"
)

type mockTrigger struct {
	callback       func(e event.Event)
	OnActivation   func()
	OnDeactivation func()
}

func (t *mockTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	t.callback = run
	if t.OnActivation != nil {
		t.OnActivation()
	}
	return nil
}

func (t *mockTrigger) deactivate(client *Client) {
	t.callback = nil
	if t.OnDeactivation != nil {
		t.OnDeactivation()
	}
}

func TestDebounce(t *testing.T) {
	var (
		counter uint64
	)

	trigger := &mockTrigger{}
	debounced := Debounce(trigger, 100*time.Millisecond)
	debounced.activate(nil, func(event.Event) {
		atomic.AddUint64(&counter, 1)
	}, RuleData{})

	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			trigger.callback(nil)
		}

		time.Sleep(200 * time.Millisecond)
	}

	assert.Equal(t, 3, int(atomic.LoadUint64(&counter)))
}

func TestDebounceConcurrentRun(t *testing.T) {
	var wg sync.WaitGroup

	var flag uint64

	trigger := &mockTrigger{}
	debounced := Debounce(trigger, 100*time.Millisecond)
	debounced.activate(nil, func(event.Event) {
		atomic.CompareAndSwapUint64(&flag, 0, 1)
	}, RuleData{})

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			trigger.callback(nil)
		}()
	}
	wg.Wait()

	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, 1, int(atomic.LoadUint64(&flag)), "Flag not set")
}

func TestDebounceDelayed(t *testing.T) {
	var (
		counter uint64
	)

	trigger := &mockTrigger{}
	debounced := Debounce(trigger, 100*time.Millisecond)
	debounced.activate(nil, func(event.Event) {
		atomic.AddUint64(&counter, 1)
	}, RuleData{})

	time.Sleep(110 * time.Millisecond)

	trigger.callback(nil)

	time.Sleep(200 * time.Millisecond)

	assert.Equal(t, 1, int(atomic.LoadUint64(&counter)))
}

func TestDebounceCancelled(t *testing.T) {
	var (
		counter uint64
	)

	trigger := &mockTrigger{}
	debounced := Debounce(trigger, 100*time.Millisecond)
	debounced.activate(nil, func(event.Event) {
		atomic.AddUint64(&counter, 1)
	}, RuleData{})

	trigger.callback(nil)

	time.Sleep(50 * time.Millisecond)
	debounced.deactivate(nil)

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 0, int(atomic.LoadUint64(&counter)))
}
