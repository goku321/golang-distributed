package main

func mergeSortedSlice(x []string, y []string) ([]string) {
	result := make([]string, len(x) + len(y))

	i := 0
	for len(x) > 0 && len(y) > 0 {
		if x[0] < y[0] {
			result[i] = x[0]
			x = x[1:]
		} else {
			result[i] = y[0]
			y = y[1:]
		}
		i++
	}

	for j := 0; j < len(x); j++ {
        result[i] = x[j]
        i++
    }
    for j := 0; j < len(y); j++ {
        result[i] = y[j]
        i++
	}
	
	return result
}