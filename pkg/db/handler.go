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

// HTTPError contains http status code
type HTTPError struct {
	StatusCode int
	Message    string
}

// NewHTTPError constructor
func NewHTTPError(statusCode int, message string) *HTTPError {
	return &HTTPError{StatusCode: statusCode, Message: message}
}

func (e *HTTPError) Error() string {
	return e.Message
}

// ErrorMiddleware wraps a normal handler and converts errors to corresponding http status codes
type ErrorMiddleware func(http.ResponseWriter, *http.Request) error

func (fn ErrorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err != nil {
		switch err.(type) {
		case *HTTPError:
			e := err.(*HTTPError)
			w.WriteHeader(e.StatusCode)
			http.Error(w, e.Message, e.StatusCode)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
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
func (h *Handler) HandleDb(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return h.getHandler(w, r)
	} else if r.Method == "POST" {
		return h.setHandler(w, r)
	} else if r.Method == "DELETE" {
		return h.deleteHandler(w, r)
	}
	return NewHTTPError(http.StatusNotFound, "NOT FOUND")
}

func (h *Handler) setHandler(w http.ResponseWriter, r *http.Request) error {
	key, err := getID(r.URL.Path)
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, "requested key not valid")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	entity := &pb.Entity{Key: key, Value: body}
	err = h.service.Set(entity)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *Handler) getHandler(w http.ResponseWriter, r *http.Request) error {
	key, err := getID(r.URL.Path)
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, "requested key not valid")
	}
	entity, err := h.service.Get(key)
	if err != nil {
		return NewHTTPError(http.StatusNotFound, "key not found")
	}
	fmt.Fprintf(w, "GET, %v", entity)
	return nil
}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request) error {
	key, err := getID(r.URL.Path)
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, "requested key not valid")
	}
	err = h.service.Delete(key)
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, "requested key does not exist")
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
