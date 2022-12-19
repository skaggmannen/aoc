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

type Stack[T any] []T

func (s *Stack[T]) Push(t T) {
	*s = append(*s, t)
}

func (s *Stack[T]) Pop() T {
	t := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return t
}

func (s *Stack[T]) Peek() T {
	return (*s)[len(*s)-1]
}

type Point struct {
	X, Y int
}
