package db

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"
)

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

func readEntry(file *os.File, data []byte, t *testing.T) []byte {
	readSize(file, data, t)
	return readData(file, data, t)
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

func clean(filename string, t *testing.T) {
	err := os.Remove(filename)

	if err != nil {
		t.Fatalf("error deleting file %v", err)
	}
}

func TestSingleAppend(t *testing.T) {
	testdb := "db.test.bin"
	data := []byte("hello")
	db := NewDb(testdb)
	err := db.Append(data)
	if err != nil {
		t.Fatalf("error append")
	}
	file, err := os.OpenFile(testdb, os.O_RDONLY, 644)
	if err != nil {
		t.Fatalf("error open file for reading %v", err)
	}
	readEntry(file, data, t)
	clean(testdb, t)
}

func TestMultiAppend(t *testing.T) {
	testdb := "db.test.bin"
	data := []byte("hello")
	data2 := []byte("foo-world")
	db := NewDb(testdb)
	err := db.Append(data)
	if err != nil {
		t.Fatalf("error append")
	}
	err = db.Append(data2)
	if err != nil {
		t.Fatalf("error append")
	}
	file, err := os.OpenFile(testdb, os.O_RDONLY, 644)
	if err != nil {
		t.Fatalf("error open file for reading %v", err)
	}
	readEntry(file, data, t)
	readEntry(file, data2, t)
	clean(testdb, t)
}
