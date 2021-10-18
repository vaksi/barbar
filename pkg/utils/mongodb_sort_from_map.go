package utils

func MongoSetSortFromMap(val map[string]bool) map[string]int {
	newVal := make(map[string]int)
	for i, v := range val {
		if v {
			newVal[i] = -1
		} else {
			newVal[i] = 1
		}
	}
	return newVal
}
