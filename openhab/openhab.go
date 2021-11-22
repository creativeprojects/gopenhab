package openhab

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/robfig/cron/v3"
)

const (
	eventStateWaiting  = 0
	eventStateBanner   = 1
	eventStateData     = 2
	eventStateFinished = 3
	eventHeader        = "event: "
	eventData          = "data: "
)

// RuleClient is an interface for a Client inside a rule
type RuleClient interface {
	GetItem(name string) (*Item, error)
	GetMembersOf(groupName string) ([]*Item, error)
	AddRule(ruleData RuleData, run Runner, triggers ...Trigger) (ruleID string)
}

var _ RuleClient = &Client{}

// Client for openHAB. It's using openHAB REST API internally.
type Client struct {
	config         Config
	baseURL        string
	client         *http.Client
	cron           *cron.Cron
	items          *Items
	rules          []*rule
	systemEventBus event.PubSub
	userEventBus   event.PubSub
	internalRules  sync.Once
	startOnce      sync.Once
	stopOnce       sync.Once
	stopChan       chan os.Signal
}

// NewClient creates a new client to connect to a openHAB instance
func NewClient(config Config) *Client {
	if config.URL == "" {
		panic("missing URL from Config")
	}
	if config.ReconnectionInitialBackoff == 0 {
		config.ReconnectionInitialBackoff = time.Second
	}
	if config.ReconnectionMultiplier == 0 {
		config.ReconnectionMultiplier = 2.0
	}
	if config.ReconnectionMaxBackoff == 0 {
		config.ReconnectionMaxBackoff = time.Minute
	}
	if config.StableConnectionDuration == 0 {
		config.StableConnectionDuration = time.Minute
	}
	if config.TimeoutHTTP == 0 {
		config.TimeoutHTTP = 5 * time.Second
	}
	baseURL := strings.ToLower(config.URL)
	if baseURL[:len(baseURL)-1] != "/" {
		baseURL += "/"
	}
	if !strings.HasSuffix(config.URL, "/rest/") {
		baseURL += "rest/"
	}
	httpClient := http.DefaultClient
	if config.Client != nil {
		httpClient = config.Client
	}
	client := &Client{
		config:  config,
		baseURL: baseURL,
		client:  httpClient,
		cron: cron.New(
			cron.WithParser(
				cron.NewParser(
					cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor))),
		systemEventBus: event.NewEventBus(false),
		userEventBus:   event.NewEventBus(true),
		stopChan:       make(chan os.Signal, 1),
	}
	client.items = newItems(client)
	return client
}

// GetItem returns an openHAB item from its name.
// The very first call of GetItem will try to load the items collection from openHAB.
func (c *Client) GetItem(name string) (*Item, error) {
	return c.items.getItem(name)
}

// GetMembersOf returns a list of items member of the group
func (c *Client) GetMembersOf(groupName string) ([]*Item, error) {
	return c.items.getMembersOf(groupName)
}

func (c *Client) get(ctx context.Context, URL string) (*http.Response, error) {
	debuglog.Printf("GET: %s", c.baseURL+URL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+URL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode >= 400 {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return resp, ErrorNotFound
		default:
			return resp, errors.New(resp.Status)
		}
	}

	return resp, nil
}

func (c *Client) getString(ctx context.Context, URL string) (string, error) {
	resp, err := c.get(ctx, URL)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) getJSON(ctx context.Context, URL string, result interface{}) error {
	resp, err := c.get(ctx, URL)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) postString(ctx context.Context, URL string, value string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+URL, strings.NewReader(value))
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	// we don't expect any body in the response
	resp.Body.Close()

	if resp.StatusCode >= 400 {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return ErrorNotFound
		default:
			return errors.New(resp.Status)
		}
	}

	return nil
}

// listenEvents listen to the events from the REST api and send them to the event bus.
// the method returns after the HTTP connection dropped
func (c *Client) listenEvents() error {
	resp, err := c.get(context.Background(), "events")
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		// send error event
		c.userEventBus.Publish(event.NewErrorEvent(err))
		return err
	}
	// send connect event
	c.userEventBus.Publish(event.NewSystemEvent(event.TypeClientConnected))
	defer func() {
		// send disconnect event
		c.userEventBus.Publish(event.NewSystemEvent(event.TypeClientDisconnected))
	}()

	state := 0
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		state++
		line := scanner.Text()
		if line == "" {
			// Move back to waiting state
			if state != eventStateFinished {
				errorlog.Printf("unexpected end of event data on state %d", state)
			}
			state = eventStateWaiting
			continue
		}
		if state == eventStateBanner {
			if !strings.HasPrefix(line, eventHeader) {
				errorlog.Printf("unexpected start of event: %q", line)
			}
			event := strings.TrimPrefix(line, eventHeader)
			if event != "message" {
				errorlog.Printf("unexpected event: %q", event)
			}
			continue
		}
		if state == eventStateData {
			if !strings.HasPrefix(line, eventData) {
				errorlog.Printf("unexpected event data: %q", line)
			}
			data := strings.TrimPrefix(line, eventData)
			if data != "" {
				c.dispatchRawEvent(data)
			}
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		// send error event
		c.userEventBus.Publish(event.NewErrorEvent(err))
		return err
	}
	return nil
}

func (c *Client) dispatchRawEvent(data string) {
	e, err := event.New(data)
	if err != nil {
		errorlog.Printf("event ignored: %s", err)
		return
	}
	c.systemEventBus.Publish(e)
	c.userEventBus.Publish(e)
}

// eventLoop listen to the events from the REST api and send them to the event bus.
// the method never returns: if the connection drops it tries to reconnect in a loop
func (c *Client) eventLoop() {
	var successTimer *time.Timer
	var successTimerMutex sync.Mutex
	backoff := c.config.ReconnectionInitialBackoff

	for {
		// run a timer in the background to reset the backoff when the connection is stable
		go func() {
			successTimerMutex.Lock()
			defer successTimerMutex.Unlock()

			successTimer = time.AfterFunc(c.config.StableConnectionDuration, func() {
				backoff = c.config.ReconnectionInitialBackoff
				debuglog.Printf("connection to the event bus looks stable")

				successTimerMutex.Lock()
				defer successTimerMutex.Unlock()
				successTimer = nil
				// ==> We could maybe publish a new system event when the connection is stable?
			})
		}()
		err := c.listenEvents()
		if err != nil {
			errorlog.Printf("error connecting or listening to openHAB events: %s", err)
		}
		// we just got logged off so we cancel any success timer
		successTimerMutex.Lock()
		if successTimer != nil {
			successTimer.Stop()
		}
		successTimerMutex.Unlock()
		debuglog.Printf("reconnecting in %s...", backoff.String())
		time.Sleep(backoff)

		// calculate next backoff
		backoff = time.Duration(float64(backoff) * c.config.ReconnectionMultiplier)
		if backoff > c.config.ReconnectionMaxBackoff {
			backoff = c.config.ReconnectionMaxBackoff
		}
	}
}

func (c *Client) subscribe(name string, eventType event.Type, callback func(e event.Event)) int {
	return c.userEventBus.Subscribe(name, eventType, func(e event.Event) {
		defer preventPanic()
		callback(e)
	})
}

// subscribeSystem is a subscription to the system (synchronous) event bus
func (c *Client) subscribeSystem(name string, eventType event.Type, callback func(e event.Event)) int {
	return c.systemEventBus.Subscribe(name, eventType, func(e event.Event) {
		defer preventPanic()
		callback(e)
	})
}

func (c *Client) unsubscribe(subID int) {
	c.userEventBus.Unsubscribe(subID)
}

// AddRule adds a rule definition
func (c *Client) AddRule(ruleData RuleData, run Runner, triggers ...Trigger) (ruleID string) {
	rule := newRule(c, ruleData, run, triggers)
	c.rules = append(c.rules, rule)
	return rule.ruleData.ID
}

// Start the handling of the defined rules.
// The function will return after the process received a Terminate, Abort or Interrupt signal,
// and after all the currently running rules have finished
//
// Please note a client can only be started once. Any other call to this method will be ignored.
func (c *Client) Start() {
	c.startOnce.Do(func() {
		c.addInternalRules()

		for _, rule := range c.rules {
			err := rule.activate(c)
			if err != nil {
				ruleName := rule.String()
				if ruleName != "" {
					ruleName = " \"" + ruleName + "\""
				}
				errorlog.Printf("error activating rule%s: %s", ruleName, err)
			}
		}
		c.cron.Start()

		// start the event bus
		go c.eventLoop()

		// listen to os signals
		signal.Notify(c.stopChan, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

		// Send the started event
		c.userEventBus.Publish(event.NewSystemEvent(event.TypeClientStarted))

		// Wait until we're politely asked to leave
		<-c.stopChan
		signal.Stop(c.stopChan)

		// Send the stopped event
		c.userEventBus.Publish(event.NewSystemEvent(event.TypeClientStopped))

		debuglog.Printf("shutting down...")
		ctx := c.cron.Stop()

		// Wait until all the cron tasks finished running
		<-ctx.Done()

		// and also all the event based rules
		c.waitFinishingRules()
	})
}

func (c *Client) Stop() {
	c.stopOnce.Do(func() {
		close(c.stopChan)
	})
}

func (c *Client) addInternalRules() {
	// make sure the internal rules are only added once
	c.internalRules.Do(func() {
		c.subscribeSystem("", event.TypeItemState, func(e event.Event) {
			c.itemStateUpdated(e)
		})
	})
}

func (c *Client) waitFinishingRules() {
	c.userEventBus.Wait()
}

func (c *Client) itemStateUpdated(e event.Event) {
	if ev, ok := e.(event.ItemReceivedState); ok {
		item, err := c.items.getItem(ev.ItemName)
		if err != nil {
			errorlog.Printf("itemStateUpdated: %s", err)
			return
		}
		item.setInternalState(ev.State)
		// debuglog.Printf("Item %s received state %s", ev.ItemName, ev.State)
	}
}
