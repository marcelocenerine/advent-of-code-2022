// Solution to https://adventofcode.com/2022/day/8
package adventofcode

import (
	"errors"
	"fmt"
	"strconv"
)

type TreetopTreeHouse struct{}

func (p TreetopTreeHouse) Details() Details {
	return Details{Day: 8, Description: "Treetop Tree House"}
}

func (p TreetopTreeHouse) Solve(input *Input) (Result, error) {
	heights, err := parseTreeHeights(input)
	if err != nil {
		return Result{}, err
	}
	visibility := heights.VisibleTrees()
	scenicScores := heights.ScenicScores()
	return Result{
		Part1: strconv.Itoa(visibility.Count),
		Part2: strconv.Itoa(scenicScores.Max),
	}, nil
}

type TreeHeights [][]int

func (m TreeHeights) Height() int {
	return len(m)
}

func (m TreeHeights) Width() int {
	if m.Height() == 0 {
		return 0
	}
	return len(m[0])
}

// TODO good use case for dynamic programming
func (m TreeHeights) VisibleTrees() VisibilityMap {
	var tallestTree func(int, int, int, int) int
	tallestTree = func(r, c, rstep, cstep int) int {
		if r < 0 || r >= m.Height() || c < 0 || c >= m.Width() {
			return 0
		}
		currTree := m[r][c]
		nextTree := tallestTree(r+rstep, c+cstep, rstep, cstep)
		if nextTree > currTree {
			return nextTree
		}
		return currTree
	}

	count := 0
	visibilityMap := make([][]bool, m.Height())

	for r := 0; r < m.Height(); r++ {
		visibilityMap[r] = make([]bool, m.Width())
		for c := 0; c < m.Width(); c++ {
			height := m[r][c]
			visible :=
				r == 0 || r == m.Height()-1 || // top/bottom edges
					c == 0 || c == m.Width()-1 || // left/right edges
					height > tallestTree(r-1, c, -1, 0) || // up
					height > tallestTree(r+1, c, 1, 0) || // down
					height > tallestTree(r, c-1, 0, -1) || // left
					height > tallestTree(r, c+1, 0, 1) // right
			if visible {
				visibilityMap[r][c] = true
				count++
			}
		}
	}
	return VisibilityMap{visibilityMap, count}
}

func (m TreeHeights) ScenicScores() ScenicScores {
	var distance func(int, int, int, int, int) int
	distance = func(height, r, c, rstep, cstep int) int {
		if r+rstep < 0 || r+rstep >= m.Height() || c+cstep < 0 || c+cstep >= m.Width() {
			return 0
		}
		if height <= m[r+rstep][c+cstep] {
			return 1
		}
		return 1 + distance(height, r+rstep, c+cstep, rstep, cstep)
	}

	max := 0
	scores := make([][]int, m.Height())
	for r := 0; r < m.Height(); r++ {
		scores[r] = make([]int, m.Width())
		for c := 0; c < m.Width(); c++ {
			height := m[r][c]
			score :=
				distance(height, r, c, -1, 0) * // up
					distance(height, r, c, 1, 0) * // down
					distance(height, r, c, 0, -1) * // left
					distance(height, r, c, 0, 1) // right
			scores[r][c] = score
			if score > max {
				max = score
			}
		}
	}
	return ScenicScores{scores, max}
}

type VisibilityMap struct {
	Visible [][]bool
	Count   int
}

type ScenicScores struct {
	Scores [][]int
	Max    int
}

func parseTreeHeights(input *Input) (TreeHeights, error) {
	lines := input.Lines()
	height := len(lines)
	if height == 0 {
		return nil, errors.New("empty input")
	}
	width := len(lines[0])
	forest := make(TreeHeights, height)
	for r, line := range lines {
		if width == 0 || len(line) != width {
			return nil, fmt.Errorf("line %d has invalid width: %s", r, line)
		}
		forest[r] = make([]int, width)
		for c, tree := range line {
			treeHeight, err := strconv.Atoi(string(tree))
			if err != nil {
				return nil, err
			}
			forest[r][c] = treeHeight
		}
	}
	return forest, nil
}
