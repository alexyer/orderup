package orderup

import "encoding/binary"

// Convert int to 8-byte big endian representation.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// Convert 8-byte big endian representation to int64.
func btoi(bytes []byte) uint64 {
	val := binary.BigEndian.Uint64(bytes)
	return val
}
