package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/genwilliam/short_code/internal"
	"github.com/genwilliam/short_code/internal/short_link"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to the service configuration file")
	flag.Parse()
	if err := run(*configPath); err != nil {
		log.Fatal(err)
	}
}

func run(configPath string) error {
	config, err := internal.LoadConfig(configPath)
	if err != nil {
		return err
	}

	db, err := internal.InitDB(config.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	server := &http.Server{
		Addr:              config.Server.Address,
		Handler:           internal.NewHandler(short_link.NewShortLinkRepository(db)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("short link service listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdownContext)
	}
}
