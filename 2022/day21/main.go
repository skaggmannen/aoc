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
	path := "test.txt"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	input, _ := os.ReadFile(path)

	fmt.Println("Part 1:", part1(bytes.NewReader(input)))
	fmt.Println("Part 2:", part2(bytes.NewReader(input)))
}

func part1(r io.Reader) int {
	monkeys := readMonkeys(r)

	return monkeys["root"].doJob(monkeys)
}

func part2(r io.Reader) int {
	monkeys := readMonkeys(r)

	stack := findMe(monkeys, monkeys["root"])

	root := stack[0]
	var expected int
	if stack[1].name == root.refs[0] {
		expected = monkeys[root.refs[1]].doJob(monkeys)
	} else {
		expected = monkeys[root.refs[0]].doJob(monkeys)
	}

	for i := 1; i < len(stack)-1; i++ {
		var operand int
		if stack[i+1].name == stack[i].refs[0] {
			// "Me" is down the left path, meaning we can calculate the right path
			operand = monkeys[stack[i].refs[1]].doJob(monkeys)

			// Now we can calculate the expected value for the left path
			switch stack[i].op {
			case "+":
				// oldExpected = newExpected + operand
				expected = expected - operand
			case "-":
				// oldExpected = newExpected - operand
				expected = expected + operand
			case "*":
				// oldExpected = newExpected * operand
				expected = expected / operand
			case "/":
				// oldExpected = newExpected / operand
				expected = expected * operand
			}
		} else {
			operand = monkeys[stack[i].refs[0]].doJob(monkeys)

			switch stack[i].op {
			case "+":
				// oldExpected = operand + newExpected
				expected = expected - operand
			case "-":
				// oldExpected = operand - newExpected
				expected = operand - expected
			case "*":
				// oldExpected = operand * oldExpected
				expected = expected / operand
			case "/":
				// oldExpected = operand / newExpected
				expected = operand / expected
			}
		}
	}

	return expected
}

func readMonkeys(r io.Reader) map[string]monkey {
	s := bufio.NewScanner(r)

	result := make(map[string]monkey)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}

		m := parseMonkey(line)
		result[m.name] = m
	}
	return result
}

func findMe(monkeys map[string]monkey, m monkey) []monkey {
	if m.name == "humn" {
		return []monkey{m}
	}
	if m.job == "num" {
		return nil
	}

	me := findMe(monkeys, monkeys[m.refs[0]])
	if me == nil {
		me = findMe(monkeys, monkeys[m.refs[1]])
	}

	if me == nil {
		return nil
	}

	return append([]monkey{m}, me...)
}

func parseMonkey(v string) monkey {
	parts := strings.Split(v, ": ")
	if len(parts) < 2 {
		panic("invalid monkey: " + v)
	}

	m := monkey{
		name: parts[0],
	}

	parts = strings.Split(parts[1], " ")
	if len(parts) == 1 {
		m.job = "num"
		m.num, _ = strconv.Atoi(parts[0])
	} else if len(parts) >= 3 {
		m.job = "math"
		m.op = parts[1]
		m.refs = []string{
			parts[0],
			parts[2],
		}
	} else {
		panic("invalid monkey: " + v)
	}

	return m
}

type monkey struct {
	name string
	job  string

	// if job is "num"
	num int

	// if job is "math"
	op   string
	refs []string
}

func (m monkey) doJob(monkeys map[string]monkey) (result int) {
	switch m.job {
	case "num":
		return m.num
	case "math":
		first := monkeys[m.refs[0]].doJob(monkeys)
		second := monkeys[m.refs[1]].doJob(monkeys)

		switch m.op {
		case "+":
			return first + second
		case "*":
			return first * second
		case "/":
			return first / second
		case "-":
			return first - second
		case "=":
			if first-second == 0 {
				return 1
			} else {
				return 0
			}
		default:
			panic("invalid op: " + m.op)
		}
	default:
		panic("invalid job: " + m.job)
	}
}
