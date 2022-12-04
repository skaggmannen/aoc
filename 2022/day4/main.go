package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, _ := io.ReadAll(os.Stdin)

	fmt.Println("Part 1:", part1(bytes.NewReader(input)))
	fmt.Println("Part 2:", part2(bytes.NewReader(input)))
}

func part1(r io.Reader) int {
	var totalCount int

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		p1, p2 := toPairs(line)

		firstP1, lastP1 := toRange(p1)
		firstP2, lastP2 := toRange(p2)

		if firstP1 <= firstP2 && lastP1 >= lastP2 {
			totalCount++
		} else if firstP2 <= firstP1 && lastP2 >= lastP1 {
			totalCount++
		}
	}

	return totalCount
}

func part2(r io.Reader) int {
	var totalCount int

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		p1, p2 := toPairs(line)

		firstP1, lastP1 := toRange(p1)
		firstP2, lastP2 := toRange(p2)

		if firstP1 <= firstP2 && lastP1 >= lastP2 {
			// P1 fully covers P2
			totalCount++
		} else if firstP2 <= firstP1 && lastP2 >= lastP1 {
			// P2 fully covers P1
			totalCount++
		} else if firstP1 <= firstP2 && lastP1 >= firstP2 {
			// P1 covers start of P2
			totalCount++
		} else if firstP1 <= lastP2 && lastP1 >= lastP2 {
			// P1 covers end of P2
			totalCount++
		}
	}

	return totalCount
}

func toPairs(v string) (string, string) {
	parts := strings.Split(v, ",")
	if len(parts) != 2 {
		log.Fatal("invalid pair: ", v)
	}

	return parts[0], parts[1]
}

func toRange(v string) (int, int) {
	parts := strings.Split(v, "-")
	if len(parts) != 2 {
		log.Fatal("invalid range: ", v)
	}

	first, _ := strconv.Atoi(parts[0])
	last, _ := strconv.Atoi(parts[1])

	return first, last
}
