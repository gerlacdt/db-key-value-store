package db

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/gerlacdt/db-key-value-store/pb"
	"github.com/mattetti/filebuffer"
)

func setup(t *testing.T) *DB {
	t.Parallel()
	return New(filebuffer.New(nil))
}

func TestSingleGet(t *testing.T) {
	db := setup(t)
	key := "foo-key"
	value := "foo-value"
	entity := &pb.Entity{Tombstone: false, Key: key, Value: []byte(value)}
	err := db.Set(entity)
	if err != nil {
		t.Fatalf("error append")
	}
	readEntity, err := db.Get(key)
	if err != nil {
		t.Fatalf("error getting entity %v", err)
	}
	if !reflect.DeepEqual(entity, readEntity) {
		t.Fatalf("expected %v, got %v", entity, readEntity)
	}
}

func TestMultipleGet(t *testing.T) {
	db := setup(t)
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
	readEntity, err := db.Get(key)
	if err != nil {
		t.Fatalf("error getting entity %v", err)
	}
	readEntity1, err := db.Get(key1)
	if !reflect.DeepEqual(entity, readEntity) {
		t.Fatalf("expected %v, got %v", entity, readEntity)
	}
	if !reflect.DeepEqual(entity1, readEntity1) {
		t.Fatalf("expected %v, got %v", entity1, readEntity1)
	}
}

func TestSingleDelete(t *testing.T) {
	db := setup(t)
	key := "foo-key"
	value := "foo-value"
	entity := &pb.Entity{Tombstone: false, Key: key, Value: []byte(value)}
	err := db.Set(entity)
	if err != nil {
		t.Fatalf("error append")
	}
	err = db.Delete(key)
	readEntity, err := db.Get(key)
	if readEntity != nil || err != nil {
		t.Fatalf("readEntity expected nil, got %v", readEntity)
	}
}

func TestSingleRecover(t *testing.T) {
	db := setup(t)
	key := "foo-key"
	value := "foo-value"
	entity := &pb.Entity{Tombstone: false, Key: key, Value: []byte(value)}
	err := db.Set(entity)
	if err != nil {
		t.Fatalf("error append")
	}

	// clear map
	db.offsets = make(map[string]int64)

	err = db.Recover()
	if err != nil {
		t.Fatalf("error recovering %v", err)
	}

	readEntity, err := db.Get(key)
	if err != nil {
		t.Fatalf("error deleting entity %v", err)
	}
	if !reflect.DeepEqual(entity, readEntity) {
		t.Fatalf("error entities not equal after recovering")
	}

}

func TestSingleRecoverWithDelete(t *testing.T) {
	db := setup(t)
	key := "foo-key"
	value := "foo-value"
	entity := &pb.Entity{Tombstone: false, Key: key, Value: []byte(value)}
	err := db.Set(entity)
	if err != nil {
		t.Fatalf("error SET")
	}

	err = db.Delete(key)
	if err != nil {
		t.Fatalf("error DELETE")
	}

	// clear map
	db.offsets = make(map[string]int64)

	err = db.Recover()
	if err != nil {
		t.Fatalf("error recovering %v", err)
	}

	readEntity, err := db.Get(key)
	if readEntity != nil || err != nil {
		t.Fatalf("readEntity expected nil, got %v", readEntity)
	}
}

func TestMultipleRecover(t *testing.T) {
	db := setup(t)

	// first item
	key := "foo-key"
	value := "foo-value"
	entity := &pb.Entity{Tombstone: false, Key: key, Value: []byte(value)}
	err := db.Set(entity)
	if err != nil {
		t.Fatalf("error setting entity 1")
	}

	// second item
	key1 := "foo-key-1"
	value1 := "foo-value-1"
	entity1 := &pb.Entity{Tombstone: false, Key: key1, Value: []byte(value1)}
	err = db.Set(entity1)
	if err != nil {
		t.Fatalf("error setting entity 2")
	}

	// third item
	key2 := "foo-key-2"
	value2 := "foo-value-2"
	entity2 := &pb.Entity{Tombstone: false, Key: key2, Value: []byte(value2)}
	err = db.Set(entity2)
	if err != nil {
		t.Fatalf("error setting entity 3")
	}

	// act
	// clear map
	db.offsets = make(map[string]int64)
	err = db.Recover()
	if err != nil {
		t.Fatalf("error recovering %v", err)
	}
	readEntity, err := db.Get(key)
	if err != nil {
		t.Fatalf("error getting entity %v", err)
	}
	readEntity1, err := db.Get(key1)
	if err != nil {
		t.Fatalf("error getting entity1 %v", err)
	}
	readEntity2, err := db.Get(key2)
	if err != nil {
		t.Fatalf("error getting entity2 %v", err)
	}

	// assert
	if !reflect.DeepEqual(entity, readEntity) {
		t.Fatalf("expected %v, got %v", entity, readEntity)
	}
	if !reflect.DeepEqual(entity1, readEntity1) {
		t.Fatalf("expected %v, got %v", entity1, readEntity1)
	}
	if !reflect.DeepEqual(entity2, readEntity2) {
		t.Fatalf("expected %v, got %v", entity2, readEntity2)
	}
}

func TestConcurrentSets(t *testing.T) {
	db := setup(t)

	var wg sync.WaitGroup
	maxItems := 1000
	buffChan := make(chan int, maxItems)
	for i := 0; i < maxItems; i++ {
		buffChan <- i
	}
	close(buffChan)

	maxConcurrency := 4
	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			for j := range buffChan {
				// set foo-key-i and foo-value-i as entry in db
				key := "foo-key-" + strconv.Itoa(j)
				value := "foo-value-" + strconv.Itoa(j)
				err := db.Set(&pb.Entity{Key: key, Value: []byte(value)})
				if err != nil {
					fmt.Printf("error inserting key-value: [T%d] %d\n", index, j)
				}
			}
		}(i)
	}

	wg.Wait() // wait for all goroutines to finish

	// check if all key-values are inserted correctly
	mapLen := len(db.offsets)
	if maxItems != mapLen {
		t.Fatalf("mapLen: expected %d, got %d", maxItems, mapLen)
	}

	for i := 0; i < maxItems; i++ {
		expectedKey := "foo-key-" + strconv.Itoa(i)
		expectedValue := "foo-value-" + strconv.Itoa(i)
		entity, err := db.Get(expectedKey)
		if err != nil {
			fmt.Printf("error getting key-value: %d\n", i)
		}
		if string(expectedValue) != string(entity.Value) {
			t.Fatalf("value expected %v, got %v", expectedValue, entity.Value)
		}
	}
}
