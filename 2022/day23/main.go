package main

import (
	"aoc2022/util"
	"bytes"
	"fmt"
	"io"
	"os"
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
	b := readBoard(r)

	directions := []int{
		0,
		1,
		2,
		3,
	}

	fmt.Println(b)
	fmt.Println()

	for i := 0; i < 10; i++ {
		proposedMoves := make(map[util.Point][]util.Point)
		for pos, e := range b {
			alts := checkAdjancent(e, b)

			if alts[0] != nil && alts[1] != nil && alts[2] != nil && alts[3] != nil {
				// no adjacent elves, stay where we are
				continue
			}

			for _, d := range directions {
				if newPos := alts[d]; newPos != nil {
					proposedMoves[*newPos] = append(proposedMoves[*newPos], pos)
					break
				}
			}
		}

		for newPos, oldPos := range proposedMoves {
			if len(oldPos) > 1 {
				continue
			}

			delete(b, oldPos[0])
			b[newPos] = &elf{pos: newPos}
		}

		shiftLeft(directions)

		fmt.Println(b)
		fmt.Println()
	}

	var score int
	min, max := b.bounds()

	for y := min.Y; y <= max.Y; y++ {
		for x := min.X; x <= max.X; x++ {
			if _, ok := b[util.Point{X: x, Y: y}]; !ok {
				score++
			}
		}
	}

	return score
}

func part2(r io.Reader) int {
	b := readBoard(r)

	directions := []int{
		0,
		1,
		2,
		3,
	}

	fmt.Println(b)
	fmt.Println()

	for i := 1; ; i++ {
		proposedMoves := make(map[util.Point][]util.Point)
		for pos, e := range b {
			alts := checkAdjancent(e, b)

			if alts[0] != nil && alts[1] != nil && alts[2] != nil && alts[3] != nil {
				// no adjacent elves, stay where we are
				continue
			}

			for _, d := range directions {
				if newPos := alts[d]; newPos != nil {
					proposedMoves[*newPos] = append(proposedMoves[*newPos], pos)
					break
				}
			}
		}

		var elfMoved bool
		for newPos, oldPos := range proposedMoves {
			if len(oldPos) > 1 {
				continue
			}

			delete(b, oldPos[0])
			b[newPos] = &elf{pos: newPos}
			elfMoved = true
		}
		if !elfMoved {
			return i
		}

		shiftLeft(directions)

		fmt.Println(b)
		fmt.Println()
	}
}

func readBoard(r io.Reader) board {
	data, _ := io.ReadAll(r)
	lines := strings.Split(string(data), "\n")

	result := make(board)
	for row, l := range lines {
		for col, b := range l {
			if b == '#' {
				pos := util.Point{X: col, Y: row}
				result[pos] = &elf{
					pos: pos,
				}
			}
		}
	}

	return result
}

type elf struct {
	pos util.Point
}

type board map[util.Point]*elf

func (b board) String() string {
	var buf strings.Builder
	fmt.Fprint(&buf, "  ")
	for i := -3; i < 10; i++ {
		if i < 0 {
			fmt.Fprint(&buf, "-")
		} else {
			fmt.Fprint(&buf, " ")
		}
	}
	fmt.Fprintln(&buf)
	fmt.Fprint(&buf, "  ")
	for i := -3; i < 10; i++ {
		if i < 0 {
			fmt.Fprint(&buf, -i)
		} else {
			fmt.Fprint(&buf, i)
		}
	}
	fmt.Fprintln(&buf)
	for y := -3; y < 10; y++ {
		fmt.Fprintf(&buf, "%2d", y)
		for x := -3; x < 10; x++ {
			if _, ok := b[util.Point{X: x, Y: y}]; ok {
				fmt.Fprint(&buf, "#")
			} else {
				fmt.Fprint(&buf, ".")
			}
		}
		fmt.Fprintln(&buf)
	}
	return buf.String()
}

func (b board) bounds() (util.Point, util.Point) {
	var (
		min = util.Point{X: 0x7fff_ffff, Y: 0x7fff_ffff}
		max = util.Point{X: -0x7fff_ffff, Y: -0x7fff_ffff}
	)

	for pos := range b {
		if pos.X < min.X {
			min.X = pos.X
		}
		if pos.X > max.X {
			max.X = pos.X
		}
		if pos.Y < min.Y {
			min.Y = pos.Y
		}
		if pos.Y > max.Y {
			max.Y = pos.Y
		}
	}

	return min, max
}

func checkAdjancent(e *elf, b board) (result [4]*util.Point) {
	var (
		ne   = util.Point{X: e.pos.X + 1, Y: e.pos.Y - 1}
		n    = util.Point{X: e.pos.X, Y: e.pos.Y - 1}
		nw   = util.Point{X: e.pos.X - 1, Y: e.pos.Y - 1}
		w    = util.Point{X: e.pos.X - 1, Y: e.pos.Y}
		sw   = util.Point{X: e.pos.X - 1, Y: e.pos.Y + 1}
		s    = util.Point{X: e.pos.X, Y: e.pos.Y + 1}
		se   = util.Point{X: e.pos.X + 1, Y: e.pos.Y + 1}
		east = util.Point{X: e.pos.X + 1, Y: e.pos.Y}

		_, neOk = b[ne]
		_, nOk  = b[n]
		_, nwOk = b[nw]
		_, wOk  = b[w]
		_, swOk = b[sw]
		_, sOk  = b[s]
		_, seOk = b[se]
		_, eOk  = b[east]
	)

	if !nOk && !neOk && !nwOk {
		result[0] = &n
	}
	if !sOk && !seOk && !swOk {
		result[1] = &s
	}
	if !wOk && !swOk && !nwOk {
		result[2] = &w
	}
	if !eOk && !seOk && !neOk {
		result[3] = &east
	}

	return
}

func shiftLeft[T any](s []T) {
	first := s[0]
	for i := 1; i < len(s); i++ {
		s[i-1] = s[i]
	}
	s[len(s)-1] = first
}
