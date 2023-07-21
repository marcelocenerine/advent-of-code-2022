// Solution to https://adventofcode.com/2022/day/15
package adventofcode

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
)

type BeaconExclusionZone struct{}

func (p BeaconExclusionZone) Details() Details {
	return Details{Day: 15, Description: "Beacon Exclusion Zone"}
}

func (p BeaconExclusionZone) Solve(input *Input) (Result, error) {
	sensors, err := p.parse(input)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Part1: strconv.Itoa(p.countCellsFreeOfBeacons(sensors, 2_000_000)),
		Part2: "TODO",
	}, nil
}

func (p BeaconExclusionZone) countCellsFreeOfBeacons(sensors []sensor, row int) int {
	segments := []rowSegment{}
	nonEmptyCells := map[int]bool{}

	for _, s := range sensors {
		rowDist := s.distanceToRow(row)
		beaconDist := s.distanceToBeacon()

		if rowDist <= beaconDist {
			diff := beaconDist - rowDist
			covered := rowSegment{
				start: s.pos.x - diff,
				end:   s.pos.x + diff,
			}
			segments = append(segments, covered)
		}
		if s.pos.y == row {
			nonEmptyCells[s.pos.x] = true
		}
		if s.closestBeacon.y == row {
			nonEmptyCells[s.closestBeacon.x] = true
		}
	}

	result := 0
	for _, seg := range p.mergeOverlappingSegments(segments) {
		for col := seg.start; col <= seg.end; col++ {
			if _, ok := nonEmptyCells[col]; !ok {
				result++
			}
		}
	}
	return result
}

func (p BeaconExclusionZone) mergeOverlappingSegments(segments []rowSegment) []rowSegment {
	result := []rowSegment{}
	sort.Slice(segments, func(i, j int) bool {
		return segments[i].start < segments[j].start
	})

	for _, curr := range segments {
		var prev *rowSegment
		if len(result) > 0 {
			prev = &result[len(result)-1]
		}
		if prev == nil || prev.end < curr.start {
			result = append(result, curr)
			continue
		}
		if curr.end > prev.end {
			prev.end = curr.end
		}
	}

	return result
}

type position struct {
	x, y int
}

// https://en.wikipedia.org/wiki/Taxicab_geometry
func (p position) distance(o position) int {
	return int(math.Abs(float64(p.x-o.x)) + math.Abs(float64(p.y-o.y)))
}

type sensor struct {
	pos, closestBeacon position
}

func (s sensor) distanceToBeacon() int {
	return s.pos.distance(s.closestBeacon)
}

func (s sensor) distanceToRow(y int) int {
	return s.pos.distance(position{x: s.pos.x, y: y})
}

type rowSegment struct {
	start, end int // both are inclusive
}

var sensorRgx = regexp.MustCompile(`^Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)$`)

func (p BeaconExclusionZone) parse(input *Input) ([]sensor, error) {
	lines := input.Lines()
	result := make([]sensor, len(lines))

	for i, line := range lines {
		if !sensorRgx.MatchString(line) {
			return nil, fmt.Errorf("invalid line %d: %s", i, line)
		}
		groups := sensorRgx.FindAllStringSubmatch(line, -1)
		sx, _ := strconv.Atoi(groups[0][1])
		sy, _ := strconv.Atoi(groups[0][2])
		bx, _ := strconv.Atoi(groups[0][3])
		by, _ := strconv.Atoi(groups[0][4])
		result[i] = sensor{
			pos:           position{x: sx, y: sy},
			closestBeacon: position{x: bx, y: by},
		}
	}

	return result, nil
}
