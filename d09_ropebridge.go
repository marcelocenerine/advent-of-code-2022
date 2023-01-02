// Solution to https://adventofcode.com/2022/day/9
package adventofcode

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

type RopeBridge struct{}

func (p RopeBridge) Details() Details {
	return Details{Day: 9, Description: "Rope Bridge"}
}

func (p RopeBridge) Solve(input *Input) (Result, error) {
	motions, err := parseMotions(input)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Part1: strconv.Itoa(countPositionsVisitedByTail(motions)),
		Part2: "?",
	}, nil
}

type Delta struct {
	X, Y int
}

type Motion struct {
	Dir   Delta
	Steps int
}

func distance(hx, hy, tx, ty int) int {
	xDist := math.Abs(float64(hx - tx))
	yDist := math.Abs(float64(hy - ty))
	return int(math.Max(xDist, yDist))
}

func countPositionsVisitedByTail(motions []Motion) int {
	var headX, headY, tailX, tailY int
	visited := make(map[string]bool)
	visit := func(x, y int) { visited[fmt.Sprintf("%d,%d", x, y)] = true }
	visit(tailX, tailY)

	for _, motion := range motions {
		for step := 1; step <= motion.Steps; step++ {
			// move head
			headX += motion.Dir.X
			headY += motion.Dir.Y

			if distance(headX, headY, tailX, tailY) <= 1 {
				continue
			}
			// move tail in the same direction
			tailX += motion.Dir.X
			tailY += motion.Dir.Y

			// same row or column?
			if headX == tailX || headY == tailY {
				visit(tailX, tailY)
				continue
			}

			// move in the diagonal
			if distance(headX, headY, tailX+motion.Dir.Y, tailY+motion.Dir.X) == 1 {
				tailX += motion.Dir.Y
				tailY += motion.Dir.X
			} else {
				tailX -= motion.Dir.Y
				tailY -= motion.Dir.X
			}
			visit(tailX, tailY)
		}
	}
	return len(visited)
}

func parseMotions(input *Input) ([]Motion, error) {
	lines := input.Lines()
	result := make([]Motion, len(lines))
	rgx := regexp.MustCompile(`^([LRUD]) (\d+)$`)
	deltas := map[string]Delta{
		"R": {0, 1},
		"L": {0, -1},
		"U": {-1, 0},
		"D": {1, 0},
	}
	for i, line := range lines {
		if !rgx.MatchString(line) {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		groups := rgx.FindAllStringSubmatch(line, -1)
		delta := deltas[groups[0][1]]
		steps, _ := strconv.Atoi(groups[0][2])
		result[i] = Motion{Dir: delta, Steps: steps}
	}
	return result, nil
}
