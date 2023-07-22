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
		Part1: strconv.Itoa(p.countBeaconFreeCells(sensors, 2_000_000)),
		Part2: strconv.Itoa(p.computeDistressBeaconTuneFrequency(sensors, 4_000_000)),
	}, nil
}

func (p BeaconExclusionZone) countBeaconFreeCells(sensors []sensor, row int) int {
	result := 0
	for _, in := range p.computeRowExclusionIntervals(row, sensors) {
		result += in.length()
	}
	return result
}

// TODO improve running time
func (p BeaconExclusionZone) computeDistressBeaconTuneFrequency(sensors []sensor, xMultiplier int) int {
	searchArea := area{
		topLeft:     position{x: 0, y: 0},
		bottomRight: position{x: xMultiplier, y: xMultiplier},
	}
	if pos, ok := p.findDistressBeacon(sensors, searchArea); ok {
		return (pos.x * xMultiplier) + pos.y
	}
	return -1
}

func (p BeaconExclusionZone) findDistressBeacon(sensors []sensor, searchArea area) (position, bool) {
	for row := searchArea.topLeft.y; row <= searchArea.bottomRight.y; row++ {
		beaconsOnRow := map[int]bool{} // TODO move this outside the loop
		for _, sen := range sensors {
			if sen.closestBeacon.y == row {
				beaconsOnRow[sen.closestBeacon.x] = true
			}
		}

		findInGap := func(lo, hi int) (position, bool) {
			for col := lo; col <= hi; col++ {
				if _, ok := beaconsOnRow[col]; !ok {
					return position{x: col, y: row}, true
				}
			}
			return position{}, false
		}

		exclusions := p.computeRowExclusionIntervals(row, sensors)
		prev := searchArea.topLeft.x - 1

		for _, excl := range exclusions {
			if excl.end < searchArea.topLeft.x {
				continue
			}
			if excl.start > searchArea.bottomRight.x {
				break
			}
			if pos, ok := findInGap(prev+1, excl.start-1); ok {
				return pos, true
			}
			prev = excl.end
		}

		if pos, ok := findInGap(prev+1, searchArea.bottomRight.x); ok {
			return pos, true
		}
	}

	return position{}, false
}

func (p BeaconExclusionZone) computeRowExclusionIntervals(row int, sensors []sensor) []interval {
	intervals := []interval{}
	beaconsOnRow := map[int]bool{}

	for _, sen := range sensors {
		if sen.closestBeacon.y == row {
			beaconsOnRow[sen.closestBeacon.x] = true
		}

		rowDist := sen.distanceToRow(row)
		beaconDist := sen.distanceToBeacon()

		if rowDist <= beaconDist {
			diff := beaconDist - rowDist
			intervals = append(intervals, interval{
				start: sen.pos.x - diff,
				end:   sen.pos.x + diff,
			})
		}
	}

	intervals = p.mergeOverlappingIntervals(intervals)

	// remove beacons found in row from the computed (ordered) intervals
	for beaconX := range beaconsOnRow {
		for i, in := range intervals {
			if beaconX < in.start {
				break // cannot be in this or in any subsequent interval
			}
			if beaconX > in.end {
				continue
			}

			var tmp []interval
			tmp = append(tmp, intervals[:i]...)
			if prefix := (interval{start: in.start, end: beaconX - 1}); prefix.length() > 0 {
				tmp = append(tmp, prefix)
			}
			if suffix := (interval{start: beaconX + 1, end: in.end}); suffix.length() > 0 {
				tmp = append(tmp, suffix)
			}
			tmp = append(tmp, intervals[i+1:]...)
			intervals = tmp
			break
		}
	}

	return intervals
}

func (p BeaconExclusionZone) mergeOverlappingIntervals(intervals []interval) []interval {
	result := []interval{}
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].start < intervals[j].start
	})

	for _, curr := range intervals {
		var prev *interval
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

type area struct {
	topLeft, bottomRight position
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

type interval struct {
	start, end int // both are inclusive
}

func (i interval) length() int {
	return 1 + i.end - i.start
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
