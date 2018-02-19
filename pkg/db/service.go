package db

import "github.com/gerlacdt/db-key-value-store/pb"

// Service provides all database methods
type Service struct {
	db                *Db
	SetRequestChannel chan *SetRequest
}

// SetRequest is input for the worker non-buffered channel. Response will be
// send to the resultChannel
type SetRequest struct {
	Entity        *pb.Entity
	ResultChannel chan error
}

// NewService creates a new db-service based on the given db reference
func NewService(db *Db) *Service {
	c := make(chan *SetRequest)
	service := &Service{db: db, SetRequestChannel: c}

	err := db.Recover()
	if err != nil {
		panic("error recovering from given filename: " + err.Error())
	}

	// start single background thread for SET request in order to prevent
	// concurrency issues during save.
	go func(queryChan <-chan *SetRequest) {
		for request := range queryChan {
			err := service.db.Set(request.Entity)
			request.ResultChannel <- err
		}
	}(c)

	return service
}

// Set and stores given entity
func (service *Service) Set(entity *pb.Entity) error {
	c := make(chan error)
	service.SetRequestChannel <- &SetRequest{Entity: entity, ResultChannel: c}
	err := <-c
	if err != nil {
		return err
	}
	return nil
}

// Get entity from database with given key
func (service *Service) Get(key string) (*pb.Entity, error) {
	entity, err := service.db.Get(key)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Delete entity with given key from database
func (service *Service) Delete(key string) error {
	err := service.db.Delete(key)
	if err != nil {
		return err
	}
	return nil
}
