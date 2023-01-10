package storageOverlay

// naiveProcessData stores as much of the input data as possible in the Block data field and cuts off the rest. This
// will later be improved by breaking down the data into several Block(s) if necessary.
func naiveProcessData(data string) [3840]byte {
	var formattedData [3840]byte
	slice := []byte(data)
	copy(formattedData[:], slice)
	return formattedData
}

// naiveProcessKeywords stores as much of the specified keywords as possible in the Block keyword subfields and cuts
// off the rest any remaining part of the keyword if it exceeds 50 characters. Requirement is that there is always 5
// keywords.
func naiveProcessKeywords(keywords []string) [5][50]byte {
	var formattedKeywords [5][50]byte
	for i := range formattedKeywords {
		var keyword [50]byte
		slice := []byte(keywords[i])
		copy(keyword[:], slice)
		formattedKeywords[i] = keyword
	}
	return formattedKeywords
}
