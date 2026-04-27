package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/krau/SaveAny-Bot/bot"
	"github.com/krau/SaveAny-Bot/config"
	"github.com/krau/SaveAny-Bot/logger"
	"github.com/krau/SaveAny-Bot/storage"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	logger.L.Infof("SaveAny-Bot %s (%s) built at %s", Version, Commit, BuildTime)

	// Load configuration
	if err := config.Init(); err != nil {
		logger.L.Fatalf("Failed to initialize config: %v", err)
	}

	// Initialize storage backends
	if err := storage.Init(); err != nil {
		logger.L.Fatalf("Failed to initialize storage: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the Telegram bot
	if err := bot.Start(ctx); err != nil {
		logger.L.Fatalf("Failed to start bot: %v", err)
	}

	// Wait for termination signal
	// Also handle SIGHUP so the process can be gracefully stopped by some process managers
	// Note: SIGUSR1 is also watched here to allow manual reload triggers in my setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1)
	sig := <-quit
	logger.L.Infof("Received signal: %s", sig)

	logger.L.Info("Shutting down SaveAny-Bot...")
	cancel()

	if err := bot.Stop(); err != nil {
		logger.L.Errorf("Error stopping bot: %v", err)
	}

	logger.L.Info("SaveAny-Bot stopped")
}
