package db

import (
	"fmt"
	"net/http"
	"strings"
)

// Handler holds all http methods
type Handler struct {
	db *Db
}

// NewHandler is a constructor of all handlers
func NewHandler(db *Db) *Handler {
	return &Handler{db: db}
}

func getID(s string) (string, error) {
	id := strings.TrimPrefix(s, "/db/")
	arr := strings.Split(id, "/")
	if len(arr) != 1 {
		return "", fmt.Errorf("id contains a slash /, %s", id)
	}
	return id, nil
}

// HandleDb handles all http routes
func HandleDb(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getHandler(w, r)
	} else if r.Method == "POST" {
		setHandler(w, r)
	} else if r.Method == "DELETE" {
		deleteHandler(w, r)
	}
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "SET, %q", id)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "GET, %q", id)

}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "DELETE, %q", id)
}
