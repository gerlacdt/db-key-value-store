package db

import "github.com/gerlacdt/db-example/pb"

// Service provides all database methods
type Service struct {
	db *Db
}

// Set and stores given entity
func (service *Service) Set(entity *pb.Entity) error {
	err := service.db.Set(entity)
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
