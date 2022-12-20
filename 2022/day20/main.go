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
	numbers := readNumbers(r)

	for _, n := range numbers {
		if n.v > 0 {
			n.shiftRight(n.v)
		}
		if n.v < 0 {
			n.shiftLeft(-n.v)
		}
	}

	var start *listNode
	for _, n := range numbers {
		if n.v == 0 {
			start = n
		}
	}

	var sum int

	n := start.next
	for i := 1; i <= 3001; i++ {
		if i == 1001 {
			sum += n.v
		}
		if i == 2001 {
			sum += n.v
		}
		if i == 3001 {
			sum += n.v
		}
		n = n.next
	}

	return sum
}

func part2(r io.Reader) int {
	numbers := readNumbers(r)

	for i := range numbers {
		numbers[i].v *= 811589153
		numbers[i].shift = numbers[i].v % (len(numbers) - 1)
	}

	for round := 0; round < 10; round++ {
		for _, n := range numbers {
			if n.v > 0 {
				n.shiftRight(n.shift)
			}
			if n.v < 0 {
				n.shiftLeft(-(n.shift))
			}
		}
	}

	var start *listNode
	for _, n := range numbers {
		if n.v == 0 {
			start = n
		}
	}

	var sum int

	n := start.next
	for i := 1; i <= 3000; i++ {
		if i == 1000 {
			sum += n.v
		}
		if i == 2000 {
			sum += n.v
		}
		if i == 3000 {
			sum += n.v
		}
		n = n.next
	}

	return sum
}

func printList(numbers []*listNode) {
	var start *listNode
	for _, n := range numbers {
		if n.v == 0 {
			start = n
		}
	}

	for {
		fmt.Printf("%d(%d)", start.v, start.shift)
		if start.next.v == 0 {
			break
		}

		fmt.Print("->")
		start = start.next
	}
	fmt.Println()
}

func readNumbers(r io.Reader) []*listNode {
	s := bufio.NewScanner(r)
	var prev *listNode

	var result []*listNode
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}

		v, _ := strconv.Atoi(line)
		n := &listNode{v: v}
		if prev != nil {
			n.prev = prev
			prev.next = n
		}
		result = append(result, n)
		prev = n
	}

	result[0].prev = result[len(result)-1]
	result[len(result)-1].next = result[0]

	return result
}

type listNode struct {
	v     int
	shift int
	prev  *listNode
	next  *listNode
}

func (n *listNode) shiftLeft(count int) {
	for i := 0; i < count; i++ {
		prev := n.prev
		next := n.next

		n.prev = prev.prev
		prev.prev.next = n

		prev.next = next
		next.prev = prev

		prev.prev = n
		n.next = prev
	}
}

func (n *listNode) shiftRight(count int) {
	for i := 0; i < count; i++ {
		prev := n.prev
		next := n.next

		n.next = next.next
		next.next.prev = n

		prev.next = next
		next.prev = prev

		next.next = n
		n.prev = next
	}
}

func (n *listNode) String() string {
	return fmt.Sprintf("%d<-%d->%d", n.prev.v, n.v, n.next.v)
}
