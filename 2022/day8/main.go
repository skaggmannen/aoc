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
	grid := readGrid(r)

	fmt.Printf("%s\n", grid)

	lookFromLeft(grid)
	lookFromRight(grid)
	lookFromTop(grid)
	lookFromBottom(grid)

	fmt.Printf("%s\n", grid)

	return countVisible(grid)
}

type Square struct {
	height  int
	visible bool
	score   int
}

func (s Square) higherThan(other Square) bool {
	return s.height > other.height
}

func (s Square) String() string {
	if s.visible {
		return fmt.Sprintf("[%d]", s.height)
	} else {
		return fmt.Sprintf(" %d ", s.height)
	}
}

func readGrid(r io.Reader) Grid {
	var grid Grid
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		var row []Square
		for _, b := range line {
			v, _ := strconv.Atoi(string(b))
			row = append(row, Square{height: v})
		}
		grid = append(grid, row)
	}
	return grid
}

type Grid [][]Square

func (g Grid) String() string {
	var buf strings.Builder
	for _, row := range g {
		fmt.Fprintf(&buf, "%s\n", row)
	}
	return buf.String()
}

func lookFromLeft(grid Grid) {
	for row := range grid {
		highest := -1
		for col := range grid[row] {
			if grid[row][col].height > highest {
				grid[row][col].visible = true
				highest = grid[row][col].height
			}
		}
	}
}

func lookFromRight(grid Grid) {
	for row := range grid {
		highest := -1
		for col := len(grid[row]) - 1; col >= 0; col-- {
			if grid[row][col].height > highest {
				grid[row][col].visible = true
				highest = grid[row][col].height
			}
		}
	}
}

func lookFromTop(grid Grid) {
	for col := 0; col < len(grid[0]); col++ {
		highest := -1
		for row := range grid {
			if grid[row][col].height > highest {
				grid[row][col].visible = true
				highest = grid[row][col].height
			}
		}
	}
}

func lookFromBottom(grid Grid) {
	for col := 0; col < len(grid[0]); col++ {
		highest := -1
		for row := len(grid) - 1; row >= 0; row-- {
			if grid[row][col].height > highest {
				grid[row][col].visible = true
				highest = grid[row][col].height
			}
		}
	}
}

func countVisible(grid Grid) int {
	var count int
	for _, row := range grid {
		for _, s := range row {
			if s.visible {
				count++
			}
		}
	}
	return count
}

func part2(r io.Reader) int {
	grid := readGrid(r)

	for row := range grid {
		for col := range grid[row] {
			grid[row][col].score = calculateScore(grid, row, col)
			fmt.Printf("%4d", grid[row][col].score)
		}
		fmt.Println()
	}

	var topScore int
	for row := range grid {
		for col := range grid[row] {
			if grid[row][col].score > topScore {
				topScore = grid[row][col].score
			}
		}
	}

	return topScore
}

func calculateScore(grid Grid, row int, col int) int {
	var (
		scoreLeft  int
		scoreRight int
		scoreUp    int
		scoreDown  int
	)

	start := grid[row][col]

	// Look left
	if col > 0 {
		for i := col - 1; i >= 0; i-- {
			scoreLeft++
			next := grid[row][i]
			if !start.higherThan(next) {
				break
			}
		}
	}

	// Look right
	if col < len(grid[row])-1 {
		for i := col + 1; i < len(grid[row]); i++ {
			scoreRight++
			next := grid[row][i]
			if !start.higherThan(next) {
				break
			}
		}
	}

	// Look up
	if row > 0 {
		for i := row - 1; i >= 0; i-- {
			scoreUp++
			next := grid[i][col]
			if !start.higherThan(next) {
				break
			}
		}
	}

	// Look down
	if row < len(grid)-1 {
		for i := row + 1; i < len(grid); i++ {
			scoreDown++
			next := grid[i][col]
			if !start.higherThan(next) {
				break
			}
		}
	}

	score := scoreLeft * scoreRight * scoreUp * scoreDown

	return score
}
