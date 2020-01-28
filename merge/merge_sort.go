package mergesort

// Sort realises the sorting of the data and returns it
func Sort(data []byte) []byte {
	if len(data) == 1 {
		// Already sorted
		return data
	}

	middle := len(data) / 2

	left, right := make([]byte, middle), make([]byte, len(data)-middle)

	for i := 0; i < len(data); i++ {
		if i < middle {
			left[i] = data[i]
		} else {
			right[i-middle] = data[i]
		}
	}

	return Merge(Sort(left), Sort(right))
}

// Merge merges the two slices and return the resulting slice
func Merge(left []byte, right []byte) (merged []byte) {
	merged = make([]byte, len(left)+len(right))

	i := 0
	for len(left) > 0 && len(right) > 0 {
		if left[0] <= right[0] {
			merged[i] = left[0]
			left = left[1:]
		} else {
			merged[i] = right[0]
			right = right[1:]
		}
		i++
	}

	for j := 0; j < len(left); j++ {
		merged[i] = left[j]
		i++
	}

	for j := 0; j < len(right); j++ {
		merged[i] = right[j]
		i++
	}

	return
}
