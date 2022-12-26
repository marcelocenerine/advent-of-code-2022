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
	forest, err := parseTreeHeightMap(input)
	if err != nil {
		return Result{}, err
	}
	visibility := forest.VisibilityMap()
	return Result{
		Part1: strconv.Itoa(visibility.Count()),
		Part2: "",
	}, nil
}

type TreeHeightMap [][]int

func (m TreeHeightMap) Height() int {
	return len(m)
}

func (m TreeHeightMap) Width() int {
	if m.Height() == 0 {
		return 0
	}
	return len(m[0])
}

// TODO good use case for dynamic programming
func (m TreeHeightMap) VisibilityMap() TreeVisibilityMap {
	visibility := make(TreeVisibilityMap, m.Height())
	for r := 0; r < m.Height(); r++ {
		visibility[r] = make([]bool, m.Width())
		for c := 0; c < m.Width(); c++ {
			visibility[r][c] = m.isTreeVisible(r, c)
		}
	}
	return visibility
}

func (m TreeHeightMap) isTreeVisible(r, c int) bool {
	if r == 0 || r == m.Height()-1 { // top/bottom edges
		return true
	}
	if c == 0 || c == m.Width()-1 { // left/right edges
		return true
	}
	tree := m[r][c]
	return tree > m.tallestTree(r-1, c, -1, 0) || // up
		tree > m.tallestTree(r+1, c, 1, 0) || // down
		tree > m.tallestTree(r, c-1, 0, -1) || // left
		tree > m.tallestTree(r, c+1, 0, 1) // right
}

func (m TreeHeightMap) tallestTree(r, c, rstep, cstep int) int {
	if r < 0 || r >= m.Height() || c < 0 || c >= m.Width() {
		return 0
	}
	currTree := m[r][c]
	if nextTree := m.tallestTree(r+rstep, c+cstep, rstep, cstep); nextTree > currTree {
		return nextTree
	}
	return currTree
}

type TreeVisibilityMap [][]bool

func (m TreeVisibilityMap) Count() int {
	count := 0
	for _, row := range m {
		for _, col := range row {
			if col {
				count++
			}
		}
	}
	return count
}

func parseTreeHeightMap(input *Input) (TreeHeightMap, error) {
	lines := input.Lines()
	height := len(lines)
	if height == 0 {
		return nil, errors.New("empty input")
	}
	width := len(lines[0])
	forest := make(TreeHeightMap, height)
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
