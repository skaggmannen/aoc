package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func main() {
	input, _ := io.ReadAll(os.Stdin)
	// input, _ := os.ReadFile("test.txt")

	fmt.Println("Part 1:", part1(bytes.NewReader(input)))
	fmt.Println("Part 2:", part2(bytes.NewReader(input)))
}

func part1(r io.Reader) int {
	packetPairs := parsePacketPairs(r)

	var sum int
	for i := range packetPairs {
		if packetPairs[i].inOrder() {
			fmt.Println(i+1, ":", packetPairs[i], "=> in order")
			sum += i + 1
		} else {
			fmt.Println(i+1, ":", packetPairs[i], "=> not in order")
		}
	}
	return sum
}

func part2(r io.Reader) int {
	var (
		firstMarker  = packet{[]any{2}}
		secondMarker = packet{[]any{6}}
	)
	packets := parsePackets(r)
	packets = append(packets, firstMarker)
	packets = append(packets, secondMarker)

	sort.Slice(packets, func(i, j int) bool {
		return compare([]any(packets[i]), []any(packets[j])) < 1
	})

	result := 1

	for i, p := range packets {
		if reflect.DeepEqual(p, firstMarker) || reflect.DeepEqual(p, secondMarker) {
			result *= (i + 1)
		}
	}

	return result
}

func parsePacketPairs(r io.Reader) []packetPair {
	data, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	pairs := strings.Split(string(data), "\n\n")

	result := make([]packetPair, 0, len(pairs))

	for _, p := range pairs {
		p = strings.TrimSpace(p)
		if len(p) == 0 {
			continue
		}

		result = append(result, parsePacketPair(p))
	}

	return result
}

func parsePackets(r io.Reader) []packet {
	data, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")

	result := make([]packet, 0, len(lines))

	for _, p := range lines {
		p = strings.TrimSpace(p)
		if len(p) == 0 {
			continue
		}

		result = append(result, parseList(p))
	}

	return result
}

func parsePacketPair(s string) packetPair {
	parts := strings.Split(s, "\n")

	return packetPair{
		left:  parseList(parts[0]),
		right: parseList(parts[1]),
	}
}

func parseList(s string) []any {
	var (
		curr  []any
		stack stack[[]any]
		work  string
	)

	for _, b := range s {
		switch b {
		case '[':
			if curr != nil {
				stack.push(curr)
			}

			curr = make([]any, 0)
		case ']':
			if len(work) > 0 {
				v, _ := strconv.Atoi(work)
				curr = append(curr, v)
				work = ""
			}
			if len(stack) > 0 {
				curr = append(stack.pop(), curr)
			}
		case ',':
			if len(work) > 0 {
				v, _ := strconv.Atoi(work)
				curr = append(curr, v)
				work = ""
			}
		default:
			work += string(b)
		}
	}

	return curr
}

type packetPair struct {
	left  packet
	right packet
}

func (p packetPair) inOrder() bool {
	return compare([]any(p.left), []any(p.right)) < 0
}

func (p packetPair) String() string {
	return fmt.Sprint("{left:", p.left, ", right:", p.right, "}")
}

type packet []any

func compare(left any, right any) int {
	var (
		leftInt, leftIsInt     = left.(int)
		rightInt, rightIsInt   = right.(int)
		leftList, leftIsList   = left.([]any)
		rightList, rightIsList = right.([]any)
	)

	if leftIsInt && rightIsList {
		return compare([]any{leftInt}, rightList)
	}

	if leftIsList && rightIsInt {
		return compare(leftList, []any{rightInt})
	}

	if leftIsList && rightIsList {
		for i := range leftList {
			if i >= len(rightList) {
				// Left is longer, meaning is it greater than right
				return len(leftList) - len(rightList)
			}

			cmp := compare(leftList[i], rightList[i])
			if cmp != 0 {
				return cmp
			}
		}

		// All values so far were equal. Compare the lengths.
		return len(leftList) - len(rightList)
	}

	// The only remaining combination is that both are ints so just compare them.
	return leftInt - rightInt
}

type stack[T any] []T

func (s *stack[T]) push(v T) {
	*s = append(*s, v)
}

func (s *stack[T]) pop() T {
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}
