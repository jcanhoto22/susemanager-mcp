package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SUSE   SUSEConfig   `yaml:"suse"`
	Server ServerConfig `yaml:"server"`
	Log    LogConfig    `yaml:"logging"`
}

type SUSEConfig struct {
	URL                string `yaml:"url"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
	Timeout            int    `yaml:"timeout"`
}

type ServerConfig struct {
	Port        int    `yaml:"port"`
	LogLevel    string `yaml:"log_level"`
	Transport   string `yaml:"transport"`
	MCPEndpoint string `yaml:"mcp_endpoint"`
}

type LogConfig struct {
	Format string `yaml:"format"`
	Level  string `yaml:"level"`
}

func Default() *Config {
	return &Config{
		SUSE: SUSEConfig{
			URL:                "https://localhost/rpc/api",
			Username:           "admin",
			Password:           "",
			InsecureSkipVerify: false,
			Timeout:            60,
		},
		Server: ServerConfig{
			Port:        8080,
			LogLevel:    "info",
			Transport:   "stdio",
			MCPEndpoint: "/mcp",
		},
		Log: LogConfig{
			Format: "text",
			Level:  "info",
		},
	}
}

func Load(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
	} else {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parsing config file: %w", err)
		}
	}

	cfg.overrideFromEnv()

	return cfg, nil
}

func (c *Config) overrideFromEnv() {
	if v := os.Getenv("SUSE_URL"); v != "" {
		c.SUSE.URL = v
	}
	if v := os.Getenv("SUSE_USERNAME"); v != "" {
		c.SUSE.Username = v
	}
	if v := os.Getenv("SUSE_PASSWORD"); v != "" {
		c.SUSE.Password = v
	}
	if v := os.Getenv("SUSE_INSECURE"); v != "" {
		c.SUSE.InsecureSkipVerify, _ = strconv.ParseBool(v)
	}
	if v := os.Getenv("SUSE_TIMEOUT"); v != "" {
		if t, err := strconv.Atoi(v); err == nil {
			c.SUSE.Timeout = t
		}
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			c.Server.Port = p
		}
	}
	if v := os.Getenv("SERVER_LOG_LEVEL"); v != "" {
		c.Server.LogLevel = v
	}
	if v := os.Getenv("SERVER_TRANSPORT"); v != "" {
		c.Server.Transport = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		c.Log.Level = v
	}
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		c.Log.Format = v
	}
}

func (c *Config) Redacted() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("SUSE URL: %s\n", c.SUSE.URL))
	b.WriteString(fmt.Sprintf("SUSE Username: %s\n", c.SUSE.Username))
	if c.SUSE.Password != "" {
		b.WriteString("SUSE Password: *****\n")
	}
	b.WriteString(fmt.Sprintf("SUSE InsecureSkipVerify: %v\n", c.SUSE.InsecureSkipVerify))
	b.WriteString(fmt.Sprintf("SUSE Timeout: %ds\n", c.SUSE.Timeout))
	b.WriteString(fmt.Sprintf("Server Port: %d\n", c.Server.Port))
	b.WriteString(fmt.Sprintf("Server LogLevel: %s\n", c.Server.LogLevel))
	b.WriteString(fmt.Sprintf("Server Transport: %s\n", c.Server.Transport))
	return b.String()
}
