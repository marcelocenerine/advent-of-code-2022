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

	return Result{
		Part1: strconv.Itoa(p.shortestPathFromStartToDest(hm)),
		Part2: strconv.Itoa(p.shortestFromLowestToDest(hm)),
	}, nil
}

func (p HillClimbingAlgorithm) shortestPathFromStartToDest(hm *heightmap) int {
	dist := p.shortestPaths(hm, hm.start, atMostOneHigher)
	return dist[hm.dest.row][hm.dest.col]
}

func (p HillClimbingAlgorithm) shortestFromLowestToDest(hm *heightmap) int {
	dist := p.shortestPaths(hm, hm.dest, atMostOneLower)
	shortest := math.MaxInt

	for r := 0; r < hm.height; r++ {
		for c := 0; c < hm.width; c++ {
			if hm.grid[r][c] != 'a' {
				continue
			}
			if d := dist[r][c]; d < shortest {
				shortest = d
			}
		}
	}

	return shortest
}

func (p HillClimbingAlgorithm) shortestPaths(hm *heightmap, from pos, canMove movePredicate) [][]int {
	visited := make([]bool, hm.size())
	dist := make([][]int, hm.height)
	for r := 0; r < hm.height; r++ {
		dist[r] = make([]int, hm.width)
		for c := 0; c < hm.width; c++ {
			dist[r][c] = math.MaxInt
		}
	}

	dist[from.row][from.col] = 0
	pq := priorityQueue{&square{from, 0}}

	for len(pq) > 0 {
		sq := heap.Pop(&pq).(*square)
		for _, adj := range hm.neighbors(sq.pos) {
			if i := hm.index(adj); visited[i] || !canMove(hm, sq.pos, adj) {
				continue
			}

			oldDist := dist[adj.row][adj.col]
			newDist := sq.dist + 1

			if newDist < oldDist {
				dist[adj.row][adj.col] = newDist
				heap.Push(&pq, &square{adj, newDist})
			}
		}
	}

	return dist
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

type movePredicate func(*heightmap, pos, pos) bool

var atMostOneHigher = func(hm *heightmap, from, to pos) bool {
	return hm.grid[to.row][to.col]-hm.grid[from.row][from.col] <= 1
}

var atMostOneLower = func(hm *heightmap, from, to pos) bool {
	return atMostOneHigher(hm, to, from)
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
