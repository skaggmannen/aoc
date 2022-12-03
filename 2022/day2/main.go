package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
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

		parts := strings.Split(line, " ")

		opponents := parts[0]
		my := parts[1]

		switch opponents {
		case "A": // Rock
			switch my {
			case "X": // Rock
				totalScore += 1 + 3 // Draw
			case "Y": // Paper
				totalScore += 2 + 6 // Win!
			case "Z": // Scissors
				totalScore += 3 + 0 // Lost :(
			}
		case "B": // Paper
			switch my {
			case "X": // Rock
				totalScore += 1 + 0 // Lost :(
			case "Y": // Paper
				totalScore += 2 + 3 // Draw
			case "Z": // Scissors
				totalScore += 3 + 6 // Won!
			}
		case "C": // Scissors
			switch my {
			case "X": // Rock
				totalScore += 1 + 6 // Won!
			case "Y": // Paper
				totalScore += 2 + 0 // Lost :(
			case "Z": // Scissors
				totalScore += 3 + 3 // Draw
			}
		}
	}

	return totalScore
}

func part2(r io.Reader) int {
	var totalScore int

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")

		opponents := parts[0]
		my := parts[1]

		switch opponents {
		case "A": // Rock
			switch my {
			case "X": // We should lose
				totalScore += 3 + 0 // Choose scissors
			case "Y": // We should draw
				totalScore += 1 + 3 // Choose rock
			case "Z": // We should win
				totalScore += 2 + 6 // Choose paper
			}
		case "B": // Paper
			switch my {
			case "X": // We should lose
				totalScore += 1 + 0 // Choose rock
			case "Y": // We should draw
				totalScore += 2 + 3 // Choose paper
			case "Z": // We should win
				totalScore += 3 + 6 // Choose scissors
			}
		case "C": // Scissors
			switch my {
			case "X": // We should lose
				totalScore += 2 + 0 // Choose paper
			case "Y": // We should draw
				totalScore += 3 + 3 // Choose scissors
			case "Z": // We should win
				totalScore += 1 + 6 // Choose rock
			}
		}
	}

	return totalScore
}
