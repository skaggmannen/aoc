package util

import "strconv"

func ToInts(ss []string) []int {
	result := make([]int, len(ss))
	for i, s := range ss {
		result[i], _ = strconv.Atoi(s)
	}
	return result
}
