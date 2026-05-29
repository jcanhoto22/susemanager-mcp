package xmlrpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	client := New("https://example.com/rpc/api",
		WithTimeout(30),
		WithLogger(nil),
	)
	assert.NotNil(t, client)
	assert.Equal(t, 30*time.Second, client.timeout)
}

func TestNewWithInsecure(t *testing.T) {
	client := New("https://example.com/rpc/api",
		WithInsecureSkipVerify(),
	)
	assert.NotNil(t, client)
}

func TestCallWithoutServer(t *testing.T) {
	client := New("https://localhost:1/rpc/api",
		WithTimeout(1),
	)
	assert.NotNil(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var result interface{}
	err := client.Call(ctx, "test.method", nil, &result)
	assert.Error(t, err)
}

func TestCallWithRetry(t *testing.T) {
	client := New("https://localhost:1/rpc/api",
		WithTimeout(1),
	)
	assert.NotNil(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result interface{}
	err := client.CallWithRetry(ctx, "test.method", nil, &result, 2)
	assert.Error(t, err)
}
