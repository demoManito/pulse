package http

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/demoManito/pulse/config"
	"github.com/demoManito/pulse/internal/server/httpsvc"
	"github.com/demoManito/pulse/pkg/logger"
	"golang.org/x/sync/errgroup"
)

// Command server
type Command struct{}

// Name command's name
func (*Command) Name() string {
	return "http"
}

// Run command
func (*Command) Run(cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	server := httpsvc.NewServer(cfg)
	g.Go(func() error {
		logger.Infof("Starting HTTP server on %v", server.Addr)
		return server.ListenAndServe()
	})
	g.Go(func() error {
		<-ctx.Done()
		logger.Info("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Errorf("Server forced to shutdown: %v", err)
			return err
		}
		logger.Info("Server exiting")
		return nil
	})

	osCh := make(chan os.Signal, 1)
	signal.Notify(osCh, os.Interrupt, syscall.SIGTERM)
	select {
	case sig := <-osCh:
		logger.Infof("Received signal: %s", sig)
		cancel()
	case <-ctx.Done():
	}

	err := g.Wait()
	if err != nil {
		logger.Errorf("Http exited with error: %v", err)
		// PASS
	}

	logger.Info("All tasks completed, exiting")

	return nil
}
