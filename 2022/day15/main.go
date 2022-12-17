package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

var (
	y     int
	maxXY int
	logf  func(fmt string, args ...any) (int, error) = fmt.Printf
)

func main() {
	path := "test.txt"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	if path == "test.txt" {
		y = 10
		maxXY = 20
	} else {
		logf = func(fmt string, args ...any) (int, error) { return 0, nil }
		y = 2_000_000
		maxXY = 4_000_000
	}

	input, _ := os.ReadFile(path)

	fmt.Println("Part 1:", part1(bytes.NewReader(input)))
	fmt.Println("Part 2:", part2(bytes.NewReader(input)))
}

func part1(r io.Reader) int {
	sensors := readSensors(r)

	min, max := getBounds(sensors)

	var result int
	logf("(%2d, %2d) ", min.x, y)
	for x := min.x; x <= max.x; x++ {
		c := findCoverage(sensors, point{x, y})
		if c.sensor != nil && c.beacon == nil {
			logf("#")
			result++
		} else if c.beacon != nil {
			logf("B")
		} else {
			logf(".")
		}
	}
	logf(" (%2d, %2d)\n", max.x, y)
	return result
}

func part2(r io.Reader) int64 {
	sensors := readSensors(r)

	min, max := getBounds(sensors)
	if min.x < 0 {
		min.x = 0
	}
	if min.y < 0 {
		min.y = 0
	}
	if max.x > maxXY {
		max.x = maxXY
	}
	if max.y > maxXY {
		max.y = maxXY
	}

	logf = fmt.Printf

	var counted int
	for y := min.y; y <= max.y; y++ {
		for x := min.x; x < max.x; x++ {
			p := point{x, y}
			c := findCoverage(sensors, p)
			if c.sensor == nil {
				return int64(x)*4_000_000 + int64(y)
			}

			xb := c.sensor.d + c.sensor.pos.x - abs(p.y-c.sensor.pos.y)

			counted += xb - x
			if counted%10_000_000 == 0 {
				fmt.Println(counted)
			}
			x = xb
		}
	}
	return 0
}

var (
	pattern, _ = regexp.Compile(`Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`)
)

func readSensors(r io.Reader) []sensor {
	s := bufio.NewScanner(r)
	var result []sensor
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			continue
		}

		result = append(result, parseSensor(line))
	}
	return result
}

func parseSensor(v string) sensor {
	matches := pattern.FindStringSubmatch(v)
	if len(matches) < 5 {
		return sensor{}
	}

	s := sensor{
		pos:           point{mustInt(matches[1]), mustInt(matches[2])},
		closestBeacon: point{mustInt(matches[3]), mustInt(matches[4])},
	}

	s.d = manhattanDistance(s.pos, s.closestBeacon)

	return s
}

type sensor struct {
	pos           point
	closestBeacon point
	d             int
}

func (s sensor) covers(p point) bool {
	d := manhattanDistance(s.pos, p)

	return d <= s.d
}

func (s sensor) String() string {
	return fmt.Sprintf("{pos:%s, d:%d}", s.pos, s.d)
}

type point struct {
	x, y int
}

func (p point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

func manhattanDistance(from, to point) int {
	d := abs(to.x - from.x)
	d += abs(to.y - from.y)
	return d
}

func abs(a int) int {
	if a < 0 {
		return -a
	}

	return a
}

func mustInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return v
}

func getBounds(sensors []sensor) (point, point) {
	min := point{0x7fff_ffff, 0x7fff_ffff}
	max := point{-0x7fff_ffff, -0x7fff_ffff}

	for _, s := range sensors {
		d := manhattanDistance(s.pos, s.closestBeacon)

		if s.pos.x-d < min.x {
			min.x = s.pos.x - d
		}
		if s.pos.y-d < min.y {
			min.y = s.pos.y - d
		}
		if s.pos.y+d > max.y {
			max.y = s.pos.y + d
		}
		if s.pos.x+d > max.x {
			max.x = s.pos.x + d
		}
	}

	return min, max
}

func findCoverage(sensors []sensor, p point) coverage {
	var c coverage
	for _, s := range sensors {
		if s.covers(p) {
			covered := s
			c.sensor = &covered
			if p == s.closestBeacon {
				c.beacon = &s.closestBeacon
			}
		}
	}

	return c
}

type coverage struct {
	beacon *point
	sensor *sensor
}

// y1 = x1 + b
// y2 = x2 + b
// y2 = -x2 + b
// y = -x + b
