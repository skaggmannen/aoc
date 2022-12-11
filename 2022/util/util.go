package util

import (
	"bufio"
	"io"
	"strconv"
)

func ToInts(ss []string) []int {
	result := make([]int, len(ss))
	for i, s := range ss {
		result[i], _ = strconv.Atoi(s)
	}
	return result
}

func ReadAllLines(r io.Reader) []string {
	lines := make([]string, 0, 1024)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		lines = append(lines, line)
	}
	return lines
}
