package main

import (
	"aoc2022/util"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println("Part 1:", part1(bytes.NewReader(input)))
	fmt.Println("Part 2:", part2(bytes.NewReader(input)))
}

func part1(r io.Reader) int {
	g := NewGame(r)
	g.worryReducer = func(w uint64) uint64 {
		return w / 3
	}

	for round := 0; round < 20; round++ {
		g.playRound()
	}

	sort.Slice(g.monkeys, func(i, j int) bool {
		return g.monkeys[i].itemsInspected > g.monkeys[j].itemsInspected
	})

	return g.monkeys[0].itemsInspected * g.monkeys[1].itemsInspected
}

func part2(r io.Reader) int {
	g := NewGame(r)

	g.worryReducer = func(w uint64) uint64 {
		return w % mod
	}

	for round := 0; round < 10000; round++ {
		g.playRound()
	}

	sort.Slice(g.monkeys, func(i, j int) bool {
		return g.monkeys[i].itemsInspected > g.monkeys[j].itemsInspected
	})

	return g.monkeys[0].itemsInspected * g.monkeys[1].itemsInspected
}

func NewGame(r io.Reader) Game {
	lines := util.ReadAllLines(r)

	var monkeys []*Monkey
	for i := 0; i < len(lines); i += 6 {
		monkeys = append(monkeys, &Monkey{
			items:       parseStartingItems(lines[i+1]),
			operation:   parseOperation(lines[i+2]),
			divisibleBy: parseTest(lines[i+3]),
			ifTrue:      parseIfTrue(lines[i+4]),
			ifFalse:     parseIfFalse(lines[i+5]),
		})
	}

	return Game{monkeys: monkeys}
}

var mod uint64

type Game struct {
	monkeys      []*Monkey
	worryReducer func(w uint64) uint64
}

func (g Game) playRound() {
	for _, m := range g.monkeys {
		if len(m.items) == 0 {
			continue
		}

		for _, item := range m.items {
			m.itemsInspected++

			worry := m.operation(item)
			worry = g.worryReducer(worry)

			if worry%uint64(m.divisibleBy) == 0 {
				g.monkeys[m.ifTrue].items = append(g.monkeys[m.ifTrue].items, worry)
			} else {
				g.monkeys[m.ifFalse].items = append(g.monkeys[m.ifFalse].items, worry)
			}

		}

		m.items = nil
	}
}

func (g Game) PrintState() {
	for i, m := range g.monkeys {
		fmt.Println("Monkey:", i)
		fmt.Println("	Items:", m.items)
		fmt.Println("	ItemsInspected:", m.itemsInspected)
	}
}

type Monkey struct {
	items          []uint64
	operation      func(w uint64) uint64
	divisibleBy    int
	ifTrue         int
	ifFalse        int
	itemsInspected int
}

func parseStartingItems(v string) []uint64 {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "Starting items: ")
	parts := strings.Split(v, ", ")
	items := make([]uint64, 0, len(parts))
	for _, p := range parts {
		item, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}

		items = append(items, uint64(item))
	}
	return items
}

func parseOperation(v string) func(w uint64) uint64 {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "Operation: new = ")
	parts := strings.Split(v, " ")

	var v1 func(w uint64) uint64
	if parts[0] == "old" {
		v1 = func(w uint64) uint64 {
			return w
		}
	} else {
		v, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}
		v1 = func(w uint64) uint64 {
			return uint64(v)
		}
	}
	var v2 func(w uint64) uint64
	if parts[2] == "old" {
		v2 = func(w uint64) uint64 {
			return w
		}
	} else {
		v, err := strconv.Atoi(parts[2])
		if err != nil {
			panic(err)
		}
		v2 = func(w uint64) uint64 {
			return uint64(v)
		}
	}

	return func(w uint64) uint64 {
		op := parts[1]

		switch op {
		case "+":
			return (v1(w) + v2(w))
		case "*":
			return (v1(w) * v2(w))
		default:
			panic("unknow op " + op)
		}
	}
}

func parseTest(v string) int {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "Test: divisible by ")
	divisor, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}

	if mod == 0 {
		mod = 1
	}

	mod = lcm(mod, uint64(divisor))

	return divisor
}

func parseIfTrue(v string) int {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "If true: throw to monkey ")
	monkey, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}

	return monkey
}

func parseIfFalse(v string) int {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "If false: throw to monkey ")
	monkey, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}

	return monkey
}

// function to calculate gcd
// or hcf of two numbers.
func gcd(a, b uint64) uint64 {
	if a == 0 {
		return b
	}
	return gcd(b%a, a)
}

// function to calculate
// lcm of two numbers.
func lcm(a, b uint64) uint64 {
	return (a * b) / gcd(a, b)
}
