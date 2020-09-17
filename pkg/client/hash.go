package client

import "hash/crc32"

func Hash(key string) (uint32, error) {
	return hashKeyCRC32(key), nil
}

func hashKeyCRC32(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}
