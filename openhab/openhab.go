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
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/robfig/cron/v3"
)

const (
	eventStateWaiting   = 0
	eventStateBanner    = 1
	eventStateData      = 2
	eventStateFinished  = 3
	eventHeader         = "event: "
	eventData           = "data: "
	eventTypeMessage    = "message"
	eventTypeEvent      = "event"
	minSupportedVersion = 3
	maxSupportedVersion = 5
)

// Client for openHAB. It's using openHAB REST API internally.
type Client struct {
	config         Config
	baseURL        string
	client         *http.Client
	user           string
	password       string
	cron           *cron.Cron
	items          *itemCollection
	rules          []*rule
	rulesMutex     sync.Mutex
	systemEventBus event.PubSub
	userEventBus   event.PubSub
	internalRules  sync.Once
	startOnce      sync.Once
	stopOnce       sync.Once
	stopChan       chan os.Signal
	running        bool
	runningMutex   sync.Mutex
	apiVersion     int
	state          ClientState
	stateMutex     sync.Mutex
	telemetry      Telemetry
	telemetryWg    sync.WaitGroup
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
	// API token takes precedence over user/password
	if config.APIToken != "" {
		config.User = config.APIToken
		config.Password = ""
	}
	client := &Client{
		config:   config,
		baseURL:  baseURL,
		client:   httpClient,
		user:     config.User,
		password: config.Password,
		cron: cron.New(
			cron.WithParser(
				cron.NewParser(
					cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor))),
		systemEventBus: event.NewEventBus(false),
		userEventBus:   event.NewEventBus(true),
		stopChan:       make(chan os.Signal, 1),
		running:        false,
		runningMutex:   sync.Mutex{},
		state:          StateStarting,
		stateMutex:     sync.Mutex{},
		telemetry:      config.Telemetry,
	}
	client.items = newItems(client)
	return client
}

// GetItem returns an openHAB item from its name.
// The very first call of GetItem will try to load the items collection from openHAB.
func (c *Client) GetItem(name string) (*Item, error) {
	return c.items.getItem(name)
}

// GetItemState returns an openHAB item state from its name. It's a shortcut of GetItem() => State().
// The very first call of GetItemState will try to load the items collection from openHAB.
func (c *Client) GetItemState(name string) (State, error) {
	item, err := c.items.getItem(name)
	if err != nil {
		return StringState(""), err
	}
	return item.State()
}

// GetMembersOf returns a list of items member of the group
func (c *Client) GetMembersOf(groupName string) ([]*Item, error) {
	return c.items.getMembersOf(groupName)
}

// SendCommand sends a command to an item. It's a shortcut for GetItem() => SendCommand().
func (c *Client) SendCommand(itemName string, command State) error {
	item, err := c.items.getItem(itemName)
	if err != nil {
		return err
	}
	return item.SendCommand(command)
}

// SendCommandWait sends a command to an item and wait until the event bus acknowledge receiving the state, or after a timeout
// It returns true if openHAB acknowledge it's setting the desired state to the item (even if it's the same value as before).
// It returns false in case the acknowledged value is different than the command, or after timeout.
// It's a shortcut for GetItem() => SendCommandWait().
func (c *Client) SendCommandWait(itemName string, command State, timeout time.Duration) (bool, error) {
	item, err := c.items.getItem(itemName)
	if err != nil {
		return false, err
	}
	return item.SendCommandWait(command, timeout)
}

func (c *Client) get(ctx context.Context, url, contentType string) (*http.Response, error) {
	debuglog.Printf("GET: %s", c.baseURL+url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.user, c.password)
	req.Header.Set("Accept", contentType)

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

func (c *Client) getString(ctx context.Context, url string) (string, error) {
	resp, err := c.get(ctx, url, "text/plain")
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

func (c *Client) getJSON(ctx context.Context, url string, result interface{}) error {
	resp, err := c.get(ctx, url, "application/json")
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

func (c *Client) postString(ctx context.Context, url, value string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+url, strings.NewReader(value))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.user, c.password)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	// we don't expect any body in the response
	resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
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
	resp, err := c.get(context.Background(), "events", "text/event-stream")
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		// send error event
		c.userEventBus.Publish(event.NewErrorEvent(err))
		return err
	}
	c.setState(StateConnected)
	// send connect event
	c.userEventBus.Publish(event.NewSystemEvent(event.TypeClientConnected))
	defer func() {
		c.setState(StateDisconnected)
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
			ev := strings.TrimPrefix(line, eventHeader)
			// "event" seems to send this data: "{\"type\":\"ALIVE\"}"
			if ev != eventTypeMessage && ev != eventTypeEvent {
				errorlog.Printf("unexpected event type %q", ev)
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
	if ev, ok := e.(event.GenericEvent); ok {
		debuglog.Printf("generic event type %q topic %q payload %q (%+v)", ev.TypeName(), ev.Topic(), ev.Payload(), data)
	}
	// debuglog.Printf("received event: %s", data)
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
		c.setState(StateConnecting)
		// run a timer in the background to reset the backoff when the connection is stable
		go func() {
			successTimerMutex.Lock()
			defer successTimerMutex.Unlock()

			successTimer = time.AfterFunc(c.config.StableConnectionDuration, func() {
				successTimerMutex.Lock()
				defer successTimerMutex.Unlock()
				backoff = c.config.ReconnectionInitialBackoff

				if !c.isState(StateConnected) {
					// still not connected, to we restart the timer
					successTimer.Reset(c.config.StableConnectionDuration)
					return
				}
				successTimer = nil
				// publish stable event
				c.userEventBus.Publish(event.NewSystemEvent(event.TypeClientConnectionStable))
				// load API version information
				c.loadIndex()
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
		backoff = nextBackoff(backoff, c.config)
	}
}

// subscribe to the user event bus (events are sent asynchronously)
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

func (c *Client) loadIndex() {
	index := restIndex{}
	ctx, cancel := context.WithTimeout(context.Background(), c.config.TimeoutHTTP)
	defer cancel()

	err := c.getJSON(ctx, "", &index)
	if err != nil {
		errorlog.Printf("cannot load information from openHAB: %s", err)
		return
	}
	apiVersion, err := strconv.Atoi(index.APIVersion)
	if err != nil {
		errorlog.Printf("invalid API version: %s", err)
		return
	}
	if apiVersion < minSupportedVersion || apiVersion > maxSupportedVersion {
		errorlog.Printf("API version %d not supported!", apiVersion)
		return
	}
	c.apiVersion = apiVersion
	debuglog.Printf("openHAB server API version %d", apiVersion)
}

// AddRule adds a rule definition
func (c *Client) AddRule(ruleData RuleData, run Runner, triggers ...Trigger) (ruleID string) {
	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()

	rule := newRule(c, ruleData, run, triggers)
	c.rules = append(c.rules, rule)
	if c.isRunning() {
		// activate it right away if the client is already running
		c.activateRule(rule)
	}
	c.addCounter(MetricRuleAdded, 1, MetricRuleID, rule.ruleData.ID)
	return rule.ruleData.ID
}

// DeleteRule deletes all the rule definition using their ruleID (it could be 0 to many)
// and returns the number of rules deleted
func (c *Client) DeleteRule(ruleID string) int {
	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()

	deleted := 0
	newRules := make([]*rule, 0, len(c.rules))
	for _, rule := range c.rules {
		if rule.ruleData.ID == ruleID {
			rule.deactivate(c)
			deleted++
			c.addCounter(MetricRuleDeleted, 1, MetricRuleID, rule.ruleData.ID)
			continue
		}
		newRules = append(newRules, rule)
	}
	c.rules = newRules
	return deleted
}

// GetRulesData returns the list of rules definition
func (c *Client) GetRulesData() []RuleData {
	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()

	rules := make([]RuleData, 0, len(c.rules))
	for _, rule := range c.rules {
		rules = append(rules, rule.ruleData)
	}
	return rules
}

// Start the handling of the defined rules.
// The function will return after the process received a Terminate, Abort or Interrupt signal,
// and after all the currently running rules have finished
//
// Please note a client can only be started once. Any other call to this method will be ignored.
func (c *Client) Start() {
	c.startOnce.Do(func() {
		c.addInternalRules()

		c.activateRules()
		c.cron.Start()

		// start the event bus
		go c.eventLoop()

		// listen to os signals
		signal.Notify(c.stopChan, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

		// Send the started event
		c.userEventBus.Publish(event.NewSystemEvent(event.TypeClientStarted))

		c.setRunning(true)

		// Wait until we're politely asked to leave
		<-c.stopChan
		signal.Stop(c.stopChan)

		c.setRunning(false)

		// Send the stopped event
		c.userEventBus.Publish(event.NewSystemEvent(event.TypeClientStopped))

		debuglog.Printf("shutting down...")
		ctx := c.cron.Stop()

		// Wait until all the cron tasks finished running
		<-ctx.Done()

		// and also all the event based rules
		c.waitFinishingRules()

		// and wait for the telemetry to finish
		c.telemetryWg.Wait()
	})
}

// Stop will send a ClientStopped event, let all the currently running rules finish, close the client, then return.
// Stop can only be called once, any subsequent call will be ignored.
func (c *Client) Stop() {
	c.stopOnce.Do(func() {
		close(c.stopChan)
	})
}

func (c *Client) setRunning(running bool) {
	c.runningMutex.Lock()
	defer c.runningMutex.Unlock()
	c.running = running
}

func (c *Client) isRunning() bool {
	c.runningMutex.Lock()
	defer c.runningMutex.Unlock()
	return c.running
}

func (c *Client) addInternalRules() {
	// make sure the internal rules are only added once
	c.internalRules.Do(func() {
		c.subscribeSystem("", event.TypeItemState, func(e event.Event) {
			c.itemStateUpdated(e)
		})
	})
}

func (c *Client) activateRules() {
	c.rulesMutex.Lock()
	defer c.rulesMutex.Unlock()

	for _, rule := range c.rules {
		c.activateRule(rule)
	}
}

func (c *Client) activateRule(rule *rule) {
	err := rule.activate(c)
	if err != nil {
		ruleName := rule.String()
		if ruleName != "" {
			ruleName = " \"" + ruleName + "\""
		}
		errorlog.Printf("error activating rule%s: %s", ruleName, err)
	}
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
		item.setInternalStateString(ev.State)
		c.addCounter(MetricItemStateUpdated, 1, MetricItemName, ev.ItemName)
	}
}

func (c *Client) setState(state ClientState) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()
	c.state = state
}

func (c *Client) isState(state ClientState) bool {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()
	return c.state == state
}

func (c *Client) getCron() *cron.Cron {
	return c.cron
}

func (c *Client) addCounter(metricName string, metricValue int64, tagName, tagValue string) {
	if c.telemetry == nil {
		return
	}
	defer preventPanic()

	c.telemetryWg.Add(1)
	go func() {
		defer c.telemetryWg.Done()
		c.telemetry.AddCounter(metricName, metricValue, getMap(tagName, tagValue))
	}()
}

func (c *Client) addGauge(metricName string, metricValue int64, tagName, tagValue string) {
	if c.telemetry == nil {
		return
	}
	defer preventPanic()

	c.telemetryWg.Add(1)
	go func() {
		defer c.telemetryWg.Done()
		c.telemetry.AddGauge(metricName, metricValue, getMap(tagName, tagValue))
	}()
}

func getMap(key, value string) map[string]string {
	if key == "" && value == "" {
		return nil
	}
	return map[string]string{key: value}
}
