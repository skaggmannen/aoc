package main

import (
	"aoc2022/util"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	path := "input.txt"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	input, _ := os.ReadFile(path)

	fmt.Println("Part 1:", part1(bytes.NewReader(input)))
	fmt.Println("Part 2:", part2(bytes.NewReader(input), path == "test.txt"))
}

func part1(r io.Reader) int {
	lines, directions := parseInput(r)

	board := newBoard(lines)

	for _, d := range directions {
		switch d {
		case "R":
			board.turnRight()
		case "L":
			board.turnLeft()
		default:
			d, _ := strconv.Atoi(d)
			for i := 0; i < d; i++ {
				board.move()
			}
		}
	}

	fmt.Println(board)
	fmt.Println()

	return board.score()
}

func part2(r io.Reader, test bool) int {
	lines, directions := parseInput(r)

	c := newCube(lines, test)

	c.visit(c.x, c.y)

	for _, d := range directions {
		switch d {
		case "R":
			c.turnRight()
			c.visit(c.x, c.y)
		case "L":
			c.turnLeft()
			c.visit(c.x, c.y)
		default:
			d, _ := strconv.Atoi(d)
			for i := 0; i < d; i++ {
				c.move()
				c.visit(c.x, c.y)
			}
		}
	}

	fmt.Println(c)
	fmt.Println()

	return c.score()

}

func parseInput(r io.Reader) ([]string, []string) {
	data, _ := io.ReadAll(r)

	parts := strings.Split(string(data), "\n\n")
	if len(parts) < 2 {
		panic("invalid input")
	}

	var (
		pattern, _ = regexp.Compile(`(\d+|[LR])`)
	)

	matches := pattern.FindAllStringSubmatch(parts[1], -1)
	if len(matches) == 0 {
		panic("invalid directions: " + parts[1])
	}

	directions := make([]string, len(matches))
	for i, m := range matches {
		directions[i] = m[0]
	}

	return strings.Split(parts[0], "\n"), directions
}

func newBoard(tiles []string) *board {
	b := &board{
		tiles:   make(map[util.Point]string),
		facing:  East,
		visited: make(map[util.Point]string),
		h:       len(tiles),
	}

	for y := 0; y < len(tiles); y++ {
		for x := 0; x < len(tiles[y]); x++ {
			if tiles[y][x] != ' ' {
				b.tiles[util.Point{x, y}] = string(tiles[y][x])
			}
			if len(tiles[y]) > b.w {
				b.w = len(tiles[y])
			}
		}
	}

	// Find the start position
	for b.get(b.x, b.y) == " " {
		b.x++
	}

	b.cubeWidth = len(strings.TrimSpace(tiles[0]))

	b.visit(b.x, b.y)

	return b
}

type board struct {
	cubeWidth int
	tiles     map[util.Point]string
	facing    string
	x, y      int
	w, h      int
	visited   map[util.Point]string
}

func (b *board) get(x, y int) string {
	v, ok := b.tiles[util.Point{x, y}]
	if ok {
		return v
	} else {
		return " "
	}
}

func (b *board) turnRight() {
	switch b.facing {
	case East:
		b.facing = South
	case North:
		b.facing = East
	case West:
		b.facing = North
	case South:
		b.facing = West
	}

	b.visit(b.x, b.y)
}

func (b *board) turnLeft() {
	switch b.facing {
	case East:
		b.facing = North
	case North:
		b.facing = West
	case West:
		b.facing = South
	case South:
		b.facing = East
	}

	b.visit(b.x, b.y)
}

func (b *board) move() {
	startX := b.x
	startY := b.y
	defer func() {
		b.visit(b.x, b.y)
	}()

	for {
		switch b.facing {
		case North:
			b.y -= 1
			if b.y < 0 {
				b.y = b.h - 1
			}
		case East:
			b.x += 1
			if b.x >= b.w {
				b.x = 0
			}
		case South:
			b.y += 1
			if b.y >= b.h {
				b.y = 0
			}
		case West:
			b.x -= 1
			if b.x < 0 {
				b.x = b.w - 1
			}
		}

		if b.get(b.x, b.y) == "#" {
			b.x = startX
			b.y = startY
		}

		if b.get(b.x, b.y) != " " {
			return
		}
	}
}

func (b *board) score() int {
	x := b.x % b.w
	if x < 0 {
		x = b.w + b.x
	}

	y := b.y % b.h
	if y < 0 {
		y = b.h + b.y
	}

	s := 1000 * (y + 1)
	s += 4 * (x + 1)

	switch b.facing {
	case South:
		s += 1
	case West:
		s += 2
	case North:
		s += 3
	}

	return s
}

func (b *board) visit(x, y int) {
	x %= b.w
	if x < 0 {
		x = b.w + x
	}

	y %= b.h
	if y < 0 {
		y = b.h + y
	}

	b.visited[util.Point{x, y}] = b.facing
}

func (b *board) String() string {
	var buf strings.Builder
	for y := 0; y < b.h; y++ {
		for x := 0; x < b.w; x++ {
			if f, ok := b.visited[util.Point{x, y}]; ok {
				fmt.Fprint(&buf, f)
			} else {
				fmt.Fprint(&buf, b.get(x, y))
			}
		}
		fmt.Fprintln(&buf)
	}
	return buf.String()
}

const (
	North = "^"
	East  = ">"
	South = "v"
	West  = "<"
)

func newCube(lines []string, test bool) *cube {
	c := &cube{
		side:    len(strings.TrimSpace(lines[0])),
		h:       len(lines),
		tiles:   make(map[string]string),
		visited: make(map[string]string),
		facing:  East,
		test:    test,
	}

	for row := range lines {
		for col := range lines[row] {
			if col > c.w {
				c.w = col
			}
			if v := lines[row][col]; v != ' ' {
				c.tiles[fmt.Sprintf("%d:%d", row, col)] = string(v)
			}
		}
	}
	if test {
		readTest(c, lines)
	} else {
		readInput(c, lines)
	}

	c.face = FaceNorth

	return c
}

type cube struct {
	side    int
	faces   [6]face
	x, y    int
	w, h    int
	face    int
	facing  string
	tiles   map[string]string
	visited map[string]string
	test    bool
}

func (b *cube) turnRight() {
	switch b.facing {
	case East:
		b.facing = South
	case North:
		b.facing = East
	case West:
		b.facing = North
	case South:
		b.facing = West
	}
}

func (b *cube) turnLeft() {
	switch b.facing {
	case East:
		b.facing = North
	case North:
		b.facing = West
	case West:
		b.facing = South
	case South:
		b.facing = East
	}
}

func (c *cube) move() {
	startX := c.x
	startY := c.y
	startFace := c.face
	startFacing := c.facing

	for {
		switch c.facing {
		case North:
			c.y -= 1
		case East:
			c.x += 1
		case South:
			c.y += 1
		case West:
			c.x -= 1
		}

		if c.test {
			foldTest(c)
		} else {
			foldInput(c)
		}

		if c.get(c.x, c.y) == "#" {
			c.x = startX
			c.y = startY
			c.face = startFace
			c.facing = startFacing
		}

		if c.get(c.x, c.y) != " " {
			return
		}
	}
}

func (c *cube) get(x, y int) string {
	return c.faces[c.face].tiles[y][x].v
}

func (c *cube) visit(x, y int) {
	t := c.faces[c.face].tiles[y][x]

	c.visited[fmt.Sprintf("%d:%d", t.row, t.col)] = c.facing
}

func (c *cube) score() int {
	t := c.faces[c.face].tiles[c.y][c.x]

	s := 1000 * (t.row + 1)
	s += 4 * (t.col + 1)

	switch c.facing {
	case South:
		s += 1
	case West:
		s += 2
	case North:
		s += 3
	}

	return s
}

func (c *cube) String() string {
	var buf strings.Builder
	for row := 0; row < c.h; row++ {
		for col := 0; col < c.w; col++ {
			if v, ok := c.visited[fmt.Sprintf("%d:%d", row, col)]; ok {
				fmt.Fprint(&buf, v)
			} else if v, ok := c.tiles[fmt.Sprintf("%d:%d", row, col)]; ok {
				fmt.Fprint(&buf, v)
			} else {
				fmt.Fprint(&buf, " ")
			}
		}
		fmt.Fprintln(&buf)
	}
	return buf.String()
}

type face struct {
	tiles [][]tile
}

type tile struct {
	v        string
	row, col int
}

const (
	FaceNorth = iota
	FaceBack
	FaceWest
	FaceFront
	FaceSouth
	FaceEast
)

func readTest(c *cube, lines []string) {
	// Read north face
	c.faces[FaceNorth].tiles = make([][]tile, c.side)
	for y := 0; y < c.side; y++ {
		c.faces[FaceNorth].tiles[y] = make([]tile, c.side)
		for x := 0; x < c.side; x++ {
			row := y
			col := x + 2*c.side
			c.faces[FaceNorth].tiles[y][x] = tile{
				v:   string(lines[row][col]),
				row: row,
				col: col,
			}
		}
	}

	// Read back face
	c.faces[FaceBack].tiles = make([][]tile, c.side)
	for y := 0; y < c.side; y++ {
		c.faces[FaceBack].tiles[y] = make([]tile, c.side)
		for x := 0; x < c.side; x++ {
			row := y + c.side
			col := x
			c.faces[FaceBack].tiles[y][x] = tile{
				v:   string(lines[row][col]),
				row: row,
				col: col,
			}
		}
	}

	// Read west face
	c.faces[FaceWest].tiles = make([][]tile, c.side)
	for y := 0; y < c.side; y++ {
		c.faces[FaceWest].tiles[y] = make([]tile, c.side)
		for x := 0; x < c.side; x++ {
			row := y + c.side
			col := x + c.side
			c.faces[FaceWest].tiles[y][x] = tile{
				v:   string(lines[row][col]),
				row: row,
				col: col,
			}
		}
	}

	// Read front face
	c.faces[FaceFront].tiles = make([][]tile, c.side)
	for y := 0; y < c.side; y++ {
		c.faces[FaceFront].tiles[y] = make([]tile, c.side)
		for x := 0; x < c.side; x++ {
			row := y + c.side
			col := x + 2*c.side
			c.faces[FaceFront].tiles[y][x] = tile{
				v:   string(lines[row][col]),
				row: row,
				col: col,
			}
		}
	}

	// Read south face
	c.faces[FaceSouth].tiles = make([][]tile, c.side)
	for y := 0; y < c.side; y++ {
		c.faces[FaceSouth].tiles[y] = make([]tile, c.side)
		for x := 0; x < c.side; x++ {
			row := y + 2*c.side
			col := x + 2*c.side
			c.faces[FaceSouth].tiles[y][x] = tile{
				v:   string(lines[row][col]),
				row: row,
				col: col,
			}
		}
	}

	// Read east face
	c.faces[FaceEast].tiles = make([][]tile, c.side)
	for y := 0; y < c.side; y++ {
		c.faces[FaceEast].tiles[y] = make([]tile, c.side)
		for x := 0; x < c.side; x++ {
			row := y + 2*c.side
			col := x + 3*c.side
			c.faces[FaceEast].tiles[y][x] = tile{
				v:   string(lines[row][col]),
				row: row,
				col: col,
			}
		}
	}
}

func foldTest(c *cube) {

	switch c.face {
	case FaceNorth:
		if c.y < 0 {
			// Move to back
			c.face = FaceBack
			c.facing = South
			c.y = 0
		} else if c.y >= c.side {
			// Move to front
			c.face = FaceFront
			c.y = 0
		} else if c.x < 0 {
			// Move to west
			c.face = FaceWest
			c.facing = South
			c.x = c.y
			c.y = 0
		} else if c.x >= c.side {
			// Move to east
			c.face = FaceEast
			c.facing = West
			c.x = c.side - 1
			c.y = c.side - c.y
		}
	case FaceWest:
		if c.y < 0 {
			// Move to north
			c.face = FaceNorth
			c.facing = East
			c.y = c.x
			c.x = 0
		} else if c.y >= c.side {
			// Move to south
			c.face = FaceSouth
			c.facing = East
			c.y = c.side - c.x
			c.x = 0
		} else if c.x < 0 {
			// Move to back
			c.face = FaceBack
			c.x = c.side - 1
		} else if c.x >= c.side {
			// Move to front
			c.face = FaceFront
			c.x = 0
		}
	case FaceFront:
		if c.y < 0 {
			// Move to north
			c.face = FaceNorth
			c.y = c.side - 1
		} else if c.y >= c.side {
			// Move to south
			c.face = FaceSouth
			c.y = 0
		} else if c.x < 0 {
			// Move to west
			c.face = FaceWest
			c.x = c.side - 1
		} else if c.x >= c.side {
			// Move to east
			c.face = FaceEast
			c.facing = South
			c.x = c.side - c.y - 1
			c.y = 0
		}
	case FaceEast:
		if c.y < 0 {
			// Move to front
			c.face = FaceFront
			c.facing = West
			c.y = c.side - c.x - 1
			c.x = c.side - 1
		} else if c.y >= c.side {
			// Move to back
			c.face = FaceBack
			c.facing = East
			c.y = c.side - c.x
			c.x = 0
		} else if c.x < 0 {
			// Move to south
			c.face = FaceSouth
			c.x = c.side - 1
		} else if c.x >= c.side {
			// Move to north
			c.face = FaceNorth
			c.facing = West
			c.x = c.side - 1
			c.y = c.side - c.y - 1
		}
	case FaceSouth:
		if c.y < 0 {
			// Move to front
			c.face = FaceFront
			c.y = c.side
		} else if c.y >= c.side {
			// Move to back
			c.face = FaceBack
			c.facing = North
			c.y = c.side - 1
			c.x = c.side - c.x - 1
		} else if c.x < 0 {
			// Move to west
			c.face = FaceWest
			c.facing = North
			c.x = c.side - c.y - 1
			c.y = 0
		} else if c.x >= c.side {
			// Move to east
			c.face = FaceEast
			c.x = 0
		}
	case FaceBack:
		if c.y < 0 {
			// Move to north
			c.face = FaceNorth
			c.facing = South
			c.y = 0
			c.x = c.side - c.x - 1
		} else if c.y >= c.side {
			// Move to south
			c.face = FaceSouth
			c.facing = North
			c.y = c.side
			c.x = c.side - c.x - 1
		} else if c.x < 0 {
			// Move to east
			c.face = FaceEast
			c.facing = North
			c.x = c.side - c.y - 1
			c.y = c.side - 1
			// Switch to north
		} else if c.x >= c.side {
			// Move to west
			c.face = FaceWest
			c.x = 0
		}
	}
}
func readInput(c *cube, lines []string) {
	c.side = 50
	// Read north face
	c.faces[FaceNorth].tiles = make([][]tile, c.side)
	for y := 0; y < c.side; y++ {
		c.faces[FaceNorth].tiles[y] = make([]tile, c.side)
		for x := 0; x < c.side; x++ {
			row := y
			col := x + c.side
			c.faces[FaceNorth].tiles[y][x] = tile{
				v:   string(lines[row][col]),
				row: row,
				col: col,
			}
		}
	}

	// Read back face
	c.faces[FaceBack].tiles = make([][]tile, c.side)
	for y := 0; y < c.side; y++ {
		c.faces[FaceBack].tiles[y] = make([]tile, c.side)
		for x := 0; x < c.side; x++ {
			row := y + 3*c.side
			col := x
			c.faces[FaceBack].tiles[y][x] = tile{
				v:   string(lines[row][col]),
				row: row,
				col: col,
			}
		}
	}

	// Read west face
	c.faces[FaceWest].tiles = make([][]tile, c.side)
	for y := 0; y < c.side; y++ {
		c.faces[FaceWest].tiles[y] = make([]tile, c.side)
		for x := 0; x < c.side; x++ {
			row := y + 2*c.side
			col := x
			c.faces[FaceWest].tiles[y][x] = tile{
				v:   string(lines[row][col]),
				row: row,
				col: col,
			}
		}
	}

	// Read front face
	c.faces[FaceFront].tiles = make([][]tile, c.side)
	for y := 0; y < c.side; y++ {
		c.faces[FaceFront].tiles[y] = make([]tile, c.side)
		for x := 0; x < c.side; x++ {
			row := y + c.side
			col := x + c.side
			c.faces[FaceFront].tiles[y][x] = tile{
				v:   string(lines[row][col]),
				row: row,
				col: col,
			}
		}
	}

	// Read south face
	c.faces[FaceSouth].tiles = make([][]tile, c.side)
	for y := 0; y < c.side; y++ {
		c.faces[FaceSouth].tiles[y] = make([]tile, c.side)
		for x := 0; x < c.side; x++ {
			row := y + 2*c.side
			col := x + c.side
			c.faces[FaceSouth].tiles[y][x] = tile{
				v:   string(lines[row][col]),
				row: row,
				col: col,
			}
		}
	}

	// Read east face
	c.faces[FaceEast].tiles = make([][]tile, c.side)
	for y := 0; y < c.side; y++ {
		c.faces[FaceEast].tiles[y] = make([]tile, c.side)
		for x := 0; x < c.side; x++ {
			row := y
			col := x + 2*c.side
			c.faces[FaceEast].tiles[y][x] = tile{
				v:   string(lines[row][col]),
				row: row,
				col: col,
			}
		}
	}
}

func foldInput(c *cube) {
	switch c.face {
	case FaceNorth:
		if c.y < 0 {
			// Move to back - OK
			c.face = FaceBack
			c.facing = East
			c.y = c.x
			c.x = 0
		} else if c.y >= c.side {
			// Move to front - OK
			c.face = FaceFront
			c.y = 0
		} else if c.x < 0 {
			// Move to west - OK
			c.face = FaceWest
			c.facing = East
			c.x = 0
			c.y = c.side - c.y - 1
		} else if c.x >= c.side {
			// Move to east - OK
			c.face = FaceEast
			c.x = 0
		}
	case FaceWest:
		if c.y < 0 {
			// Move to front - OK
			c.face = FaceFront
			c.facing = East
			c.y = c.x
			c.x = 0
		} else if c.y >= c.side {
			// Move to back - OK
			c.face = FaceBack
			c.y = 0
		} else if c.x < 0 {
			// Move to north - OK
			c.face = FaceNorth
			c.facing = East
			c.x = 0
			c.y = c.side - c.y - 1
		} else if c.x >= c.side {
			// Move to south - OK
			c.face = FaceSouth
			c.x = 0
		}
	case FaceFront:
		if c.y < 0 {
			// Move to north - OK
			c.face = FaceNorth
			c.y = c.side - 1
		} else if c.y >= c.side {
			// Move to south - OK
			c.face = FaceSouth
			c.y = 0
		} else if c.x < 0 {
			// Move to west - OK
			c.face = FaceWest
			c.facing = South
			c.x = c.y
			c.y = 0
		} else if c.x >= c.side {
			// Move to east - OK
			c.face = FaceEast
			c.facing = North
			c.x = c.y
			c.y = c.side - 1
		}
	case FaceEast:
		if c.y < 0 {
			// Move to back
			c.face = FaceBack
			c.y = c.side - 1
		} else if c.y >= c.side {
			// Move to front
			c.face = FaceFront
			c.facing = West
			c.y = c.x
			c.x = c.side - 1
		} else if c.x < 0 {
			// Move to north
			c.face = FaceNorth
			c.x = c.side - 1
		} else if c.x >= c.side {
			// Move to south
			c.face = FaceSouth
			c.facing = West
			c.x = c.side - 1
			c.y = c.side - c.y - 1
		}
	case FaceSouth:
		if c.y < 0 {
			// Move to front
			c.face = FaceFront
			c.y = c.side - 1
		} else if c.y >= c.side {
			// Move to back
			c.face = FaceBack
			c.facing = West
			c.y = c.x
			c.x = c.side - 1
		} else if c.x < 0 {
			// Move to west
			c.face = FaceWest
			c.x = c.side - 1
		} else if c.x >= c.side {
			// Move to east
			c.face = FaceEast
			c.facing = West
			c.x = c.side - 1
			c.y = c.side - c.y - 1
		}
	case FaceBack:
		if c.y < 0 {
			// Move to west
			c.face = FaceWest
			c.y = c.side - 1
		} else if c.y >= c.side {
			// Move to east
			c.face = FaceEast
			c.y = 0
		} else if c.x < 0 {
			// Move to north
			c.face = FaceNorth
			c.facing = South
			c.x = c.y
			c.y = 0
		} else if c.x >= c.side {
			// Move to south
			c.face = FaceSouth
			c.facing = North
			c.x = c.y
			c.y = c.side - 1
		}
	}
}
