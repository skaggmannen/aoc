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
	input, _ := io.ReadAll(os.Stdin)

	fmt.Println("Part 1:", part1(bytes.NewReader(input)))
	fmt.Println("Part 2:", part2(bytes.NewReader(input)))
}

func part1(r io.Reader) string {
	var result string

	s := bufio.NewScanner(r)

	// Read the starting state
	state := readStartingState(s)

	var (
		pattern, _ = regexp.Compile(`move (\d+) from (\d+) to (\d+)`)
	)

	for s.Scan() {
		line := s.Text()

		matches := pattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		count, _ := strconv.Atoi(matches[1])
		from, _ := strconv.Atoi(matches[2])
		to, _ := strconv.Atoi(matches[3])

		for i := 0; i < count; i++ {
			state[to].Add(state[from].Remove(1)...)
		}
	}

	for i := range state {
		crates := state[i].Remove(1)
		if len(crates) > 0 {
			result += crates[0]
		}
	}

	return result
}

func readStartingState(s *bufio.Scanner) []Stack {
	var state []Stack

	for s.Scan() {
		line := s.Text()

		// The starting state ends with an empty line
		if line == "" {
			break
		}

		stacks := chunkify(line, 4)

		if state == nil {
			state = make([]Stack, len(stacks)+1)
		}

		for i, s := range stacks {
			s := strings.TrimSpace(s)
			if s == "" {
				// No crate for this stack on this level
				continue
			} else if _, err := strconv.Atoi(s); err == nil {
				// This was a number, meaning we are at the lowest level.
				continue
			} else {
				// Add the crate to the stack. Skip position 0 since the instructions are 1 based.
				state[i+1].Add(strings.Trim(s, "[]"))
			}
		}
	}

	for i := range state {
		state[i].Reverse()
	}

	return state
}

type Stack struct {
	entries []string
}

func (s *Stack) Remove(count int) []string {
	if len(s.entries) == 0 || count == 0 {
		return nil
	}

	if count > len(s.entries) {
		count = len(s.entries)
	}

	e := s.entries[len(s.entries)-count:]
	s.entries = s.entries[:len(s.entries)-count]
	return e
}

func (s *Stack) Add(entries ...string) {
	s.entries = append(s.entries, entries...)
}

func (s *Stack) Reverse() {
	for i := len(s.entries)/2 - 1; i >= 0; i-- {
		opp := len(s.entries) - 1 - i
		s.entries[i], s.entries[opp] = s.entries[opp], s.entries[i]
	}
}

func (s Stack) String() string {
	return "[" + strings.Join(s.entries, ",") + "]"
}

func chunkify(v string, chunkSize int) []string {
	result := make([]string, 0, len(v)/chunkSize+1)
	for i := 0; i < len(v); i += chunkSize {
		end := i + chunkSize
		if end >= len(v) {
			end = len(v) - 1
		}

		chunk := v[i:end]
		result = append(result, chunk)
	}
	return result
}

func part2(r io.Reader) string {
	var result string

	s := bufio.NewScanner(r)

	// Read the starting state
	state := readStartingState(s)

	var (
		pattern, _ = regexp.Compile(`move (\d+) from (\d+) to (\d+)`)
	)

	for s.Scan() {
		line := s.Text()

		matches := pattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		count, _ := strconv.Atoi(matches[1])
		from, _ := strconv.Atoi(matches[2])
		to, _ := strconv.Atoi(matches[3])

		state[to].Add(state[from].Remove(count)...)
	}

	for i := range state {
		crates := state[i].Remove(1)
		if len(crates) > 0 {
			result += crates[0]
		}
	}

	return result
}
