package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gerlacdt/db-key-value-store/pb"
	"github.com/gerlacdt/db-key-value-store/pkg/db"
)

// handler holds all http methods.
type handler struct{ db *db.DB }

// New creates all http handlers.
func New(db *db.DB) (http.Handler, error) {
	r := http.NewServeMux()

	if err := db.Recover(); err != nil {
		return nil, fmt.Errorf("could not recover database: %v", err)
	}

	h := &handler{db}
	r.Handle("/db/", errorMiddleware(h.handleDb))
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Handle("/version", errorMiddleware(versionHandler))
	return r, nil
}

func versionHandler(w http.ResponseWriter, r *http.Request) error {
	info := struct {
		BuildTime string `json:"buildTime"`
		Commit    string `json:"commit"`
		Release   string `json:"release"`
	}{
		BuildTime, Commit, Release,
	}
	body, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("Could not encode version data: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		return fmt.Errorf("error writing to http response: %v", err)
	}
	return nil
}

// httpError contains http status code
type httpError struct {
	err  error
	code int
	msg  string
}

// errorf constructor
func errorf(err error, code int, msg string) error {
	return &httpError{err: err, code: code, msg: msg}
}

func (e *httpError) Error() string { return fmt.Sprintf("%v: %v [%d]", e.err, e.msg, e.code) }

// errorMiddleware wraps a normal handler and converts errors to corresponding http status codes
type errorMiddleware func(http.ResponseWriter, *http.Request) error

func (fn errorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil {
		return
	}

	msg := err.Error()
	jsonError := struct {
		msg string `json:"msg"`
	}{msg}

	w.Header().Set("Content-Type", "application/json")

	if herr, ok := err.(*httpError); ok {
		w.WriteHeader(herr.code)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	body, err := json.Marshal(jsonError)
	if err != nil {
		fmt.Printf("Could not encode error data: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
	}
	if _, err = w.Write(body); err != nil {
		http.Error(w, msg, http.StatusInternalServerError)
	}
}

func getID(s string) (string, error) {
	id := strings.TrimPrefix(s, "/db/")
	arr := strings.Split(id, "/")
	if len(arr) != 1 {
		return "", fmt.Errorf("id contains a slash /, %s", id)
	}
	return id, nil
}

func (h *handler) handleDb(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return h.getHandler(w, r)
	case http.MethodPost:
		return h.setHandler(w, r)
	case http.MethodDelete:
		return h.deleteHandler(w, r)
	default:
		return errorf(fmt.Errorf(""), http.StatusMethodNotAllowed, r.Method+": method not allowed")
	}
}

func (h *handler) setHandler(w http.ResponseWriter, r *http.Request) error {
	key, err := getID(r.URL.Path)
	if err != nil {
		return errorf(err, http.StatusBadRequest, "requested key not valid")
	}

	if r.Header.Get("Content-Type") != "application/octet-stream" {
		return errorf(fmt.Errorf(""), http.StatusBadRequest, "Mime-Type not supported, application/octet-stream is supported")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	entity := &pb.Entity{Key: key, Value: body}
	err = h.db.Set(entity)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *handler) getHandler(w http.ResponseWriter, r *http.Request) error {
	key, err := getID(r.URL.Path)
	if err != nil {
		return errorf(err, http.StatusBadRequest, "requested key not valid")
	}
	entity, err := h.db.Get(key)
	if err != nil {
		return errorf(err, http.StatusInternalServerError, "error GET the requested key")
	}
	if entity == nil {
		return errorf(fmt.Errorf(""), http.StatusNotFound, "key does not exist")
	}

	if r.URL.Query().Get("format") == "json" {
		w.Header().Set("Content-Type", "application/json")
	}
	fmt.Fprintf(w, "%s", entity.Value)
	return nil
}

func (h *handler) deleteHandler(w http.ResponseWriter, r *http.Request) error {
	key, err := getID(r.URL.Path)
	if err != nil {
		return errorf(err, http.StatusBadRequest, "requested key not valid")
	}
	err = h.db.Delete(key)
	if err != nil {
		return errorf(err, http.StatusBadRequest, "requested key does not exist")
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
