package db

import (
	"os"
	"reflect"
	"testing"

	"github.com/gerlacdt/db-example/pb"
)

var testdb = "db.test.bin"

func clean(filename string) error {
	err := os.Remove(filename)

	if err != nil {
		return err
	}
	return nil
}

func before(filename string) {
	err := clean(filename)
	if err != nil {
		// panic("could not delete db file")
	}
}

func teardown(filename string) {
	err := clean(filename)
	if err != nil {
		// panic("could not delete db file")
	}
}

func TestSinglePbGet(t *testing.T) {
	// prepare
	before(testdb)
	defer teardown(testdb)
	db := NewDb(testdb)

	key := "foo-key"
	value := "foo-value"
	entity := &pb.Entity{Tombstone: false, Key: key, Value: []byte(value)}

	// act
	err := db.Set(entity)
	if err != nil {
		t.Fatalf("error append")
	}

	readEntity, err := db.Get(key)
	if err != nil {
		t.Fatalf("error getting entity %v", err)
	}

	// assert
	if !reflect.DeepEqual(entity, readEntity) {
		t.Fatalf("expected %v, got %v", entity, readEntity)
	}
}

func TestMultiplePbGet(t *testing.T) {
	// prepare
	before(testdb)
	defer teardown(testdb)
	db := NewDb(testdb)
	key := "foo-key"
	value := "foo-value"
	entity := &pb.Entity{Tombstone: false, Key: key, Value: []byte(value)}
	err := db.Set(entity)
	if err != nil {
		t.Fatalf("error setting entity 1")
	}

	key1 := "foo-key-1"
	value1 := "foo-value-1"
	entity1 := &pb.Entity{Tombstone: false, Key: key1, Value: []byte(value1)}
	err = db.Set(entity1)
	if err != nil {
		t.Fatalf("error setting entity 2")
	}

	// act
	readEntity, err := db.Get(key)
	if err != nil {
		t.Fatalf("error getting entity %v", err)
	}

	readEntity1, err := db.Get(key1)

	// assert
	if !reflect.DeepEqual(entity, readEntity) {
		t.Fatalf("expected %v, got %v", entity, readEntity)
	}
	if !reflect.DeepEqual(entity1, readEntity1) {
		t.Fatalf("expected %v, got %v", entity1, readEntity1)
	}
}
