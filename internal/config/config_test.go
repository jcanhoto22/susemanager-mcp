package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	assert.Equal(t, "https://localhost/rpc/api", cfg.SUSE.URL)
	assert.Equal(t, "admin", cfg.SUSE.Username)
	assert.Equal(t, "", cfg.SUSE.Password)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "stdio", cfg.Server.Transport)
}

func TestLoadNonExistentFile(t *testing.T) {
	cfg, err := Load("/tmp/nonexistent-file.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, "https://localhost/rpc/api", cfg.SUSE.URL)
}

func TestLoadFromFile(t *testing.T) {
	content := []byte(`
suse:
  url: "https://suma.test/rpc/api"
  username: "testuser"
  password: "testpass"
  insecure_skip_verify: true
  timeout: 120
server:
  port: 9090
  log_level: debug
  transport: http
logging:
  format: json
  level: debug
`)
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(content)
	require.NoError(t, err)
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	require.NoError(t, err)
	assert.Equal(t, "https://suma.test/rpc/api", cfg.SUSE.URL)
	assert.Equal(t, "testuser", cfg.SUSE.Username)
	assert.Equal(t, "testpass", cfg.SUSE.Password)
	assert.True(t, cfg.SUSE.InsecureSkipVerify)
	assert.Equal(t, 120, cfg.SUSE.Timeout)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Server.LogLevel)
	assert.Equal(t, "http", cfg.Server.Transport)
}

func TestEnvOverride(t *testing.T) {
	os.Setenv("SUSE_URL", "https://env-override/rpc/api")
	os.Setenv("SUSE_USERNAME", "envuser")
	os.Setenv("SUSE_PASSWORD", "envpass")
	os.Setenv("SUSE_INSECURE", "true")
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Unsetenv("SUSE_URL")
	defer os.Unsetenv("SUSE_USERNAME")
	defer os.Unsetenv("SUSE_PASSWORD")
	defer os.Unsetenv("SUSE_INSECURE")
	defer os.Unsetenv("SERVER_PORT")
	defer os.Unsetenv("LOG_LEVEL")

	cfg, err := Load("/tmp/nonexistent.yaml")
	require.NoError(t, err)
	assert.Equal(t, "https://env-override/rpc/api", cfg.SUSE.URL)
	assert.Equal(t, "envuser", cfg.SUSE.Username)
	assert.Equal(t, "envpass", cfg.SUSE.Password)
	assert.True(t, cfg.SUSE.InsecureSkipVerify)
	assert.Equal(t, 3000, cfg.Server.Port)
}
