package node

// mbToBytes convert an input value in mb to a value in bytes
func mbToBytes(mb uint64) uint64 {
	return mb * 1024 * 1024
}
