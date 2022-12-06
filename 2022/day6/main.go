package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	input, _ := io.ReadAll(os.Stdin)

	fmt.Println("Part 1:", part1(bytes.NewReader(input)))
	fmt.Println("Part 2:", part2(bytes.NewReader(input)))
}

func part1(r io.Reader) int {
	data, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	for i := 0; i+4 <= len(data); i++ {
		marker := data[i : i+4]
		if isValidMarker(marker) {
			return i + 4
		}
	}

	return -1
}

func part2(r io.Reader) int {
	data, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	for i := 0; i+14 <= len(data); i++ {
		marker := data[i : i+14]
		if isValidMarker(marker) {
			return i + 14
		}
	}

	return -1
}

func isValidMarker(m []byte) bool {
	for i, b := range m {
		if bytes.Contains(m[i+1:], []byte{b}) {
			return false
		}
	}

	return true
}
