package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
	var (
		totalScore int
	)

	// Cycles are 1 based...
	cycle := 1
	// Register X starts with value 1
	x := 1

	program := readAllLines(r)

	for _, line := range program {
		parts := strings.Split(line, " ")

		switch parts[0] {
		case "noop":
			cycle++
			if (cycle+20)%40 == 0 {
				fmt.Println(cycle, ": ", x)
				totalScore += cycle * x
			}
		case "addx":
			v, _ := strconv.Atoi(parts[1])
			if (cycle+20+2)%40 == 0 {
				fmt.Println(cycle+2, ": ", x+v)
				totalScore += (cycle + 2) * (x + v)
			} else if (cycle+20+1)%40 == 0 {
				fmt.Println(cycle+1, ": ", x)
				totalScore += (cycle + 1) * x
			}
			x += v
			cycle += 2
		}
	}

	return totalScore
}

func part2(r io.Reader) int {
	var (
		x  int
		pc int
	)

	x = 1

	s := bufio.NewScanner(r)

	for s.Scan() {
		line := s.Text()
		parts := strings.Split(line, " ")

		switch parts[0] {
		case "noop":
			draw(pc%40, x)
			pc++
		case "addx":
			v, _ := strconv.Atoi(parts[1])
			draw(pc%40, x)
			pc++
			draw(pc%40, x)
			pc++
			x += v
		}
	}

	return 0
}

func draw(pixel int, x int) {
	if pixel < x-1 || pixel > x+1 {
		fmt.Print(".")
	} else {
		fmt.Print("#")
	}
	if pixel == 39 {
		fmt.Println()
	}
}

func readAllLines(r io.Reader) []string {
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
