package db

import (
	"fmt"

	"github.com/gerlacdt/db-key-value-store/pb"
)

// Service provides all database methods
type Service struct {
	*DB
}

// request is input for the worker non-buffered channel. Response will be
// send to the result
type request struct {
	entity *pb.Entity
	result chan error
}

// NewService creates a new db-service based on the given db reference
func NewService(db *DB) (*Service, error) {
	if err := db.Recover(); err != nil {
		return nil, fmt.Errorf("error recovering from given filename: %v", err)
	}
	svc := &Service{DB: db}
	return svc, nil
}
