package utils

import "sort"

func Include(arr []string, val string) bool {
	sort.Strings(arr)
	i := sort.SearchStrings(arr, val)
	if i >= len(arr) || arr[i] != val {
		return false
	}
	return true
}
