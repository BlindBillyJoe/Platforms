package main

import (
	"io"
	"unsafe"
)

func readByte(reader io.Reader) uint8 {
	var buf = make([]byte, 1)
	var _, err = reader.Read(buf)
	check(err)
	return uint8(buf[0])
}

func readHWord(reader io.Reader) uint16 {
	return uint16(readByte(reader)) | uint16(readByte(reader))<<4
}

func readWord(reader io.Reader) uint32 {
	return uint32(readHWord(reader)) | uint32(readHWord(reader))<<8
}

func readDWord(reader io.Reader) uint64 {
	return uint64(readWord(reader)) | uint64(readWord(reader))<<16
}

func readUShort(reader io.Reader) uint16 {
	return uint16(readByte(reader))<<8 | uint16(readByte(reader)<<0)
}

func readShort(reader io.Reader) int16 {
	return int16(readUShort(reader))
}

func readULong(reader io.Reader) uint32 {
	return uint32(readByte(reader))<<24 | uint32(readByte(reader))<<16 | uint32(readByte(reader))<<8 | uint32(readByte(reader))<<0
}

func readFixed(reader io.Reader) uint32 {
	return readULong(reader)
}

func readTag(reader io.Reader) uint32 {
	return readWord(reader)
}

func readLongDateTime(reader io.Reader) int64 {
	return int64(readULong(reader))<<32 | int64(readULong(reader))<<0
}

func fixedToFloat(fixed uint32) float32 {
	var ptr = uintptr(unsafe.Pointer(&fixed))
	var to16 = func(p uintptr) uint16 {
		return *(*uint16)(unsafe.Pointer(p))
	}
	return float32(to16(ptr+2)) + float32(to16(ptr))/float32(1<<2)
}

func sToTag(s []uint8) uint32 {
	return uint32(s[0])<<0 | uint32(s[1])<<8 | uint32(s[2])<<16 | uint32(s[3])<<24
}
