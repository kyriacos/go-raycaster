package main

import (
	"encoding/binary"
)

func uintsToBytes(vs []uint32) []byte {
	buf := make([]byte, len(vs)*4)
	for i, v := range vs {
		binary.LittleEndian.PutUint32(buf[i*4:], v)
	}
	return buf
}

func bytesToUints(vs []byte) []uint32 {
	out := make([]uint32, len(vs)/4)
	for i := range out {
		out[i] = binary.LittleEndian.Uint32(vs[i*4:])
	}
	return out
}
