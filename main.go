package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/time-sheet/config"
	"github.com/time-sheet/mongo"
	"github.com/time-sheet/service"
	"github.com/time-sheet/store"
)

// TimeSheetAPIStore is a wrapper which embeds a Mongo stuct which satisfies the store.Storer interface.
type TimeSheetAPIStore struct {
	*mongo.Mongo
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.Get()
	if err != nil {
		fmt.Printf("errored getting configuration: [%v]\n", err)
		os.Exit(1)
	}

	fmt.Printf("config on startup: [%v]\n", cfg)

	mongodb := &mongo.Mongo{
		Collection: cfg.MongoConfig.Collection,
		Database:   cfg.MongoConfig.Database,
		URI:        cfg.MongoConfig.BindAddr,
	}

	session, err := mongodb.Init()
	if err != nil {
		fmt.Printf("failed to initialise mongo: [%v]", err)
		os.Exit(1)
	}

	mongodb.Session = session

	store := store.DataStore{Backend: TimeSheetAPIStore{Mongo: mongodb}}

	apiErrors := make(chan error, 1)

	service.CreateTimeSheetAPI(*cfg, store, apiErrors)

	// Gracefully shutdown the application closing any open resources.
	gracefulShutdown := func() {
		fmt.Printf("shutdown with timeout: %s\n", cfg.GracefulShutdownTimeout)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)

		// stop any incoming requests before closing any outbound connections
		service.Close(ctx)

		fmt.Println("shutdown complete")

		cancel()
		os.Exit(1)
	}

	for {
		select {
		case err := <-apiErrors:
			fmt.Printf("api error received: [%v]\n", err)
			gracefulShutdown()
		case <-signals:
			fmt.Println("os signal received")
			gracefulShutdown()
		}
	}
}
