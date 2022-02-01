package node

// mbToBytes converts a megabyte value to bytes.
func mbToBytes(mb uint64) uint64 {
	return mb * 1024 * 1024
}

// naiveProcessData stores as much of the input data as possible in the Block data field and cuts off the rest. This
// will later be improved by breaking down the data into several Block(s) if necessary.
func naiveProcessData(data string) [3840]byte {
	var formattedData [3840]byte
	for i := range formattedData {
		formattedData[i] = data[i]
	}
	return formattedData
}

// naiveProcessKeywords stores as much of the specified keywords as possible in the Block keyword subfields and cuts
// off the rest any remaining part of the keyword if it exceeds 50 characters.
func naiveProcessKeywords(keywords []string) [5][50]byte {
	var formattedKeywords [5][50]byte
	for i := range formattedKeywords {
		var word [50]byte
		for j := range word {
			formattedKeywords[i][j] = keywords[i][j]
		}
	}
	return formattedKeywords
}
