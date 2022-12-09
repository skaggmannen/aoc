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
		head Pos
		tail Pos
	)
	visitedByTail := make(map[string]struct{})

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			continue
		}

		delta, _ := strconv.Atoi(parts[1])

		for i := 0; i < delta; i++ {
			switch parts[0] {
			case "U":
				head.y++
			case "D":
				head.y--
			case "R":
				head.x++
			case "L":
				head.x--
			}

			if head.x > tail.x+1 {
				tail.x++
				if head.y != tail.y {
					tail.y = head.y
				}
			} else if head.x < tail.x-1 {
				tail.x--
				if head.y != tail.y {
					tail.y = head.y
				}
			} else if head.y > tail.y+1 {
				tail.y++
				if head.x != tail.x {
					tail.x = head.x
				}
			} else if head.y < tail.y-1 {
				tail.y--
				if head.x != tail.x {
					tail.x = head.x
				}
			}

			visitedByTail[fmt.Sprintf("%d,%d", tail.x, tail.y)] = struct{}{}
		}
	}

	return len(visitedByTail)
}

func part2(r io.Reader) int {
	visitedByTail := make(map[Pos]struct{})
	rope := make([]Pos, 10)
	// b := NewBoard(1000)
	// b.Draw(rope)

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			continue
		}

		delta, _ := strconv.Atoi(parts[1])

		for i := 0; i < delta; i++ {
			head := &rope[0]
			switch parts[0] {
			case "U":
				head.y++
			case "D":
				head.y--
			case "R":
				head.x++
			case "L":
				head.x--
			}

			for i := 1; i < len(rope); i++ {
				tail := &rope[i]
				head := &rope[i-1]
				dx := head.x - tail.x
				dy := head.y - tail.y

				if dx == 2 {
					tail.x++
					if dy == 2 {
						tail.y++
					} else if dy == -2 {
						tail.y--
					} else {
						tail.y = head.y
					}
				} else if dx == -2 {
					tail.x--
					if dy == 2 {
						tail.y++
					} else if dy == -2 {
						tail.y--
					} else {
						tail.y = head.y
					}
				} else if dy == 2 {
					tail.y++
					if dx == 2 {
						tail.x++
					} else if dx == -2 {
						tail.x--
					} else {
						tail.x = head.x
					}
				} else if dy == -2 {
					tail.y--
					if dx == 2 {
						tail.x++
					} else if dx == -2 {
						tail.x--
					} else {
						tail.x = head.x
					}
				} else {
					break
				}
			}

			visitedByTail[rope[len(rope)-1]] = struct{}{}

			// b.Draw(rope)
			// fmt.Println(b.String())
		}

	}

	return len(visitedByTail)
}

type Pos struct {
	x int
	y int
}

func NewBoard(size int) Board {
	tiles := make([][]string, size)
	for i := range tiles {
		tiles[i] = make([]string, size)
		for j := range tiles[i] {
			tiles[i][j] = "."
		}
	}
	return Board{tiles: tiles, size: size}
}

type Board struct {
	size  int
	tiles [][]string
}

func (b Board) Draw(rope []Pos) {
	b.Clear()
	for i := len(rope) - 1; i >= 0; i-- {
		b.tiles[rope[i].y+b.size/2][rope[i].x+b.size/2] = fmt.Sprint(i)
	}
}

func (b Board) Clear() {
	for row := range b.tiles {
		for col := range b.tiles[row] {
			b.tiles[row][col] = "."
		}
	}
}

func (b Board) String() string {
	var buf strings.Builder
	for _, row := range b.tiles {
		for _, col := range row {
			fmt.Fprint(&buf, col)
		}
		fmt.Fprintln(&buf)
	}
	return buf.String()
}
