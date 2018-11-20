package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/s-take/goecho-nats-postgre-sample/db"
	"github.com/s-take/goecho-nats-postgre-sample/event"
	"github.com/s-take/goecho-nats-postgre-sample/retry"
)

// Config of Other Service
type Config struct {
	PostgresHost     string `envconfig:"POSTGRES_HOST"`
	PostgresPort     string `envconfig:"POSTGRES_PORT"`
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
	NatsClient       string `envconfig:"NATS_CLIENT"`
}

func main() {
	// Read Config
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to PostgreSQL
	retry.ForeverSleep(2*time.Second, func(attempt int) error {
		addr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)
		repo, err := db.NewPostgres(addr)
		if err != nil {
			log.Println(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

	// Connect to Nats
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := event.NewStan(fmt.Sprintf("nats://%s", cfg.NatsAddress), "test-cluster", cfg.NatsClient)
		if err != nil {
			log.Println(err)
			return err
		}
		err = es.OnTaskPublished(onTaskPublished)
		if err != nil {
			log.Println(err)
		}
		event.SetEventStore(es)
		return nil
	})
	defer event.Close()

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			event.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
