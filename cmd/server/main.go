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
	handler := db.NewMainHandler(config.Filename)
	err := http.ListenAndServe(":"+config.Port, handler)
	if err != nil {
		log.Fatalf("error listening %v", err)
	}
}
