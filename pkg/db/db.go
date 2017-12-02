package db

import (
	"bytes"
	"encoding/binary"
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
	_, err := db.fileWrite.Write(keySizeBuf.Bytes())
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
	key := []byte(item.Key)
	err := db.Append(key, item.Value)
	if err != nil {
		return err
	}
	offset, err := db.fileWrite.Seek(0, 1)
	if err != nil {
		return err
	}
	db.offsetMap[item.Key] = offset
	return nil
}

// Get a key-value pair from the database
func (db *Db) Get(key string) *Entity {

	return nil
}

// ReadAll of a file and return all entries in the database
func (db *Db) ReadAll() []*Entity {
	return nil
}

// NewDb return a new intialized Db
func NewDb(filename string) *Db {
	fileWrite, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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
