package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"

	"github.com/gerlacdt/db-key-value-store/pkg/db"
)

func main() {
	var config struct {
		Port     string `required:"true"`
		Filename string `required:"true"`
	}
	if err := envconfig.Process("db_key_value_store", &config); err != nil {
		fmt.Fprintln(os.Stderr, err)
		envconfig.Usage("db_key_value_store", &config)
		os.Exit(1)
	}

	dbs, err := db.New(config.Filename)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: db.NewMainHandler(dbs),
	}

	go func() {
		// graceful shutdown
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
		<-interrupt
		log.Println("App is shutting down...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down: %v\n", err)
		}
	}()

	log.Printf("App is ready to listen and serve on port %s\n", config.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Println("Good bye")
}
