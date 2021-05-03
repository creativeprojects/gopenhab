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

	"github.com/robfig/cron/v3"
)

type Client struct {
	config  Config
	baseURL string
	client  *http.Client
	cron    *cron.Cron
	items   *Items
	rules   []*Rule
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

func (c *Client) Subscribe(topic string) error {
	resp, err := c.get(context.Background(), "events")
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (c *Client) AddRule(config RuleConfig, run func(), triggers ...Trigger) error {
	rule := NewRule(config, run, triggers)
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

	// Wait until we're politely asked to leave
	<-stop

	log.Printf("shutting down...")
	ctx := c.cron.Stop()
	// Wait until all the cron tasks finished running
	<-ctx.Done()
}
