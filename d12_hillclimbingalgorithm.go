// Solution to https://adventofcode.com/2022/day/12
package adventofcode

import (
	"container/heap"
	"errors"
	"math"
	"strconv"
)

type HillClimbingAlgorithm struct{}

func (p HillClimbingAlgorithm) Details() Details {
	return Details{Day: 12, Description: "Hill Climbing Algorithm"}
}

func (p HillClimbingAlgorithm) Solve(input *Input) (Result, error) {
	hm, err := p.parseHeightmap(input)
	if err != nil {
		return Result{}, err
	}

	part1 := strconv.Itoa(p.shortestPath(hm))
	return Result{Part1: part1, Part2: "TODO"}, nil
}

func (p HillClimbingAlgorithm) shortestPath(hm *heightmap) int {
	dist := make([]int, hm.size())
	for i := 0; i < len(dist); i++ {
		dist[i] = math.MaxInt
	}
	dist[hm.index(hm.start)] = 0

	visited := make([]bool, hm.size())
	pq := priorityQueue{
		&square{hm.start, 0},
	}

	for len(pq) > 0 {
		sq := heap.Pop(&pq).(*square)
		sqi := hm.index(sq.pos)

		if visited[sqi] {
			continue
		}

		visited[sqi] = true

		for _, adj := range hm.neighbors(sq.pos) {
			adji := hm.index(adj)

			if visited[adji] || !hm.canMove(sq.pos, adj) {
				continue
			}

			oldDist := dist[adji]
			newDist := sq.dist + 1

			if newDist < oldDist {
				dist[adji] = newDist
				heap.Push(&pq, &square{adj, newDist})
			}
		}
	}

	return dist[hm.index(hm.dest)]
}

type pos struct {
	row, col int
}

type heightmap struct {
	start, dest   pos
	width, height int
	grid          [][]rune
}

func (h *heightmap) size() int {
	return h.height * h.width
}

func (h *heightmap) index(p pos) int {
	return p.row*h.width + p.col
}

func (h *heightmap) canMove(from, to pos) bool {
	return h.grid[to.row][to.col]-h.grid[from.row][from.col] <= 1
}

func (h *heightmap) neighbors(p pos) []pos {
	res := make([]pos, 0, 4)
	if p.row > 0 {
		res = append(res, pos{p.row - 1, p.col})
	}
	if p.row < h.height-1 {
		res = append(res, pos{p.row + 1, p.col})
	}
	if p.col > 0 {
		res = append(res, pos{p.row, p.col - 1})
	}
	if p.col < h.width-1 {
		res = append(res, pos{p.row, p.col + 1})
	}
	return res
}

func (p HillClimbingAlgorithm) parseHeightmap(input *Input) (*heightmap, error) {
	lines := input.Lines()
	height := len(lines)
	if height == 0 {
		return nil, errors.New("empty input")
	}

	width := len(lines[0])
	if width == 0 {
		return nil, errors.New("empty input")
	}

	hm := &heightmap{
		height: height,
		width:  width,
		grid:   make([][]rune, height),
	}

	for r, line := range lines {
		if len(line) != width {
			return nil, errors.New("all lines must have the same length")
		}

		hm.grid[r] = make([]rune, width)

		for c, square := range line {
			if square == 'S' { // start
				hm.start = pos{r, c}
				hm.grid[r][c] = 'a'
				continue
			}

			if square == 'E' { // destination
				hm.dest = pos{r, c}
				hm.grid[r][c] = 'z'
				continue
			}

			hm.grid[r][c] = square
		}
	}

	return hm, nil
}

type square struct {
	pos  pos
	dist int
}

type priorityQueue []*square

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq priorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *priorityQueue) Push(x any) {
	item := x.(*square)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}
