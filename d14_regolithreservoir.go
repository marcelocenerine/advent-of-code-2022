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
	rockPaths, err := p.parse(input)
	if err != nil {
		return Result{}, err
	}

	source := point{x: 500, y: 0}
	return Result{
		Part1: strconv.Itoa(p.countPouredUnitsOfSand(rockPaths, source, -1)),
		Part2: strconv.Itoa(p.countPouredUnitsOfSand(rockPaths, source, 1)),
	}, nil
}

func (p RegolithReservoir) countPouredUnitsOfSand(paths []rockPath, source point, floorPadding int) int {
	result := 0
	cave := p.draw(paths, floorPadding)
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

type rockPath []point

func (p rockPath) full() rockPath {
	if len(p) <= 1 {
		return p
	}

	var result rockPath

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
	tiles    [][]tile
	abyss    []bool
	hasFloor bool
}

func (c *cave) tile(p point) tile {
	c.mustBeValid(p)

	if p.y >= len(c.tiles) {
		if c.hasFloor {
			return rock
		}
		return air
	}

	if p.x >= len(c.tiles[p.y]) {
		return air
	}

	return c.tiles[p.y][p.x]
}

func (c *cave) pourSand(source point) (point, bool) {
	if c.tile(source) != air || c.leadToAbyss(source) {
		return point{}, false
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
				break OUTER
			}

			if c.tile(move) == air {
				cur = move
				continue OUTER
			}
		}

		if cur.x == len(c.tiles[cur.y]) { // row needs to be expanded
			c.tiles[cur.y] = append(c.tiles[cur.y], sand)
		} else {
			c.tiles[cur.y][cur.x] = sand
		}
		return cur, true
	}

	return point{}, false
}

func (c *cave) mustBeValid(p point) {
	if p.x < 0 || p.y < 0 {
		// Infinite left/up is not supported.
		panic(fmt.Sprintf("invalid point %v", p))
	}
}

func (c *cave) leadToAbyss(p point) bool {
	c.mustBeValid(p)

	if c.hasFloor {
		return false
	}

	return p.y >= len(c.tiles) || p.x >= len(c.abyss) || c.abyss[p.x]
}

var pointRgx = regexp.MustCompile(`^(\d+),(\d+)$`)

func (p RegolithReservoir) parse(input *Input) ([]rockPath, error) {
	lines := input.Lines()
	result := make([]rockPath, len(lines))

	for i, line := range lines {
		var path rockPath
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

func (p RegolithReservoir) draw(paths []rockPath, floorPadding int) cave {
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

	hasFloor := floorPadding >= 0
	if hasFloor {
		rows += floorPadding
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
		abyss[c] = !hasFloor
	}

	// Draw rocks and mark off columns that don't lead to the abyss.
	for _, path := range paths {
		for _, point := range path.full() {
			tiles[point.y][point.x] = rock
			abyss[point.x] = false
		}
	}

	return cave{tiles, abyss, hasFloor}
}
