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
	heightMap, err := parseTreeHeights(input)
	if err != nil {
		return Result{}, err
	}
	visibilityMap := heightMap.VisibleTrees()
	return Result{
		Part1: strconv.Itoa(visibilityMap.Count),
		Part2: "",
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

func (m TreeHeights) VisibleTrees() VisibilityMap {
	isEdge := func(row, col int) bool {
		return row == 0 || row == m.Height()-1 || col == 0 || col == m.Width()-1
	}
	count := 0
	visibilityMap := make([][]bool, m.Height())
	topDown, bottomUp, leftRight, rightLeft := m.tallestTrees()

	for row := 0; row < m.Height(); row++ {
		visibilityMap[row] = make([]bool, m.Width())
		for col := 0; col < m.Width(); col++ {
			if visible := isEdge(row, col) ||
				m[row][col] > topDown[row-1][col] ||
				m[row][col] > bottomUp[row+1][col] ||
				m[row][col] > leftRight[row][col-1] ||
				m[row][col] > rightLeft[row][col+1]; visible {
				visibilityMap[row][col] = true
				count++
			}
		}
	}

	return VisibilityMap{visibilityMap, count}
}

func (m TreeHeights) tallestTrees() (TreeHeights, TreeHeights, TreeHeights, TreeHeights) {
	topDown := make(TreeHeights, m.Height())
	bottomUp := make(TreeHeights, m.Height())
	leftRight := make(TreeHeights, m.Height())
	rightLeft := make(TreeHeights, m.Height())

	for r := 0; r < m.Height(); r++ {
		topDown[r] = make([]int, m.Width())
		leftRight[r] = make([]int, m.Width())

		for c := 0; c < m.Width(); c++ {
			if r == 0 || m[r][c] > topDown[r-1][c] {
				topDown[r][c] = m[r][c]
			} else {
				topDown[r][c] = topDown[r-1][c]
			}

			if c == 0 || m[r][c] > leftRight[r][c-1] {
				leftRight[r][c] = m[r][c]
			} else {
				leftRight[r][c] = leftRight[r][c-1]
			}
		}
	}

	for r := m.Height() - 1; r >= 0; r-- {
		bottomUp[r] = make([]int, m.Width())
		rightLeft[r] = make([]int, m.Width())

		for c := m.Width() - 1; c >= 0; c-- {
			if r == m.Height()-1 || m[r][c] > bottomUp[r+1][c] {
				bottomUp[r][c] = m[r][c]
			} else {
				bottomUp[r][c] = bottomUp[r+1][c]
			}

			if c == m.Width()-1 || m[r][c] > rightLeft[r][c+1] {
				rightLeft[r][c] = m[r][c]
			} else {
				rightLeft[r][c] = rightLeft[r][c+1]
			}
		}
	}

	return topDown, bottomUp, leftRight, rightLeft
}

type VisibilityMap struct {
	Visible [][]bool
	Count   int
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
