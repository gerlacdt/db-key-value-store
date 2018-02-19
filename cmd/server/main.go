package main

import (
	"context"
	"fmt"
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
	fmt.Println("start app...")
	handler := db.NewMainHandler(config.Filename)

	// graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: handler,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	fmt.Println("App is ready to listen and serve.")

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		fmt.Println("Got SIGINT...")
	case syscall.SIGTERM:
		fmt.Println("got SIGTERM...")
	}
	fmt.Println("App is shutting down...")
	err := srv.Shutdown(context.Background())
	if err != nil {
		fmt.Printf("Error shutting down: %v\n", err)
	}
	fmt.Println("Done")
}
