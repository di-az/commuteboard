package main

import (
	"commuteboard/internal/config"
	"commuteboard/internal/db"
	"commuteboard/internal/engine"
	"commuteboard/internal/server"
	"commuteboard/internal/store"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// const UpdateRate = 10 * time.Minute
// const tickRate = 10 * time.Second

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	sqliteDB, err := db.NewSQLite(cfg.SQLiteDB)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting engine")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store := store.NewRouteStore(sqliteDB)

	engine, err := engine.NewRouteEngine(
		ctx,
		store,
		cfg.UpdateRate,
		cfg.TickRate,
		cfg.GoogleMapsAPIKey,
	)
	if err != nil {
		log.Fatal(err)
	}
	server := server.NewHttpServer(store, engine)

	// Run HTTP server
	go func() {
		if err := server.Run(ctx); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Run engine
	go engine.Run(ctx)

	<-ctx.Done()
	log.Println("Shutting down application...")
}
