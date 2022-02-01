package utils

func ParsePacket(packet []byte) ([]byte, []byte) {
	for i, v := range packet {
		if v == '/' { // looking for separator
			return packet[:i+1], packet[i+1:]
		}
	}
	return nil, nil
}
