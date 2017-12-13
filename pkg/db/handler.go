package db

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gerlacdt/db-example/pb"
)

// Handler holds all http methods
type Handler struct {
	service *Service
}

// NewHandler is a constructor of all handlers
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
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
func (h *Handler) HandleDb(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.getHandler(w, r)
	} else if r.Method == "POST" {
		h.setHandler(w, r)
	} else if r.Method == "DELETE" {
		h.deleteHandler(w, r)
	}
}

func (h *Handler) setHandler(w http.ResponseWriter, r *http.Request) {
	key, err := getID(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	entity := &pb.Entity{Key: key, Value: body}
	err = h.service.Set(entity)
	if err != nil {
		fmt.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getHandler(w http.ResponseWriter, r *http.Request) {
	key, err := getID(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	entity, err := h.service.Get(key)
	if err != nil {
		fmt.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "GET, %v", entity)
}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	key, err := getID(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.service.Delete(key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
