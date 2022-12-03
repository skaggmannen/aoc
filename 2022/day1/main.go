package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func main() {
	input, _ := io.ReadAll(os.Stdin)

	fmt.Println("Part 1:", part1(bytes.NewReader(input)))
	fmt.Println("Part 2:", part2(bytes.NewReader(input)))
}

func part1(r io.Reader) int {
	var totals []int
	var current int

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			totals = append(totals, current)
			current = 0
		} else {
			calories, err := strconv.Atoi(line)
			if err != nil {
				continue
			}

			current += calories
		}
	}

	totals = append(totals, current)

	sort.Ints(totals)

	return totals[len(totals)-1]
}

func part2(r io.Reader) int {
	var totals []int
	var current int

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			totals = append(totals, current)
			current = 0
		} else {
			calories, err := strconv.Atoi(line)
			if err != nil {
				continue
			}

			current += calories
		}
	}

	totals = append(totals, current)

	sort.Ints(totals)

	return totals[len(totals)-1] + totals[len(totals)-2] + totals[len(totals)-3]
}
