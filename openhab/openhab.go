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

	"github.com/creativeprojects/gopenhab/api"
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

type Client struct {
	config         Config
	baseURL        string
	client         *http.Client
	cron           *cron.Cron
	items          *Items
	rules          []*Rule
	eventBus       event.PubSub
	internalRules  sync.Once
	rulesWaitGroup sync.WaitGroup
	rulesLocker    sync.Mutex
}

func NewClient(config Config) *Client {
	if config.URL == "" {
		panic("missing URL from Config")
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
		eventBus: event.NewEventBus(),
	}
	client.items = newItems(client)
	return client
}

func (c *Client) Items() *Items {
	return c.items
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
		return err
	}
	// send connect event
	c.eventBus.Publish(event.NewSystemEvent(event.TypeClientConnected))
	defer func() {
		// send disconnect event
		c.eventBus.Publish(event.NewSystemEvent(event.TypeClientDisconnected))
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
	c.eventBus.Publish(e)
}

// eventLoop listen to the events from the REST api and send them to the event bus.
// the method never returns: if the connection drops it tries to reconnect in a loop
func (c *Client) eventLoop() {
	for {
		err := c.listenEvents()
		if err != nil {
			errorlog.Printf("error listening to openhab events: %s", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func (c *Client) subscribe(topic string, eventType event.Type, callback func(e event.Event)) int {
	return c.eventBus.Subscribe(topic, eventType, func(e event.Event) {
		c.ruleExecutionStarted()
		defer c.rulesWaitGroup.Done()
		defer preventPanic()
		callback(e)
	})
}

func (c *Client) unsubscribe(subID int) {
	c.eventBus.Unsubscribe(subID)
}

func (c *Client) AddRule(ruleData RuleData, run Runner, triggers ...Trigger) error {
	rule := newRule(c, ruleData, run, triggers)
	c.rules = append(c.rules, rule)
	return nil
}

// Start the handling of the defined rules.
// The function will return after the process received a Terminate, Abort or Interrupt signal,
// and after all the currently running rules have finished
func (c *Client) Start() {
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

	// Wait until we're politely asked to leave
	<-stop

	debuglog.Printf("shutting down...")
	ctx := c.cron.Stop()

	// Wait until all the cron tasks finished running
	<-ctx.Done()

	// and also all the event based rules
	c.waitFinishingRules()
}

func (c *Client) addInternalRules() {
	// make sure the internal rules are only added once
	c.internalRules.Do(func() {
		c.subscribe("", event.TypeItemState, func(e event.Event) {
			c.itemStateChanged(e)
		})
	})
}

func (c *Client) ruleExecutionStarted() {
	c.rulesLocker.Lock()
	defer c.rulesLocker.Unlock()

	c.rulesWaitGroup.Add(1)
}

func (c *Client) waitFinishingRules() {
	c.rulesLocker.Lock()
	defer c.rulesLocker.Unlock()

	c.rulesWaitGroup.Wait()
}

func (c *Client) itemStateChanged(e event.Event) {
	if ev, ok := e.(event.ItemReceivedState); ok {
		itemName := strings.TrimPrefix(ev.Topic(), itemTopicPrefix)
		itemName = strings.TrimSuffix(itemName, "/"+api.TopicEventState)

		item, err := c.items.GetItem(itemName)
		if err != nil {
			errorlog.Printf("itemStateChanged: %w", err)
			return
		}
		item.setInternalState(ev.State)
		debuglog.Printf("Item %s received state %s", itemName, ev.State)
	}
}
