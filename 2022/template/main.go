package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	path := "test.txt"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	input, _ := os.ReadFile(path)

	fmt.Println("Part 1:", part1(bytes.NewReader(input)))
	fmt.Println("Part 2:", part2(bytes.NewReader(input)))
}

func part1(r io.Reader) int {
	var totalScore int

	s := bufio.NewScanner(r)
	for s.Scan() {

	}

	return totalScore
}

func part2(r io.Reader) int {
	var totalScore int

	s := bufio.NewScanner(r)
	for s.Scan() {

	}

	return totalScore
}
