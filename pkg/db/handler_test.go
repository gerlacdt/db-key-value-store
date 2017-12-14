package db

import "testing"

func TestSingleHttpGet(t *testing.T) {
	//
}

// func TestSingleGet(t *testing.T) {
// 	before(testdb)
// 	defer teardown(testdb)
// 	db := NewDb(testdb)
// 	key := "foo-key"
// 	value := "foo-value"
// 	entity := &pb.Entity{Tombstone: false, Key: key, Value: []byte(value)}
// 	err := db.Set(entity)
// 	if err != nil {
// 		t.Fatalf("error append")
// 	}
// 	readEntity, err := db.Get(key)
// 	if err != nil {
// 		t.Fatalf("error getting entity %v", err)
// 	}
// 	if !reflect.DeepEqual(entity, readEntity) {
// 		t.Fatalf("expected %v, got %v", entity, readEntity)
// 	}
// }
