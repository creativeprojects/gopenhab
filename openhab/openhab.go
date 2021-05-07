package openhab

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	config   Config
	baseURL  string
	client   *http.Client
	cron     *cron.Cron
	items    *Items
	rules    []*Rule
	eventBus eventBus
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
		eventBus: newEventBus(),
	}
	client.items = newItems(client)
	return client
}

func (c *Client) Items() *Items {
	return c.items
}

func (c *Client) get(ctx context.Context, URL string) (*http.Response, error) {
	log.Printf("GET: %s", c.baseURL+URL)
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
	c.eventBus.publish(event.NewSystemEvent(event.ClientConnected))
	defer func() {
		// send disconnect event
		c.eventBus.publish(event.NewSystemEvent(event.ClientDisconnected))
	}()

	state := 0
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		state++
		line := scanner.Text()
		if line == "" {
			// Move back to waiting state
			if state != eventStateFinished {
				log.Printf("unexpected end of event data on state %d", state)
			}
			state = eventStateWaiting
			continue
		}
		if state == eventStateBanner {
			if !strings.HasPrefix(line, eventHeader) {
				log.Printf("unexpected start of event: %q", line)
			}
			event := strings.TrimPrefix(line, eventHeader)
			if event != "message" {
				log.Printf("unexpected event: %q", event)
			}
			continue
		}
		if state == eventStateData {
			if !strings.HasPrefix(line, eventData) {
				log.Printf("unexpected event data: %q", line)
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
	decoder := json.NewDecoder(strings.NewReader(data))
	message := api.EventMessage{}
	err := decoder.Decode(&message)
	if err != nil {
		log.Printf("invalid event data: %s", err)
		return
	}
	switch message.Type {
	case api.EventItemCommand:
		e, err := event.NewItemReceivedCommand(message.Topic, message.Payload)
		if err != nil {
			log.Printf("error decoding message: %s", err)
			break
		}
		c.eventBus.publish(e)

	case api.EventItemState:
		e, err := event.NewItemReceivedState(message.Topic, message.Payload)
		if err != nil {
			log.Printf("error decoding message: %s", err)
			break
		}
		c.eventBus.publish(e)

	case api.EventItemStateChanged:
		e, err := event.NewItemChanged(message.Topic, message.Payload)
		if err != nil {
			log.Printf("error decoding message: %s", err)
			break
		}
		c.eventBus.publish(e)

	default:
		log.Printf("EVENT: %s on %s", message.Type, message.Topic)
	}
}

// eventLoop listen to the events from the REST api and send them to the event bus.
// the method never returns: if the connection drops it tries to reconnect in a loop
func (c *Client) eventLoop() {
	for {
		err := c.listenEvents()
		if err != nil {
			log.Printf("error listening to openhab events: %s", err)
		}
		time.Sleep(1 * time.Second)
	}
}

// func (c *Client) Subscribe(eventType event.Type, topic string, callback func(e event.Event)) {
// 	c.eventBus.subscribe(topic, eventType, callback)
// }

func (c *Client) AddRule(ruleData RuleData, run Runner, triggers ...Trigger) error {
	rule := newRule(c, ruleData, run, triggers)
	c.rules = append(c.rules, rule)
	return nil
}

// Start the handling of the defined rules.
// The function will return after the process received a Terminate, Abort or Interrupt signal,
// and after all the currently running rules have finished
func (c *Client) Start() {
	for _, rule := range c.rules {
		err := rule.activate(c)
		if err != nil {
			ruleName := rule.String()
			if ruleName != "" {
				ruleName = " \"" + ruleName + "\""
			}
			log.Printf("error activating rule%s: %s", ruleName, err)
		}
	}
	c.cron.Start()

	// start the event bus
	go c.eventLoop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

	// Wait until we're politely asked to leave
	<-stop

	log.Printf("shutting down...")
	ctx := c.cron.Stop()
	// Wait until all the cron tasks finished running
	<-ctx.Done()
}

func (c *Client) Subscribers() []string {
	subs := make([]string, len(c.eventBus.subs))
	for i, sub := range c.eventBus.subs {
		subs[i] = fmt.Sprintf("id=%d; topic=%q", sub.id, sub.topic)
	}
	return subs
}
