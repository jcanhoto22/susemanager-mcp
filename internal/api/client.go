package api

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"mcp-susemanager/internal/config"
	"mcp-susemanager/internal/models"
	xrpc "mcp-susemanager/pkg/xmlrpc"
)

type Client struct {
	mu        sync.RWMutex
	xmlrpc    *xrpc.Client
	session   string
	logger    *slog.Logger
	cfg       config.SUSEConfig
	lastLogin time.Time
}

func New(cfg *config.Config, logger *slog.Logger) (*Client, error) {
	opts := []xrpc.Option{
		xrpc.WithLogger(logger),
		xrpc.WithTimeout(cfg.SUSE.Timeout),
	}
	if cfg.SUSE.InsecureSkipVerify {
		opts = append(opts, xrpc.WithInsecureSkipVerify())
	}

	xmlCli := xrpc.New(cfg.SUSE.URL, opts...)

	cli := &Client{
		xmlrpc: xmlCli,
		logger: logger.With("component", "api-client"),
		cfg:    cfg.SUSE,
	}

	return cli, nil
}

func (c *Client) Login(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var sessionKey string
	args := []interface{}{c.cfg.Username, c.cfg.Password}
	err := c.xmlrpc.Call(ctx, "auth.login", args, &sessionKey)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	c.session = sessionKey
	c.lastLogin = time.Now()
	c.logger.Info("login successful")
	return nil
}

func (c *Client) Logout(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.session == "" {
		return nil
	}

	var result int
	err := c.xmlrpc.Call(ctx, "auth.logout", []interface{}{c.session}, &result)
	if err != nil {
		c.logger.Warn("logout warning", "error", err)
	}
	c.session = ""
	return nil
}

func (c *Client) SessionKey() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.session
}

func (c *Client) IsLoggedIn() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.session != ""
}

func (c *Client) doCall(ctx context.Context, method string, extraArgs []interface{}, reply interface{}) error {
	c.mu.RLock()
	session := c.session
	c.mu.RUnlock()

	if session == "" {
		return fmt.Errorf("not authenticated: call auth.login first")
	}

	args := []interface{}{session}
	args = append(args, extraArgs...)

	err := c.xmlrpc.Call(ctx, method, args, reply)
	if err != nil {
		if isSessionExpired(err) {
			c.logger.Warn("session expired, re-authenticating")
			if loginErr := c.Login(ctx); loginErr != nil {
				return fmt.Errorf("re-authentication failed: %w", loginErr)
			}
			return c.RetryCall(ctx, method, extraArgs, reply, 1)
		}
		return err
	}
	return nil
}

func (c *Client) RetryCall(ctx context.Context, method string, args []interface{}, reply interface{}, maxRetries int) error {
	return c.xmlrpc.CallWithRetry(ctx, method, append([]interface{}{c.SessionKey()}, args...), reply, maxRetries)
}

func (c *Client) Call(ctx context.Context, method string, extraArgs []interface{}, reply interface{}) error {
	return c.doCall(ctx, method, extraArgs, reply)
}

func (c *Client) GetVersion(ctx context.Context) (*models.ApiVersion, error) {
	var result struct {
		Version string `xmlrpc:"version"`
	}
	err := c.Call(ctx, "api.getVersion", nil, &result)
	if err != nil {
		return nil, err
	}
	return &models.ApiVersion{Version: result.Version}, nil
}

func (c *Client) TestConnection(ctx context.Context) (bool, string, error) {
	version, err := c.GetVersion(ctx)
	if err != nil {
		return false, "", fmt.Errorf("connection test failed: %w", err)
	}
	return true, version.Version, nil
}

func isSessionExpired(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	keywords := []string{"session", "expired", "invalid session", "not authenticated", "IwmAuthFault"}
	for _, kw := range keywords {
		if contains(errMsg, kw) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && search(s, substr)
}

func search(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
