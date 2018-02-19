package db

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/gerlacdt/db-key-value-store/pb"
)

func TestSingleServiceGet(t *testing.T) {
	db, teardown := setup(t)
	defer teardown()

	service := NewService(db)
	key := "foo-key"
	value := "foo-value"
	entity := &pb.Entity{Tombstone: false, Key: key, Value: []byte(value)}
	err := service.Set(entity)
	if err != nil {
		t.Fatalf("error append")
	}
	readEntity, err := service.Get(key)
	if err != nil {
		t.Fatalf("error getting entity %v", err)
	}
	if !reflect.DeepEqual(entity, readEntity) {
		t.Fatalf("expected %v, got %v", entity, readEntity)
	}
}

func TestSingleServiceDelete(t *testing.T) {
	// prepare
	db, teardown := setup(t)
	defer teardown()

	service := NewService(db)
	key := "foo-key"
	value := "foo-value"
	entity := &pb.Entity{Tombstone: false, Key: key, Value: []byte(value)}
	err := service.Set(entity)
	if err != nil {
		t.Fatalf("error append")
	}
	err = service.Delete(key)
	readEntity, err := service.Get(key)
	if readEntity != nil || err != nil {
		t.Fatalf("readEntity expected nil, got %v", readEntity)
	}
}

func TestMultipleServiceSet(t *testing.T) {
	// prepare testint
	db, teardown := setup(t)
	defer teardown()

	service := NewService(db)

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
				err := service.Set(&pb.Entity{Key: key, Value: []byte(value)})
				if err != nil {
					fmt.Printf("error inserting key-value: [T%d] %d\n", index, j)
				}
			}
		}(i)
	}

	wg.Wait() // wait for all goroutines to finish

	// check if all key-values are inserted correctly
	mapLen := len(service.db.offsets)
	if maxItems != mapLen {
		t.Fatalf("mapLen: expected %d, got %d", maxItems, mapLen)
	}

	for i := 0; i < maxItems; i++ {
		expectedKey := "foo-key-" + strconv.Itoa(i)
		expectedValue := "foo-value-" + strconv.Itoa(i)
		entity, err := service.Get(expectedKey)
		if err != nil {
			fmt.Printf("error getting key-value: %d\n", i)
		}
		if string(expectedValue) != string(entity.Value) {
			t.Fatalf("value expected %v, got %v", expectedValue, entity.Value)
		}
	}
}
