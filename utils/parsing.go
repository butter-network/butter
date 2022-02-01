package utils

import "errors"

// ParsePacket into its component route URI and payload
func ParsePacket(packet []byte) ([]byte, []byte, error) {
	for i, v := range packet {
		if v == '/' { // looking for separator
			return packet[:i+1], packet[i+1:], nil
		}
	}
	return nil, nil, errors.New("incorrectly formatted packet")
}
