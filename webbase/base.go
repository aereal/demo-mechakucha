package webbase

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	HostPort        string
	ShutdownTimeout time.Duration
}

func Run(config *Config, handler http.Handler) error {
	server := &http.Server{
		Addr:    config.HostPort,
		Handler: handler,
	}
	go graceful(server, config.ShutdownTimeout)

	log.Printf("starting server (config:%#v) ...", config)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func graceful(server *http.Server, timeout time.Duration) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	log.Printf("shutting down server (%v) ...", sig)
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown: %v", err)
	}
}
