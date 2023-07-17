// Solution to https://adventofcode.com/2022/day/14
package adventofcode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type RegolithReservoir struct{}

func (p RegolithReservoir) Details() Details {
	return Details{Day: 14, Description: "Regolith Reservoir"}
}

func (p RegolithReservoir) Solve(input *Input) (Result, error) {
	paths, err := p.parse(input)
	if err != nil {
		return Result{}, err
	}

	cave := p.draw(paths)
	return Result{
		Part1: strconv.Itoa(p.countUnitsOfSandBeforeFlowingIntoAbyss(cave, point{x: 500, y: 0})),
		Part2: "TODO",
	}, nil
}

func (p RegolithReservoir) countUnitsOfSandBeforeFlowingIntoAbyss(cave cave, source point) int {
	result := 0
	for {
		if _, ok := cave.pourSand(source); !ok {
			break
		}
		result++
	}
	return result
}

const (
	air  tile = '.'
	rock tile = '#'
	sand tile = 'o'
)

type point struct {
	x, y int
}

type path []point

func (p path) full() path {
	if len(p) <= 1 {
		return p
	}

	var result path

	for i, curr := range p {
		if i == 0 {
			result = append(result, curr)
			continue
		}

		prev := p[i-1]

		if curr.y != prev.y { // horizontal path segment
			loY, hiY := curr.y, curr.y
			if loY > prev.y {
				loY = prev.y
			}
			if hiY < prev.y {
				hiY = prev.y
			}
			for y := loY; y <= hiY; y++ {
				result = append(result, point{curr.x, y})
			}
		} else { // vertical path segment
			loX, hiX := curr.x, curr.x
			if loX > prev.x {
				loX = prev.x
			}
			if hiX < prev.x {
				hiX = prev.x
			}
			for x := loX; x <= hiX; x++ {
				result = append(result, point{x, curr.y})
			}
		}
	}

	return result
}

type tile rune

type cave struct {
	rows, cols int
	tiles      [][]tile
	abyss      []bool
}

func (c *cave) tile(p point) tile {
	if !c.isWithinBounds(p) {
		return air
	}
	return c.tiles[p.y][p.x]
}

func (c *cave) pourSand(source point) (point, bool) {
	if c.tile(source) != air || c.leadToAbyss(source) {
		return point{-1, -1}, false
	}

	cur := source

OUTER:
	for {
		candidateMoves := []point{
			{x: cur.x, y: cur.y + 1},     // down
			{x: cur.x - 1, y: cur.y + 1}, // down left
			{x: cur.x + 1, y: cur.y + 1}, // down right
		}

		for _, move := range candidateMoves {
			if c.leadToAbyss(move) {
				c.abyss[cur.x] = true
				break OUTER
			}

			if c.tile(move) == air {
				cur = move
				continue OUTER
			}
		}

		c.tiles[cur.y][cur.x] = sand
		return cur, true
	}

	return point{-1, -1}, false
}

func (c *cave) leadToAbyss(p point) bool {
	return !c.isWithinBounds(p) || c.abyss[p.x]
}

func (c *cave) isWithinBounds(p point) bool {
	return p.x >= 0 && p.x < c.cols && p.y >= 0 && p.y < c.rows
}

var pointRgx = regexp.MustCompile(`^(\d+),(\d+)$`)

func (p RegolithReservoir) parse(input *Input) ([]path, error) {
	lines := input.Lines()
	result := make([]path, len(lines))

	for i, line := range lines {
		var path path
		for j, segment := range strings.Split(line, " -> ") {
			if !pointRgx.MatchString(segment) {
				return nil, fmt.Errorf("invalid segment on line %d: %s", i, segment)
			}
			groups := pointRgx.FindAllStringSubmatch(segment, -1)
			x, _ := strconv.Atoi(groups[0][1])
			y, _ := strconv.Atoi(groups[0][2])

			if len(path) > 0 && (path[j-1].x != x && path[j-1].y != y) {
				return nil, fmt.Errorf("invalid segment on line %d: %s", i, segment)
			}

			path = append(path, point{x, y})
		}
		result[i] = path
	}

	return result, nil
}

func (p RegolithReservoir) draw(paths []path) cave {
	if len(paths) == 0 {
		return cave{}
	}

	// Initialise data structures.
	var rows, cols int
	for _, path := range paths {
		for _, point := range path {
			if point.x >= cols {
				cols = point.x + 1
			}
			if point.y >= rows {
				rows = point.y + 1
			}
		}
	}

	tiles := make([][]tile, rows)
	for r := 0; r < rows; r++ {
		tiles[r] = make([]tile, cols)
		for c := 0; c < cols; c++ {
			tiles[r][c] = air
		}
	}

	abyss := make([]bool, cols)
	for c := 0; c < cols; c++ {
		abyss[c] = true
	}

	// Draw rocks and mark off columns that don't lead to the abyss.
	for _, path := range paths {
		for _, point := range path.full() {
			tiles[point.y][point.x] = rock
			abyss[point.x] = false
		}
	}

	return cave{rows, cols, tiles, abyss}
}
