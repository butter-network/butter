package node

func mbToBytes(mb uint64) uint64 {
	return mb * 1024 * 1024
}

// just store as much of the data as possible - cut off the rest
func naiveProcessData(data string) [3840]byte {
	var formattedData [3840]byte
	for i, _ := range formattedData {
		formattedData[i] = data[i]
	}
	return formattedData
}

func processKeywords(keywords []string) [5][50]byte {
	var formattedKeywords [5][50]byte
	for i, _ := range formattedKeywords {
		var word [50]byte
		for j, _ := range word {
			formattedKeywords[i][j] = keywords[i][j]
		}
	}
	return formattedKeywords
}
