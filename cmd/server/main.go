package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gerlacdt/db-key-value-store/pkg/config"
	"github.com/gerlacdt/db-key-value-store/pkg/db"
)

func main() {
	config := config.NewConfig()
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: db.NewMainHandler(config.Filename),
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
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Println("Good bye")
}
