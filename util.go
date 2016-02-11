package main

import "encoding/binary"

// Convert int to 8-byte big endian representation.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
