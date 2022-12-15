package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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
	board := newBoard(r, false)

	start := point{500, 0}
	board.set(start, "+")

	fmt.Println(board.String())
	fmt.Println()

	for i := 0; ; i++ {
		sand := board.pour(start)
		if sand.y > board.bottom {
			fmt.Println(board.String())
			fmt.Println()
			return i
		}
	}
}

func part2(r io.Reader) int {
	board := newBoard(r, true)

	start := point{500, 0}
	board.set(start, "+")

	fmt.Println(board.String())
	fmt.Println()

	for i := 0; ; i++ {
		sand := board.pour(start)
		if sand.y == start.y {
			fmt.Println(board.String())
			fmt.Println()
			return i + 1
		}
	}
}

func newBoard(r io.Reader, infiniteFloor bool) board {
	rocks := parseRocks(r)
	b := board{
		elems:         make(map[point]string),
		infiniteFloor: infiniteFloor,
	}

	for _, r := range rocks {
		r.draw(b)
	}

	min, max := b.findBounds()
	b.bottom = max.y

	if infiniteFloor {
		b.drawLine(point{min.x, max.y + 2}, point{max.x, max.y + 2}, "#")
	}

	return b
}

func parseRocks(r io.Reader) []rock {
	data, _ := io.ReadAll(r)
	lines := strings.Split(string(data), "\n")

	rocks := make([]rock, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		rocks = append(rocks, parseRock(l))
	}
	return rocks
}

func parseRock(l string) rock {
	parts := strings.Split(l, " -> ")

	path := make([]point, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if len(p) == 0 {
			continue
		}

		path = append(path, parsePoint(p))
	}

	return rock{path: path}
}

func parsePoint(p string) point {
	parts := strings.Split(p, ",")
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	return point{x, y}
}

type board struct {
	elems         map[point]string
	bottom        int
	infiniteFloor bool
}

func (b *board) set(p point, v string) {
	b.elems[p] = v
}

func (b *board) get(p point) string {
	v, ok := b.elems[p]
	if !ok {
		return "."
	}
	return v
}

func (b *board) drawLine(from point, to point, v string) {
	if from.x > to.x {
		from.x, to.x = to.x, from.x
	}
	if from.y > to.y {
		from.y, to.y = to.y, from.y
	}

	for x := from.x; x <= to.x; x++ {
		for y := from.y; y <= to.y; y++ {
			b.set(point{x, y}, v)
		}
	}
}

func (b *board) pour(sand point) point {
	if sand.y > b.bottom {
		// The sand fell off the bottom of the board...
		if b.infiniteFloor {
			// ...or to the floor
			b.set(sand, "o")
			b.set(point{sand.x, sand.y + 1}, "#")
			b.set(point{sand.x + 1, sand.y + 1}, "#")
			b.set(point{sand.x - 1, sand.y + 1}, "#")
		}
		return sand
	}

	// Check if sand falls straight down
	down := point{sand.x, sand.y + 1}
	if b.get(down) == "." {
		return b.pour(down)
	}

	// Check if sand rolls to the left
	left := point{sand.x - 1, sand.y + 1}
	if b.get(left) == "." {
		return b.pour(left)
	}

	// Check if sand rolls to the right
	right := point{sand.x + 1, sand.y + 1}
	if b.get(right) == "." {
		return b.pour(right)
	}

	// The sand could go nowhere, so it lands here
	b.set(sand, "o")
	return sand
}

func (b *board) String() string {
	min, max := b.findBounds()

	var buf strings.Builder
	for y := min.y; y <= max.y+1; y++ {
		for x := min.x - 1; x <= max.x+1; x++ {
			fmt.Fprint(&buf, b.get(point{x, y}))
		}
		fmt.Fprintln(&buf)
	}
	return buf.String()
}

func (b *board) findBounds() (point, point) {
	min := point{0x7fff_ffff, 0x7fff_ffff}
	max := point{-0x7fff_ffff, -0x7fff_ffff}

	for e := range b.elems {
		if e.x < min.x {
			min.x = e.x
		}
		if e.x > max.x {
			max.x = e.x
		}
		if e.y < min.y {
			min.y = e.y
		}
		if e.y > max.y {
			max.y = e.y
		}
	}

	return min, max
}

type rock struct {
	path []point
}

func (r rock) draw(b board) {
	start := r.path[0]
	for i := 1; i < len(r.path); i++ {
		stop := r.path[i]
		b.drawLine(start, stop, "#")
		start = stop
	}
}

type point struct {
	x, y int
}
