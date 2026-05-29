package xmlrpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/kolo/xmlrpc"
)

type Client struct {
	endpoint string
	client   *xmlrpc.Client
	httpCli  *http.Client
	logger   *slog.Logger
	timeout  time.Duration
}

type Option func(*Client)

func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

func WithTimeout(seconds int) Option {
	return func(c *Client) {
		c.timeout = time.Duration(seconds) * time.Second
	}
}

func WithInsecureSkipVerify() Option {
	return func(c *Client) {
		tr, ok := c.httpCli.Transport.(*http.Transport)
		if !ok {
			tr = &http.Transport{}
			c.httpCli.Transport = tr
		}
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{}
		}
		tr.TLSClientConfig.InsecureSkipVerify = true
	}
}

func New(endpoint string, opts ...Option) *Client {
	c := &Client{
		endpoint: endpoint,
		timeout:  60 * time.Second,
		logger:   slog.Default(),
		httpCli: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{},
			},
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.httpCli.Timeout = c.timeout

	var err error
	c.client, err = xmlrpc.NewClient(c.endpoint, c.httpCli.Transport)
	if err != nil {
		c.logger.Warn("xmlrpc client creation warning", "error", err)
	}

	return c
}

func (c *Client) Call(ctx context.Context, method string, args interface{}, reply interface{}) error {
	start := time.Now()

	c.logger.Debug("xmlrpc call",
		"method", method,
		"endpoint", c.endpoint,
	)

	err := c.client.Call(method, args, reply)

	elapsed := time.Since(start)
	if err != nil {
		c.logger.Warn("xmlrpc call failed",
			"method", method,
			"error", err,
			"elapsed", elapsed,
		)
		return fmt.Errorf("xmlrpc %s: %w", method, err)
	}

	c.logger.Debug("xmlrpc call success",
		"method", method,
		"elapsed", elapsed,
	)

	return nil
}

func (c *Client) CallWithRetry(ctx context.Context, method string, args interface{}, reply interface{}, maxRetries int) error {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Debug("retrying xmlrpc call",
				"method", method,
				"attempt", attempt,
				"max_retries", maxRetries,
			)
			backoff := time.Duration(attempt*attempt) * 200 * time.Millisecond
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		if err := c.Call(ctx, method, args, reply); err != nil {
			lastErr = err
			continue
		}
		return nil
	}
	return fmt.Errorf("xmlrpc call failed after %d retries: %w", maxRetries, lastErr)
}

func (c *Client) Close() error {
	return nil
}
