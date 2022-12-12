package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
)

func main() {
	input, _ := io.ReadAll(os.Stdin)

	fmt.Println("Part 1:", part1(bytes.NewReader(input)))
	fmt.Println("Part 2:", part2(bytes.NewReader(input)))
}

func part1(r io.Reader) int {
	h := readMap(r)

	for _, n := range h {
		if n.start {
			n.state.distance = 0
		}
	}

	return findShortestPath(h)
}

func part2(r io.Reader) int {
	h := readMap(r)

	var startPositions []*Node
	for _, n := range h {
		if n.height == 'a' {
			startPositions = append(startPositions, n)
		}
	}

	shortest := 1_000_000_000

	for _, n := range startPositions {
		curr := h
		for _, n := range curr {
			n.Reset()
		}

		n.state.distance = 0

		distance := findShortestPath(curr)

		if distance < shortest {
			shortest = distance
		}
	}

	return shortest
}

func findShortestPath(h []*Node) int {
	update := func(height byte, n *Node, d int) {
		if n == nil || n.state.visited || n.state.distance <= d {
			return
		}

		if n.height > height+1 {
			return
		}

		n.state.distance = d + 1
		sort.Slice(h, func(i, j int) bool {
			return h[i].state.distance > h[j].state.distance
		})
	}

	sort.Slice(h, func(i, j int) bool {
		return h[i].state.distance > h[j].state.distance
	})

	var (
		end  *Node
		curr *Node
	)

	for len(h) > 0 {
		curr = pop(&h)
		if curr.end {
			end = curr
			break
		} else if curr.state.distance == 1_000_000_000 {
			break
		}

		update(curr.height, curr.north, curr.state.distance)
		update(curr.height, curr.south, curr.state.distance)
		update(curr.height, curr.east, curr.state.distance)
		update(curr.height, curr.west, curr.state.distance)

		curr.state.visited = true
	}

	if end == nil {
		return 1_000_000_000
	}

	return end.state.distance
}

func readMap(r io.Reader) []*Node {
	type pos struct {
		x, y int
	}

	var (
		curr          pos
		h             []*Node
		width, height int
	)
	nodes := make(map[pos]*Node)

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		curr.x = 0
		width = 0

		for _, b := range line {
			n := &Node{height: byte(b)}
			n.state.distance = 1_000_000_000 // "infinity"
			if b == 'S' {
				n.height = 'a'
				n.start = true
			} else if b == 'E' {
				n.height = 'z'
				n.end = true
			}
			nodes[curr] = n
			curr.x++
			width++
		}

		curr.y++
		height++
	}

	for p, n := range nodes {
		if p.y > 0 {
			n.north = nodes[pos{x: p.x, y: p.y - 1}]
		}
		if p.x < width-1 {
			n.east = nodes[pos{x: p.x + 1, y: p.y}]
		}
		if p.y < height-1 {
			n.south = nodes[pos{x: p.x, y: p.y + 1}]
		}
		if p.x > 0 {
			n.west = nodes[pos{x: p.x - 1, y: p.y}]
		}

		h = append(h, n)
	}

	return h
}

func pop(h *[]*Node) *Node {
	n := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return n
}

type Node struct {
	height byte
	end    bool
	start  bool

	state NodeState

	north *Node
	east  *Node
	south *Node
	west  *Node
}

type NodeState struct {
	distance int
	visited  bool
}

func (n *Node) String() string {
	return fmt.Sprintf("{%16p, %s, %2d, [N:%16p, E:%16p, S:%16p, W:%16p]}", n, string(n.height), n.state.distance, n.north, n.east, n.south, n.west)
}

func (n *Node) Reset() {
	n.state = NodeState{
		distance: 1_000_000_000,
	}
}
