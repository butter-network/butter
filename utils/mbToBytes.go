package utils

// MbToBytes converts a megabyte value to bytes.
func MbToBytes(mb uint64) uint64 {
	return mb * 1024 * 1024
}
