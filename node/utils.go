package node

// mbToBytes converts a megabyte value to bytes.
func mbToBytes(mb uint64) uint64 {
	return mb * 1024 * 1024
}
