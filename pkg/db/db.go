package db

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

// Db type
type Db struct {
	filename  string
	fileWrite *os.File
	fileRead  *os.File
	offsetMap map[string]int64
}

// Entity is the default structure which is used for the database api
type Entity struct {
	Key   string
	Value []byte
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

func writeBinaryBuffer(data []byte) *bytes.Buffer {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, data)
	if err != nil {
		log.Fatalf("error writing data length %v", err)
	}
	return buf
}

// Append the given byte-array to file
func (db *Db) Append(key []byte, data []byte) error {
	keySizeBuf := writeBinaryBufferLength(key)
	keyBuf := writeBinaryBuffer(key)
	sizeBuf := writeBinaryBufferLength(data)
	dataBuf := writeBinaryBuffer(data)

	_, err := db.fileWrite.Seek(0, 2)
	if err != nil {
		return nil
	}
	_, err = db.fileWrite.Write(keySizeBuf.Bytes())
	if err != nil {
		return err
	}
	_, err = db.fileWrite.Write(keyBuf.Bytes())
	if err != nil {
		return nil
	}
	_, err = db.fileWrite.Write(sizeBuf.Bytes())
	if err != nil {
		return err
	}
	_, err = db.fileWrite.Write(dataBuf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// Set / stores a key-value pair in the database
func (db *Db) Set(item *Entity) error {
	offset, err := db.fileWrite.Seek(0, 2)
	key := []byte(item.Key)
	err = db.Append(key, item.Value)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	db.offsetMap[item.Key] = offset
	return nil
}

// Get a key-value pair from the database
func (db *Db) Get(key string) (*Entity, error) {
	offset, ok := db.offsetMap[key]
	if !ok {
		return nil, fmt.Errorf("Key not in database")
	}
	db.fileRead.Seek(offset, 0)
	keySize, err := db.readSize()
	if err != nil {
		return nil, fmt.Errorf("key readSize error, %v", err)
	}
	readKey, err := db.readData(keySize)
	if err != nil {
		return nil, fmt.Errorf("key readData error, %v", err)
	}
	dataSize, err := db.readSize()
	if err != nil {
		return nil, fmt.Errorf("data readSize error, %v", err)
	}
	readData, err := db.readData(dataSize)
	if err != nil {
		return nil, fmt.Errorf("data readData error, %v", err)
	}
	return &Entity{Key: string(readKey), Value: readData}, nil
}

func (db *Db) readSize() (uint64, error) {
	intsize := 8
	sizeBuf := make([]byte, intsize)
	_, err := db.fileRead.Read(sizeBuf)
	if err != nil {
		return 0, err
	}
	var b = bytes.NewReader(sizeBuf)
	var readSize uint64
	err = binary.Read(b, binary.LittleEndian, &readSize)
	return readSize, nil
}

func (db *Db) readData(lengthOf uint64) ([]byte, error) {
	dataBuf := make([]byte, lengthOf)
	_, err := db.fileRead.Read(dataBuf)
	if err != nil {
		return nil, err
	}
	var b = bytes.NewReader(dataBuf)
	bbuf := make([]byte, lengthOf)
	err = binary.Read(b, binary.LittleEndian, bbuf)
	if err != nil {
		return nil, err
	}
	return dataBuf, nil
}

// ReadAll of a file and return all entries in the database
func (db *Db) ReadAll() []*Entity {
	return nil
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
