package openhab

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	config  Config
	baseURL string
	client  *http.Client
	items   *Items
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
