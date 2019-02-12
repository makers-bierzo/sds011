package sds011

import "encoding/binary"

func sum(input []byte) byte {
	sum := byte(0)

	for i := range input {
		sum += input[i]
	}
	return sum
}

func parseId(id uint16) []byte {
	idBuffer := make([]byte, 2)
	binary.LittleEndian.PutUint16(idBuffer, id)
	return idBuffer
}