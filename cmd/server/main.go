package main

import (
	"context"
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
	if err := envconfig.Process("", &config); err != nil {
		log.Print(err)
		envconfig.Usage("", &config)
		os.Exit(1)
	}

	f, err := os.Create(config.Filename)
	if err != nil {
		log.Printf("could not create file %s: %v", config.Filename, err)
		os.Exit(1)
	}
	defer f.Close()
	defer os.Remove(config.Filename)

	h, err := handler.New(db.New(f))
	if err != nil {
		log.Printf("could not create handler: %v", err)
		os.Exit(1)
	}

	srv := &http.Server{Addr: ":" + config.Port, Handler: h}
	go func() {
		// graceful shutdown
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
		<-interrupt
		log.Print("app is shutting down...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("could not shutdown: %v\n", err)
		}
	}()

	log.Printf("app is ready to listen and serve on port %s", config.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("server failed: %v", err)
		os.Exit(1)
	}

	log.Print("good bye!")
}
