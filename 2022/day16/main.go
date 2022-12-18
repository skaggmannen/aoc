package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
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
	valves = readValves(r)
	fmt.Println(valves)

	cache = make(map[string]int)
	return maxScore("AA", 0, 30, 0)
}

func part2(r io.Reader) int {
	valves = readValves(r)
	fmt.Println(valves)

	cache = make(map[string]int)
	return maxScore("AA", 0, 26, 1)
}

var (
	valves map[string]valve
	cache  map[string]int
)

func maxScore(v string, openedValves int64, timeLeft int, elephants int) int {
	if timeLeft == 0 {
		if elephants > 0 {
			return maxScore("AA", openedValves, 26, elephants-1)
		}
		return 0
	}

	key := fmt.Sprintf("%s:%d:%d:%d", v, openedValves, timeLeft, elephants)
	if v, ok := cache[key]; ok {
		return v
	}

	currentValve := valves[v]

	var score int
	if ((openedValves & currentValve.id) == 0) && currentValve.flowRate > 0 {
		score = (timeLeft-1)*currentValve.flowRate + maxScore(v, openedValves|currentValve.id, timeLeft-1, elephants)
	}

	for _, n := range currentValve.leadsTo {
		s := maxScore(n, openedValves, timeLeft-1, elephants)
		if s > score {
			score = s
		}
	}

	cache[key] = score

	return score
}

func readValves(r io.Reader) map[string]valve {
	pattern := regexp.MustCompile(`^Valve ([A-Z]+) has flow rate=(\d+); tunnels? leads? to valves? (.*)$`)

	s := bufio.NewScanner(r)

	result := make(map[string]valve)

	var i int
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			continue
		}

		matches := pattern.FindStringSubmatch(line)
		name := matches[1]
		flowRate, _ := strconv.Atoi(matches[2])
		leadsTo := strings.Split(matches[3], ", ")

		result[name] = valve{id: 1 << i, name: name, flowRate: flowRate, leadsTo: leadsTo}
		i++
	}
	return result
}

type valve struct {
	name     string
	id       int64
	flowRate int
	leadsTo  []string
}

func (v valve) String() string {
	return fmt.Sprintf("{name:%s, flow:%d, leadsTo:[%s]", v.name, v.flowRate, strings.Join(v.leadsTo, ","))
}
