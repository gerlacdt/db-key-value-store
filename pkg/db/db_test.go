package db

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"
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

func TestWriteBinaryLength(t *testing.T) {

	data := []byte("hello")
	buf := writeBinaryBufferLength(data)

	var i uint64
	err := binary.Read(buf, binary.LittleEndian, &i)
	if err != nil {
		t.Fatalf("error reading binary err: %v, data: %v", err, buf.Bytes())
	}

	if i != uint64(len(data)) {
		t.Fatalf("Expected %d, got %d", len(data), i)
	}
}

func TestWriteBinary(t *testing.T) {
	data := []byte("hello")
	buf := writeBinaryBuffer(data)

	var readData = make([]byte, len(data))
	err := binary.Read(buf, binary.LittleEndian, readData)
	if err != nil {
		t.Fatalf("error reading binary err: %v, data: %v", err, buf.Bytes())
	}

	if !bytes.Equal(data, readData) {
		t.Fatalf("expected %s, got %s", data, readData)
	}
}

type testEntry struct {
	Key       []byte
	KeySize   int
	Value     []byte
	ValueSize int
}

func readKeyValue(file *os.File, key []byte, data []byte, t *testing.T) *testEntry {
	keySize := readSize(file, key, t)
	readKey := readData(file, key, t)

	dataSize := readSize(file, data, t)
	readData := readData(file, data, t)

	return &testEntry{Key: readKey, KeySize: keySize, Value: readData, ValueSize: dataSize}
}

func readSize(file *os.File, data []byte, t *testing.T) int {
	intsize := 8
	sizeBuf := make([]byte, intsize)
	n, err := file.Read(sizeBuf)
	if err != nil {
		t.Fatalf("error reading file for size %v", err)
	}
	var b = bytes.NewReader(sizeBuf)
	var readSize uint64
	err = binary.Read(b, binary.LittleEndian, &readSize)
	if n != intsize {
		t.Fatalf("error reading file, n byte not matching for size, expected %d, got %d", intsize, n)
	}

	if readSize != uint64(len(data)) {
		t.Fatalf("readSize expected %d, got %d", len(data), readSize)
	}
	return n
}

func readData(file *os.File, data []byte, t *testing.T) []byte {
	dataBuf := make([]byte, len(data))
	n, err := file.Read(dataBuf)
	if err != nil {
		t.Fatalf("error reading file for data %v", err)
	}

	var b = bytes.NewReader(dataBuf)
	bbuf := make([]byte, len(data))
	err = binary.Read(b, binary.LittleEndian, bbuf)
	if err != nil {
		t.Fatalf("error reading binary err: %v, data: %v", err, bbuf)
	}

	if n != len(data) {
		t.Fatalf("size expected %d, got %d", len(data), n)
	}
	if !bytes.Equal(data, dataBuf) {
		t.Fatalf("data expected %s, got %s", data, dataBuf)
	}
	return dataBuf
}

func TestSingleAppend(t *testing.T) {
	before(testdb)
	defer teardown(testdb)

	key := []byte("foo-key")
	data := []byte("foo-value")
	db := NewDb(testdb)
	err := db.Append(key, data)
	if err != nil {
		t.Fatalf("error append")
	}
	file, err := os.OpenFile(testdb, os.O_RDONLY, 644)
	if err != nil {
		t.Fatalf("error open file for reading %v", err)
	}
	entry := readKeyValue(file, key, data, t)
	if !bytes.Equal(key, entry.Key) {
		t.Fatalf("data expected %s, got %s", key, entry.Key)
	}
	if !bytes.Equal(data, entry.Value) {
		t.Fatalf("data expected %s, got %s", data, entry.Value)
	}
}

func TestMultiAppend(t *testing.T) {
	before(testdb)
	defer teardown(testdb)
	key := []byte("k1")
	data := []byte("h1")
	key2 := []byte("k2")
	data2 := []byte("h2")
	db := NewDb(testdb)
	err := db.Append(key, data)
	if err != nil {
		t.Fatalf("error append")
	}
	err = db.Append(key2, data2)
	if err != nil {
		t.Fatalf("error append")
	}
	file, err := os.OpenFile(testdb, os.O_RDONLY, 644)
	if err != nil {
		t.Fatalf("error open file for reading %v", err)
	}
	entry := readKeyValue(file, key, data, t)
	if !bytes.Equal(key, entry.Key) {
		t.Fatalf("data expected %s, got %s", key, entry.Key)
	}
	if !bytes.Equal(data, entry.Value) {
		t.Fatalf("data expected %s, got %s", data, entry.Value)
	}
	entry = readKeyValue(file, key2, data2, t)
	if !bytes.Equal(key2, entry.Key) {
		t.Fatalf("data expected %s, got %s", key, entry.Key)
	}
	if !bytes.Equal(data2, entry.Value) {
		t.Fatalf("data expected %s, got %s", data, entry.Value)
	}
}

func TestSingleSet(t *testing.T) {
	before(testdb)
	defer teardown(testdb)

	key := "k1"
	data := []byte("v1")
	db := NewDb(testdb)
	err := db.Set(&Entity{Key: key, Value: data})
	if err != nil {
		t.Fatalf("error db.set: %v", err)
	}
	offset := db.offsetMap[key]
	expectedOffset := int64(0)
	if offset != expectedOffset {
		t.Fatalf("expected offset %d, got %d", expectedOffset, offset)
	}
}

func TestMultiSet(t *testing.T) {
	before(testdb)
	defer teardown(testdb)

	key := "k1"
	data := []byte("v1")
	db := NewDb(testdb)
	err := db.Set(&Entity{Key: key, Value: data})
	if err != nil {
		t.Fatalf("error db.set: %v", err)
	}
	offset := db.offsetMap[key]
	expectedOffset := int64(0)
	if offset != expectedOffset {
		t.Fatalf("expected offset %d, got %d", expectedOffset, offset)
	}

	key = "k2"
	data = []byte("v2")
	err = db.Set(&Entity{Key: key, Value: data})
	if err != nil {
		t.Fatalf("error db.set: %v", err)
	}
	offset = db.offsetMap[key]
	expectedOffset = int64(20)
	if offset != expectedOffset {
		t.Fatalf("expected offset %d, got %d", expectedOffset, offset)
	}

}

func TestSingleGet(t *testing.T) {
	before(testdb)
	defer teardown(testdb)

	key := "k12"
	data := []byte("v12")
	db := NewDb(testdb)
	err := db.Set(&Entity{Key: key, Value: data})
	if err != nil {
		t.Fatalf("error db.set: %v", err)
	}

	item, err := db.Get(key)
	if err != nil {
		t.Fatalf("error GETTING from database %v", err)
	}
	if item.Key != key {
		t.Fatalf("expected %v, got %v", key, item.Key)
	}

	if !bytes.Equal(data, item.Value) {
		t.Fatalf("expected %s, got %s", data, item.Value)
	}
}

func TestMultiGet(t *testing.T) {
	before(testdb)
	defer teardown(testdb)

	key := "foo"
	data := []byte("bar")
	db := NewDb(testdb)
	err := db.Set(&Entity{Key: key, Value: data})
	if err != nil {
		t.Fatalf("error db.set: %v", err)
	}

	item, err := db.Get(key)
	if err != nil {
		t.Fatalf("error GETTING from database %v", err)
	}
	if item.Key != key {
		t.Fatalf("expected %v, got %v", key, item.Key)
	}

	if !bytes.Equal(data, item.Value) {
		t.Fatalf("expected %s, got %s", data, item.Value)
	}

	key = "key1"
	data = []byte("value1")
	err = db.Set(&Entity{Key: key, Value: data})
	if err != nil {
		t.Fatalf("error db.set: %v", err)
	}

	item, err = db.Get(key)
	if err != nil {
		t.Fatalf("error GETTING from database %v", err)
	}
	if item.Key != key {
		t.Fatalf("expected %v, got %v", key, item.Key)
	}

	if !bytes.Equal(data, item.Value) {
		t.Fatalf("expected %s, got %s", data, item.Value)
	}
}
