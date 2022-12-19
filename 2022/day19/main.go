package main

import (
	"aoc2022/util"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
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
	blueprints := readBlueprints(r)

	var total int
	for _, bp := range blueprints {
		score := newFactory(bp).run(24, stock{}, robots{1, 0, 0, 0}, map[int]struct{}{}, 0)
		total += bp.id * score
	}

	return total
}

func part2(r io.Reader) int {
	blueprints := readBlueprints(r)

	count := 3
	if len(blueprints) < count {
		count = len(blueprints)
	}

	total := 1
	for _, bp := range blueprints[:count] {
		score := newFactory(bp).run(32, stock{}, robots{1, 0, 0, 0}, map[int]struct{}{}, 0)
		total *= score
	}

	return total
}

func readBlueprints(r io.Reader) []blueprint {
	pattern := regexp.MustCompile(`Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.`)

	s := bufio.NewScanner(r)

	var result []blueprint
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}

		matches := pattern.FindStringSubmatch(line)

		values := util.ToInts(matches[1:])

		result = append(result, blueprint{
			id: values[0],
			costs: map[int]cost{
				Ore:      {ore: values[1]},
				Clay:     {ore: values[2]},
				Obsidian: {ore: values[3], clay: values[4]},
				Geode:    {ore: values[5], obsidian: values[6]},
			},
		})
	}
	return result
}

type blueprint struct {
	id    int
	costs map[int]cost
}

type cost struct {
	ore      int
	clay     int
	obsidian int
}

func newFactory(b blueprint) *factory {
	return &factory{
		blueprint: b,
		cache:     make(map[string]int),
	}
}

type factory struct {
	blueprint blueprint
	cache     map[string]int
}

func (f *factory) run(time int, stock stock, robots robots, previouslySkipped map[int]struct{}, max int) int {
	if time == 1 {
		// There is not enough time for any new robots to produce material
		return stock[Geode] + robots[Geode]
	}

	// If there is no chance at all for this branch to produce any more Geodes
	// than the currently known best we can ignore it altoghether
	optimistic_max := stock[Geode] + robots[Geode]*time + time*(time-1)/2
	if optimistic_max < max {
		return 0
	}

	// If there's no chance we can produce enough obsidian to produce another Geode
	// we can calculate the final score right now.
	optimistic_max = stock[Obsidian] + robots[Obsidian]*time + time*(time-1)/2
	if optimistic_max < f.blueprint.costs[Geode].obsidian {
		return stock[Geode] + robots[Geode]*time
	}

	// The robots will have collected some resources until next iteration.
	nextStock := stock
	for r, c := range robots {
		nextStock[r] += c
	}

	// If we can afford a Geode robot it's always the best choice
	if stock.canAfford(f.blueprint.costs[Geode]) {
		return f.run(
			time-1,
			nextStock.withdraw(f.blueprint.costs[Geode]),
			robots.build(Geode),
			map[int]struct{}{},
			max,
		)
	}

	// Now check the remaining options.
	skipped := make(map[int]struct{})
	for i := 0; i <= Obsidian; i++ {
		if stock.canAfford(f.blueprint.costs[i]) {
			if _, ok := previouslySkipped[i]; ok {
				// Building something now that we could've afforded previous round
				// will never be a good choice
				continue
			}

			score := f.run(
				time-1,
				nextStock.withdraw(f.blueprint.costs[i]),
				robots.build(i),
				map[int]struct{}{},
				max,
			)
			if score > max {
				max = score
			}
			skipped[i] = struct{}{}
		}
	}

	// We should also consider not building any robots at all
	score := f.run(time-1, nextStock, robots, skipped, max)
	if score > max {
		max = score
	}

	return max
}

const (
	Ore = iota
	Clay
	Obsidian
	Geode
)

type stock [4]int

func (s stock) canAfford(c cost) bool {
	return s[Ore] >= c.ore && s[Clay] >= c.clay && s[Obsidian] >= c.obsidian
}

func (s stock) withdraw(c cost) stock {
	return stock{
		s[Ore] - c.ore,
		s[Clay] - c.clay,
		s[Obsidian] - c.obsidian,
		s[Geode],
	}
}

func (s stock) String() string {
	return fmt.Sprintf("[%d,%d,%d,%d]", s[0], s[1], s[2], s[3])
}

type robots [4]int

func (r robots) build(b int) robots {
	next := r
	next[b]++
	return next
}

func (r robots) String() string {
	return fmt.Sprintf("[%d,%d,%d,%d]", r[0], r[1], r[2], r[3])
}
