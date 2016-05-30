package util

// Chunk separates a slice into multiple slices of a given size
func Chunk(slice []string, size int) [][]string {
	ret := [][]string{}
	l := len(slice)
	n := l / size
	for i := 0; i <= n; i++ {
		first := i * size
		last := min(l, (i+1)*size)
		ret = append(ret, slice[first:last])
	}
	return ret
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
