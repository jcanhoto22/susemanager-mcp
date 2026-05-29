package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"mcp-susemanager/internal/api"
	"mcp-susemanager/internal/config"
	"mcp-susemanager/internal/logger"
	"mcp-susemanager/internal/mcp"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(cfg.Log.Level, cfg.Log.Format)
	log.Info("starting mcp-susemanager")
	log.Debug("configuration loaded", "config", cfg.Redacted())

	if cfg.SUSE.Password == "" {
		log.Warn("SUSE password is empty. Set SUSE_PASSWORD env var or config.yaml")
	}

	apiClient, err := api.New(cfg, log)
	if err != nil {
		log.Error("failed to create API client", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30)
	defer cancel()

	if err := apiClient.Login(ctx); err != nil {
		log.Error("failed to login to SUSE Manager", "error", err)
		os.Exit(1)
	}
	log.Info("authenticated to SUSE Manager")

	mcpServer, err := mcp.New(cfg, apiClient, log)
	if err != nil {
		log.Error("failed to create MCP server", "error", err)
		apiClient.Logout(context.Background())
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := mcpServer.Serve(context.Background()); err != nil {
			log.Error("MCP server error", "error", err)
			os.Exit(1)
		}
	}()

	sig := <-sigCh
	log.Info("received signal, shutting down", "signal", sig)

	logoutCtx, logoutCancel := context.WithTimeout(context.Background(), 10)
	defer logoutCancel()
	if err := apiClient.Logout(logoutCtx); err != nil {
		log.Warn("logout warning", "error", err)
	}

	log.Info("shutdown complete")
}

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))
}
