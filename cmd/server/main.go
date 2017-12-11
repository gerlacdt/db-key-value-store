package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gerlacdt/db-example/pkg/db"
)

func main() {
	fmt.Println("start app...")
	http.HandleFunc("/db/", db.HandleDb)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
