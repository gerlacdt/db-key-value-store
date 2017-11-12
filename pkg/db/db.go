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
func (db *Db) Append(data []byte) error {
	sizeBuf := writeBinaryBufferLength(data)
	dataBuf := writeBinaryBuffer(data)
	_, err := db.fileWrite.Write(sizeBuf.Bytes())
	if err != nil {
		return err
	}
	_, err = db.fileWrite.Write(dataBuf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// ReadAll of file and return as string
func (db *Db) ReadAll() string {
	return ""
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
	db := &Db{filename: filename, fileWrite: fileWrite, fileRead: fileRead}
	return db
}
