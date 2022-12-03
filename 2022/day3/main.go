package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	input, _ := io.ReadAll(os.Stdin)

	fmt.Println("Part 1:", part1(bytes.NewReader(input)))
	fmt.Println("Part 2:", part2(bytes.NewReader(input)))
}

func part1(r io.Reader) int {
	var totalScore int

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		firstCompartement := make(map[rune]struct{})
		for _, r := range line[:len(line)/2] {
			firstCompartement[r] = struct{}{}
		}

		for _, r := range line[len(line)/2:] {
			if _, ok := firstCompartement[r]; ok {
				totalScore += itemPrio(r)
				break
			}
		}
	}

	return totalScore
}

func part2(r io.Reader) int {
	var totalScore int

	s := bufio.NewScanner(r)
	for s.Scan() {
		first := s.Text()
		s.Scan()

		second := s.Text()
		s.Scan()

		third := s.Text()

		firstItems := make(map[rune]struct{})
		for _, r := range first {
			firstItems[r] = struct{}{}
		}

		alsoInSecond := make(map[rune]struct{})
		for _, r := range second {
			if _, ok := firstItems[r]; ok {
				alsoInSecond[r] = struct{}{}
			}
		}

		for _, r := range third {
			if _, ok := alsoInSecond[r]; ok {
				totalScore += itemPrio(r)
				break
			}
		}
	}

	return totalScore
}

func itemPrio(r rune) int {
	if r >= 'a' && r <= 'z' {
		return int(r) - int('a') + 1
	} else {
		return int(r) - int('A') + 27
	}
}
