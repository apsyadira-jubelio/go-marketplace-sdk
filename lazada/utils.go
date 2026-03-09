package lazada

func SplitFileToBlocks(file []byte, maxBlockSize int) [][]byte {
	var blocks [][]byte
	for start := 0; start < len(file); start += maxBlockSize {
		end := start + maxBlockSize
		if end > len(file) {
			end = len(file)
		}
		blocks = append(blocks, file[start:end])
	}
	return blocks
}
