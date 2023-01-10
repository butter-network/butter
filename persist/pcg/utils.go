package pcg

// MbToBytes converts a megabyte value to bytes.
func MbToBytes(mb uint64) uint64 {
	return uint64(mb * 1024 * 1024)
}

// MaxStorage returns the maximum amount of Groups a node can naivePersist.
func MaxStorage(maxMemory uint64) uint64 {
	return maxMemory / uint64(GroupStructSize)
}
