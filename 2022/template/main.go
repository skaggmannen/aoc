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
