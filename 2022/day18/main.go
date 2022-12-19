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
	cubes := readCubes(r)

	var area int
	for c := range cubes {
		if _, ok := cubes[cube{c.x + 1, c.y, c.z}]; !ok {
			area++
		}
		if _, ok := cubes[cube{c.x - 1, c.y, c.z}]; !ok {
			area++
		}
		if _, ok := cubes[cube{c.x, c.y + 1, c.z}]; !ok {
			area++
		}
		if _, ok := cubes[cube{c.x, c.y - 1, c.z}]; !ok {
			area++
		}
		if _, ok := cubes[cube{c.x, c.y, c.z + 1}]; !ok {
			area++
		}
		if _, ok := cubes[cube{c.x, c.y, c.z - 1}]; !ok {
			area++
		}
	}

	return area
}

func part2(r io.Reader) int {
	s := newSpace(readCubes(r))

	var area int
	for c := range s.cubes {
		if s.isOpenAir(cube{c.x + 1, c.y, c.z}, map[cube]struct{}{}) {
			area++
		}
		if s.isOpenAir(cube{x: c.x - 1, y: c.y, z: c.z}, map[cube]struct{}{}) {
			area++
		}
		if s.isOpenAir(cube{x: c.x, y: c.y + 1, z: c.z}, map[cube]struct{}{}) {
			area++
		}
		if s.isOpenAir(cube{x: c.x, y: c.y - 1, z: c.z}, map[cube]struct{}{}) {
			area++
		}
		if s.isOpenAir(cube{x: c.x, y: c.y, z: c.z + 1}, map[cube]struct{}{}) {
			area++
		}
		if s.isOpenAir(cube{x: c.x, y: c.y, z: c.z - 1}, map[cube]struct{}{}) {
			area++
		}
	}

	return area
}

func readCubes(r io.Reader) map[cube]struct{} {
	s := bufio.NewScanner(r)

	result := make(map[cube]struct{})
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, ",")
		result[cube{
			x: mustInt(parts[0]),
			y: mustInt(parts[1]),
			z: mustInt(parts[2]),
		}] = struct{}{}
	}
	return result
}

func newSpace(cubes map[cube]struct{}) *space {
	s := &space{
		cubes:   cubes,
		xPlane:  make(map[cube][]cube),
		yPlane:  make(map[cube][]cube),
		zPlane:  make(map[cube][]cube),
		openAir: make(map[cube]bool),
	}

	for c := range cubes {
		s.xPlane[cube{y: c.y, z: c.z}] = append(s.xPlane[cube{y: c.y, z: c.z}], c)
		s.yPlane[cube{x: c.x, z: c.z}] = append(s.yPlane[cube{x: c.x, z: c.z}], c)
		s.zPlane[cube{x: c.x, y: c.y}] = append(s.zPlane[cube{x: c.x, y: c.y}], c)
	}

	return s
}

type space struct {
	cubes   map[cube]struct{}
	xPlane  map[cube][]cube
	yPlane  map[cube][]cube
	zPlane  map[cube][]cube
	openAir map[cube]bool
}

func (s space) isOpenAir(c cube, body map[cube]struct{}) (result bool) {
	// Check if this cube has already been evaluated
	if _, ok := body[c]; ok {
		return false
	}
	body[c] = struct{}{}

	// Check so it's not part of the shape
	if _, ok := s.cubes[c]; ok {
		return false
	}

	// Check if we have a cached result
	if result, ok := s.openAir[c]; ok {
		return result
	}
	defer func() {
		s.openAir[c] = result
	}()

	// Check all directions for a hole
	if s.canSeeHoleX(c) || s.canSeeHoleY(c) || s.canSeeHoleZ(c) {
		return true
	}

	// Expand through the space and see if we can find a hole
	return s.isOpenAir(cube{c.x + 1, c.y, c.z}, body) ||
		s.isOpenAir(cube{c.x - 1, c.y, c.z}, body) ||
		s.isOpenAir(cube{c.x, c.y + 1, c.z}, body) ||
		s.isOpenAir(cube{c.x, c.y - 1, c.z}, body) ||
		s.isOpenAir(cube{c.x, c.y, c.z + 1}, body) ||
		s.isOpenAir(cube{c.x, c.y, c.z - 1}, body)
}

func (s space) canSeeHoleX(c cube) bool {
	var (
		posBlocked bool
		negBlocked bool
	)

	for _, o := range s.xPlane[cube{y: c.y, z: c.z}] {
		if o.x > c.x {
			posBlocked = true
		}
		if o.x < c.x {
			negBlocked = true
		}
	}

	return !(posBlocked && negBlocked)
}

func (s space) canSeeHoleY(c cube) bool {
	var (
		posBlocked bool
		negBlocked bool
	)

	for _, o := range s.yPlane[cube{x: c.x, z: c.z}] {
		if o.y > c.y {
			posBlocked = true
		}
		if o.y < c.y {
			negBlocked = true
		}
	}

	return !(posBlocked && negBlocked)
}

func (s space) canSeeHoleZ(c cube) bool {
	var (
		posBlocked bool
		negBlocked bool
	)

	for _, o := range s.zPlane[cube{x: c.x, y: c.y}] {
		if o.z > c.z {
			posBlocked = true
		}
		if o.z < c.z {
			negBlocked = true
		}
	}

	return !(posBlocked && negBlocked)
}

type cube struct {
	x, y, z int
}

func mustInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return v
}
