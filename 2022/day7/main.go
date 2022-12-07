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
	input, _ := io.ReadAll(os.Stdin)

	fmt.Println("Part 1:", part1(bytes.NewReader(input)))
	fmt.Println("Part 2:", part2(bytes.NewReader(input)))
}

func part1(r io.Reader) int {
	var totalSize int

	dirSizes := make(map[string]int)

	var cwd Stack[string]

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		fmt.Println(line)

		if strings.HasPrefix(line, "$ cd ") {
			dir := strings.TrimPrefix(line, "$ cd ")
			switch dir {
			case "/":
				cwd = Stack[string]{}
			case "..":
				cwd.Pop()
			default:
				cwd.Push(dir)
			}

			fmt.Println(cwd.String())
		} else if strings.HasPrefix(line, "$ ls") {
			// Nothing to do really...
		} else if strings.HasPrefix(line, "dir ") {
			// Not here either...
		} else {
			parts := strings.Split(line, " ")
			size, _ := strconv.Atoi(parts[0])

			dirSizes["/"] += size

			var s Stack[string]
			for _, dir := range cwd {
				s.Push(dir)
				dirSizes[s.String()] += size
			}

			fmt.Println(dirSizes)
		}
	}

	for dir, size := range dirSizes {
		fmt.Printf("%-60s %10d\n", dir, size)
		if size <= 100000 {
			totalSize += size
		}
	}

	return totalSize
}

func part2(r io.Reader) int {
	dirSizes := make(map[string]int)

	var cwd Stack[string]

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		fmt.Println(line)

		if strings.HasPrefix(line, "$ cd ") {
			dir := strings.TrimPrefix(line, "$ cd ")
			switch dir {
			case "/":
				cwd = Stack[string]{}
			case "..":
				cwd.Pop()
			default:
				cwd.Push(dir)
			}

			fmt.Println(cwd.String())
		} else if strings.HasPrefix(line, "$ ls") {
			// Nothing to do really...
		} else if strings.HasPrefix(line, "dir ") {
			// Not here either...
		} else {
			parts := strings.Split(line, " ")
			size, _ := strconv.Atoi(parts[0])

			dirSizes["/"] += size

			var s Stack[string]
			for _, dir := range cwd {
				s.Push(dir)
				dirSizes[s.String()] += size
			}

			fmt.Println(dirSizes)
		}
	}

	rootSize := dirSizes["/"]
	freeSpace := 70000000 - rootSize
	deletedSize := rootSize

	fmt.Println(freeSpace)

	for _, size := range dirSizes {
		if (freeSpace+size) > 30000000 && size < deletedSize {
			deletedSize = size
		}
	}

	return deletedSize
}

type Stack[T any] []T

func (s *Stack[T]) Push(v T) {
	*s = append(*s, v)
}

func (s *Stack[T]) Pop() {
	*s = (*s)[:len(*s)-1]
}

func (s *Stack[T]) Peek() T {
	return (*s)[len(*s)-1]
}

func (s *Stack[T]) String() string {
	result := "/"

	for i, v := range *s {
		if i != 0 {
			result += "/"
		}

		result += fmt.Sprint(v)
	}
	return result
}
