package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jekabolt/grbpwr-manager/app"
	"github.com/jekabolt/grbpwr-manager/config"
	"github.com/jekabolt/grbpwr-manager/internal/store"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func run(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("cannot load a config %v", err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.Level(cfg.Logger.Level),
		AddSource: cfg.Logger.AddSource,
	}))
	slog.SetDefault(logger)

	repo, err := store.New(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("cannot create a repository %v", err.Error())
	}

	a := app.New(cfg, repo)
	if err := a.Start(ctx); err != nil {
		return fmt.Errorf("cannot start the application %v", err.Error())
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case s := <-sigCh:
		logger.With("signal", s.String()).Warn("signal received, exiting")
		a.Stop(ctx)
		logger.Info("application exited")
	case <-a.Done():
		logger.Error("application exited")
	}

	return nil
}
