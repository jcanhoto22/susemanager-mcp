package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	logger := New("info", "text")
	assert.NotNil(t, logger)

	loggerJSON := New("debug", "json")
	assert.NotNil(t, loggerJSON)
}

func TestNewWriter(t *testing.T) {
	logger := New("info", "text")
	writer := NewWriter(logger, slog.LevelInfo)
	assert.NotNil(t, writer)

	n, err := writer.Write([]byte("test message"))
	assert.NoError(t, err)
	assert.Greater(t, n, 0)
}

func TestNewWriterEmptyMessage(t *testing.T) {
	logger := New("info", "text")
	writer := NewWriter(logger, slog.LevelInfo)

	n, err := writer.Write([]byte("\n"))
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
}
