package main

import (
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
	gusts := readGusts(r)

	b := newBoard(gusts)

	for i := 0; i < 2022; i++ {
		b.drop()
	}

	return b.height()
}

func part2(r io.Reader) int {
	gusts := readGusts(r)

	b := newBoard(gusts)

	for i := 0; i < 1000000000000; i++ {
		if state, ok := b.detectCycle(i); ok {
			currHeight := b.height()
			cycleLength := i - state.i
			heightIncrease := currHeight - state.height
			cycleCount := (1000000000000 - state.i) / cycleLength
			remainder := 1000000000000 - cycleCount*cycleLength - state.i

			for i := 0; i < remainder; i++ {
				b.drop()
			}

			finalHeight := state.height
			finalHeight += heightIncrease * cycleCount
			finalHeight += b.height() - currHeight

			return finalHeight
		}
		b.drop()
		if i%10_000_000 == 0 {
			fmt.Printf("%d: %d\n", i, b.height())
		}
	}

	return b.height()
}

func readGusts(r io.Reader) string {
	data, _ := io.ReadAll(r)
	return string(data)
}

func newBoard(gusts string) *board {
	s := &board{
		gusts: gusts,
	}

	return s
}

type board struct {
	rocks       []uint8
	nextPattern int
	nextGust    int
	gusts       string
}

var cache map[string]cachedState

type cachedState struct {
	i      int
	height int
}

func (b *board) detectCycle(i int) (cachedState, bool) {
	if i < 10_000 {
		return cachedState{}, false
	}

	var relHeights []int
	for i := 0; i < 7; i++ {
		for j := len(b.rocks) - 1; j >= 0; j-- {
			if b.rocks[j]&(0b1000000>>i) != 0 {
				relHeights = append(relHeights, j)
				break
			}
		}
	}

	var base int
	for _, v := range relHeights {
		if v > base {
			base = v
		}
	}

	for i := range relHeights {
		relHeights[i] -= base
	}

	key := fmt.Sprintf("%d:%d", b.nextGust, b.nextPattern)
	for _, h := range relHeights {
		key += fmt.Sprintf(":%d", h)
	}

	if state, ok := cache[key]; ok {
		return state, ok
	}

	if cache == nil {
		cache = make(map[string]cachedState)
	}

	cache[key] = cachedState{
		i:      i,
		height: b.height(),
	}

	return cachedState{}, false
}

func (b *board) drop() {
	p := rocks[b.nextPattern]
	b.nextPattern = (b.nextPattern + 1) % len(rocks)

	y := b.height() + 3

	for y >= 0 && !b.overlaps(y, p) {
		g := b.gusts[b.nextGust]
		b.nextGust = (b.nextGust + 1) % len(b.gusts)

		switch g {
		case '<':
			if p.canMoveLeft() {
				p.moveLeft()
			}
			if b.overlaps(y, p) {
				p.moveRight()
			}
		case '>':
			if p.canMoveRight() {
				p.moveRight()
			}
			if b.overlaps(y, p) {
				p.moveLeft()
			}
		}

		y -= 1
	}

	b.addRocks(y+1, p)
}

func (b *board) overlaps(y int, p pattern) bool {
	if y >= len(b.rocks) {
		return false
	}

	for i := 0; i < len(p); i++ {
		if y+i >= len(b.rocks) {
			return false
		}
		if b.rocks[y+i]&p[len(p)-i-1] != 0 {
			return true
		}
	}

	return false
}

func (b *board) addRocks(y int, p pattern) {
	for i := 0; i < len(p); i++ {
		if y+i >= len(b.rocks) {
			b.rocks = append(b.rocks, p[len(p)-i-1])
		} else {
			b.rocks[y+i] |= p[len(p)-i-1]
		}
	}
}

func (b *board) height() int {
	for i := len(b.rocks) - 1; i >= 0; i-- {
		if b.rocks[i] != 0 {
			return i + 1
		}
	}

	return 0
}

func (b *board) String() string {
	var buf strings.Builder

	for i := len(b.rocks) - 1; i >= 0; i-- {
		for j := 0; j < 7; j++ {
			if b.rocks[i]&(0b1000000>>j) != 0 {
				fmt.Fprint(&buf, "#")
			} else {
				fmt.Fprint(&buf, ".")
			}
		}
		fmt.Fprintln(&buf)
	}

	return buf.String()
}

type pattern [4]uint8

func (p *pattern) canMoveLeft() bool {
	for i := range p {
		if p[i]&0b1000000 != 0 {
			return false
		}
	}
	return true
}

func (p *pattern) moveLeft() {
	for i := range p {
		p[i] <<= 1
	}
}

func (p *pattern) moveRight() {
	for i := range p {
		p[i] >>= 1
	}
}

func (p *pattern) canMoveRight() bool {
	for i := range p {
		if p[i]&0b0000001 != 0 {
			return false
		}
	}
	return true
}

var rocks = []pattern{
	{
		0b0000000,
		0b0000000,
		0b0000000,
		0b0011110,
	}, {
		0b0000000,
		0b0001000,
		0b0011100,
		0b0001000,
	}, {
		0b0000000,
		0b0000100,
		0b0000100,
		0b0011100,
	}, {
		0b0010000,
		0b0010000,
		0b0010000,
		0b0010000,
	}, {
		0b0000000,
		0b0000000,
		0b0011000,
		0b0011000,
	},
}
