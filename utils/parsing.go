package utils

//func ParsePacket(packet []byte) (string, string) {
//	// Convert the packet byte buffer to a string
//	packetString := string(packet)
//
//	// Get the uri by splitting the packet at the first space
//	uri := strings.Split(packetString, " ")[0]
//
//	// Calculate the start of the payload (based on the length of the URI)
//	uriLength := len(uri)
//	startOfPayload := uriLength + 1
//
//	// Get the payload by getting everything after the first space
//	payload := packet[startOfPayload:]
//
//	return uri, payload
//}

func ParsePacket(packet []byte) (byte, []byte) {
	//for i, v := range packet {
	//	if v == 32 {
	//		return string(packet[:i]), packet[i+1:]
	//	}
	//}
	//return "", []byte{1}
	return packet[0], packet[1:]
}
