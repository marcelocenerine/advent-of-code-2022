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
		Part1: strconv.Itoa(countPositionsVisitedByTail(2, motions)),
		Part2: strconv.Itoa(countPositionsVisitedByTail(10, motions)),
	}, nil
}

type Motion struct {
	Dir   Delta
	Steps int
}

type Delta struct {
	X, Y int
}

type Rope []Knot

func (r Rope) tail() Knot {
	return r[len(r)-1]
}

type Knot struct {
	X, Y int
}

func (p Knot) delta(that Knot) Delta {
	return Delta{X: that.X - p.X, Y: that.Y - p.Y}
}

func (p Knot) isAdjacent(that Knot) bool {
	delta := p.delta(that)
	return math.Abs(float64(delta.X)) <= 1 && math.Abs(float64(delta.Y)) <= 1
}

func (p Knot) move(delta Delta) Knot {
	return Knot{X: p.X + delta.X, Y: p.Y + delta.Y}
}

func (p Knot) moveCloser(that Knot) Knot {
	delta := p.delta(that)
	move := Delta{
		X: int(math.Max(-1, math.Min(1, float64(delta.X)))),
		Y: int(math.Max(-1, math.Min(1, float64(delta.Y)))),
	}
	return p.move(move)
}

func (p Knot) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

func countPositionsVisitedByTail(knots int, motions []Motion) int {
	rope := makeRope(knots)
	visited := make(map[string]bool)
	visit := func(k Knot) { visited[k.String()] = true }
	visit(rope.tail())

	for _, motion := range motions {
	sloop:
		for step := 0; step < motion.Steps; step++ {
			// move head
			rope[0] = rope[0].move(motion.Dir)

			// move remaining knots if needed
			for k := 1; k < knots; k++ {
				curr := rope[k]
				prev := rope[k-1]

				if curr.isAdjacent(prev) {
					continue sloop
				}
				rope[k] = curr.moveCloser(prev)
			}
			visit(rope.tail())
		}
	}
	return len(visited)
}

func makeRope(knots int) Rope {
	rope := make([]Knot, knots)
	for i := 0; i < knots; i++ {
		rope[i] = Knot{}
	}
	return rope
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
