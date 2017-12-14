package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gerlacdt/db-example/pkg/config"
	"github.com/gerlacdt/db-example/pkg/db"
)

func main() {
	config := config.NewConfig()
	fmt.Println("start app...")
	err := http.ListenAndServe(":"+config.Port, handler(config))
	if err != nil {
		log.Fatalf("error listening %v", err)
	}
}

func handler(config *config.Config) http.Handler {
	r := http.NewServeMux()
	mydb := db.NewDb(config.Filename)
	service := db.NewService(mydb)
	myhandler := db.NewHandler(service)
	r.Handle("/db/", db.ErrorMiddleware(myhandler.HandleDb))
	return r
}
