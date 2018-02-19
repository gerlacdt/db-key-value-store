package db

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/gerlacdt/db-key-value-store/pb"
	"github.com/golang/protobuf/proto"
)

// DB type
type DB struct {
	f       io.ReadWriteSeeker
	offsets map[string]int64
}

// New return a new intialized DB.
func New(f io.ReadWriteSeeker) *DB {
	offsetMap := make(map[string]int64)
	db := &DB{f: f, offsets: offsetMap}
	return db
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

func (db *DB) pbAppend(entity *pb.Entity) (int64, error) {
	entityBytes, err := proto.Marshal(entity)
	if err != nil {
		return 0, fmt.Errorf("pb marshall error %v", err)
	}
	byteBuffer := writeBinaryBufferLength(entityBytes)
	offset, err := db.f.Seek(0, 2)
	if err != nil {
		return 0, fmt.Errorf("file seek error %v", err)
	}
	_, err = byteBuffer.Write(entityBytes)
	if err != nil {
		return 0, fmt.Errorf("error writing byte buffer %v", err)
	}
	_, err = db.f.Write(byteBuffer.Bytes())
	if err != nil {
		return 0, fmt.Errorf("entity size file write error %v", err)
	}
	if err != nil {
		return 0, fmt.Errorf("entity data file write error %v", err)
	}
	return offset, nil
}

// Set / stores a key-value pair in the database
func (db *DB) Set(entity *pb.Entity) error {
	offset, err := db.pbAppend(entity)
	if err != nil {
		return err
	}
	db.offsets[entity.Key] = offset
	return nil
}

// Delete an entry for given key from database
func (db *DB) Delete(key string) error {
	entity := &pb.Entity{Tombstone: true, Key: key}
	offset, err := db.pbAppend(entity)
	if err != nil {
		return err
	}
	db.offsets[key] = offset
	return nil
}

// Get a key-value pair from the database
func (db *DB) Get(key string) (*pb.Entity, error) {
	offset, ok := db.offsets[key]
	if !ok {
		return nil, nil
	}
	_, err := db.f.Seek(offset, 0)
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
		return nil, nil
	}
	return entity, nil
}

func (db *DB) readSize() (uint64, error) {
	intsize := 8
	byteBuffer := make([]byte, intsize)
	_, err := db.f.Read(byteBuffer)
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
func (db *DB) Recover() error {
	// start reading file at beginning
	offset := int64(0)
	_, err := db.f.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("file seek error %v", err)
	}
	// run through all key-value pairs and populate in-memory hashmap
	for i := 0; i < 8; i++ {
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
		db.offsets[entity.Key] = offset
		offset += int64(size) + int64(8) // calculate next offset
	}
	return nil
}

func (db *DB) readPbData(lengthOf uint64) (*pb.Entity, error) {
	dataBuf := make([]byte, lengthOf)
	_, err := db.f.Read(dataBuf)
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
