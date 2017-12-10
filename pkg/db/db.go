package db

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gerlacdt/db-example/pb"
	"github.com/golang/protobuf/proto"
)

// Db type
type Db struct {
	filename  string
	fileWrite *os.File
	fileRead  *os.File
	offsetMap map[string]int64
}

func writeBinaryBufferLength(data []byte) *bytes.Buffer {
	var length = uint64(len(data))
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, length)
	if err != nil {
		log.Fatalf("error writing data length %v", err)
	}
	return buf
}

func (db *Db) pbAppend(entity *pb.Entity) (int64, error) {
	entityBytes, err := proto.Marshal(entity)
	if err != nil {
		return 0, fmt.Errorf("pb marshall error %v", err)
	}
	byteBuffer := writeBinaryBufferLength(entityBytes)
	offset, err := db.fileWrite.Seek(0, 2)
	if err != nil {
		return 0, fmt.Errorf("file seek error %v", err)
	}
	_, err = byteBuffer.Write(entityBytes)
	if err != nil {
		return 0, fmt.Errorf("error writing byte buffer %v", err)
	}
	_, err = db.fileWrite.Write(byteBuffer.Bytes())
	if err != nil {
		return 0, fmt.Errorf("entity size file write error %v", err)
	}
	if err != nil {
		return 0, fmt.Errorf("entity data file write error %v", err)
	}
	return offset, nil
}

// Set / stores a key-value pair in the database
func (db *Db) Set(entity *pb.Entity) error {
	offset, err := db.pbAppend(entity)
	if err != nil {
		return err
	}
	db.offsetMap[entity.Key] = offset
	return nil
}

// Delete an entry for given key from database
func (db *Db) Delete(key string) error {
	entity := &pb.Entity{Tombstone: true, Key: key}
	offset, err := db.pbAppend(entity)
	if err != nil {
		return err
	}
	db.offsetMap[key] = offset
	return nil
}

// Get a key-value pair from the database
func (db *Db) Get(key string) (*pb.Entity, error) {
	offset, ok := db.offsetMap[key]
	if !ok {
		return nil, fmt.Errorf("Key not in database")
	}
	_, err := db.fileRead.Seek(offset, 0)
	if err != nil {
		return nil, fmt.Errorf("file seek error %v", err)
	}
	size, err := db.readSize()
	if err != nil {
		return nil, fmt.Errorf("read size error, %v", err)
	}
	entity, err := db.readPbData(size)
	if err != nil {
		return nil, fmt.Errorf("key readData error, %v", err)
	}
	if entity.Tombstone {
		return nil, fmt.Errorf("Key not in database (already deleted)")
	}
	return entity, nil
}

func (db *Db) readSize() (uint64, error) {
	intsize := 8
	byteBuffer := make([]byte, intsize)
	_, err := db.fileRead.Read(byteBuffer)
	if err != nil {
		return 0, err
	}
	var b = bytes.NewReader(byteBuffer)
	var readSize uint64
	err = binary.Read(b, binary.LittleEndian, &readSize)
	if err != nil {
		return 0, err
	}
	return readSize, nil
}

// Recover from a crash and populate in-memory hashmap from existing file
func (db *Db) Recover() error {
	// start reading file at beginning
	offset := int64(0)
	_, err := db.fileRead.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("file seek error %v", err)
	}
	// run through all key-value pairs and populate in-memory hashmap
	for i := 0; i < 8; i++ {
		fmt.Printf("begin offset %d\n", offset)
		size, err := db.readSize()
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read size error, %v", err)
		}
		entity, err := db.readPbData(size)
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("key readData error, %v", err)
		}
		db.offsetMap[entity.Key] = offset
		offset += int64(size) + int64(8) // calculate next offset
		fmt.Printf("end offset %d\n", offset)
	}
	return nil
}

func (db *Db) readPbData(lengthOf uint64) (*pb.Entity, error) {
	dataBuf := make([]byte, lengthOf)
	_, err := db.fileRead.Read(dataBuf)
	if err != nil {
		return nil, err
	}

	entity := &pb.Entity{}
	err = proto.Unmarshal(dataBuf, entity)

	if err != nil {
		return nil, fmt.Errorf("proto unmarshal error %v", err)
	}
	return entity, nil
}

// NewDb return a new intialized Db
func NewDb(filename string) *Db {
	fileWrite, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error file opening for write")
	}
	fileRead, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("error file opening for read")
	}
	offsetMap := make(map[string]int64)
	db := &Db{filename: filename, fileWrite: fileWrite, fileRead: fileRead, offsetMap: offsetMap}
	return db
}
