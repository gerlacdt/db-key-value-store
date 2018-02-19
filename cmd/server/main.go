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
	"github.com/gerlacdt/db-key-value-store/pkg/handler"
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

	f, err := os.Create(config.Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	defer os.Remove(config.Filename)

	h, err := handler.New(db.New(f))
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{Addr: ":" + config.Port, Handler: h}
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
