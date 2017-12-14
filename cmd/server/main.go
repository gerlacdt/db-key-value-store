package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gerlacdt/db-example/pkg/db"
)

func main() {
	fmt.Println("start app...")
	err := http.ListenAndServe(":8080", handler())
	if err != nil {
		log.Fatalf("error listening %v", err)
	}
}

func handler() http.Handler {
	r := http.NewServeMux()

	mydb := db.NewDb("app.db.bin")
	service := db.NewService(mydb)
	myhandler := db.NewHandler(service)
	r.Handle("/db/", db.ErrorMiddleware(myhandler.HandleDb))
	return r
}
